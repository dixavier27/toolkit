// Package tupa fornece uma abstração mínima de servidor HTTP em Go usando
// apenas a stdlib (net/http 1.22+, com path params nativos no ServeMux).
//
// O ponto de partida é um Servidor: registre rotas com Rota, encadeie
// middlewares com Usar e suba com Iniciar (shutdown gracioso ao cancelar o
// contexto). Sem dependências externas — esse é o núcleo reutilizável.
package tupa

import (
	"context"
	"net/http"
	"time"
)

// Servidor envolve um http.ServeMux e a configuração de um http.Server.
type Servidor struct {
	mux         *http.ServeMux
	addr        string
	middlewares []Middleware
	timeoutPar  time.Duration // tempo máximo para o shutdown gracioso
}

// Opcao configura um Servidor na criação (Functional Options).
type Opcao func(*Servidor)

// ComTimeoutDeParada define quanto tempo Iniciar espera as conexões em voo
// terminarem durante o shutdown gracioso. Default: 10s.
func ComTimeoutDeParada(d time.Duration) Opcao {
	return func(s *Servidor) { s.timeoutPar = d }
}

// Novo cria um Servidor que escutará em addr (ex.: ":8080").
func Novo(addr string, opts ...Opcao) *Servidor {
	s := &Servidor{
		mux:        http.NewServeMux(),
		addr:       addr,
		timeoutPar: 10 * time.Second,
	}
	for _, opt := range opts {
		opt(s)
	}
	return s
}

// Rota registra um handler para o par método+padrão. O padrão segue a sintaxe
// do ServeMux 1.22+, então path params são suportados: Rota("GET",
// "/usuarios/{id}", h). metodo vazio casa qualquer método.
func (s *Servidor) Rota(metodo, padrao string, h http.HandlerFunc) {
	if metodo == "" {
		s.mux.HandleFunc(padrao, h)
		return
	}
	s.mux.HandleFunc(metodo+" "+padrao, h)
}

// Usar adiciona middlewares aplicados ao handler raiz, na ordem de registro
// (o primeiro a entrar é o mais externo).
func (s *Servidor) Usar(mw ...Middleware) {
	s.middlewares = append(s.middlewares, mw...)
}

// Handler devolve o http.Handler final (mux + middlewares aplicados). Útil
// para testes com httptest sem subir um socket.
func (s *Servidor) Handler() http.Handler {
	var h http.Handler = s.mux
	// Aplica de trás pra frente para que o primeiro registrado fique externo.
	for i := len(s.middlewares) - 1; i >= 0; i-- {
		h = s.middlewares[i](h)
	}
	return h
}

// Iniciar sobe o servidor e bloqueia até ctx ser cancelado, quando faz um
// shutdown gracioso limitado por ComTimeoutDeParada. Devolve nil em parada
// limpa; caso contrário, o erro do servidor.
func (s *Servidor) Iniciar(ctx context.Context) error {
	srv := &http.Server{
		Addr:    s.addr,
		Handler: s.Handler(),
	}

	errc := make(chan error, 1)
	go func() {
		err := srv.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			errc <- err
			return
		}
		errc <- nil
	}()

	select {
	case err := <-errc:
		return err
	case <-ctx.Done():
		shutCtx, cancel := context.WithTimeout(context.Background(), s.timeoutPar)
		defer cancel()
		return srv.Shutdown(shutCtx)
	}
}
