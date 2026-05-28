package usuarios

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/dixavier27/eco/pkg/inmemdb"
	"github.com/dixavier27/eco/pkg/tupa"
)

func montar() (*tupa.Servidor, *Servico) {
	repo := inmemdb.NovaMemoria(
		func(u Usuario) string { return u.ID },
		inmemdb.ComDefinirID(func(u *Usuario, id string) { u.ID = id }),
	)
	svc := NovoServico(repo, BcryptHasher{Custo: 4}) // custo baixo: testes rápidos
	srv := tupa.Novo(":0")
	Registrar(srv, svc)
	return srv, svc
}

// TestServicoComMemoria prova o serviço contra inmemdb.Memoria (contexto dados).
func TestServicoComMemoria(t *testing.T) {
	_, svc := montar()
	u, err := svc.Criar(context.Background(), entradaValida())
	if err != nil {
		t.Fatal(err)
	}
	if u.ID == "" {
		t.Error("id não gerado")
	}
	if u.SenhaCriptografada == "segredo123" || u.SenhaCriptografada == "" {
		t.Error("senha não foi hasheada")
	}
}

// TestHTTPCRUD prova os endpoints com inmemdb por trás (contexto HTTP).
func TestHTTPCRUD(t *testing.T) {
	srv, _ := montar()
	h := srv.Handler()

	// POST cria
	body := `{"nome":"Ana","sobrenome":"Lima","email":"ana@ex.com","whatsapp":"+5511999990000","senha":"segredo123"}`
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, httptest.NewRequest("POST", "/usuarios", strings.NewReader(body)))
	if rec.Code != http.StatusCreated {
		t.Fatalf("POST status = %d, quer 201 (corpo: %s)", rec.Code, rec.Body)
	}
	if strings.Contains(rec.Body.String(), "segredo") || strings.Contains(rec.Body.String(), "Criptografada") {
		t.Errorf("senha vazou no JSON: %s", rec.Body)
	}
	var criado Usuario
	if err := json.Unmarshal(rec.Body.Bytes(), &criado); err != nil {
		t.Fatal(err)
	}

	// GET lista
	rec = httptest.NewRecorder()
	h.ServeHTTP(rec, httptest.NewRequest("GET", "/usuarios", nil))
	if rec.Code != http.StatusOK {
		t.Errorf("GET lista status = %d, quer 200", rec.Code)
	}

	// GET busca por id
	rec = httptest.NewRecorder()
	h.ServeHTTP(rec, httptest.NewRequest("GET", "/usuarios/"+criado.ID, nil))
	if rec.Code != http.StatusOK {
		t.Errorf("GET id status = %d, quer 200", rec.Code)
	}

	// DELETE
	rec = httptest.NewRecorder()
	h.ServeHTTP(rec, httptest.NewRequest("DELETE", "/usuarios/"+criado.ID, nil))
	if rec.Code != http.StatusNoContent {
		t.Errorf("DELETE status = %d, quer 204", rec.Code)
	}
}

func TestHTTPErros(t *testing.T) {
	srv, _ := montar()
	h := srv.Handler()

	// id inexistente → 404
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, httptest.NewRequest("GET", "/usuarios/999", nil))
	if rec.Code != http.StatusNotFound {
		t.Errorf("GET inexistente = %d, quer 404", rec.Code)
	}

	// email inválido → 400
	body := `{"nome":"Ana","sobrenome":"Lima","email":"invalido","whatsapp":"+5511999990000","senha":"segredo123"}`
	rec = httptest.NewRecorder()
	h.ServeHTTP(rec, httptest.NewRequest("POST", "/usuarios", strings.NewReader(body)))
	if rec.Code != http.StatusBadRequest {
		t.Errorf("POST email inválido = %d, quer 400", rec.Code)
	}
}
