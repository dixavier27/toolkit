package usuarios

import (
	"context"
	"errors"

	"github.com/dixavier27/eco/pkg/repo"
)

// PapelPadrao é o papel atribuído a usuários recém-criados.
const PapelPadrao = "user"

// ErrCredenciais é devolvido por Autenticar quando email/senha não conferem.
var ErrCredenciais = errors.New("credenciais inválidas")

// Servico orquestra validação, hashing e persistência de usuários. Depende da
// interface repo.Repositorio (não da impl) e de um Hasher.
type Servico struct {
	repo   repo.Repositorio[Usuario]
	hasher Hasher
}

// NovoServico cria um serviço sobre o repositório e hasher fornecidos.
func NovoServico(repo repo.Repositorio[Usuario], h Hasher) *Servico {
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
		Papel:              PapelPadrao,
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
// preservando o id e o papel atual.
func (s *Servico) Atualizar(ctx context.Context, id string, e EntradaCriar) (Usuario, error) {
	if err := validarEntrada(e); err != nil {
		return Usuario{}, err
	}
	atual, err := s.repo.Buscar(ctx, id)
	if err != nil {
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
		Papel:              atual.Papel,
		SenhaCriptografada: hash,
	}
	return s.repo.Atualizar(ctx, id, u)
}

// Autenticar localiza o usuário pelo email e confere a senha. Devolve
// ErrCredenciais se o email não existir ou a senha não bater.
//
// Gotcha (PoC): a busca por email é feita varrendo Listar (O(n), sem índice).
// Em produção, use um índice único de email e uma query nativa por backend.
func (s *Servico) Autenticar(ctx context.Context, email, senha string) (Usuario, error) {
	todos, err := s.repo.Listar(ctx)
	if err != nil {
		return Usuario{}, err
	}
	for _, u := range todos {
		if u.Email == email {
			if s.hasher.Conferir(u.SenhaCriptografada, senha) {
				return u, nil
			}
			return Usuario{}, ErrCredenciais
		}
	}
	return Usuario{}, ErrCredenciais
}

// Deletar remove o usuário de id.
func (s *Servico) Deletar(ctx context.Context, id string) error {
	return s.repo.Deletar(ctx, id)
}
