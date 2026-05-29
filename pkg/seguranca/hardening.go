package seguranca

import (
	"net"
	"net/http"
	"sync"

	"golang.org/x/time/rate"

	"github.com/dixavier27/eco/pkg/tupa"
)

// CabecalhosSeguranca devolve um middleware que adiciona cabeçalhos de segurança
// padrão a toda resposta: bloqueia sniffing de MIME, framing (clickjacking),
// vazamento de referrer e restringe a origem de recursos via CSP.
func CabecalhosSeguranca() tupa.Middleware {
	return func(prox http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			h := w.Header()
			h.Set("X-Content-Type-Options", "nosniff")
			h.Set("X-Frame-Options", "DENY")
			h.Set("Referrer-Policy", "no-referrer")
			h.Set("Content-Security-Policy", "default-src 'none'")
			prox.ServeHTTP(w, r)
		})
	}
}

// CORS devolve um middleware que responde a preflight (OPTIONS) e libera as
// origens informadas. Sem origens, libera todas ("*"). Com origens, só ecoa o
// Access-Control-Allow-Origin se a origem da requisição estiver na lista.
func CORS(origens ...string) tupa.Middleware {
	permitidas := make(map[string]bool, len(origens))
	for _, o := range origens {
		permitidas[o] = true
	}
	return func(prox http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			h := w.Header()
			if len(permitidas) == 0 {
				h.Set("Access-Control-Allow-Origin", "*")
			} else if origem := r.Header.Get("Origin"); origem != "" && permitidas[origem] {
				h.Set("Access-Control-Allow-Origin", origem)
				h.Add("Vary", "Origin")
			}
			h.Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
			h.Set("Access-Control-Allow-Headers", "Authorization, Content-Type")
			if r.Method == http.MethodOptions {
				w.WriteHeader(http.StatusNoContent)
				return
			}
			prox.ServeHTTP(w, r)
		})
	}
}

// LimitarTaxa devolve um middleware de rate limiting por IP (token bucket):
// rps tokens por segundo com capacidade de rajada burst. Ao estourar, responde
// 429.
//
// Gotcha (PoC): o mapa de limitadores cresce por IP sem expiração. Em produção,
// adicione coleta periódica dos IPs ociosos para evitar crescimento ilimitado.
func LimitarTaxa(rps float64, burst int) tupa.Middleware {
	var mu sync.Mutex
	limitadores := make(map[string]*rate.Limiter)
	limitadorDe := func(ip string) *rate.Limiter {
		mu.Lock()
		defer mu.Unlock()
		l, ok := limitadores[ip]
		if !ok {
			l = rate.NewLimiter(rate.Limit(rps), burst)
			limitadores[ip] = l
		}
		return l
	}
	return func(prox http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if !limitadorDe(ipDe(r)).Allow() {
				tupa.EscreverErro(w, http.StatusTooManyRequests, "muitas requisições")
				return
			}
			prox.ServeHTTP(w, r)
		})
	}
}

// ipDe extrai o IP do cliente a partir de RemoteAddr (host:porta).
func ipDe(r *http.Request) string {
	host, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return r.RemoteAddr
	}
	return host
}
