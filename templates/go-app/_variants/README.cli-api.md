# {{name}}

Aplicação Go criada com [eco](https://dixavier27.github.io/eco/), com **CLI** (Cobra) e **API REST** no mesmo binário.

## Quickstart

```bash
go mod tidy
go run ./cmd/{{name}} hello mundo       # subcomando CLI
go run ./cmd/{{name}} serve             # sobe servidor HTTP
curl http://localhost:3000/hello/mundo
```

Em desenvolvimento com hot-reload do servidor:

```bash
make dev                                # equivalente a `air -- serve`
```

## Subcomandos

| Subcomando | Descrição |
|---|---|
| `hello [nome]` | Imprime saudação no stdout |
| `serve` | Sobe o servidor HTTP |

## Rotas (`serve`)

| Método | Rota | Descrição |
|---|---|---|
| GET | `/hello/{nome}` | Saudação com validação do nome |
| GET | `/healthz` | Health check com versão embutida |

## Comandos

| Make target | Descrição |
|---|---|
| `make dev` | Hot-reload via `air` rodando `serve` |
| `make build` | Compila binário em `dist/` |
| `make test` | Executa testes |
| `make lint` | Roda golangci-lint |
| `make fmt` | Formata com gofmt + goimports |
| `make cover` | Gera relatório HTML de cobertura |
| `make tidy` | `go mod tidy` |

## Configuração

| Env | Default | Descrição |
|---|---|---|
| `ADDR` | `:3000` | Endereço de bind do servidor (subcomando `serve`) |

## Release

```bash
eco release --platforms=linux,win,macos
```

## Estrutura

```
cmd/{{name}}/    # entrypoint que registra `hello` + `serve`
internal/cli/    # subcomando hello
internal/api/    # servidor, rotas, validação e ServeCmd
VERSION          # versão embutida no binário via -ldflags
```

## Dependências dev

- `air` — `go install github.com/air-verse/air@latest`
- `golangci-lint` — `go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest`
- `goimports` — `go install golang.org/x/tools/cmd/goimports@latest`
