package posts

import (
	"errors"
	"strings"
)

// ErrValidacao indica que a entrada é inválida. A mensagem detalha o campo.
var ErrValidacao = errors.New("validação")

// validarEntrada checa os campos obrigatórios de um post.
func validarEntrada(e EntradaCriar) error {
	if strings.TrimSpace(e.Titulo) == "" {
		return erroValidacao("título é obrigatório")
	}
	if strings.TrimSpace(e.Conteudo) == "" {
		return erroValidacao("conteúdo é obrigatório")
	}
	return nil
}

func erroValidacao(msg string) error {
	return errEmbrulhado{msg: msg}
}

// errEmbrulhado carrega a mensagem específica e satisfaz errors.Is(_, ErrValidacao).
type errEmbrulhado struct{ msg string }

func (e errEmbrulhado) Error() string { return e.msg }
func (e errEmbrulhado) Unwrap() error { return ErrValidacao }
