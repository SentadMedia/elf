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
	server := mdhttp.NewServer(logger, tracer)
	server.HandleFunc(config.GraphqlPath, handler)
	if config.Graphiql.Include {
		server.HandleFunc(string(config.Graphiql.Path), http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { w.Write(config.Graphiql.Page) }))
	}
	return GraphGophers{
		logger: logger,
		server: &server,
	}
}
