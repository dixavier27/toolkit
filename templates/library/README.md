# {{name}}

TypeScript library scaffolded com [eco](https://github.com/dixavier27/eco).

## Desenvolvimento

```bash
bun install
bun test
bun run build
```

## Uso

```ts
import { greet } from '{{name}}'

greet('mundo')                   // "Olá, mundo!"
greet('world', { locale: 'en' }) // "Hello, world!"
```

## Estrutura

- `src/index.ts` — entry point com exports públicos
- `src/index.test.ts` — testes com `bun test`
- `dist/` — build de produção (gerado por `bun run build`)
