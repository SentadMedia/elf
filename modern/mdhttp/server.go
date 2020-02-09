package mdhttp

import (
	"context"
	"fmt"
	"net"
	"net/http"

	"github.com/sentadmedia/elf/fw"
)

type Server struct {
	mux    *http.ServeMux
	server *http.Server
	tracer fw.Tracer
	logger fw.Logger
}

func (s *Server) ListenAndServe(port int) error {
	addr := fmt.Sprintf(":%d", port)

	s.server = &http.Server{Addr: addr, Handler: s.mux}
	err := s.server.ListenAndServe()

	if err == nil || err == http.ErrServerClosed {
		return nil
	}

	return err
}

func (s Server) Shutdown() error {
	return s.server.Shutdown(context.Background())
}

func (s Server) HandleFunc(pattern string, handler http.Handler) {
	s.mux.HandleFunc(pattern, func(w http.ResponseWriter, r *http.Request) {
		w = setupPreFlight(w)
		if (*r).Method == "OPTIONS" {
			return
		}

		// w = enableCors(w)
		handler.ServeHTTP(w, r)
	})
}

func setupPreFlight(w http.ResponseWriter) http.ResponseWriter {
	w.Header().Set("Access-Control-Allow-Credentials", "true")
	w.Header().Set("Access-Control-Allow-Origin", "localhost:3000")
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	w.Header().Set("Access-Control-Allow-Headers",
		"Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
	return w
}

func enableCors(w http.ResponseWriter) http.ResponseWriter {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	return w
}

// GetClientIP gets a requests IP address by reading off the forwarded-for
// header (for proxies) and falls back to use the remote address.
func GetClientIP(r *http.Request) string {
	forwarded := r.Header.Get("X-FORWARDED-FOR")
	if forwarded != "" {
		return forwarded
	}
	if host, _, err := net.SplitHostPort(r.RemoteAddr); err != nil {
		return host
	}

	return r.RemoteAddr
}

func NewServer(logger fw.Logger, tracer fw.Tracer) Server {
	mux := http.NewServeMux()
	return Server{
		mux:    mux,
		tracer: tracer,
		logger: logger,
	}
}
