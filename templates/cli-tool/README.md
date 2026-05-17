# {{name}}

CLI tool criado com [eco](https://github.com/dixavier27/toolkit).

## Desenvolvimento

```bash
bun install
bun run dev hello
```

## Build e release

```bash
bun run package          # gera dist/bundle.js
bun run release          # compila binários (linux + win)
```

Para configurar workflows CI/CD no GitHub:

```bash
bunx eco ci generate
```

## Estrutura

- `src/main.ts` — entry point do CLI
- `eco.config.js` — configuração do eco (plataformas, bundle, etc.)
- `obfuscator.config.cjs` — gerado por `eco doctor --fix` quando rodar `release`
