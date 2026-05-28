// Comando de exemplo: sobe uma API integrando os módulos tupa + inmemdb com
// os domínios usuarios e posts. Use para verificação manual com curl.
package main

import (
	"context"
	"log"
	"os"
	"os/signal"

	"github.com/dixavier27/eco/internal/posts"
	"github.com/dixavier27/eco/internal/usuarios"
	"github.com/dixavier27/eco/pkg/inmemdb"
	"github.com/dixavier27/eco/pkg/tupa"
)

func main() {
	repoUsuarios := inmemdb.NovaMemoria(
		func(u usuarios.Usuario) string { return u.ID },
		inmemdb.ComDefinirID(func(u *usuarios.Usuario, id string) { u.ID = id }),
	)
	svcUsuarios := usuarios.NovoServico(repoUsuarios, usuarios.BcryptHasher{})

	repoPosts := inmemdb.NovaMemoria(
		func(p posts.Post) string { return p.ID },
		inmemdb.ComDefinirID(func(p *posts.Post, id string) { p.ID = id }),
	)
	svcPosts := posts.NovoServico(repoPosts, svcUsuarios)

	srv := tupa.Novo(":8080")
	srv.Usar(tupa.LogRequisicoes(log.Default()))
	usuarios.Registrar(srv, svcUsuarios)
	posts.Registrar(srv, svcPosts)

	ctx, parar := signal.NotifyContext(context.Background(), os.Interrupt)
	defer parar()

	log.Println("escutando em :8080 (Ctrl+C para parar)")
	if err := srv.Iniciar(ctx); err != nil {
		log.Fatal(err)
	}
}
