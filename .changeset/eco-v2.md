---
"@dixavier27/eco": major
---

🌱 **Rebrand para `@dixavier27/eco`** — toolkit unificado de build/package/obfuscate/release para apps Bun.

## Breaking changes

- Pacote renomeado: `toolkit` → `@dixavier27/eco`
- CLI renomeado: `biglaw-scripts` → `eco`
- Config renomeado: `biglaw.config.js` → `eco.config.{js,ts,mjs}`
- `ToolkitConfig` → `EcoConfig`
- Bundle único `dist/bundle.js` em vez de N bundles idênticos por plataforma

## Novos comandos

- **`eco init`** — gera `eco.config.js` com defaults inferidos do `package.json`
- **`eco check`** — valida config + ambiente (entry, obfuscator, Bun, cross-compile)
- **`eco info`** — exibe versão do eco, Bun, javascript-obfuscator e config path
- **`eco config show`** — imprime config resolvido (com defaults) em JSON colorizado

## Novas features

- **Validação com Zod** — erros legíveis em vez de stack traces
- **Hooks de pipeline** — `afterPackage`, `afterObfuscate`, `afterRelease`
- **Help contextual** — `eco <comando> --help` mostra flags e exemplos específicos
- **Flags globais** — `--help`, `--version`, `--config`, `--platforms`, `--verbose`, `--quiet`, `--dry-run`
- **Suporte a `eco.config.ts`** — TypeScript nativo via Bun
- **Logger colorido** — cores semânticas (success verde, warn amarelo, error vermelho)
- **Spinner** durante operações longas
- **Tabela final no release** — plataforma | tamanho | tempo | path
- **Timing por etapa** — `📦 Bundled → dist/bundle.js (1.2 KB, 80ms)`
- **Watch mode** — `eco package --watch` rebundla em mudanças

## Dependências

- `javascript-obfuscator`: 4.x → **5.4.2**
- `@biomejs/biome`: 1.x → **2.4.15**
- `typescript`: 5.8 → **5.9.3**
- Adicionado: `zod` ^3.23, `picocolors` ^1.1

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

Renomear `biglaw.config.js` → `eco.config.js`. Use hook `afterPackage` para cópia de assets (substitui post-build commands ad-hoc).
