package posts

import (
	"errors"
	"net/http"

	"github.com/dixavier27/eco/pkg/repo"
	"github.com/dixavier27/eco/pkg/seguranca"
	"github.com/dixavier27/eco/pkg/tupa"
)

// Registrar registra as rotas REST de posts no servidor tupa. As rotas de
// criação e listagem são aninhadas no usuário dono; as individuais são planas.
// As rotas mutadoras exigem token válido e que o autenticado seja o dono (ou
// admin); as de leitura são públicas.
//
//	POST   /usuarios/{usuarioID}/posts  cria post do usuário (dono/admin)
//	GET    /usuarios/{usuarioID}/posts  lista posts do usuário (público)
//	GET    /posts/{id}                  busca (público)
//	PUT    /posts/{id}                  atualiza (dono/admin)
//	DELETE /posts/{id}                  remove (dono/admin)
func Registrar(s *tupa.Servidor, svc *Servico, emissor *seguranca.Emissor) {
	s.Rota("POST", "/usuarios/{usuarioID}/posts", seguranca.Autenticar(emissor, func(w http.ResponseWriter, r *http.Request) {
		usuarioID := r.PathValue("usuarioID")
		if !seguranca.EhDono(r.Context(), usuarioID) {
			tupa.EscreverErro(w, http.StatusForbidden, "acesso negado")
			return
		}
		var e EntradaCriar
		if err := tupa.LerJSON(r, &e); err != nil {
			tupa.EscreverErro(w, http.StatusBadRequest, "JSON inválido")
			return
		}
		p, err := svc.Criar(r.Context(), usuarioID, e)
		if err != nil {
			responderErro(w, err)
			return
		}
		_ = tupa.EscreverJSON(w, http.StatusCreated, p)
	}))

	s.Rota("GET", "/usuarios/{usuarioID}/posts", func(w http.ResponseWriter, r *http.Request) {
		ps, err := svc.ListarDoUsuario(r.Context(), r.PathValue("usuarioID"))
		if err != nil {
			responderErro(w, err)
			return
		}
		_ = tupa.EscreverJSON(w, http.StatusOK, ps)
	})

	s.Rota("GET", "/posts/{id}", func(w http.ResponseWriter, r *http.Request) {
		p, err := svc.Buscar(r.Context(), r.PathValue("id"))
		if err != nil {
			responderErro(w, err)
			return
		}
		_ = tupa.EscreverJSON(w, http.StatusOK, p)
	})

	s.Rota("PUT", "/posts/{id}", seguranca.Autenticar(emissor, func(w http.ResponseWriter, r *http.Request) {
		if _, ok := autorizarDono(w, r, svc, r.PathValue("id")); !ok {
			return
		}
		var e EntradaCriar
		if err := tupa.LerJSON(r, &e); err != nil {
			tupa.EscreverErro(w, http.StatusBadRequest, "JSON inválido")
			return
		}
		p, err := svc.Atualizar(r.Context(), r.PathValue("id"), e)
		if err != nil {
			responderErro(w, err)
			return
		}
		_ = tupa.EscreverJSON(w, http.StatusOK, p)
	}))

	s.Rota("DELETE", "/posts/{id}", seguranca.Autenticar(emissor, func(w http.ResponseWriter, r *http.Request) {
		if _, ok := autorizarDono(w, r, svc, r.PathValue("id")); !ok {
			return
		}
		if err := svc.Deletar(r.Context(), r.PathValue("id")); err != nil {
			responderErro(w, err)
			return
		}
		w.WriteHeader(http.StatusNoContent)
	}))
}

// autorizarDono carrega o post de id e confirma que o autenticado é o dono (ou
// admin). Em falha já escreve a resposta (404 se não existe, 403 se não é dono)
// e devolve ok=false.
func autorizarDono(w http.ResponseWriter, r *http.Request, svc *Servico, id string) (Post, bool) {
	p, err := svc.Buscar(r.Context(), id)
	if err != nil {
		responderErro(w, err)
		return Post{}, false
	}
	if !seguranca.EhDono(r.Context(), p.UsuarioID) {
		tupa.EscreverErro(w, http.StatusForbidden, "acesso negado")
		return Post{}, false
	}
	return p, true
}

// responderErro mapeia erros de domínio para status HTTP.
func responderErro(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, repo.ErrNaoEncontrado):
		tupa.EscreverErro(w, http.StatusNotFound, "post não encontrado")
	case errors.Is(err, ErrUsuarioInexistente):
		tupa.EscreverErro(w, http.StatusNotFound, "usuário dono não encontrado")
	case errors.Is(err, ErrValidacao):
		tupa.EscreverErro(w, http.StatusBadRequest, err.Error())
	default:
		tupa.EscreverErro(w, http.StatusInternalServerError, "erro interno")
	}
}
