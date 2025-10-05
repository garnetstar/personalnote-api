package middleware

import (
	"net/http"
	"os"
	"strings"
)

var (
	allowedOrigins        = parseAllowedOrigins()
	allowAllOrigins       = len(allowedOrigins) == 1 && allowedOrigins[0] == "*"
	defaultAllowedMethods = strings.Join([]string{
		http.MethodGet,
		http.MethodPost,
		http.MethodPut,
		http.MethodPatch,
		http.MethodDelete,
		http.MethodOptions,
	}, ", ")
	defaultAllowedHeaders = "Content-Type, Authorization, X-Requested-With"
	defaultExposedHeaders = "Content-Length, Content-Disposition"
)

// WithCORS adds the standard CORS headers and handles preflight requests
func WithCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("Origin")

		if allowAllOrigins {
			if origin != "" {
				w.Header().Set("Access-Control-Allow-Origin", origin)
				w.Header().Add("Vary", "Origin")
			} else {
				w.Header().Set("Access-Control-Allow-Origin", "*")
			}
		} else if origin != "" && isOriginAllowed(origin) {
			w.Header().Set("Access-Control-Allow-Origin", origin)
			w.Header().Add("Vary", "Origin")
		} else if origin != "" {
			http.Error(w, "CORS origin not allowed", http.StatusForbidden)
			return
		}

		w.Header().Add("Vary", "Access-Control-Request-Method")
		w.Header().Add("Vary", "Access-Control-Request-Headers")

		w.Header().Set("Access-Control-Allow-Methods", defaultAllowedMethods)

		if requestHeaders := r.Header.Get("Access-Control-Request-Headers"); requestHeaders != "" {
			w.Header().Set("Access-Control-Allow-Headers", requestHeaders)
		} else {
			w.Header().Set("Access-Control-Allow-Headers", defaultAllowedHeaders)
		}

		w.Header().Set("Access-Control-Expose-Headers", defaultExposedHeaders)
		w.Header().Set("Access-Control-Max-Age", "600")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func parseAllowedOrigins() []string {
	raw := strings.TrimSpace(os.Getenv("CORS_ALLOWED_ORIGINS"))
	if raw == "" {
		return []string{"*"}
	}

	parts := strings.Split(raw, ",")
	var origins []string
	for _, part := range parts {
		origin := strings.TrimSpace(part)
		if origin != "" {
			origins = append(origins, origin)
		}
	}

	if len(origins) == 0 {
		return []string{"*"}
	}

	return origins
}

func isOriginAllowed(origin string) bool {
	if origin == "" {
		return false
	}

	for _, allowed := range allowedOrigins {
		if allowed == "*" || strings.EqualFold(allowed, origin) {
			return true
		}
	}

	return false
}
