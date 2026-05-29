// Package posts é um CRUD de posts pertencentes a usuários. Segue o mesmo
// padrão de internal/usuarios: domínio + validação + Servico sobre
// repo.Repositorio[Post] + rotas HTTP via tupa.
//
// As rotas são aninhadas no usuário dono (/usuarios/{usuarioID}/posts) e o
// serviço valida, na criação, que o usuário existe — via a interface
// BuscadorDeUsuario, satisfeita por *usuarios.Servico.
package posts

// Post é o recurso do domínio. UsuarioID referencia o dono; vem da URL
// (rota aninhada), não do corpo da requisição.
type Post struct {
	ID        string `json:"id" bson:"_id,omitempty"`
	UsuarioID string `json:"usuarioId" bson:"usuarioId"`
	Titulo    string `json:"titulo" bson:"titulo"`
	Conteudo  string `json:"conteudo" bson:"conteudo"`
}

// EntradaCriar é o payload de criação/atualização de um post.
type EntradaCriar struct {
	Titulo   string `json:"titulo"`
	Conteudo string `json:"conteudo"`
}
