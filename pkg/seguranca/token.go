// Package seguranca reúne os blocos de autenticação e autorização da aplicação.
//
// Nesta iteração cobre a emissão e verificação de tokens JWT (HS256) via
// Emissor. As camadas de middleware HTTP (autenticação por request,
// autorização por dono/papel e hardening) estão planejadas e serão adicionadas
// em iterações seguintes.
package seguranca

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// ErrTokenInvalido é devolvido por Verificar quando o token está ausente,
// expirado, mal-formado ou com assinatura inválida.
var ErrTokenInvalido = errors.New("seguranca: token inválido")

// Claims são as reivindicações carregadas pelo token. Subject (sub) é o id do
// usuário; Papel transporta o papel para autorização posterior.
type Claims struct {
	Papel string `json:"papel"`
	jwt.RegisteredClaims
}

// Emissor assina e verifica tokens JWT com um segredo simétrico (HS256).
type Emissor struct {
	segredo []byte
	ttl     time.Duration
}

// NovoEmissor cria um Emissor. segredo nunca deve ser hard-coded — injete via
// variável de ambiente/secret. ttl é a validade de cada token emitido.
func NovoEmissor(segredo string, ttl time.Duration) *Emissor {
	return &Emissor{segredo: []byte(segredo), ttl: ttl}
}

// Emitir gera um token assinado para o usuário sub com o papel informado,
// válido por ttl a partir de agora.
func (e *Emissor) Emitir(sub, papel string) (string, error) {
	agora := time.Now()
	c := Claims{
		Papel: papel,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   sub,
			IssuedAt:  jwt.NewNumericDate(agora),
			ExpiresAt: jwt.NewNumericDate(agora.Add(e.ttl)),
		},
	}
	return jwt.NewWithClaims(jwt.SigningMethodHS256, c).SignedString(e.segredo)
}

// Verificar valida a assinatura e a expiração do token e devolve suas Claims.
// Qualquer falha resulta em ErrTokenInvalido (embrulhando a causa).
func (e *Emissor) Verificar(token string) (*Claims, error) {
	var c Claims
	_, err := jwt.ParseWithClaims(token, &c, func(*jwt.Token) (any, error) {
		return e.segredo, nil
	}, jwt.WithValidMethods([]string{"HS256"}))
	if err != nil {
		return nil, errors.Join(ErrTokenInvalido, err)
	}
	return &c, nil
}
