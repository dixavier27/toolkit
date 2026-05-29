// Comando de exemplo: sobe uma API integrando os módulos tupa + repo (inmemdb
// ou mongodb) com os domínios usuarios e posts. Use para verificação manual com
// curl.
//
// O backend de persistência é escolhido por env var na inicialização:
//
//	DB_DRIVER=memory                 # padrão: in-memory, zero dependências
//	DB_DRIVER=mongo MONGO_URI=...    # MongoDB (MONGO_DB opcional, default "eco")
//
// Autenticação (login + JWT) exige um segredo:
//
//	JWT_SECRET=...                   # obrigatório; nunca hard-code/commit
//	JWT_TTL=15m                      # validade do token (default 15m)
package main

import (
	"bufio"
	"context"
	"log"
	"os"
	"os/signal"
	"strings"
	"time"

	"github.com/dixavier27/eco/internal/posts"
	"github.com/dixavier27/eco/internal/usuarios"
	"github.com/dixavier27/eco/pkg/id"
	"github.com/dixavier27/eco/pkg/inmemdb"
	"github.com/dixavier27/eco/pkg/mongodb"
	"github.com/dixavier27/eco/pkg/repo"
	"github.com/dixavier27/eco/pkg/seguranca"
	"github.com/dixavier27/eco/pkg/tupa"
)

func main() {
	carregarEnv(".env")

	ctx, parar := signal.NotifyContext(context.Background(), os.Interrupt)
	defer parar()

	repoUsuarios, repoPosts, fechar := abrirRepos(ctx)
	defer fechar()

	svcUsuarios := usuarios.NovoServico(repoUsuarios, usuarios.BcryptHasher{})
	svcPosts := posts.NovoServico(repoPosts, svcUsuarios)

	emissor := abrirEmissor()

	srv := tupa.Novo(":8080")
	srv.Usar(
		tupa.LogRequisicoes(log.Default()),
		seguranca.CabecalhosSeguranca(),
		seguranca.CORS(),
		seguranca.LimitarTaxa(50, 100), // 50 req/s por IP, rajada de 100
	)
	usuarios.Registrar(srv, svcUsuarios, emissor)
	usuarios.RegistrarAuth(srv, svcUsuarios, emissor)
	posts.Registrar(srv, svcPosts, emissor)

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

// carregarEnv lê um arquivo no formato KEY=VALUE e popula os.Environ com as
// variáveis ainda não definidas (env explícita do shell tem precedência).
// Ignora linhas em branco e comentários (#). Silencioso se o arquivo não existe.
func carregarEnv(arquivo string) {
	f, err := os.Open(arquivo)
	if err != nil {
		return
	}
	defer f.Close()
	sc := bufio.NewScanner(f)
	for sc.Scan() {
		linha := strings.TrimSpace(sc.Text())
		if linha == "" || strings.HasPrefix(linha, "#") {
			continue
		}
		chave, valor, ok := strings.Cut(linha, "=")
		if !ok {
			continue
		}
		chave = strings.TrimSpace(chave)
		valor = strings.TrimSpace(valor)
		if os.Getenv(chave) == "" {
			os.Setenv(chave, valor)
		}
	}
}

// abrirEmissor monta o emissor de tokens JWT a partir das env vars JWT_SECRET
// (obrigatória) e JWT_TTL (default 15m).
func abrirEmissor() *seguranca.Emissor {
	segredo := os.Getenv("JWT_SECRET")
	if segredo == "" {
		log.Fatal("JWT_SECRET é obrigatório (defina via env var; nunca hard-code)")
	}
	ttl := 15 * time.Minute
	if v := os.Getenv("JWT_TTL"); v != "" {
		d, err := time.ParseDuration(v)
		if err != nil {
			log.Fatalf("JWT_TTL inválido (%q): %v", v, err)
		}
		ttl = d
	}
	return seguranca.NovoEmissor(segredo, ttl)
}
