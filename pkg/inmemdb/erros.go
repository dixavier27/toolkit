package inmemdb

import "errors"

var (
	// ErrNaoEncontrado é devolvido quando um item com o id pedido não existe.
	ErrNaoEncontrado = errors.New("inmemdb: item não encontrado")
	// ErrJaExiste é devolvido por Criar quando o id já está em uso.
	ErrJaExiste = errors.New("inmemdb: item já existe")
)
