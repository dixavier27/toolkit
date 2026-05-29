package seguranca

import (
	"context"
	"net/http"
	"strings"

	"github.com/dixavier27/eco/pkg/tupa"
)

// PapelAdmin é o papel com acesso irrestrito nas checagens de dono.
const PapelAdmin = "admin"

// chaveContexto é a chave privada sob a qual as Claims do request autenticado
// ficam guardadas no context.Context.
type chaveContexto struct{}

// comClaims devolve um contexto derivado carregando as Claims.
func comClaims(ctx context.Context, c *Claims) context.Context {
	return context.WithValue(ctx, chaveContexto{}, c)
}

// DoContexto recupera as Claims injetadas por Autenticar. ok é false se a
// requisição não passou por Autenticar.
func DoContexto(ctx context.Context) (*Claims, bool) {
	c, ok := ctx.Value(chaveContexto{}).(*Claims)
	return c, ok
}

// Autenticar envolve um handler exigindo um token Bearer válido. Em sucesso,
// injeta as Claims no contexto (recuperáveis com DoContexto) e segue; caso
// contrário responde 401.
func Autenticar(e *Emissor, prox http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		token, ok := strings.CutPrefix(r.Header.Get("Authorization"), "Bearer ")
		if !ok || token == "" {
			tupa.EscreverErro(w, http.StatusUnauthorized, "token ausente")
			return
		}
		c, err := e.Verificar(token)
		if err != nil {
			tupa.EscreverErro(w, http.StatusUnauthorized, "token inválido")
			return
		}
		prox(w, r.WithContext(comClaims(r.Context(), c)))
	}
}

// ExigirPapel envolve um handler exigindo que o usuário autenticado tenha o
// papel informado. Pressupõe Autenticar antes; responde 403 caso contrário.
func ExigirPapel(papel string, prox http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		c, ok := DoContexto(r.Context())
		if !ok || c.Papel != papel {
			tupa.EscreverErro(w, http.StatusForbidden, "acesso negado")
			return
		}
		prox(w, r)
	}
}

// EhDono informa se o usuário autenticado é o dono do recurso (Subject == donoID)
// ou tem papel admin. Use nos handlers para autorizar acesso a recursos de
// outrem. Retorna false se não houver Claims no contexto.
func EhDono(ctx context.Context, donoID string) bool {
	c, ok := DoContexto(ctx)
	if !ok {
		return false
	}
	return c.Subject == donoID || c.Papel == PapelAdmin
}
