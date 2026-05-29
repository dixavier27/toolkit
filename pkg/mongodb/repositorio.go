// Package mongodb fornece uma implementação do contrato repo.Repositorio[T]
// sobre o MongoDB, usando o driver oficial (go.mongodb.org/mongo-driver/v2).
//
// Espelha a API de construção de pkg/inmemdb (NovaColecao/ComDefinirID) para que
// o wiring de um backend seja simétrico ao do outro. O id do domínio é uma
// string (mapeada para o campo _id); quando ausente na criação, geramos um
// ObjectID em hex e o gravamos de volta no item via ComDefinirID.
package mongodb

import (
	"context"
	"errors"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"

	"github.com/dixavier27/eco/pkg/repo"
)

// Colecao é uma implementação de repo.Repositorio[T] sobre uma coleção Mongo.
type Colecao[T any] struct {
	col      *mongo.Collection
	idDe     func(T) string   // extrai o id de um item
	defineID func(*T, string) // opcional: grava o id gerado de volta no item
	gerar    func() string    // opcional: estratégia de geração de id
}

// Opcao configura uma *Colecao na criação.
type Opcao[T any] func(*Colecao[T])

// ComDefinirID registra um setter que grava o id gerado de volta no item quando
// Criar recebe um item sem id. Sem ele, Criar exige id não-vazio.
func ComDefinirID[T any](definir func(*T, string)) Opcao[T] {
	return func(c *Colecao[T]) { c.defineID = definir }
}

// ComGerarID define a estratégia de geração de id usada quando Criar recebe um
// item sem id (ex.: id.UUIDv7{}.Novo). Sem ela, é gerado um ObjectID em hex.
// Requer ComDefinirID para gravar o id no item.
func ComGerarID[T any](gerar func() string) Opcao[T] {
	return func(c *Colecao[T]) { c.gerar = gerar }
}

// NovaColecao cria um repositório sobre col. idDe extrai o id de cada item.
// Veja ComDefinirID para geração automática de id (ObjectID hex).
func NovaColecao[T any](col *mongo.Collection, idDe func(T) string, opts ...Opcao[T]) *Colecao[T] {
	c := &Colecao[T]{col: col, idDe: idDe}
	for _, opt := range opts {
		opt(c)
	}
	return c
}

// Criar insere um item. Se o id vier vazio e ComDefinirID tiver sido fornecido,
// um ObjectID em hex é gerado e gravado no item; caso contrário, devolve
// repo.ErrNaoEncontrado. Id duplicado devolve repo.ErrJaExiste.
func (c *Colecao[T]) Criar(ctx context.Context, item T) (T, error) {
	id := c.idDe(item)
	if id == "" {
		if c.defineID == nil {
			var zero T
			return zero, repo.ErrNaoEncontrado
		}
		if c.gerar != nil {
			id = c.gerar()
		} else {
			id = bson.NewObjectID().Hex()
		}
		c.defineID(&item, id)
	}
	if _, err := c.col.InsertOne(ctx, item); err != nil {
		if mongo.IsDuplicateKeyError(err) {
			var zero T
			return zero, repo.ErrJaExiste
		}
		var zero T
		return zero, err
	}
	return item, nil
}

// Buscar devolve o item de id, ou repo.ErrNaoEncontrado.
func (c *Colecao[T]) Buscar(ctx context.Context, id string) (T, error) {
	var out T
	err := c.col.FindOne(ctx, bson.M{"_id": id}).Decode(&out)
	if errors.Is(err, mongo.ErrNoDocuments) {
		var zero T
		return zero, repo.ErrNaoEncontrado
	}
	if err != nil {
		var zero T
		return zero, err
	}
	return out, nil
}

// Listar devolve todos os itens da coleção (ordem não garantida).
func (c *Colecao[T]) Listar(ctx context.Context) ([]T, error) {
	cur, err := c.col.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	out := make([]T, 0)
	if err := cur.All(ctx, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// Atualizar substitui o item de id. Devolve repo.ErrNaoEncontrado se não existir.
func (c *Colecao[T]) Atualizar(ctx context.Context, id string, item T) (T, error) {
	res, err := c.col.ReplaceOne(ctx, bson.M{"_id": id}, item)
	if err != nil {
		var zero T
		return zero, err
	}
	if res.MatchedCount == 0 {
		var zero T
		return zero, repo.ErrNaoEncontrado
	}
	return item, nil
}

// Deletar remove o item de id. Devolve repo.ErrNaoEncontrado se não existir.
func (c *Colecao[T]) Deletar(ctx context.Context, id string) error {
	res, err := c.col.DeleteOne(ctx, bson.M{"_id": id})
	if err != nil {
		return err
	}
	if res.DeletedCount == 0 {
		return repo.ErrNaoEncontrado
	}
	return nil
}

// Conectar abre uma conexão com o Mongo em uri, valida com um ping e devolve o
// banco dbName junto de uma função de desligamento. Chame a função devolvida
// (tipicamente com defer) para encerrar a conexão.
func Conectar(ctx context.Context, uri, dbName string) (*mongo.Database, func(context.Context) error, error) {
	cliente, err := mongo.Connect(options.Client().ApplyURI(uri))
	if err != nil {
		return nil, nil, err
	}
	if err := cliente.Ping(ctx, nil); err != nil {
		_ = cliente.Disconnect(ctx)
		return nil, nil, err
	}
	return cliente.Database(dbName), cliente.Disconnect, nil
}
