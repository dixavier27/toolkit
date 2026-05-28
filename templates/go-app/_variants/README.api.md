# {{name}}

API REST Go criada com [eco](https://dixavier27.github.io/eco/).

## Quickstart

```bash
go mod tidy
make dev                                  # sobe servidor com hot-reload (air)
curl http://localhost:3000/hello/mundo
curl http://localhost:3000/healthz
```

## Rotas

| Método | Rota | Descrição |
|---|---|---|
| GET | `/hello/{nome}` | Saudação com validação do nome |
| GET | `/healthz` | Health check com versão embutida |

## Comandos

| Make target | Descrição |
|---|---|
| `make dev` | Hot-reload via `air` |
| `make build` | Compila binário em `dist/` |
| `make test` | Executa testes |
| `make lint` | Roda golangci-lint |
| `make fmt` | Formata com gofmt + goimports |
| `make cover` | Gera relatório HTML de cobertura |
| `make tidy` | `go mod tidy` |

## Configuração

| Env | Default | Descrição |
|---|---|---|
| `ADDR` | `:3000` | Endereço de bind do servidor |

## Release

```bash
eco release --platforms=linux,win,macos
```

## Estrutura

```
cmd/{{name}}/    # entrypoint
internal/api/    # servidor, rotas e validação
VERSION          # versão embutida no binário via -ldflags
```

## Dependências dev

- `air` — `go install github.com/air-verse/air@latest`
- `golangci-lint` — `go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest`
- `goimports` — `go install golang.org/x/tools/cmd/goimports@latest`
