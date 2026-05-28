package tupa

import (
	"log"
	"net/http"
	"time"
)

// Middleware envolve um http.Handler, podendo executar antes/depois dele.
type Middleware func(http.Handler) http.Handler

// LogRequisicoes registra método, caminho e duração de cada requisição no
// logger fornecido. Útil como exemplo e ponto de partida.
func LogRequisicoes(l *log.Logger) Middleware {
	return func(prox http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			inicio := time.Now()
			prox.ServeHTTP(w, r)
			l.Printf("%s %s (%s)", r.Method, r.URL.Path, time.Since(inicio))
		})
	}
}
