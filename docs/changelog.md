# Changelog

Histórico de releases do eco (`@dixavier27/eco`). Para o changelog completo gerado pelo changeset, veja [CHANGELOG.md no repo](https://github.com/dixavier27/eco/blob/main/CHANGELOG.md).

## v2.8 — Docs site

- 📚 Site de documentação com VitePress
- Deploy automático para GitHub Pages
- Guias: início rápido, instalação, migração v1→v2, pipeline, hooks
- Referência completa: comandos, templates, config, hooks, Composite Action

## v2.7 — Ecossistema

- 🤝 **Composite GitHub Action** `dixavier27/eco/composite-action@v2.7.0` — substitui 30 linhas de YAML por 3
- 🐚 `eco completion <bash|zsh|fish>` — autocomplete no shell
- Dogfooding: `eco ci generate` agora usa o Composite Action nos workflows gerados

## v2.6 — Template frontend-angular-tauri

- 🖥️ Template completo Angular 21 standalone + Tauri 2 + Tailwind + Playwright
- Bridge Rust ↔ Angular demonstrado com command `greet`

## v2.5 — Template backend-fastify

- 🌐 Template Fastify 5 + Zod + `fastify-type-provider-zod` + helmet + CORS
- Validação Zod de env vars
- Error handler global, teste com `app.inject()`

## v2.4 — Templates

- 🌱 **`eco new <nome> --template=<tipo>`** — scaffolda projetos novos
- Templates `cli-tool` e `library`
- Substituição de `{{name}}` em todos os arquivos

## v2.3 — Platform engineering

- 🔧 **`eco scripts inject`** — adiciona scripts no `package.json`
- ⚙️ **`eco ci generate`** — gera workflows GitHub Actions
- 🩺 **`eco doctor --fix`** — diagnóstico com autocorreção

## v2.2 — Build features

- 🗺️ Sourcemaps (`'inline' | 'external'`)
- 📁 `assets: [{ from, to }]` — cópia declarativa
- 🔧 `define` — env-var injection
- 🏷️ `embedVersion` — injeta `__VERSION__` automático
- ⚡ Paralelização do release (~3× mais rápido)
- 🔐 `checksums: true` — SHA256 automático
- `--keep-going`, `--no-parallel` no release

## v2.1 — DX polish

- 🎨 Cores semânticas via `picocolors`
- ⏳ Spinner braille durante operações longas
- 📊 Tabela final no `release` (plataforma | tamanho | tempo | path)
- ⏱️ Timing por etapa
- 📋 Help contextual (`eco <cmd> --help`)
- 🔍 `eco info` — versões + plataforma
- 🔧 `eco config show` — JSON colorizado do config resolvido
- 👀 `eco package --watch`

## v2.0 — Fundação (BREAKING)

- 🌱 **Rebrand**: `toolkit` → `@dixavier27/eco`, `biglaw-scripts` → `eco`
- Config: `biglaw.config.js` → `eco.config.{js,ts,mjs}` (suporta TypeScript nativo)
- Bundle único `dist/bundle.js` em vez de N idênticos por plataforma
- ✅ **Validação Zod** com erros legíveis
- 🪝 **Hooks**: `afterPackage`, `afterObfuscate`, `afterRelease`
- **Comandos novos**: `eco init`, `eco check`
- **Flags globais**: `--help`, `--version`, `--config`, `--platforms`, `--verbose`, `--quiet`, `--dry-run`
- Deps atualizadas: `javascript-obfuscator` v5, Biome v2, TypeScript 5.9

[Migração v1→v2](/guide/migration-v1-to-v2)

## v1.x — Histórico

Versões antigas (`toolkit` / `biglaw-scripts`). Recomenda-se migrar para v2.x.
