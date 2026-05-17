---
"@dixavier27/eco": minor
---

## v2.4 — Templates curados

Comando `eco new` scaffolda projetos completos a partir de templates curados. Vai do zero ao primeiro `bun run dev` em <30 segundos.

### Comando novo

- **`eco new <nome> [--template=<tipo>]`** — copia template, substitui `{{name}}`, deixa pronto para `bun install`

### Templates disponíveis

- **`cli-tool`** (default) — CLI Bun com `bin/`, `src/main.ts` exemplo, integração com `eco package`/`release`, biome, tsconfig estrito
- **`library`** — biblioteca TypeScript com `bun test`, declarações `.d.ts` via `tsc --emitDeclarationOnly`, exports map, exemplo `greet()`

### Fluxo

```bash
bunx @dixavier27/eco new minha-app
cd minha-app
bun install
bun run dev hello
```

### Próximas releases

- **v2.5** — templates pesados: `backend-fastify`, `frontend-angular-tauri` (extraídos do big-law)
- **v2.6** — ecossistema: Composite GitHub Action, code signing, docs
