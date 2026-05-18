# Hooks

Hooks permitem injetar lógica customizada depois de cada etapa do pipeline.

## Disponíveis

| Hook | Quando dispara |
|------|----------------|
| `afterPackage` | Após `eco package` (e dentro de `release`) |
| `afterObfuscate` | Após `eco obfuscate` (e dentro de `release`) |
| `afterRelease` | Após `eco release` (depois de compilar todas plataformas) |

Cada hook recebe o `EcoConfig` resolvido como parâmetro.

## Assinatura

```ts
type EcoHook = (config: EcoConfig) => void | Promise<void>
```

## Exemplo: copiar assets dinâmicos

```js
// eco.config.js
import { cp } from 'node:fs/promises'

/** @type {import('@dixavier27/eco').EcoConfig} */
export default {
  entry: 'src/main.ts',
  afterPackage: async (cfg) => {
    await cp(
      'node_modules/@fastify/swagger-ui/static',
      `${cfg.outDir}/static`,
      { recursive: true }
    )
  },
}
```

::: tip Quando usar hook vs `assets`
O campo declarativo `assets: [{ from, to }]` resolve 90% dos casos. Use hook para lógica condicional (ex: copiar só em produção) ou transformações.
:::

## Exemplo: invalidar CDN após release

```js
import { $ } from 'bun'

export default {
  entry: 'src/main.ts',
  afterRelease: async (cfg) => {
    if (process.env.CI) {
      await $`aws cloudfront create-invalidation --distribution-id ABC --paths "/*"`
      console.log('CloudFront invalidado.')
    }
  },
}
```

## Exemplo: gerar release notes

```js
import { writeFile } from 'node:fs/promises'
import { $ } from 'bun'

export default {
  entry: 'src/main.ts',
  afterRelease: async (cfg) => {
    const log = await $`git log --oneline -20`.text()
    await writeFile(`release/CHANGELOG.txt`, log, 'utf8')
  },
}
```

## Exemplo: notificar Slack

```js
export default {
  entry: 'src/main.ts',
  afterRelease: async (cfg) => {
    if (!process.env.SLACK_WEBHOOK) return
    await fetch(process.env.SLACK_WEBHOOK, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({
        text: `🚀 Release ${cfg.releaseName} concluído em ${cfg.platforms.length} plataformas`,
      }),
    })
  },
}
```

## Ordem de execução

```
runRelease
  ├─ runPackage
  │   ├─ bun build → dist/bundle.js
  │   ├─ copia assets declarativos
  │   └─ afterPackage()       ← hook 1
  ├─ runObfuscate (se não --skip-obfuscate)
  │   ├─ javascript-obfuscator
  │   └─ afterObfuscate()     ← hook 2
  ├─ compila plataformas (paralelo)
  ├─ gera checksums (se configurado)
  └─ afterRelease()           ← hook 3
```

## Validação

Os hooks são validados pelo schema Zod no `loadConfig()`. Se o campo não for uma função, o erro é claro:

```
Configuração inválida em eco.config.js:
  - afterPackage: Expected function, received string
```
