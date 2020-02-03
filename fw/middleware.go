package fw

import "net/http"

type Middleware func(h http.Handler) http.Handler
