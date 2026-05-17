---
"@dixavier27/eco": major
---

🌱 **Rebrand para `@dixavier27/eco`** — toolkit unificado de build/package/obfuscate/release para apps Bun.

## Breaking changes

- Pacote renomeado: `toolkit` → `@dixavier27/eco`
- CLI renomeado: `biglaw-scripts` → `eco`
- Config renomeado: `biglaw.config.js` → `eco.config.{js,ts,mjs}` (suporta TS nativo)
- `ToolkitConfig` → `EcoConfig`
- Bundle único `dist/bundle.js` em vez de N bundles idênticos por plataforma

## Novidades

- **Comando `eco init`** — gera `eco.config.js` com defaults inferidos do `package.json`
- **Comando `eco check`** — valida config + ambiente (entry, obfuscator no PATH, Bun >= 1.2.0, cross-compile)
- **Validação com Zod** — erros legíveis em vez de stack traces crus
- **Hooks de pipeline** — `afterPackage`, `afterObfuscate`, `afterRelease` para customizar cada etapa
- **Flags globais** — `--help`, `--version`, `--config`, `--platforms`, `--verbose`, `--quiet`, `--dry-run`
- **Logger com níveis** — silent / normal / verbose

## Dependências atualizadas

- `javascript-obfuscator`: 4.x → **5.4.2**
- `@biomejs/biome`: 1.x → **2.4.15**
- `typescript`: 5.8 → **5.9.3**
- Adicionado: `zod` ^3.23

## Migração de v1.x

```diff
- "toolkit": "github:dixavier27/toolkit#v1.0.1"
+ "@dixavier27/eco": "github:dixavier27/toolkit#v2.0.0"

- "package":   "biglaw-scripts package"
- "obfuscate": "biglaw-scripts obfuscate"
- "release":   "biglaw-scripts release"
+ "package":   "eco package"
+ "obfuscate": "eco obfuscate"
+ "release":   "eco release"
```

Renomear `biglaw.config.js` → `eco.config.js`. Hook `afterPackage` substitui post-build commands de cópia de assets.
