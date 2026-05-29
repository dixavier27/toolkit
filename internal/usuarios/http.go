package usuarios

import (
	"errors"
	"net/http"

	"github.com/dixavier27/eco/pkg/repo"
	"github.com/dixavier27/eco/pkg/tupa"
)

// Registrar registra as rotas REST de usuários no servidor tupa.
//
//	POST   /usuarios       cria
//	GET    /usuarios       lista
//	GET    /usuarios/{id}  busca
//	PUT    /usuarios/{id}  atualiza
//	DELETE /usuarios/{id}  remove
func Registrar(s *tupa.Servidor, svc *Servico) {
	s.Rota("POST", "/usuarios", func(w http.ResponseWriter, r *http.Request) {
		var e EntradaCriar
		if err := tupa.LerJSON(r, &e); err != nil {
			tupa.EscreverErro(w, http.StatusBadRequest, "JSON inválido")
			return
		}
		u, err := svc.Criar(r.Context(), e)
		if err != nil {
			responderErro(w, err)
			return
		}
		_ = tupa.EscreverJSON(w, http.StatusCreated, u)
	})

	s.Rota("GET", "/usuarios", func(w http.ResponseWriter, r *http.Request) {
		us, err := svc.Listar(r.Context())
		if err != nil {
			responderErro(w, err)
			return
		}
		_ = tupa.EscreverJSON(w, http.StatusOK, us)
	})

	s.Rota("GET", "/usuarios/{id}", func(w http.ResponseWriter, r *http.Request) {
		u, err := svc.Buscar(r.Context(), r.PathValue("id"))
		if err != nil {
			responderErro(w, err)
			return
		}
		_ = tupa.EscreverJSON(w, http.StatusOK, u)
	})

	s.Rota("PUT", "/usuarios/{id}", func(w http.ResponseWriter, r *http.Request) {
		var e EntradaCriar
		if err := tupa.LerJSON(r, &e); err != nil {
			tupa.EscreverErro(w, http.StatusBadRequest, "JSON inválido")
			return
		}
		u, err := svc.Atualizar(r.Context(), r.PathValue("id"), e)
		if err != nil {
			responderErro(w, err)
			return
		}
		_ = tupa.EscreverJSON(w, http.StatusOK, u)
	})

	s.Rota("DELETE", "/usuarios/{id}", func(w http.ResponseWriter, r *http.Request) {
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
		tupa.EscreverErro(w, http.StatusNotFound, "usuário não encontrado")
	case errors.Is(err, repo.ErrJaExiste):
		tupa.EscreverErro(w, http.StatusConflict, "usuário já existe")
	case errors.Is(err, ErrValidacao):
		tupa.EscreverErro(w, http.StatusBadRequest, err.Error())
	default:
		tupa.EscreverErro(w, http.StatusInternalServerError, "erro interno")
	}
}
