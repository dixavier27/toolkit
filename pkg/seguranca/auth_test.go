package seguranca

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func tokenDe(t *testing.T, e *Emissor, sub, papel string) string {
	t.Helper()
	tok, err := e.Emitir(sub, papel)
	if err != nil {
		t.Fatal(err)
	}
	return tok
}

func TestAutenticar(t *testing.T) {
	e := NovoEmissor("segredo", time.Hour)
	protegido := Autenticar(e, func(w http.ResponseWriter, r *http.Request) {
		c, ok := DoContexto(r.Context())
		if !ok || c.Subject != "user-1" {
			t.Errorf("claims não injetadas: ok=%v c=%+v", ok, c)
		}
		w.WriteHeader(http.StatusOK)
	})

	casos := []struct {
		nome     string
		header   string
		querCode int
	}{
		{"sem header", "", http.StatusUnauthorized},
		{"sem Bearer", tokenDe(t, e, "user-1", "user"), http.StatusUnauthorized},
		{"token inválido", "Bearer abc.def.ghi", http.StatusUnauthorized},
		{"válido", "Bearer " + tokenDe(t, e, "user-1", "user"), http.StatusOK},
	}
	for _, c := range casos {
		t.Run(c.nome, func(t *testing.T) {
			rec := httptest.NewRecorder()
			req := httptest.NewRequest("GET", "/", nil)
			if c.header != "" {
				req.Header.Set("Authorization", c.header)
			}
			protegido(rec, req)
			if rec.Code != c.querCode {
				t.Errorf("status = %d, quer %d", rec.Code, c.querCode)
			}
		})
	}
}

func TestExigirPapel(t *testing.T) {
	e := NovoEmissor("segredo", time.Hour)
	final := func(w http.ResponseWriter, r *http.Request) { w.WriteHeader(http.StatusOK) }
	protegido := Autenticar(e, ExigirPapel(PapelAdmin, final))

	t.Run("admin passa", func(t *testing.T) {
		rec := httptest.NewRecorder()
		req := httptest.NewRequest("GET", "/", nil)
		req.Header.Set("Authorization", "Bearer "+tokenDe(t, e, "u", PapelAdmin))
		protegido(rec, req)
		if rec.Code != http.StatusOK {
			t.Errorf("status = %d, quer 200", rec.Code)
		}
	})
	t.Run("user é barrado", func(t *testing.T) {
		rec := httptest.NewRecorder()
		req := httptest.NewRequest("GET", "/", nil)
		req.Header.Set("Authorization", "Bearer "+tokenDe(t, e, "u", "user"))
		protegido(rec, req)
		if rec.Code != http.StatusForbidden {
			t.Errorf("status = %d, quer 403", rec.Code)
		}
	})
}

func TestEhDono(t *testing.T) {
	claimsDe := func(sub, papel string) *Claims {
		c := &Claims{Papel: papel}
		c.Subject = sub
		return c
	}
	ctxCom := func(sub, papel string) context.Context {
		return comClaims(context.Background(), claimsDe(sub, papel))
	}

	dono := ctxCom("user-1", "user")
	admin := ctxCom("outro", PapelAdmin)
	estranho := ctxCom("user-2", "user")

	if !EhDono(dono, "user-1") {
		t.Error("dono deveria ser autorizado")
	}
	if !EhDono(admin, "user-1") {
		t.Error("admin deveria ser autorizado")
	}
	if EhDono(estranho, "user-1") {
		t.Error("estranho não deveria ser autorizado")
	}
	if EhDono(context.Background(), "user-1") {
		t.Error("sem claims não deveria ser autorizado")
	}
}
