package mdgraphql

import (
	"net/http"

	"github.com/sentadmedia/elf/fw"
	"github.com/sentadmedia/elf/modern/mdhttp"
)

var _ fw.Server = (*GraphGophers)(nil)

var (
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

func NewGraphGophers(config fw.GraphGopherConfig, handler http.Handler, logger fw.Logger, tracer fw.Tracer) fw.Server {
	// schema := graphql.MustParseSchema(
	// 	config.G.GetSchema(),
	// 	config.G.GetResolver(),
	// 	graphql.UseStringDescriptions(),
	// )

	// rootHandler := handlers.NewRelayHandler(schema)
	// sessionStore := sessions.NewCookieStore(authKey, encryptionKey)
	// authMiddleWare := handlers.NewAuthMiddleWare(schema, sessionStore, cookieName, logger)
	// wrapped := handlers.Chain(rootHandler, authMiddleWare)

	server := mdhttp.NewServer(logger, tracer)
	server.HandleFunc(config.GraphqlPath, handler)
	if config.Graphiql.Include {
		server.HandleFunc(config.Graphiql.Path, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { w.Write(config.Graphiql.Page) }))
	}
	return GraphGophers{
		logger: logger,
		server: &server,
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
