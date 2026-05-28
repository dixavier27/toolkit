package tupa

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func TestRotaComPathParam(t *testing.T) {
	s := Novo(":0")
	s.Rota("GET", "/itens/{id}", func(w http.ResponseWriter, r *http.Request) {
		_ = EscreverJSON(w, http.StatusOK, map[string]string{"id": r.PathValue("id")})
	})

	req := httptest.NewRequest("GET", "/itens/42", nil)
	rec := httptest.NewRecorder()
	s.Handler().ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, quer 200", rec.Code)
	}
	var got map[string]string
	if err := json.Unmarshal(rec.Body.Bytes(), &got); err != nil {
		t.Fatal(err)
	}
	if got["id"] != "42" {
		t.Errorf("id = %q, quer 42", got["id"])
	}
}

func TestMetodoErrado(t *testing.T) {
	s := Novo(":0")
	s.Rota("GET", "/x", func(w http.ResponseWriter, r *http.Request) {})

	rec := httptest.NewRecorder()
	s.Handler().ServeHTTP(rec, httptest.NewRequest("POST", "/x", nil))
	if rec.Code != http.StatusMethodNotAllowed {
		t.Errorf("status = %d, quer 405", rec.Code)
	}
}

func TestMiddlewareOrdem(t *testing.T) {
	s := Novo(":0")
	var ordem []string
	marca := func(nome string) Middleware {
		return func(prox http.Handler) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				ordem = append(ordem, nome)
				prox.ServeHTTP(w, r)
			})
		}
	}
	s.Usar(marca("a"), marca("b"))
	s.Rota("GET", "/", func(w http.ResponseWriter, r *http.Request) {})

	s.Handler().ServeHTTP(httptest.NewRecorder(), httptest.NewRequest("GET", "/", nil))
	if strings.Join(ordem, ",") != "a,b" {
		t.Errorf("ordem = %v, quer [a b] (primeiro registrado mais externo)", ordem)
	}
}

func TestLerJSONRejeitaCampoDesconhecido(t *testing.T) {
	req := httptest.NewRequest("POST", "/", strings.NewReader(`{"x":1}`))
	var dst struct {
		Y int `json:"y"`
	}
	if err := LerJSON(req, &dst); err == nil {
		t.Error("esperava erro por campo desconhecido")
	}
}

func TestIniciarShutdownGracioso(t *testing.T) {
	s := Novo("127.0.0.1:0", ComTimeoutDeParada(time.Second))
	s.Rota("GET", "/", func(w http.ResponseWriter, r *http.Request) {})

	ctx, cancel := context.WithCancel(context.Background())
	done := make(chan error, 1)
	go func() { done <- s.Iniciar(ctx) }()

	time.Sleep(50 * time.Millisecond)
	cancel()

	select {
	case err := <-done:
		if err != nil {
			t.Errorf("Iniciar devolveu erro: %v", err)
		}
	case <-time.After(3 * time.Second):
		t.Fatal("Iniciar não retornou após cancelamento")
	}
}
