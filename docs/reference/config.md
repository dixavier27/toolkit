# Referência: `eco.config.js`

Schema completo do arquivo de configuração. Todos os campos têm defaults sensatos — você pode começar com `export default {}` e ajustar conforme necessário.

## Campos básicos

| Campo | Tipo | Default | Descrição |
|-------|------|---------|-----------|
| `entry` | `string` | `'src/main.ts'` | Arquivo de entrada para o `bun build` |
| `outDir` | `string` | `'dist'` | Diretório de saída do bundle |
| `bundleName` | `string` | `'bundle.js'` | Nome do arquivo do bundle |
| `releaseName` | `string` | `'app'` | Prefixo dos binários (`<releaseName>-<platform>`) |
| `obfuscatorConfig` | `string` | `'obfuscator.config.cjs'` | Path da config do `javascript-obfuscator` |
| `platforms` | `Platform[]` | `['linux', 'win']` | Plataformas alvo do `release` |

## Plataformas suportadas

```ts
type Platform = 'linux' | 'win' | 'macos' | 'macos-arm64'
```

| Platform | Bun target | Extensão do binário |
|----------|------------|---------------------|
| `linux` | `bun-linux-x64` | — |
| `win` | `bun-windows-x64` | `.exe` |
| `macos` | `bun-darwin-x64` | — |
| `macos-arm64` | `bun-darwin-arm64` | — |

::: warning Cross-compile
Compilar `macos-*` em runners Linux/Windows pode falhar. Use `macos-latest` no CI. `eco doctor` alerta sobre isso.
:::

## Build features (v2.2+)

| Campo | Tipo | Default | Descrição |
|-------|------|---------|-----------|
| `sourcemap` | `false \| 'inline' \| 'external'` | `false` | Sourcemaps |
| `assets` | `Array<{ from: string; to: string }>` | `[]` | Cópia declarativa pós-bundle |
| `define` | `Record<string, string>` | `{}` | `--define key=value` para `bun build` |
| `embedVersion` | `boolean` | `true` | Injeta `__VERSION__` do `package.json` |
| `parallel` | `boolean` | `true` | Paralelização do `release` |
| `checksums` | `boolean` | `false` | Gera `release/checksums.txt` (SHA256) |

## Hooks (v2.0+)

| Campo | Tipo | Descrição |
|-------|------|-----------|
| `afterPackage` | `(cfg) => void \| Promise<void>` | Chamado após bundle (e em `release`) |
| `afterObfuscate` | `(cfg) => void \| Promise<void>` | Chamado após ofuscação (e em `release`) |
| `afterRelease` | `(cfg) => void \| Promise<void>` | Chamado após compilar plataformas |

Veja [Hooks](/guide/hooks) para exemplos.

## Exemplo completo

```js
// eco.config.js
import { cp } from 'node:fs/promises'

/** @type {import('@dixavier27/eco').EcoConfig} */
export default {
  // básicos
  entry: 'src/main.ts',
  outDir: 'dist',
  bundleName: 'bundle.js',
  releaseName: 'meu-app',
  obfuscatorConfig: 'obfuscator.config.cjs',
  platforms: ['linux', 'win', 'macos'],

  // build features
  sourcemap: 'external',
  assets: [
    { from: 'node_modules/@fastify/swagger-ui/static', to: 'dist/static' },
  ],
  define: {
    'process.env.NODE_ENV': '"production"',
  },
  embedVersion: true,
  parallel: true,
  checksums: true,

  // hooks
  afterRelease: async (cfg) => {
    console.log(`✅ Release ${cfg.releaseName} em ${cfg.platforms.length} plataformas`)
  },
}
```

## TypeScript

Use `.ts` em vez de `.js` para autocompletar nativo:

```ts
// eco.config.ts
import type { EcoConfig } from '@dixavier27/eco'

const config: EcoConfig = {
  entry: 'src/main.ts',
  platforms: ['linux', 'win'],
}

export default config
```

Bun resolve `eco.config.ts` direto sem build separado.

## Resolução do arquivo

O eco procura, em ordem:

1. Path passado via `--config <path>`
2. `eco.config.ts`
3. `eco.config.js`
4. `eco.config.mjs`

Se nenhum for encontrado, usa só os defaults.

## Validação

A config é validada por schema Zod ao carregar. Erros são reportados de forma legível:

```
Configuração inválida em /home/user/projeto/eco.config.js:
  - platforms.0: Invalid enum value. Expected 'linux' | 'win' | 'macos' | 'macos-arm64', received 'linux64'
  - releaseName: Expected string, received number

Consulte 'eco --help' ou rode 'eco check'.
```
