package mdgraphql

import (
	"net/http"

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

func NewGraphGophers(config fw.GraphGopherConfig, handler http.Handler, logger fw.Logger, tracer fw.Tracer) fw.Server {
	server := mdhttp.NewServer(logger, tracer)
	server.HandleFunc(string(config.GraphqlPath), handler)
	if config.Graphiql.Include {
		server.HandleFunc(string(config.Graphiql.Path), http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { w.Write(config.Graphiql.Page) }))
	}
	return GraphGophers{
		logger: logger,
		server: &server,
	}
}
