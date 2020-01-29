package mdgraphql

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/graph-gophers/graphql-go"
	"github.com/graph-gophers/graphql-go/log"
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
		graphql.Logger(logger.(log.Logger)),
	)
	return RelayHandler{
		handler: relay.Handler{
			Schema: schema,
		},
	}
}

func middlewareOne(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var params struct {
			Query         string                 `json:"query"`
			OperationName string                 `json:"operationName"`
			Variables     map[string]interface{} `json:"variables"`
		}
		if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		fmt.Print(fmt.Sprintf("params=%+v", params))
		next.ServeHTTP(w, r)
	})
}
