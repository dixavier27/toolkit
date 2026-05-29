package posts

import (
	"errors"
	"net/http"

	"github.com/dixavier27/eco/pkg/repo"
	"github.com/dixavier27/eco/pkg/tupa"
)

// Registrar registra as rotas REST de posts no servidor tupa. As rotas de
// criação e listagem são aninhadas no usuário dono; as individuais são planas.
//
//	POST /usuarios/{usuarioID}/posts  cria post do usuário
//	GET  /usuarios/{usuarioID}/posts  lista posts do usuário
//	GET    /posts/{id}                busca
//	PUT    /posts/{id}                atualiza
//	DELETE /posts/{id}                remove
func Registrar(s *tupa.Servidor, svc *Servico) {
	s.Rota("POST", "/usuarios/{usuarioID}/posts", func(w http.ResponseWriter, r *http.Request) {
		var e EntradaCriar
		if err := tupa.LerJSON(r, &e); err != nil {
			tupa.EscreverErro(w, http.StatusBadRequest, "JSON inválido")
			return
		}
		p, err := svc.Criar(r.Context(), r.PathValue("usuarioID"), e)
		if err != nil {
			responderErro(w, err)
			return
		}
		_ = tupa.EscreverJSON(w, http.StatusCreated, p)
	})

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

	s.Rota("PUT", "/posts/{id}", func(w http.ResponseWriter, r *http.Request) {
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
	})

	s.Rota("DELETE", "/posts/{id}", func(w http.ResponseWriter, r *http.Request) {
		if err := svc.Deletar(r.Context(), r.PathValue("id")); err != nil {
			responderErro(w, err)
			return
		}
		w.WriteHeader(http.StatusNoContent)
	})
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
