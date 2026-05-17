---
"@dixavier27/eco": minor
---

## v2.2 — Build features

Novos campos de config para customizar o pipeline:

- **`sourcemap`** — `false | 'inline' | 'external'` (default: `false`)
- **`assets`** — cópia declarativa de arquivos/dirs (`[{ from, to }]`)
- **`define`** — env-var injection (passado como `--define` para `bun build`)
- **`embedVersion`** — injeta automaticamente `__VERSION__` do `package.json` (default: `true`)
- **`parallel`** — compila plataformas em paralelo no `release` (default: `true`)
- **`checksums`** — gera `release/checksums.txt` com SHA256 (default: `false`)

Novas flags do `release`:

- `--keep-going` — continua mesmo se uma plataforma falhar
- `--no-parallel` — força execução sequencial

Outros:

- Tabela final do `release` agora inclui coluna `Status` (ok / falhou)
- Erros por plataforma não interrompem outras se `--keep-going` ativo
