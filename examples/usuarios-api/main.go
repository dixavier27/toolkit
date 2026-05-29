// Comando de exemplo: sobe uma API integrando os módulos tupa + repo (inmemdb
// ou mongodb) com os domínios usuarios e posts. Use para verificação manual com
// curl.
//
// O backend de persistência é escolhido por env var na inicialização:
//
//	DB_DRIVER=memory                 # padrão: in-memory, zero dependências
//	DB_DRIVER=mongo MONGO_URI=...    # MongoDB (MONGO_DB opcional, default "eco")
package main

import (
	"context"
	"log"
	"os"
	"os/signal"

	"github.com/dixavier27/eco/internal/posts"
	"github.com/dixavier27/eco/internal/usuarios"
	"github.com/dixavier27/eco/pkg/id"
	"github.com/dixavier27/eco/pkg/inmemdb"
	"github.com/dixavier27/eco/pkg/mongodb"
	"github.com/dixavier27/eco/pkg/repo"
	"github.com/dixavier27/eco/pkg/tupa"
)

func main() {
	ctx, parar := signal.NotifyContext(context.Background(), os.Interrupt)
	defer parar()

	repoUsuarios, repoPosts, fechar := abrirRepos(ctx)
	defer fechar()

	svcUsuarios := usuarios.NovoServico(repoUsuarios, usuarios.BcryptHasher{})
	svcPosts := posts.NovoServico(repoPosts, svcUsuarios)

	srv := tupa.Novo(":8080")
	srv.Usar(tupa.LogRequisicoes(log.Default()))
	usuarios.Registrar(srv, svcUsuarios)
	posts.Registrar(srv, svcPosts)

	log.Println("escutando em :8080 (Ctrl+C para parar)")
	if err := srv.Iniciar(ctx); err != nil {
		log.Fatal(err)
	}
}

// abrirRepos seleciona o backend de persistência conforme DB_DRIVER e devolve os
// repositórios de usuários e posts, mais uma função de desligamento.
func abrirRepos(ctx context.Context) (repo.Repositorio[usuarios.Usuario], repo.Repositorio[posts.Post], func()) {
	idUsuario := func(u usuarios.Usuario) string { return u.ID }
	defUsuario := func(u *usuarios.Usuario, id string) { u.ID = id }
	idPost := func(p posts.Post) string { return p.ID }
	defPost := func(p *posts.Post, id string) { p.ID = id }

	// Estratégia de identidade única para todos os backends: UUIDv7.
	gerar := id.UUIDv7{}.Novo

	switch os.Getenv("DB_DRIVER") {
	case "mongo":
		uri := os.Getenv("MONGO_URI")
		if uri == "" {
			log.Fatal("DB_DRIVER=mongo exige MONGO_URI")
		}
		nomeDB := os.Getenv("MONGO_DB")
		if nomeDB == "" {
			nomeDB = "eco"
		}
		db, desconectar, err := mongodb.Conectar(ctx, uri, nomeDB)
		if err != nil {
			log.Fatalf("conectar ao mongo: %v", err)
		}
		log.Printf("backend: mongo (%s/%s)", uri, nomeDB)
		ru := mongodb.NovaColecao(db.Collection("usuarios"), idUsuario, mongodb.ComDefinirID(defUsuario), mongodb.ComGerarID[usuarios.Usuario](gerar))
		rp := mongodb.NovaColecao(db.Collection("posts"), idPost, mongodb.ComDefinirID(defPost), mongodb.ComGerarID[posts.Post](gerar))
		return ru, rp, func() { _ = desconectar(context.Background()) }
	default:
		log.Println("backend: memory")
		ru := inmemdb.NovaMemoria(idUsuario, inmemdb.ComDefinirID(defUsuario), inmemdb.ComGerarID[usuarios.Usuario](gerar))
		rp := inmemdb.NovaMemoria(idPost, inmemdb.ComDefinirID(defPost), inmemdb.ComGerarID[posts.Post](gerar))
		return ru, rp, func() {}
	}
}
