package usuarios

import (
	"errors"
	"testing"
)

func entradaValida() EntradaCriar {
	return EntradaCriar{
		Nome:      "Ana",
		Sobrenome: "Lima",
		Email:     "ana@ex.com",
		Whatsapp:  "+5511999990000",
		Senha:     "segredo123",
	}
}

func TestValidacao(t *testing.T) {
	casos := []struct {
		nome     string
		mutar    func(*EntradaCriar)
		invalido bool
	}{
		{"válida", func(*EntradaCriar) {}, false},
		{"sem nome", func(e *EntradaCriar) { e.Nome = "" }, true},
		{"sem sobrenome", func(e *EntradaCriar) { e.Sobrenome = "" }, true},
		{"email ruim", func(e *EntradaCriar) { e.Email = "ana#ex" }, true},
		{"whatsapp ruim", func(e *EntradaCriar) { e.Whatsapp = "abc" }, true},
		{"senha curta", func(e *EntradaCriar) { e.Senha = "123" }, true},
	}
	for _, c := range casos {
		t.Run(c.nome, func(t *testing.T) {
			e := entradaValida()
			c.mutar(&e)
			err := validarEntrada(e)
			if c.invalido && !errors.Is(err, ErrValidacao) {
				t.Errorf("erro = %v, quer ErrValidacao", err)
			}
			if !c.invalido && err != nil {
				t.Errorf("erro inesperado: %v", err)
			}
		})
	}
}

func TestBcryptHasher(t *testing.T) {
	h := BcryptHasher{}
	hash, err := h.Gerar("segredo123")
	if err != nil {
		t.Fatal(err)
	}
	if hash == "segredo123" {
		t.Fatal("hash igual ao texto plano")
	}
	if !h.Conferir(hash, "segredo123") {
		t.Error("Conferir falhou para senha correta")
	}
	if h.Conferir(hash, "errada") {
		t.Error("Conferir aceitou senha errada")
	}
}
