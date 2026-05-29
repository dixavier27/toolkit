//go:build mongo

// Testes de integração do repositório Mongo. Compilam apenas com -tags mongo e
// pulam se MONGO_URI não estiver definida. Requerem um Mongo acessível, ex.:
//
//	docker run -d -p 27017:27017 mongo
//	MONGO_URI=mongodb://localhost:27017 go test -tags mongo ./pkg/mongodb/...
package mongodb

import (
	"context"
	"errors"
	"os"
	"testing"
	"time"

	"github.com/dixavier27/eco/pkg/repo"
)

type item struct {
	ID   string `bson:"_id,omitempty"`
	Nome string `bson:"nome"`
}

func novoRepo(t *testing.T) (*Colecao[item], func()) {
	t.Helper()
	uri := os.Getenv("MONGO_URI")
	if uri == "" {
		t.Skip("MONGO_URI não definida; pulando teste de integração do Mongo")
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	db, fechar, err := Conectar(ctx, uri, "eco_test")
	if err != nil {
		t.Fatalf("Conectar: %v", err)
	}
	col := db.Collection("itens_" + t.Name())
	_ = col.Drop(ctx)

	r := NovaColecao(
		col,
		func(i item) string { return i.ID },
		ComDefinirID(func(i *item, id string) { i.ID = id }),
	)
	return r, func() {
		_ = col.Drop(context.Background())
		_ = fechar(context.Background())
	}
}

func TestCRUD(t *testing.T) {
	r, fechar := novoRepo(t)
	defer fechar()
	ctx := context.Background()

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
	if _, err := r.Buscar(ctx, criado.ID); !errors.Is(err, repo.ErrNaoEncontrado) {
		t.Errorf("após deletar, erro = %v, quer ErrNaoEncontrado", err)
	}
}

func TestErros(t *testing.T) {
	r, fechar := novoRepo(t)
	defer fechar()
	ctx := context.Background()

	if _, err := r.Atualizar(ctx, "nada", item{ID: "nada"}); !errors.Is(err, repo.ErrNaoEncontrado) {
		t.Errorf("atualizar inexistente: erro = %v, quer ErrNaoEncontrado", err)
	}
	if err := r.Deletar(ctx, "nada"); !errors.Is(err, repo.ErrNaoEncontrado) {
		t.Errorf("deletar inexistente: erro = %v, quer ErrNaoEncontrado", err)
	}
}
