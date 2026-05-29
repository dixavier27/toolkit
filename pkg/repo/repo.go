// Package repo define o contrato de persistência neutro do projeto: a interface
// Repositorio[T] e os erros sentinela que os domínios usam para mapear falhas a
// status HTTP. É implementado por pkg/inmemdb (em memória) e pkg/mongodb
// (MongoDB), sem que os domínios dependam de nenhuma impl específica.
package repo

import (
	"context"
	"errors"
)

var (
	// ErrNaoEncontrado é devolvido quando um item com o id pedido não existe.
	ErrNaoEncontrado = errors.New("repo: item não encontrado")
	// ErrJaExiste é devolvido por Criar quando o id já está em uso.
	ErrJaExiste = errors.New("repo: item já existe")
)

// Repositorio é a abstração mínima de persistência para um tipo T.
type Repositorio[T any] interface {
	Criar(ctx context.Context, item T) (T, error)
	Buscar(ctx context.Context, id string) (T, error)
	Listar(ctx context.Context) ([]T, error)
	Atualizar(ctx context.Context, id string, item T) (T, error)
	Deletar(ctx context.Context, id string) error
}
