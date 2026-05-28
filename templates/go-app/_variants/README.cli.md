# {{name}}

CLI Go criada com [eco](https://dixavier27.github.io/eco/).

## Quickstart

```bash
go mod tidy
go run ./cmd/{{name}} hello mundo
```

## Comandos

| Make target | Descrição |
|---|---|
| `make dev` | Roda `hello` em modo de desenvolvimento |
| `make build` | Compila binário em `dist/` |
| `make test` | Executa testes |
| `make lint` | Roda golangci-lint |
| `make fmt` | Formata com gofmt + goimports |
| `make cover` | Gera relatório HTML de cobertura |
| `make tidy` | `go mod tidy` |

## Release

A partir do diretório do projeto, com o eco instalado:

```bash
eco release --platforms=linux,win,macos
```

Binários cross-compilados ficam em `release/`.

## Estrutura

```
cmd/{{name}}/    # entrypoint
internal/cli/    # subcomandos Cobra
VERSION          # versão embutida no binário via -ldflags
```

## Dependências dev

- `golangci-lint` — `go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest`
- `goimports` — `go install golang.org/x/tools/cmd/goimports@latest`
