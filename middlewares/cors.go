package middlewares

import (
	"net/http"

	"github.com/MilindGour/jellyfin-media-renamer/util"
)

func CorsMW(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		if !util.IsProduction() {
			// Set CORS headers for every request
			w.Header().Set("Access-Control-Allow-Origin", "http://localhost:5173")
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
			w.Header().Set("Access-Control-Allow-Credentials", "true") // if needed

			// Handle preflight
			if r.Method == http.MethodOptions {
				w.WriteHeader(http.StatusOK)
				return
			}
		}

		next.ServeHTTP(w, r)
	})
}
