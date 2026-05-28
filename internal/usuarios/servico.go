package usuarios

import (
	"context"

	"github.com/dixavier27/eco/pkg/inmemdb"
)

// Servico orquestra validação, hashing e persistência de usuários. Depende da
// interface inmemdb.Repositorio (não da impl) e de um Hasher.
type Servico struct {
	repo   inmemdb.Repositorio[Usuario]
	hasher Hasher
}

// NovoServico cria um serviço sobre o repositório e hasher fornecidos.
func NovoServico(repo inmemdb.Repositorio[Usuario], h Hasher) *Servico {
	return &Servico{repo: repo, hasher: h}
}

// Criar valida a entrada, hasheia a senha e persiste o usuário.
func (s *Servico) Criar(ctx context.Context, e EntradaCriar) (Usuario, error) {
	if err := validarEntrada(e); err != nil {
		return Usuario{}, err
	}
	hash, err := s.hasher.Gerar(e.Senha)
	if err != nil {
		return Usuario{}, err
	}
	u := Usuario{
		Nome:               e.Nome,
		Sobrenome:          e.Sobrenome,
		Email:              e.Email,
		Whatsapp:           e.Whatsapp,
		SenhaCriptografada: hash,
	}
	return s.repo.Criar(ctx, u)
}

// Listar devolve todos os usuários.
func (s *Servico) Listar(ctx context.Context) ([]Usuario, error) {
	return s.repo.Listar(ctx)
}

// Buscar devolve o usuário de id.
func (s *Servico) Buscar(ctx context.Context, id string) (Usuario, error) {
	return s.repo.Buscar(ctx, id)
}

// Atualizar revalida a entrada, re-hasheia a senha e substitui o usuário de id,
// preservando o id.
func (s *Servico) Atualizar(ctx context.Context, id string, e EntradaCriar) (Usuario, error) {
	if err := validarEntrada(e); err != nil {
		return Usuario{}, err
	}
	hash, err := s.hasher.Gerar(e.Senha)
	if err != nil {
		return Usuario{}, err
	}
	u := Usuario{
		ID:                 id,
		Nome:               e.Nome,
		Sobrenome:          e.Sobrenome,
		Email:              e.Email,
		Whatsapp:           e.Whatsapp,
		SenhaCriptografada: hash,
	}
	return s.repo.Atualizar(ctx, id, u)
}

// Deletar remove o usuário de id.
func (s *Servico) Deletar(ctx context.Context, id string) error {
	return s.repo.Deletar(ctx, id)
}
