---
"@dixavier27/eco": minor
---

## v2.6 — Template `frontend-angular-tauri`

Adiciona template para apps desktop com Angular 21 standalone + Tauri 2 + Tailwind, scaffoldável via `eco new meu-app --template=frontend-angular-tauri`.

### Stack incluído

- **Angular 21** standalone (sem NgModule), com `provideRouter` e `signal()` no exemplo
- **Tauri 2** (Rust) com command `greet` de exemplo demonstrando bridge Angular ↔ Rust
- **Tailwind CSS 3** já configurado (postcss + tailwind.config.js)
- **Playwright** para e2e com `webServer` integrado
- **Biome 2** para lint+format (sem ESLint+Prettier)
- **TypeScript estrito** com `strictTemplates`, `noPropertyAccessFromIndexSignature`, etc.

### Conteúdo (22 arquivos)

```
meu-app/
├── package.json
├── angular.json
├── tsconfig.json / tsconfig.app.json
├── biome.json
├── tailwind.config.js / postcss.config.js
├── playwright.config.ts
├── README.md
├── .gitignore
├── src/
│   ├── main.ts                  ← bootstrapApplication
│   ├── index.html
│   ├── styles.css               ← @tailwind directives
│   └── app/
│       ├── app.config.ts        ← provideRouter
│       ├── app.routes.ts
│       └── app.component.ts     ← standalone, signal-based counter
├── src-tauri/
│   ├── Cargo.toml
│   ├── tauri.conf.json
│   ├── build.rs
│   └── src/
│       ├── main.rs
│       └── lib.rs               ← greet command
└── e2e/
    └── example.spec.ts          ← Playwright testando o contador
```

### Fluxo

```bash
bunx @dixavier27/eco new meu-app --template=frontend-angular-tauri
cd meu-app
bun install
bun run dev          # web em http://localhost:4200
bun run tauri:dev    # janela desktop (requer Rust)
bun run e2e          # Playwright
```

Antes do `tauri:build`: `bunx @tauri-apps/cli icon logo.png` para gerar os ícones requeridos.

### Próximas releases

- **v2.7** — ecossistema: Composite GitHub Action, code signing, docs site
