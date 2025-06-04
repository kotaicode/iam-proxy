package server

import (
	"net"
	"net/http"
	"strings"

	"github.com/kotaicode/iam-proxy/internal/config"
)

func AuthMiddleware(cfg *config.Config, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Skip auth for health check
		if r.URL.Path == "/healthz" {
			next.ServeHTTP(w, r)
			return
		}

		// Check IP if allowed IPs are configured
		if len(cfg.AllowedIPs) > 0 {
			clientIP := net.ParseIP(strings.Split(r.RemoteAddr, ":")[0])
			allowed := false
			for _, ip := range cfg.AllowedIPs {
				if ip.Equal(clientIP) {
					allowed = true
					break
				}
			}
			if !allowed {
				http.Error(w, "Forbidden", http.StatusForbidden)
				return
			}
		}

		// Check token if configured
		if cfg.SecurityToken != "" {
			token := r.Header.Get("Authorization")
			if !strings.HasPrefix(token, "Bearer ") {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}
			token = strings.TrimPrefix(token, "Bearer ")
			if token != cfg.SecurityToken {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}
		}

		next.ServeHTTP(w, r)
	})
}
