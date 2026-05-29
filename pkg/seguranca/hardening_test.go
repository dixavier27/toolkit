package seguranca

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func okHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
}

func TestCabecalhosSeguranca(t *testing.T) {
	h := CabecalhosSeguranca()(okHandler())
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, httptest.NewRequest("GET", "/", nil))

	quer := map[string]string{
		"X-Content-Type-Options":  "nosniff",
		"X-Frame-Options":         "DENY",
		"Referrer-Policy":         "no-referrer",
		"Content-Security-Policy": "default-src 'none'",
	}
	for k, v := range quer {
		if got := rec.Header().Get(k); got != v {
			t.Errorf("%s = %q, quer %q", k, got, v)
		}
	}
}

func TestCORSPreflight(t *testing.T) {
	h := CORS()(okHandler())
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, httptest.NewRequest("OPTIONS", "/", nil))

	if rec.Code != http.StatusNoContent {
		t.Errorf("preflight status = %d, quer 204", rec.Code)
	}
	if got := rec.Header().Get("Access-Control-Allow-Origin"); got != "*" {
		t.Errorf("Allow-Origin = %q, quer *", got)
	}
}

func TestCORSOrigemRestrita(t *testing.T) {
	h := CORS("https://app.exemplo.com")(okHandler())

	// origem permitida → ecoada
	rec := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/", nil)
	req.Header.Set("Origin", "https://app.exemplo.com")
	h.ServeHTTP(rec, req)
	if got := rec.Header().Get("Access-Control-Allow-Origin"); got != "https://app.exemplo.com" {
		t.Errorf("Allow-Origin = %q, quer a origem permitida", got)
	}

	// origem não listada → não ecoada
	rec = httptest.NewRecorder()
	req = httptest.NewRequest("GET", "/", nil)
	req.Header.Set("Origin", "https://intruso.com")
	h.ServeHTTP(rec, req)
	if got := rec.Header().Get("Access-Control-Allow-Origin"); got != "" {
		t.Errorf("Allow-Origin = %q, quer vazio para origem não listada", got)
	}
}

func TestLimitarTaxa(t *testing.T) {
	// 1 token/s com rajada de 2: as 2 primeiras passam, a 3ª é barrada.
	h := LimitarTaxa(1, 2)(okHandler())
	req := httptest.NewRequest("GET", "/", nil)
	req.RemoteAddr = "203.0.113.7:5555"

	for i := 1; i <= 2; i++ {
		rec := httptest.NewRecorder()
		h.ServeHTTP(rec, req)
		if rec.Code != http.StatusOK {
			t.Fatalf("requisição %d status = %d, quer 200", i, rec.Code)
		}
	}
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)
	if rec.Code != http.StatusTooManyRequests {
		t.Errorf("3ª requisição status = %d, quer 429", rec.Code)
	}
}

func TestLimitarTaxaPorIP(t *testing.T) {
	// Estourar um IP não deve afetar outro.
	h := LimitarTaxa(1, 1)(okHandler())

	reqA := httptest.NewRequest("GET", "/", nil)
	reqA.RemoteAddr = "203.0.113.1:1111"
	reqB := httptest.NewRequest("GET", "/", nil)
	reqB.RemoteAddr = "203.0.113.2:2222"

	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, reqA) // consome o token de A
	rec = httptest.NewRecorder()
	h.ServeHTTP(rec, reqA) // A barrado
	if rec.Code != http.StatusTooManyRequests {
		t.Fatalf("A 2ª status = %d, quer 429", rec.Code)
	}
	rec = httptest.NewRecorder()
	h.ServeHTTP(rec, reqB) // B ainda tem token
	if rec.Code != http.StatusOK {
		t.Errorf("B status = %d, quer 200 (IPs independentes)", rec.Code)
	}
}
