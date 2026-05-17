---
"@dixavier27/eco": minor
---

## v2.5 — Template `backend-fastify`

Adiciona template completo para API Fastify, scaffoldável via `eco new minha-api --template=backend-fastify`.

### Conteúdo do template

- **Fastify 5** com type-provider-zod para validação tipo-segura de schemas
- **Helmet + CORS** configurados
- **Validação de env vars com Zod** em `src/env.ts` (falha cedo com mensagens legíveis)
- **Error handler global** que mapeia `ZodError` → 400, outros erros → status do erro ou 500
- **Logger Pino** com pino-pretty em desenvolvimento
- **Exemplo** `GET /health` + `GET /hello/:name` com schema Zod
- **Teste** em `tests/hello.test.ts` usando `app.inject()` (sem precisar de servidor real)
- **Integração com eco**: scripts `package`, `obfuscate`, `release`, `check` no `package.json`; `eco.config.js` com `embedVersion` e `checksums`
- **Biome v2** para lint+formato (sem ESLint+Prettier)

### Fluxo

```bash
bunx @dixavier27/eco new minha-api --template=backend-fastify
cd minha-api
cp .env.example .env
bun install
bun run dev
curl http://localhost:3000/hello/mundo  # { "greeting": "Olá, mundo!" }
```

### Próximas releases

- **v2.6** — template `frontend-angular-tauri` (Angular 21 + Tauri + Tailwind + Playwright)
- **v2.7** — ecossistema: Composite GitHub Action, code signing, docs site
