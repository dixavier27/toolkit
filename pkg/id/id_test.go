package id

import (
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestUUIDv7VersaoEVariante(t *testing.T) {
	s := UUIDv7{}.Novo()
	u, err := uuid.Parse(s)
	if err != nil {
		t.Fatalf("Parse(%q): %v", s, err)
	}
	if u.Version() != 7 {
		t.Errorf("versão = %d, quer 7", u.Version())
	}
	if u.Variant() != uuid.RFC4122 {
		t.Errorf("variante = %v, quer RFC4122", u.Variant())
	}
}

func TestUUIDv7Unicidade(t *testing.T) {
	const n = 10000
	g := UUIDv7{}
	vistos := make(map[string]struct{}, n)
	for i := 0; i < n; i++ {
		s := g.Novo()
		if _, dup := vistos[s]; dup {
			t.Fatalf("id duplicado gerado: %s", s)
		}
		vistos[s] = struct{}{}
	}
}

func TestUUIDv7OrdenadoPorTempo(t *testing.T) {
	g := UUIDv7{}
	antes := g.Novo()
	time.Sleep(2 * time.Millisecond)
	depois := g.Novo()
	if !(depois > antes) {
		t.Errorf("UUIDv7 deveria crescer com o tempo: antes=%s depois=%s", antes, depois)
	}
}

func TestGeradorFunc(t *testing.T) {
	var g Gerador = GeradorFunc(func() string { return "fixo" })
	if g.Novo() != "fixo" {
		t.Errorf("GeradorFunc.Novo() = %q, quer \"fixo\"", g.Novo())
	}
}
