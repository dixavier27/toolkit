package seguranca

import (
	"errors"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func TestEmitirVerificar(t *testing.T) {
	e := NovoEmissor("segredo", time.Hour)
	tok, err := e.Emitir("user-1", "admin")
	if err != nil {
		t.Fatal(err)
	}
	c, err := e.Verificar(tok)
	if err != nil {
		t.Fatalf("Verificar: %v", err)
	}
	if c.Subject != "user-1" {
		t.Errorf("sub = %q, quer user-1", c.Subject)
	}
	if c.Papel != "admin" {
		t.Errorf("papel = %q, quer admin", c.Papel)
	}
}

func TestTokenExpirado(t *testing.T) {
	e := NovoEmissor("segredo", -time.Minute) // já nasce expirado
	tok, err := e.Emitir("user-1", "user")
	if err != nil {
		t.Fatal(err)
	}
	if _, err := e.Verificar(tok); !errors.Is(err, ErrTokenInvalido) {
		t.Errorf("erro = %v, quer ErrTokenInvalido", err)
	}
}

func TestAssinaturaAdulterada(t *testing.T) {
	tok, err := NovoEmissor("segredo-certo", time.Hour).Emitir("user-1", "user")
	if err != nil {
		t.Fatal(err)
	}
	// Outro emissor (segredo diferente) não deve validar o token.
	if _, err := NovoEmissor("segredo-errado", time.Hour).Verificar(tok); !errors.Is(err, ErrTokenInvalido) {
		t.Errorf("erro = %v, quer ErrTokenInvalido", err)
	}
}

func TestAlgoritmoRejeitado(t *testing.T) {
	// Token "none" (sem assinatura) deve ser rejeitado pela lista de métodos válidos.
	tok, err := jwt.NewWithClaims(jwt.SigningMethodNone, Claims{
		RegisteredClaims: jwt.RegisteredClaims{Subject: "intruso"},
	}).SignedString(jwt.UnsafeAllowNoneSignatureType)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := NovoEmissor("segredo", time.Hour).Verificar(tok); !errors.Is(err, ErrTokenInvalido) {
		t.Errorf("erro = %v, quer ErrTokenInvalido", err)
	}
}
