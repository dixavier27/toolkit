// Package inmemdb fornece uma implementação em memória, thread-safe, do
// contrato repo.Repositorio[T], pensada para MVPs, demos e testes. Só stdlib —
// nenhuma dependência externa.
//
// O domínio depende da interface repo.Repositorio[T]; nos testes, injete uma
// *Memoria[T] em vez de um banco real.
package inmemdb

import (
	"context"
	"strconv"
	"sync"

	"github.com/dixavier27/eco/pkg/repo"
)

// Memoria é uma implementação in-memory thread-safe de repo.Repositorio[T].
type Memoria[T any] struct {
	idDe     func(T) string   // extrai o id de um item
	defineID func(*T, string) // opcional: grava o id gerado de volta no item
	mu       sync.RWMutex
	itens    map[string]T
	seq      int
}

// OpcaoMemoria configura uma *Memoria na criação.
type OpcaoMemoria[T any] func(*Memoria[T])

// ComDefinirID registra um setter que grava o id gerado de volta no item
// quando Criar recebe um item sem id. Sem ele, Criar exige id não-vazio.
func ComDefinirID[T any](definir func(*T, string)) OpcaoMemoria[T] {
	return func(m *Memoria[T]) { m.defineID = definir }
}

// NovaMemoria cria um repositório em memória vazio. idDe extrai o id de cada
// item. Veja ComDefinirID para geração automática de id sequencial.
func NovaMemoria[T any](idDe func(T) string, opts ...OpcaoMemoria[T]) *Memoria[T] {
	m := &Memoria[T]{
		idDe:  idDe,
		itens: make(map[string]T),
	}
	for _, opt := range opts {
		opt(m)
	}
	return m
}

// Criar insere um item. Se o id vier vazio e ComDefinirID tiver sido fornecido,
// um id sequencial é gerado e gravado no item; caso contrário, devolve
// repo.ErrNaoEncontrado/repo.ErrJaExiste conforme o caso.
func (m *Memoria[T]) Criar(ctx context.Context, item T) (T, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	id := m.idDe(item)
	if id == "" {
		if m.defineID == nil {
			var zero T
			return zero, repo.ErrNaoEncontrado
		}
		m.seq++
		id = strconv.Itoa(m.seq)
		m.defineID(&item, id)
	}
	if _, dup := m.itens[id]; dup {
		var zero T
		return zero, repo.ErrJaExiste
	}
	m.itens[id] = item
	return item, nil
}

// Buscar devolve o item de id, ou repo.ErrNaoEncontrado.
func (m *Memoria[T]) Buscar(ctx context.Context, id string) (T, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	v, ok := m.itens[id]
	if !ok {
		var zero T
		return zero, repo.ErrNaoEncontrado
	}
	return v, nil
}

// Listar devolve todos os itens (ordem não garantida).
func (m *Memoria[T]) Listar(ctx context.Context) ([]T, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	out := make([]T, 0, len(m.itens))
	for _, v := range m.itens {
		out = append(out, v)
	}
	return out, nil
}

// Atualizar substitui o item de id. Devolve repo.ErrNaoEncontrado se não existir.
func (m *Memoria[T]) Atualizar(ctx context.Context, id string, item T) (T, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	if _, ok := m.itens[id]; !ok {
		var zero T
		return zero, repo.ErrNaoEncontrado
	}
	m.itens[id] = item
	return item, nil
}

// Deletar remove o item de id. Devolve repo.ErrNaoEncontrado se não existir.
func (m *Memoria[T]) Deletar(ctx context.Context, id string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	if _, ok := m.itens[id]; !ok {
		return repo.ErrNaoEncontrado
	}
	delete(m.itens, id)
	return nil
}
