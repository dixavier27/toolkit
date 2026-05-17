# {{name}}

API Fastify scaffolded com [eco](https://github.com/dixavier27/toolkit).

## Desenvolvimento

```bash
cp .env.example .env
bun install
bun run dev
```

Servidor sobe em `http://localhost:3000`. Rotas disponíveis:

- `GET /health` — health check
- `GET /hello/:name` — exemplo com validação Zod

## Estrutura

```
src/
  main.ts            ← entry point
  server.ts          ← cria a instância Fastify
  env.ts             ← validação de env vars com Zod
  plugins/
    error-handler.ts ← handler global de erros
  routes/
    hello.ts         ← rotas de exemplo
tests/
  hello.test.ts      ← exemplo de teste com fastify.inject()
```

## Build e release

```bash
bun run test
bun run typecheck
bun run lint
bun run release       # gera binários nativos (linux + win)
```

Para configurar CI/CD no GitHub:

```bash
bunx eco ci generate
```

## Próximos passos

1. Adicione mais rotas em `src/routes/`
2. Para conectar banco: instale `mongoose` (ou outro driver), crie um plugin em `src/plugins/`
3. Para auth: instale `@fastify/jwt`, crie middleware
4. Estruture módulos como `src/modules/<feature>/{domain,application,infrastructure,presentation}` quando o app crescer
