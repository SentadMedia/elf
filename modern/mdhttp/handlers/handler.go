package handlers

import (
	"net/http"

	"github.com/sentadmedia/elf/fw"
)

func Chain(rootHandler http.Handler, m ...fw.Middleware) http.Handler {
	if len(m) < 1 {
		return rootHandler
	}
	wrapped := rootHandler
	// loop in reverse to preserve middleware order
	for i := len(m) - 1; i >= 0; i-- {
		wrapped = m[i](wrapped)
	}
	return wrapped

}
