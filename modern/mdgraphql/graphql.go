package mdgraphql

import (
	"context"
	"net/http"

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
	relayHandler := NewRelayHandler(g)

	server := mdhttp.NewServer(logger, tracer)
	server.HandleFunc(graphqlPath, testMiddleWare(relayHandler))

	return GraphGophers{
		logger: logger,
		server: &server,
	}
}

func testMiddleWare(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		next.ServeHTTP(w, r.WithContext(context.WithValue(ctx, "Soufiane", "KAWTAR")))
	})
}

var _ http.Handler = (*RelayHandler)(nil)

type RelayHandler struct {
	handler relay.Handler
}

func (r RelayHandler) ServeHTTP(writer http.ResponseWriter, reader *http.Request) {
	r.handler.ServeHTTP(writer, reader)
}

func NewRelayHandler(g fw.GraphQLAPI) RelayHandler {
	schema := graphql.MustParseSchema(
		g.GetSchema(),
		g.GetResolver(),
		graphql.UseStringDescriptions(),
	)
	return RelayHandler{handler: relay.Handler{
		Schema: schema,
	}}
}
