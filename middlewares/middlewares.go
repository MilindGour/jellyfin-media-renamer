package middlewares

import "net/http"

type Middleware func(http.Handler) http.Handler

func Pipe(mws ...Middleware) Middleware {
	return func(next http.Handler) http.Handler {
		for _, mw := range mws {
			next = mw(next)
		}

		return next
	}
}
