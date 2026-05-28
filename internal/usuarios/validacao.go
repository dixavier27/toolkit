package usuarios

import (
	"errors"
	"regexp"
	"strings"
)

// ErrValidacao indica que a entrada é inválida. A mensagem detalha o campo.
var ErrValidacao = errors.New("validação")

// formatos mínimos — não exaustivos, suficientes para um MVP enxuto.
var (
	reEmail    = regexp.MustCompile(`^[^@\s]+@[^@\s]+\.[^@\s]+$`)
	reWhatsapp = regexp.MustCompile(`^\+?[0-9]{10,15}$`)
)

// validarEntrada checa campos obrigatórios e formatos de email/whatsapp.
// Devolve um erro que embrulha ErrValidacao quando algo está inválido.
func validarEntrada(e EntradaCriar) error {
	if strings.TrimSpace(e.Nome) == "" {
		return erroValidacao("nome é obrigatório")
	}
	if strings.TrimSpace(e.Sobrenome) == "" {
		return erroValidacao("sobrenome é obrigatório")
	}
	if !reEmail.MatchString(e.Email) {
		return erroValidacao("email inválido")
	}
	if !reWhatsapp.MatchString(e.Whatsapp) {
		return erroValidacao("whatsapp inválido")
	}
	if len(e.Senha) < 8 {
		return erroValidacao("senha deve ter ao menos 8 caracteres")
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
