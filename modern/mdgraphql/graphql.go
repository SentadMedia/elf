package mdgraphql

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/graph-gophers/graphql-go"
	"github.com/graph-gophers/graphql-go/relay"
	"github.com/sentadmedia/elf/fw"
	"github.com/sentadmedia/elf/modern/mdhttp"
)

var _ fw.Server = (*GraphGophers)(nil)

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
	relayHandler := NewRelayHandler(g, logger)

	server := mdhttp.NewServer(logger, tracer)
	server.HandleFunc(graphqlPath, middlewareOne(relayHandler))
	server.HandleFunc("/graphiql", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "graphiql.html")
	}))

	return GraphGophers{
		logger: logger,
		server: &server,
	}
}

var _ http.Handler = (*RelayHandler)(nil)

type RelayHandler struct {
	handler relay.Handler
}

func (r RelayHandler) ServeHTTP(writer http.ResponseWriter, reader *http.Request) {
	r.handler.ServeHTTP(writer, reader)
}

func NewRelayHandler(g fw.GraphQLAPI, logger fw.Logger) RelayHandler {
	schema := graphql.MustParseSchema(
		g.GetSchema(),
		g.GetResolver(),
		graphql.UseStringDescriptions(),
		// graphql.Logger(logger.(log.Logger)),
	)
	return RelayHandler{
		handler: relay.Handler{
			Schema: schema,
		},
	}
}

func middlewareOne(next RelayHandler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
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

		if strings.HasPrefix(params.Query, "mutation") && strings.Contains(params.Query, "createAccount") {
			ctx := r.Context()
			response := next.handler.Schema.Exec(ctx, params.Query, params.OperationName, params.Variables)
			responseJSON, err := json.Marshal(response)
			fmt.Print(fmt.Sprintf("Query=%s response=%s err=%s", params.Query, responseJSON, err.Error()))
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			w.Write(responseJSON)
			return
		}

		r.Body = bd2
		next.ServeHTTP(w, r)
	})
}
