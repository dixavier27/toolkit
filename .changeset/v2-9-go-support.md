---
"@dixavier27/eco": minor
---

## v2.9 — Suporte a aplicações Go

Expansão do eco para criar e empacotar apps Go com a mesma UX dos templates Bun. Pipeline unificado com detecção automática de linguagem.

### `eco new --go`

Flag componível com `--cli` e `--api`:

- `eco new minha-cli --go` — CLI puro (Cobra)
- `eco new minha-api --go --api` — API REST stdlib (Go 1.22+ pattern routing) + go-playground/validator
- `eco new meu-app --go --cli --api` — CLI + API no mesmo binário (subcomandos `hello` + `serve`)
- `--module=<path>` para path do módulo (default: `example.com/<nome>`)

### Template `templates/go-app/` componível

Um único template com variantes em `_variants/`. O scaffold filtra `internal/cli/`, `internal/api/` e escolhe a variante correta de `main.go`, `go.mod`, `Makefile` e `README.md` conforme as flags.

Tooling incluído: `air` (hot-reload), `golangci-lint`, `gofmt` + `goimports`, `go test -v -cover`, Makefile com targets `dev`/`build`/`test`/`lint`/`fmt`/`cover`/`tidy`, workflows `ci.yml` e `release.yml` com matriz GOOS/GOARCH.

### Pipeline unificado

- `eco build` — novo comando, detecta `go.mod` vs `package.json` e despacha (Bun → bundle JS; Go → `go build`)
- `eco release` — detecta linguagem e usa `go build` cross-compile para Go. Mesmas flags do Bun (`--platforms`, `--keep-going`, `--no-parallel`, `--dry-run`)
- `eco package` e `eco obfuscate` permanecem Bun-only com erro claro se rodados em projeto Go

### Diagnóstico Go-aware

- `eco info` mostra Go, air e golangci-lint além do Bun
- `eco doctor` em projeto Go valida `go.mod`, Go ≥ 1.22, air e golangci-lint (apenas avisos — sem `--fix` instalando binários)
- `eco check` em projeto Go valida `go.mod`, plataformas vs host e Go ≥ 1.22

### Versão embutida

Equivalente Go do `__VERSION__` do Bun: arquivo `VERSION` na raiz, lido pelo Makefile e injetado via `-ldflags "-X main.version=$(cat VERSION)"`.

### Decisões de escopo

Fora desta release (por opção explícita):

- Docker/Dockerfile (Go produz binário standalone)
- Devcontainer
- Template `go-library` (foco é CLI/API)
- Frameworks HTTP externos (chi/gin/echo) — stdlib Go 1.22+ basta
- Aliases pt-br para comandos do eco — discussão separada
