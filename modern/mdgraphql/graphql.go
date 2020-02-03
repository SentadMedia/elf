package mdgraphql

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/gorilla/sessions"
	"github.com/graph-gophers/graphql-go"
	"github.com/graph-gophers/graphql-go/relay"
	"github.com/sentadmedia/elf/fw"
	"github.com/sentadmedia/elf/modern/mdhttp"
)

var _ fw.Server = (*GraphGophers)(nil)

var (
	// key must be 16, 24 or 32 bytes long (AES-128, AES-192 or AES-256)
	authKey       = []byte("CBD8DC165293A3F32084DACFF2B0D1522095AD3741919D82")
	encryptionKey = []byte("A7A2F5A4707904314079C61EF6E49CEE")
	cookieName    = "SANTA_SESSION"
)

type GraphGophers struct {
	logger fw.Logger
	server fw.Server
}

func (g GraphGophers) Shutdown() error {
	return g.server.Shutdown()
}

func (g GraphGophers) ListenAndServe(port int) error {
	return g.server.ListenAndServe(port)
}

func NewGraphGophers(graphqlPath string, logger fw.Logger, tracer fw.Tracer, g fw.GraphQLAPI) fw.Server {
	rootHandler := NewRelayHandler(g)
	sessionStore := sessions.NewCookieStore(authKey, encryptionKey)
	authMiddleWare := NewAuthMiddleWare(sessionStore)
	testMiddleWare := NewMiddleWareTest()
	logMiddleWare := NewMiddleWareLog()
	wrapped := MultipleMiddleware(rootHandler, testMiddleWare, authMiddleWare, logMiddleWare)

	server := mdhttp.NewServer(logger, tracer)
	server.HandleFunc(graphqlPath, wrapped)
	server.HandleFunc("/graphql", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { w.Write(page) }))

	return GraphGophers{
		logger: logger,
		server: &server,
	}
}

var _ http.Handler = (*relay.Handler)(nil)

// type RelayHandler struct {
// 	handler relay.Handler
// }

// func (r RelayHandler) ServeHTTP(writer http.ResponseWriter, reader *http.Request) {
// 	r.handler.ServeHTTP(writer, reader)
// }

// func NewRelayHandler(g fw.GraphQLAPI) RelayHandler {
// 	schema := graphql.MustParseSchema(
// 		g.GetSchema(),
// 		g.GetResolver(),
// 		graphql.UseStringDescriptions(),
// 	)
// 	return RelayHandler{
// 		handler: relay.Handler{
// 			Schema: schema,
// 		},
// 	}
// }

func NewRelayHandler(g fw.GraphQLAPI) http.Handler {
	schema := graphql.MustParseSchema(
		g.GetSchema(),
		g.GetResolver(),
		graphql.UseStringDescriptions(),
	)
	return &relay.Handler{schema}
}

func NewMiddleWareTest() Middleware {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			fmt.Print("MiddleWare Test Start\n")
			defer fmt.Print("MiddleWare Test End\n")
			h.ServeHTTP(w, r)
		})
	}
}

func NewMiddleWareLog() Middleware {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			fmt.Print("MiddleWare Log Start\n")
			defer fmt.Print("MiddleWare Log End\n")
			h.ServeHTTP(w, r)
		})
	}
}

func NewAuthMiddleWare(store sessions.Store) Middleware {
	return func(handler http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			relayHandler, ok := handler.(*relay.Handler)
			if !ok {
				fmt.Printf("Unknown http handler. expected 'relay.Handler', got %T\n", handler)
				// http.Error(w, err.Error(), http.StatusBadRequest)
				handler.ServeHTTP(w, r)
				return
			}

			session, _ := store.Get(r, cookieName)
			fmt.Print(fmt.Sprintf("Session Before=%+v\n", session))
			var params struct {
				Query         string                 `json:"query"`
				OperationName string                 `json:"operationName"`
				Variables     map[string]interface{} `json:"variables"`
			}

			buf, _ := ioutil.ReadAll(r.Body)
			bd1 := ioutil.NopCloser(bytes.NewBuffer(buf))
			bd2 := ioutil.NopCloser(bytes.NewBuffer(buf))

			if err := json.NewDecoder(bd1).Decode(&params); err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}

			if strings.HasPrefix(params.Query, "mutation") && strings.Contains(params.Query, "signIn") {
				ctx := r.Context()
				session.Values["authenticated"] = true
				if err := session.Save(r, w); err != nil {
					fmt.Print(fmt.Sprintf("ERROR Saving session. err=%s\n", err.Error()))
				} else {
					fmt.Print(fmt.Sprintf("Session After=%+v\n", session))
				}
				response := relayHandler.Schema.Exec(ctx, params.Query, params.OperationName, params.Variables)
				responseJSON, err := json.Marshal(response)
				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}
				w.Write(responseJSON)
				return
			}

			r.Body = bd2
			relayHandler.ServeHTTP(w, r)
		})
	}

}

var page = []byte(`
<!DOCTYPE html>
<html>
	<head>
		<link href="https://cdnjs.cloudflare.com/ajax/libs/graphiql/0.11.11/graphiql.min.css" rel="stylesheet" />
		<script src="https://cdnjs.cloudflare.com/ajax/libs/es6-promise/4.1.1/es6-promise.auto.min.js"></script>
		<script src="https://cdnjs.cloudflare.com/ajax/libs/fetch/2.0.3/fetch.min.js"></script>
		<script src="https://cdnjs.cloudflare.com/ajax/libs/react/16.2.0/umd/react.production.min.js"></script>
		<script src="https://cdnjs.cloudflare.com/ajax/libs/react-dom/16.2.0/umd/react-dom.production.min.js"></script>
		<script src="https://cdnjs.cloudflare.com/ajax/libs/graphiql/0.11.11/graphiql.min.js"></script>
	</head>
	<body style="width: 100%; height: 100%; margin: 0; overflow: hidden;">
		<div id="graphiql" style="height: 100vh;">Loading...</div>
		<script>
			function graphQLFetcher(graphQLParams) {
				return fetch("/query", {
					method: "post",
					body: JSON.stringify(graphQLParams),
					credentials: "include",
				}).then(function (response) {
					return response.text();
				}).then(function (responseBody) {
					try {
						return JSON.parse(responseBody);
					} catch (error) {
						return responseBody;
					}
				});
			}
			ReactDOM.render(
				React.createElement(GraphiQL, {fetcher: graphQLFetcher}),
				document.getElementById("graphiql")
			);
		</script>
	</body>
</html>
`)

type Middleware func(h http.Handler) http.Handler

func MultipleMiddleware(h http.Handler, m ...Middleware) http.Handler {

	if len(m) < 1 {
		return h
	}

	wrapped := h

	// loop in reverse to preserve middleware order
	for i := len(m) - 1; i >= 0; i-- {
		wrapped = m[i](wrapped)
	}

	return wrapped

}
