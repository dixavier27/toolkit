package posts

import (
	"errors"
	"testing"
)

func TestValidacao(t *testing.T) {
	casos := []struct {
		nome     string
		entrada  EntradaCriar
		invalido bool
	}{
		{"válida", EntradaCriar{Titulo: "Olá", Conteudo: "mundo"}, false},
		{"sem título", EntradaCriar{Conteudo: "mundo"}, true},
		{"sem conteúdo", EntradaCriar{Titulo: "Olá"}, true},
		{"só espaços", EntradaCriar{Titulo: "  ", Conteudo: "  "}, true},
	}
	for _, c := range casos {
		t.Run(c.nome, func(t *testing.T) {
			err := validarEntrada(c.entrada)
			if c.invalido && !errors.Is(err, ErrValidacao) {
				t.Errorf("erro = %v, quer ErrValidacao", err)
			}
			if !c.invalido && err != nil {
				t.Errorf("erro inesperado: %v", err)
			}
		})
	}
}
