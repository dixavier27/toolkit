package usuarios

import "golang.org/x/crypto/bcrypt"

// Hasher abstrai o algoritmo de hash de senha, mantendo o domínio agnóstico
// quanto à escolha concreta.
type Hasher interface {
	Gerar(senha string) (string, error)
	Conferir(hash, senha string) bool
}

// BcryptHasher implementa Hasher usando bcrypt (golang.org/x/crypto).
type BcryptHasher struct {
	// Custo do bcrypt; 0 usa bcrypt.DefaultCost.
	Custo int
}

func (h BcryptHasher) Gerar(senha string) (string, error) {
	custo := h.Custo
	if custo == 0 {
		custo = bcrypt.DefaultCost
	}
	b, err := bcrypt.GenerateFromPassword([]byte(senha), custo)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

func (h BcryptHasher) Conferir(hash, senha string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(senha)) == nil
}
