// Package usrserv provides the router/server handlers for the user service.
package usrserv

import (
	"net/http"
	"time"

	"github.com/anchorage-store/anchorage/serv"
)

// Router creates an HTTP router for the user service routes.
func Router() *http.ServeMux {
	m := http.NewServeMux()

	m.HandleFunc("GET /healthz", healthCheck(time.Now()))

	return m
}

func healthCheck(startTime time.Time) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		serv.WriteJSON(w, struct {
			Uptime uint64 `json:"uptime"`
		}{
			Uptime: uint64(time.Since(startTime).Seconds()),
		})
	}
}
