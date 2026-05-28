package inmemdb

import (
	"context"
	"errors"
	"sync"
	"testing"
)

type item struct {
	ID   string
	Nome string
}

func novoRepo() *Memoria[item] {
	return NovaMemoria(
		func(i item) string { return i.ID },
		ComDefinirID(func(i *item, id string) { i.ID = id }),
	)
}

func TestCRUD(t *testing.T) {
	ctx := context.Background()
	r := novoRepo()

	criado, err := r.Criar(ctx, item{Nome: "a"})
	if err != nil {
		t.Fatal(err)
	}
	if criado.ID == "" {
		t.Fatal("id não foi gerado")
	}

	got, err := r.Buscar(ctx, criado.ID)
	if err != nil || got.Nome != "a" {
		t.Fatalf("Buscar = %+v, %v", got, err)
	}

	if _, err := r.Atualizar(ctx, criado.ID, item{ID: criado.ID, Nome: "b"}); err != nil {
		t.Fatal(err)
	}
	got, _ = r.Buscar(ctx, criado.ID)
	if got.Nome != "b" {
		t.Errorf("após atualizar, nome = %q, quer b", got.Nome)
	}

	lista, _ := r.Listar(ctx)
	if len(lista) != 1 {
		t.Errorf("len(lista) = %d, quer 1", len(lista))
	}

	if err := r.Deletar(ctx, criado.ID); err != nil {
		t.Fatal(err)
	}
	if _, err := r.Buscar(ctx, criado.ID); !errors.Is(err, ErrNaoEncontrado) {
		t.Errorf("após deletar, erro = %v, quer ErrNaoEncontrado", err)
	}
}

func TestErros(t *testing.T) {
	ctx := context.Background()

	t.Run("sem id e sem definirID", func(t *testing.T) {
		r := NovaMemoria(func(i item) string { return i.ID })
		if _, err := r.Criar(ctx, item{}); !errors.Is(err, ErrNaoEncontrado) {
			t.Errorf("erro = %v, quer ErrNaoEncontrado", err)
		}
	})

	t.Run("id duplicado", func(t *testing.T) {
		r := novoRepo()
		if _, err := r.Criar(ctx, item{ID: "x"}); err != nil {
			t.Fatal(err)
		}
		if _, err := r.Criar(ctx, item{ID: "x"}); !errors.Is(err, ErrJaExiste) {
			t.Errorf("erro = %v, quer ErrJaExiste", err)
		}
	})

	t.Run("atualizar inexistente", func(t *testing.T) {
		r := novoRepo()
		if _, err := r.Atualizar(ctx, "nada", item{ID: "nada"}); !errors.Is(err, ErrNaoEncontrado) {
			t.Errorf("erro = %v, quer ErrNaoEncontrado", err)
		}
	})
}

func TestConcorrencia(t *testing.T) {
	ctx := context.Background()
	r := novoRepo()
	var wg sync.WaitGroup
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			_, _ = r.Criar(ctx, item{Nome: "x"})
		}()
	}
	wg.Wait()
	lista, _ := r.Listar(ctx)
	if len(lista) != 100 {
		t.Errorf("len = %d, quer 100", len(lista))
	}
}
