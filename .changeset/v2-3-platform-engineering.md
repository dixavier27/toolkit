---
"@dixavier27/eco": minor
---

## v2.3 — Platform engineering (scaffolding helpers)

Comandos novos que transformam o eco de "biblioteca de build" em **plataforma interna de engenharia**. Setup de um projeto novo em 30 segundos.

### Novos comandos

- **`eco scripts inject`** — adiciona scripts (`package`, `obfuscate`, `release`, `check`) no `package.json` sem sobrescrever os existentes. Use `--force` para sobrescrever.
- **`eco ci generate`** — gera `.github/workflows/ci.yml` e `release.yml` baseado no `eco.config.js`. Matrix automática por plataforma (linux→ubuntu, win→windows, macos*→macos). Use `--only=ci|release` ou `--update`.
- **`eco doctor`** — superset de `eco check` com `--fix` para autocorrigir:
  - Cria `.gitignore` com entries do eco se ausente
  - Adiciona `dist/`, `release/` ao `.gitignore` existente
  - Cria `obfuscator.config.cjs` com preset `medium` se ausente
  - Roda `scripts inject` se scripts faltam no `package.json`

### Setup de projeto novo

```bash
cd meu-projeto
bun add -D github:dixavier27/toolkit#v2.3.0
eco init                 # cria eco.config.js
eco scripts inject       # adiciona scripts no package.json
eco ci generate          # gera workflows
eco doctor --fix         # cria gitignore + obfuscator config
```

Resultado: projeto pronto para `bun run release` em qualquer ambiente, com CI funcionando.

### Não incluído (v2.4)

- `eco new <template>` — scaffolding completo a partir de templates curados (backend-fastify, frontend-angular-tauri, cli-tool, library) virá na próxima minor.
