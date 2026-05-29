package posts

import (
	"context"
	"errors"

	"github.com/dixavier27/eco/internal/usuarios"
	"github.com/dixavier27/eco/pkg/repo"
)

// ErrUsuarioInexistente é devolvido por Criar quando o usuário dono não existe.
var ErrUsuarioInexistente = errors.New("usuário dono não existe")

// BuscadorDeUsuario é o mínimo que posts precisa do domínio de usuários:
// confirmar que o dono existe. *usuarios.Servico satisfaz esta interface.
type BuscadorDeUsuario interface {
	Buscar(ctx context.Context, id string) (usuarios.Usuario, error)
}

// Servico orquestra validação e persistência de posts. Depende da interface
// repo.Repositorio (não da impl) e de um BuscadorDeUsuario para validar o dono.
type Servico struct {
	repo     repo.Repositorio[Post]
	usuarios BuscadorDeUsuario
}

// NovoServico cria um serviço sobre o repositório e o buscador de usuários.
func NovoServico(repo repo.Repositorio[Post], donos BuscadorDeUsuario) *Servico {
	return &Servico{repo: repo, usuarios: donos}
}

// Criar valida a entrada, confirma que o usuário dono existe e persiste o post.
func (s *Servico) Criar(ctx context.Context, usuarioID string, e EntradaCriar) (Post, error) {
	if err := validarEntrada(e); err != nil {
		return Post{}, err
	}
	if _, err := s.usuarios.Buscar(ctx, usuarioID); err != nil {
		if errors.Is(err, repo.ErrNaoEncontrado) {
			return Post{}, ErrUsuarioInexistente
		}
		return Post{}, err
	}
	p := Post{
		UsuarioID: usuarioID,
		Titulo:    e.Titulo,
		Conteudo:  e.Conteudo,
	}
	return s.repo.Criar(ctx, p)
}

// ListarDoUsuario devolve os posts de um usuário (filtra em memória).
func (s *Servico) ListarDoUsuario(ctx context.Context, usuarioID string) ([]Post, error) {
	todos, err := s.repo.Listar(ctx)
	if err != nil {
		return nil, err
	}
	out := make([]Post, 0)
	for _, p := range todos {
		if p.UsuarioID == usuarioID {
			out = append(out, p)
		}
	}
	return out, nil
}

// Buscar devolve o post de id.
func (s *Servico) Buscar(ctx context.Context, id string) (Post, error) {
	return s.repo.Buscar(ctx, id)
}

// Atualizar revalida a entrada e substitui o post de id, preservando id e dono.
func (s *Servico) Atualizar(ctx context.Context, id string, e EntradaCriar) (Post, error) {
	if err := validarEntrada(e); err != nil {
		return Post{}, err
	}
	atual, err := s.repo.Buscar(ctx, id)
	if err != nil {
		return Post{}, err
	}
	atual.Titulo = e.Titulo
	atual.Conteudo = e.Conteudo
	return s.repo.Atualizar(ctx, id, atual)
}

// Deletar remove o post de id.
func (s *Servico) Deletar(ctx context.Context, id string) error {
	return s.repo.Deletar(ctx, id)
}
