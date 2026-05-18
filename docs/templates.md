# Templates

O eco vem com 4 templates curados, scaffoldáveis via `eco new <nome> --template=<tipo>`.

## `cli-tool` (default)

CLI Bun pronta pra distribuir como binário.

```bash
eco new minha-cli
# ou explicitamente:
eco new minha-cli --template=cli-tool
```

**Stack:** TypeScript estrito + Bun + Biome.

**Estrutura gerada:**

```
minha-cli/
├── package.json       ← bin: minha-cli → dist/bundle.js
├── tsconfig.json
├── biome.json
├── eco.config.js      ← embedVersion: true, checksums: true
├── .gitignore
├── README.md
└── src/
    └── main.ts        ← --version, --help, exemplo 'hello'
```

**Próximos passos pós-scaffold:**

```bash
cd minha-cli
bun install
bun run dev hello                # executa CLI em watch
bunx eco ci generate             # opcional: workflows CI/CD
```

---

## `library`

Biblioteca TypeScript publicável no npm.

```bash
eco new minha-lib --template=library
```

**Stack:** TypeScript + Bun + `bun test` + tsc para `.d.ts`.

**Estrutura gerada:**

```
minha-lib/
├── package.json       ← exports map, types: dist/index.d.ts
├── tsconfig.json
├── biome.json
├── .gitignore
├── README.md
└── src/
    ├── index.ts       ← greet() bilingue de exemplo
    └── index.test.ts  ← bun test
```

**Próximos passos:**

```bash
cd minha-lib
bun install
bun test
bun run build         # bun build + tsc --emitDeclarationOnly
```

---

## `backend-fastify`

API Fastify + Zod, pronta para validação de schemas e env vars.

```bash
eco new minha-api --template=backend-fastify
```

**Stack:** Fastify 5 + `fastify-type-provider-zod` + Helmet + CORS + Biome.

**Estrutura gerada (13 arquivos):**

```
minha-api/
├── package.json
├── tsconfig.json
├── biome.json
├── eco.config.js
├── .env.example
├── .gitignore
├── README.md
├── src/
│   ├── main.ts                ← bootstrap
│   ├── server.ts              ← cria Fastify
│   ├── env.ts                 ← validação Zod das env vars
│   ├── plugins/
│   │   └── error-handler.ts   ← ZodError → 400
│   └── routes/
│       └── hello.ts           ← GET /hello/:name com schema
└── tests/
    └── hello.test.ts          ← bun test + app.inject()
```

**Próximos passos:**

```bash
cd minha-api
cp .env.example .env
bun install
bun run dev
curl http://localhost:3000/hello/mundo
# { "greeting": "Olá, mundo!" }
```

---

## `frontend-angular-tauri`

App desktop Angular 21 standalone + Tauri 2 + Tailwind + Playwright.

```bash
eco new meu-app --template=frontend-angular-tauri
```

**Stack:** Angular 21 (standalone) + Tauri 2 (Rust) + Tailwind 3 + Playwright + Biome.

**Estrutura gerada (22 arquivos):**

```
meu-app/
├── package.json
├── angular.json
├── tsconfig.json, tsconfig.app.json
├── biome.json
├── tailwind.config.js, postcss.config.js
├── playwright.config.ts
├── .gitignore
├── README.md
├── src/
│   ├── main.ts             ← bootstrapApplication
│   ├── index.html
│   ├── styles.css          ← @tailwind directives
│   └── app/
│       ├── app.config.ts   ← provideRouter
│       ├── app.routes.ts
│       └── app.component.ts ← standalone, signal counter
├── src-tauri/
│   ├── Cargo.toml
│   ├── tauri.conf.json
│   ├── build.rs
│   └── src/
│       ├── main.rs
│       └── lib.rs          ← command 'greet'
└── e2e/
    └── example.spec.ts     ← Playwright
```

**Próximos passos:**

```bash
cd meu-app
bun install
bun run dev                # web em http://localhost:4200
bun run tauri:dev          # janela desktop (requer Rust)
bun run e2e                # Playwright headless
```

::: warning Ícones do Tauri
Antes do primeiro `tauri:build`, gere ícones:

```bash
bunx @tauri-apps/cli icon caminho/para/logo.png
```

Os ícones não estão no template (binários PNG/ICO/ICNS).
:::

---

## Adicionar template próprio

O eco resolve templates por nome em `templates/<nome>/`. Para adicionar um template novo:

1. Crie `templates/meu-template/` com os arquivos
2. Use `{{name}}` como placeholder do nome do projeto
3. Adicione `meu-template` à whitelist em `src/commands/new.ts`
4. Documente aqui

PRs welcome — veja [contribuindo](https://github.com/dixavier27/eco/issues).
