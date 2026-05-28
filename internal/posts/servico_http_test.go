package posts

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/dixavier27/eco/internal/usuarios"
	"github.com/dixavier27/eco/pkg/inmemdb"
	"github.com/dixavier27/eco/pkg/tupa"
)

// monta cria um servidor com usuarios + posts ligados, e devolve também o
// serviço de usuários para semear um dono nos testes.
func monta(t *testing.T) (*tupa.Servidor, *usuarios.Servico) {
	t.Helper()

	repoU := inmemdb.NovaMemoria(
		func(u usuarios.Usuario) string { return u.ID },
		inmemdb.ComDefinirID(func(u *usuarios.Usuario, id string) { u.ID = id }),
	)
	svcU := usuarios.NovoServico(repoU, usuarios.BcryptHasher{Custo: 4})

	repoP := inmemdb.NovaMemoria(
		func(p Post) string { return p.ID },
		inmemdb.ComDefinirID(func(p *Post, id string) { p.ID = id }),
	)
	svcP := NovoServico(repoP, svcU)

	srv := tupa.Novo(":0")
	Registrar(srv, svcP)
	return srv, svcU
}

func semearUsuario(t *testing.T, svc *usuarios.Servico) string {
	t.Helper()
	u, err := svc.Criar(context.Background(), usuarios.EntradaCriar{
		Nome: "Ana", Sobrenome: "Lima", Email: "ana@ex.com",
		Whatsapp: "+5511999990000", Senha: "segredo123",
	})
	if err != nil {
		t.Fatalf("semear usuário: %v", err)
	}
	return u.ID
}

// TestHTTPCRUD prova o fluxo aninhado completo com usuarios+inmemdb por trás.
func TestHTTPCRUD(t *testing.T) {
	srv, svcU := monta(t)
	h := srv.Handler()
	uid := semearUsuario(t, svcU)

	// POST cria post do usuário
	body := `{"titulo":"Primeiro","conteudo":"olá mundo"}`
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, httptest.NewRequest("POST", "/usuarios/"+uid+"/posts", strings.NewReader(body)))
	if rec.Code != http.StatusCreated {
		t.Fatalf("POST status = %d, quer 201 (corpo: %s)", rec.Code, rec.Body)
	}
	var criado Post
	if err := json.Unmarshal(rec.Body.Bytes(), &criado); err != nil {
		t.Fatal(err)
	}
	if criado.UsuarioID != uid {
		t.Errorf("usuarioId = %q, quer %q", criado.UsuarioID, uid)
	}

	// GET lista do usuário
	rec = httptest.NewRecorder()
	h.ServeHTTP(rec, httptest.NewRequest("GET", "/usuarios/"+uid+"/posts", nil))
	if rec.Code != http.StatusOK {
		t.Errorf("GET lista = %d, quer 200", rec.Code)
	}
	var lista []Post
	_ = json.Unmarshal(rec.Body.Bytes(), &lista)
	if len(lista) != 1 {
		t.Errorf("len(lista) = %d, quer 1", len(lista))
	}

	// GET individual
	rec = httptest.NewRecorder()
	h.ServeHTTP(rec, httptest.NewRequest("GET", "/posts/"+criado.ID, nil))
	if rec.Code != http.StatusOK {
		t.Errorf("GET id = %d, quer 200", rec.Code)
	}

	// PUT atualiza
	upd := `{"titulo":"Editado","conteudo":"novo texto"}`
	rec = httptest.NewRecorder()
	h.ServeHTTP(rec, httptest.NewRequest("PUT", "/posts/"+criado.ID, strings.NewReader(upd)))
	if rec.Code != http.StatusOK {
		t.Errorf("PUT = %d, quer 200", rec.Code)
	}

	// DELETE
	rec = httptest.NewRecorder()
	h.ServeHTTP(rec, httptest.NewRequest("DELETE", "/posts/"+criado.ID, nil))
	if rec.Code != http.StatusNoContent {
		t.Errorf("DELETE = %d, quer 204", rec.Code)
	}
}

func TestHTTPErros(t *testing.T) {
	srv, svcU := monta(t)
	h := srv.Handler()
	uid := semearUsuario(t, svcU)

	// post para usuário inexistente → 404
	rec := httptest.NewRecorder()
	body := `{"titulo":"x","conteudo":"y"}`
	h.ServeHTTP(rec, httptest.NewRequest("POST", "/usuarios/999/posts", strings.NewReader(body)))
	if rec.Code != http.StatusNotFound {
		t.Errorf("POST dono inexistente = %d, quer 404", rec.Code)
	}

	// entrada inválida → 400
	rec = httptest.NewRecorder()
	h.ServeHTTP(rec, httptest.NewRequest("POST", "/usuarios/"+uid+"/posts", strings.NewReader(`{"titulo":"","conteudo":""}`)))
	if rec.Code != http.StatusBadRequest {
		t.Errorf("POST inválido = %d, quer 400", rec.Code)
	}

	// post inexistente → 404
	rec = httptest.NewRecorder()
	h.ServeHTTP(rec, httptest.NewRequest("GET", "/posts/999", nil))
	if rec.Code != http.StatusNotFound {
		t.Errorf("GET post inexistente = %d, quer 404", rec.Code)
	}
}
