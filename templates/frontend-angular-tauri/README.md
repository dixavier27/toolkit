# {{name}}

App desktop Angular 21 + Tauri 2 + Tailwind, scaffolded com [eco](https://github.com/dixavier27/toolkit).

## Pré-requisitos

- [Bun](https://bun.sh) >= 1.2.0
- [Rust](https://rustup.rs/) (para Tauri)
- Plataforma-específicas para Tauri: ver [pré-requisitos do Tauri](https://tauri.app/start/prerequisites/)

## Desenvolvimento

### Web (browser)

```bash
bun install
bun run dev          # http://localhost:4200
```

### Desktop (Tauri)

```bash
bun run tauri:dev    # janela nativa + dev server
```

## Build

### Web

```bash
bun run build        # gera dist/
```

### Desktop

```bash
bun run tauri:build  # binário nativo em src-tauri/target/release/
```

> ⚠️ Antes do primeiro `tauri:build`, gere ícones: `bunx @tauri-apps/cli icon path/para/logo.png`.
> Isso cria os arquivos em `src-tauri/icons/` (PNG, ICO, ICNS) requeridos pelo build.

## Testes E2E

```bash
bun run e2e          # roda Playwright headless
bun run e2e:ui       # modo interativo
```

## Estrutura

```
src/
  main.ts              ← bootstrap Angular
  index.html           ← shell HTML
  styles.css           ← Tailwind diretivas + estilos globais
  app/
    app.config.ts      ← providers (router, etc.)
    app.routes.ts      ← rotas
    app.component.ts   ← root component (standalone)
src-tauri/
  Cargo.toml           ← deps Rust
  tauri.conf.json      ← config Tauri (janela, bundle, security)
  src/
    main.rs            ← entry binário
    lib.rs             ← setup Tauri (commands, plugins)
  build.rs             ← build script
e2e/
  example.spec.ts      ← teste Playwright
```

## Bridge Angular ↔ Rust

Adicione comandos Tauri em `src-tauri/src/lib.rs` com `#[tauri::command]`, depois invoke do Angular:

```ts
import { invoke } from '@tauri-apps/api/core'

const result = await invoke<string>('meu_comando', { arg: 'valor' })
```
