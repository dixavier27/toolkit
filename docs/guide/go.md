# Aplicações Go

A partir da v2.9, o eco também cria e empacota apps **Go**. A flag `--go` no `eco new` muda a família de templates; `eco build` e `eco release` detectam automaticamente a linguagem do projeto (`go.mod` vs `package.json`) e despacham o pipeline correto.

## Quando usar

Use Go quando precisar de:

- Binário standalone compilado para múltiplas plataformas
- CLI de sistema (`cobra`), API REST de baixa latência, ou ambos no mesmo binário
- Ecossistema maduro para infraestrutura, observabilidade, networking

## Criar projeto

```bash
eco new minha-cli --go                              # CLI puro (default)
eco new minha-api --go --api                        # API REST stdlib
eco new meu-app --go --cli --api                    # CLI + API no mesmo binário
eco new meu-app --go --cli --api \
  --module=github.com/user/meu-app                  # módulo Go customizado
```

Sem flag adicional, `--go` assume `--cli`. Sem `--module`, default é `example.com/<nome>` (edite depois em `go.mod`).

## Estrutura gerada

```
meu-app/
  cmd/meu-app/main.go     # entrypoint
  internal/cli/           # subcomandos Cobra (se --cli)
  internal/api/           # servidor stdlib + validator (se --api)
  Makefile                # dev, build, test, lint, fmt, cover, tidy
  .air.toml               # hot-reload (se --api)
  .golangci.yml           # config mínima de lint
  VERSION                 # versão embutida no binário via -ldflags
  .github/workflows/      # ci.yml + release.yml
```

## Comandos

```bash
go mod tidy               # resolve dependências
make dev                  # CLI: roda hello | API: hot-reload com air
make build                # binário em dist/
make test                 # go test -v -cover
make lint                 # golangci-lint
make fmt                  # gofmt + goimports
```

## Build e release

`eco build` e `eco release` detectam `go.mod` e usam `go build` nativo. As mesmas flags do pipeline Bun valem:

```bash
eco build                                    # binário em dist/
eco release --platforms=linux,win,macos      # cross-compile em release/
eco release --keep-going --no-parallel       # tolerância a falha + sequencial
eco release --dry-run                        # mostra comandos sem executar
```

### Mapeamento `--platforms` → `GOOS`/`GOARCH`

| `--platforms` | `GOOS` | `GOARCH` | Sufixo |
|---|---|---|---|
| `linux` | linux | amd64 | `-linux` |
| `win` | windows | amd64 | `-win.exe` |
| `macos` | darwin | amd64 | `-macos` |
| `macos-arm64` | darwin | arm64 | `-macos-arm64` |

A versão é embutida no binário via `-ldflags "-X main.version=$(cat VERSION)"`. O arquivo `VERSION` na raiz é o equivalente Go ao `__VERSION__` do pipeline Bun.

## Diagnóstico

`eco info` mostra versões de Go, air e golangci-lint independentemente do projeto detectado.

`eco doctor` em projeto Go valida:

- `go.mod` presente
- Go ≥ 1.22 no PATH
- `air` instalado (warn + sugestão `go install github.com/air-verse/air@latest`)
- `golangci-lint` instalado (warn + sugestão de instalação)
- `.gitignore` cobrindo `dist/` e `release/`

`eco check` em projeto Go valida `go.mod`, plataformas vs host e Go ≥ 1.22.

## Dependências dev

Instaláveis via `go install`:

```bash
go install github.com/air-verse/air@latest
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
go install golang.org/x/tools/cmd/goimports@latest
```

## O que **não** está incluído

- **Docker/Dockerfile** — Go produz binário standalone; adicione conforme seu deploy
- **Devcontainer** — projeto desenvolve nativo com Go instalado
- **`eco package` e `eco obfuscate`** — específicos do pipeline Bun (binários Go já são compilados, sem JS para ofuscar)
- **Template `go-library`** — esta versão foca em CLI/API; biblioteca importável fica para iteração futura

## DX unificada com Bun

| Aspecto | Como funciona |
|---|---|
| Criação | `eco new <nome>` (com ou sem `--go`) |
| Diagnóstico | `eco info`, `doctor`, `check` detectam ambas linguagens |
| Build | `eco build` despacha por `go.mod`/`package.json` |
| Release | `eco release` despacha por `go.mod`/`package.json` |
| Saída | `dist/` (binário) e `release/` (cross-compile) idênticos |
| Flags globais | `--platforms`, `--keep-going`, `--no-parallel`, `--dry-run`, `--verbose`, `--quiet` valem para ambas |
