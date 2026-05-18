# Migração v1.x → v2.x

A v2 trouxe um rebranding completo. Se você tem um projeto usando a v1.x (`toolkit` / `biglaw-scripts`), siga este guia.

## Quebras

| v1.x | v2.x |
|------|------|
| Pacote `toolkit` | `@dixavier27/eco` |
| CLI `biglaw-scripts` | `eco` |
| Config `biglaw.config.js` | `eco.config.{js,ts,mjs}` |
| Type `ToolkitConfig` | `EcoConfig` |
| N bundles por plataforma (`bundle-linux.js`, `bundle-win.js`, ...) | 1 bundle único (`bundle.js`) |

## Passo a passo

### 1. `package.json`

```diff
- "toolkit": "github:dixavier27/toolkit#v1.0.1"
+ "@dixavier27/eco": "github:dixavier27/toolkit#v2.8.0"
```

E os scripts:

```diff
- "package":   "biglaw-scripts package"
- "obfuscate": "biglaw-scripts obfuscate"
- "release":   "biglaw-scripts release"
+ "package":   "eco package"
+ "obfuscate": "eco obfuscate"
+ "release":   "eco release"
```

Ou rode o autofix:

```bash
eco scripts inject --force
```

### 2. Config

Renomeie `biglaw.config.js` → `eco.config.js`:

```bash
mv biglaw.config.js eco.config.js
```

E ajuste o JSDoc type reference:

```diff
- /** @type {import('toolkit/dist/index.js').ToolkitConfig} */
+ /** @type {import('@dixavier27/eco').EcoConfig} */
```

### 3. Hooks no lugar de pós-build inline

Se você tinha cópia de assets em scripts shell, migre para hook:

```js
// eco.config.js
import { cp } from 'node:fs/promises'

export default {
  entry: 'src/main.ts',
  // antes: cp ... no script bun:build:prod
  afterPackage: async (cfg) => {
    await cp(
      'node_modules/@fastify/swagger-ui/static',
      `${cfg.outDir}/static`,
      { recursive: true }
    )
  },
}
```

Ou use o campo declarativo `assets` (v2.2+):

```js
export default {
  entry: 'src/main.ts',
  assets: [
    { from: 'node_modules/@fastify/swagger-ui/static', to: 'dist/static' },
  ],
}
```

### 4. Workflows GitHub Actions

Substitua YAML manual pelo Composite Action:

```diff
  steps:
    - uses: actions/checkout@v4
-   - uses: oven-sh/setup-bun@v2
-     with: { bun-version: '1.2.0' }
-   - run: bun install --frozen-lockfile
-   - run: bunx biglaw-scripts release
+   - uses: dixavier27/toolkit/composite-action@v2.8.0
+     with:
+       command: release
+       platforms: linux,win
```

Ou regenere do zero:

```bash
eco ci generate --update
```

### 5. Validação

```bash
eco doctor
```

Verifica config, entry, ambiente. Use `--fix` para autocorrigir warnings.

## Novidades que valem aproveitar

- **`eco init`** — gera config inferindo do `package.json`
- **`eco new <template>`** — scaffolda projetos novos
- **Hooks** — `afterPackage`, `afterObfuscate`, `afterRelease`
- **Validação Zod** — erros legíveis no config
- **Flags** `--platforms`, `--config`, `--dry-run`, `--verbose`, `--quiet`
- **Bundle único** + paralelização do release (~3× mais rápido)
- **Checksums SHA256** opcionais (`checksums: true` no config)
- **Composite Action** GitHub
- **Autocomplete** bash/zsh/fish
