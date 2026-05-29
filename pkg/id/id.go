// Package id fornece a estratégia de geração de identificadores da aplicação,
// isolada do storage. O padrão é UUIDv7 (RFC 9562): aleatório o bastante para
// não ser enumerável, porém ordenado por tempo — o que preserva boa localidade
// de índice nas escritas.
//
// A interface Gerador é injetada nos repositórios (inmemdb, mongodb) via a
// opção ComGerarID, de modo que ambos os backends compartilhem a mesma
// estratégia de identidade em vez de cada um inventar a sua.
//
// Lembrete de segurança: o id NÃO é mecanismo de controle de acesso. Mesmo
// não-enumerável, a proteção contra acesso indevido é sempre autorização.
package id

import "github.com/google/uuid"

// Gerador produz identificadores únicos como string.
type Gerador interface {
	Novo() string
}

// GeradorFunc adapta uma função simples para a interface Gerador. Útil em
// testes e para estratégias customizadas.
type GeradorFunc func() string

// Novo satisfaz Gerador chamando a própria função.
func (f GeradorFunc) Novo() string { return f() }

// UUIDv7 gera identificadores UUID versão 7 (time-ordered).
type UUIDv7 struct{}

// Novo devolve um novo UUIDv7 em formato canônico (string). Em pânico apenas se
// a fonte de entropia do sistema falhar — caso em que não há como continuar.
func (UUIDv7) Novo() string {
	return uuid.Must(uuid.NewV7()).String()
}
