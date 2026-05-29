package posts

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/dixavier27/eco/internal/usuarios"
	"github.com/dixavier27/eco/pkg/inmemdb"
	"github.com/dixavier27/eco/pkg/seguranca"
	"github.com/dixavier27/eco/pkg/tupa"
)

// monta cria um servidor com usuarios + posts ligados, e devolve também o
// serviço de usuários (para semear um dono) e o emissor (para mintar tokens).
func monta(t *testing.T) (*tupa.Servidor, *usuarios.Servico, *seguranca.Emissor) {
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

	emissor := seguranca.NovoEmissor("teste", time.Hour)
	srv := tupa.Novo(":0")
	Registrar(srv, svcP, emissor)
	return srv, svcU, emissor
}

// comToken devolve um request com o header Authorization preenchido.
func comToken(t *testing.T, emissor *seguranca.Emissor, metodo, url, corpo, sub, papel string) *http.Request {
	t.Helper()
	tok, err := emissor.Emitir(sub, papel)
	if err != nil {
		t.Fatal(err)
	}
	var r *http.Request
	if corpo != "" {
		r = httptest.NewRequest(metodo, url, strings.NewReader(corpo))
	} else {
		r = httptest.NewRequest(metodo, url, nil)
	}
	r.Header.Set("Authorization", "Bearer "+tok)
	return r
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

// TestHTTPCRUD prova o fluxo aninhado completo com usuarios+inmemdb por trás,
// já com autenticação e autorização por dono.
func TestHTTPCRUD(t *testing.T) {
	srv, svcU, emissor := monta(t)
	h := srv.Handler()
	uid := semearUsuario(t, svcU)

	// POST cria post do usuário (autenticado como o dono)
	body := `{"titulo":"Primeiro","conteudo":"olá mundo"}`
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, comToken(t, emissor, "POST", "/usuarios/"+uid+"/posts", body, uid, usuarios.PapelPadrao))
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

	// GET lista do usuário (público)
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

	// GET individual (público)
	rec = httptest.NewRecorder()
	h.ServeHTTP(rec, httptest.NewRequest("GET", "/posts/"+criado.ID, nil))
	if rec.Code != http.StatusOK {
		t.Errorf("GET id = %d, quer 200", rec.Code)
	}

	// PUT atualiza (dono)
	upd := `{"titulo":"Editado","conteudo":"novo texto"}`
	rec = httptest.NewRecorder()
	h.ServeHTTP(rec, comToken(t, emissor, "PUT", "/posts/"+criado.ID, upd, uid, usuarios.PapelPadrao))
	if rec.Code != http.StatusOK {
		t.Errorf("PUT = %d, quer 200", rec.Code)
	}

	// DELETE (dono)
	rec = httptest.NewRecorder()
	h.ServeHTTP(rec, comToken(t, emissor, "DELETE", "/posts/"+criado.ID, "", uid, usuarios.PapelPadrao))
	if rec.Code != http.StatusNoContent {
		t.Errorf("DELETE = %d, quer 204", rec.Code)
	}
}

// TestHTTPAutorizacao cobre os caminhos de negação: sem token e dono errado.
func TestHTTPAutorizacao(t *testing.T) {
	srv, svcU, emissor := monta(t)
	h := srv.Handler()
	uid := semearUsuario(t, svcU)

	// POST sem token → 401
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, httptest.NewRequest("POST", "/usuarios/"+uid+"/posts", strings.NewReader(`{"titulo":"x","conteudo":"y"}`)))
	if rec.Code != http.StatusUnauthorized {
		t.Errorf("POST sem token = %d, quer 401", rec.Code)
	}

	// POST como outro usuário → 403
	rec = httptest.NewRecorder()
	h.ServeHTTP(rec, comToken(t, emissor, "POST", "/usuarios/"+uid+"/posts", `{"titulo":"x","conteudo":"y"}`, "intruso", usuarios.PapelPadrao))
	if rec.Code != http.StatusForbidden {
		t.Errorf("POST outro usuário = %d, quer 403", rec.Code)
	}

	// cria um post legítimo (como dono) para testar PUT/DELETE de terceiros
	rec = httptest.NewRecorder()
	h.ServeHTTP(rec, comToken(t, emissor, "POST", "/usuarios/"+uid+"/posts", `{"titulo":"meu","conteudo":"texto"}`, uid, usuarios.PapelPadrao))
	var p Post
	_ = json.Unmarshal(rec.Body.Bytes(), &p)

	// PUT por terceiro → 403
	rec = httptest.NewRecorder()
	h.ServeHTTP(rec, comToken(t, emissor, "PUT", "/posts/"+p.ID, `{"titulo":"hack","conteudo":"x"}`, "intruso", usuarios.PapelPadrao))
	if rec.Code != http.StatusForbidden {
		t.Errorf("PUT terceiro = %d, quer 403", rec.Code)
	}

	// DELETE por admin → 204 (admin ignora dono)
	rec = httptest.NewRecorder()
	h.ServeHTTP(rec, comToken(t, emissor, "DELETE", "/posts/"+p.ID, "", "admin-id", seguranca.PapelAdmin))
	if rec.Code != http.StatusNoContent {
		t.Errorf("DELETE admin = %d, quer 204", rec.Code)
	}
}

func TestHTTPErros(t *testing.T) {
	srv, svcU, emissor := monta(t)
	h := srv.Handler()
	uid := semearUsuario(t, svcU)

	// post para usuário inexistente (autenticado como admin p/ passar a autz) → 404
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, comToken(t, emissor, "POST", "/usuarios/999/posts", `{"titulo":"x","conteudo":"y"}`, "admin-id", seguranca.PapelAdmin))
	if rec.Code != http.StatusNotFound {
		t.Errorf("POST dono inexistente = %d, quer 404", rec.Code)
	}

	// entrada inválida (autenticado como dono) → 400
	rec = httptest.NewRecorder()
	h.ServeHTTP(rec, comToken(t, emissor, "POST", "/usuarios/"+uid+"/posts", `{"titulo":"","conteudo":""}`, uid, usuarios.PapelPadrao))
	if rec.Code != http.StatusBadRequest {
		t.Errorf("POST inválido = %d, quer 400", rec.Code)
	}

	// post inexistente → 404 (rota pública de leitura)
	rec = httptest.NewRecorder()
	h.ServeHTTP(rec, httptest.NewRequest("GET", "/posts/999", nil))
	if rec.Code != http.StatusNotFound {
		t.Errorf("GET post inexistente = %d, quer 404", rec.Code)
	}
}
