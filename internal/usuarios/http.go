package usuarios

import (
	"errors"
	"net/http"

	"github.com/dixavier27/eco/pkg/repo"
	"github.com/dixavier27/eco/pkg/seguranca"
	"github.com/dixavier27/eco/pkg/tupa"
)

// Registrar registra as rotas REST de usuários no servidor tupa. Criação
// (registro) e leitura são públicas; atualização e remoção exigem token válido
// e que o autenticado seja o próprio usuário (ou admin).
//
//	POST   /usuarios       cria (público — registro)
//	GET    /usuarios       lista (público)
//	GET    /usuarios/{id}  busca (público)
//	PUT    /usuarios/{id}  atualiza (self/admin)
//	DELETE /usuarios/{id}  remove (self/admin)
func Registrar(s *tupa.Servidor, svc *Servico, emissor *seguranca.Emissor) {
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

	s.Rota("PUT", "/usuarios/{id}", seguranca.Autenticar(emissor, func(w http.ResponseWriter, r *http.Request) {
		if !seguranca.EhDono(r.Context(), r.PathValue("id")) {
			tupa.EscreverErro(w, http.StatusForbidden, "acesso negado")
			return
		}
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
	}))

	s.Rota("DELETE", "/usuarios/{id}", seguranca.Autenticar(emissor, func(w http.ResponseWriter, r *http.Request) {
		if !seguranca.EhDono(r.Context(), r.PathValue("id")) {
			tupa.EscreverErro(w, http.StatusForbidden, "acesso negado")
			return
		}
		if err := svc.Deletar(r.Context(), r.PathValue("id")); err != nil {
			responderErro(w, err)
			return
		}
		w.WriteHeader(http.StatusNoContent)
	}))
}

// EmissorToken é o mínimo que o login precisa: emitir um token para o usuário
// autenticado. Satisfeito por *seguranca.Emissor — assim o domínio não importa
// o pacote seguranca nem conhece JWT.
type EmissorToken interface {
	Emitir(sub, papel string) (string, error)
}

// EntradaLogin é o payload de autenticação.
type EntradaLogin struct {
	Email string `json:"email"`
	Senha string `json:"senha"`
}

// RegistrarAuth registra a rota de autenticação:
//
//	POST /login  valida email+senha e devolve {"token": "..."}
func RegistrarAuth(s *tupa.Servidor, svc *Servico, emissor EmissorToken) {
	s.Rota("POST", "/login", func(w http.ResponseWriter, r *http.Request) {
		var e EntradaLogin
		if err := tupa.LerJSON(r, &e); err != nil {
			tupa.EscreverErro(w, http.StatusBadRequest, "JSON inválido")
			return
		}
		u, err := svc.Autenticar(r.Context(), e.Email, e.Senha)
		if err != nil {
			responderErro(w, err)
			return
		}
		token, err := emissor.Emitir(u.ID, u.Papel)
		if err != nil {
			tupa.EscreverErro(w, http.StatusInternalServerError, "erro interno")
			return
		}
		_ = tupa.EscreverJSON(w, http.StatusOK, map[string]string{"token": token})
	})
}

// responderErro mapeia erros de domínio para status HTTP.
func responderErro(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, ErrCredenciais):
		tupa.EscreverErro(w, http.StatusUnauthorized, "credenciais inválidas")
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
