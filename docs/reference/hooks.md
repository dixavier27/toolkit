# Referência: Hooks

::: tip
Conceito explicado em [Guia · Hooks](/guide/hooks). Esta página foca em assinaturas e edge cases.
:::

## Assinatura

```ts
type EcoHook = (config: EcoConfig) => void | Promise<void>
```

O hook recebe o **EcoConfig resolvido** (com defaults aplicados), incluindo overrides de flags como `--platforms`.

## Ordem de execução

```
runRelease(config)
  ├─ runPackage(config)
  │   ├─ bun build → dist/bundle.js
  │   ├─ copia assets declarativos (config.assets)
  │   └─ ↪ afterPackage(config)
  ├─ runObfuscate(config)              [se não --skip-obfuscate]
  │   ├─ javascript-obfuscator
  │   └─ ↪ afterObfuscate(config)
  ├─ Promise.all(config.platforms.map(compileOne))
  ├─ writeChecksumsFile()              [se config.checksums]
  └─ ↪ afterRelease(config)
```

## Edge cases

### Hook em `--dry-run`

Hooks **não executam** em dry-run. O log apenas informa que seriam chamados.

### Hook lança erro

Se um hook lança erro, a etapa atual falha e a próxima não é executada. Em `--keep-going` (apenas para release), outras plataformas continuam compilando, mas o hook final é skipado.

### Hook async

Hooks podem retornar `void` ou `Promise<void>`. O eco aguarda a Promise antes de prosseguir.

### Hook em comandos individuais

| Comando | Hooks chamados |
|---------|----------------|
| `eco package` | `afterPackage` |
| `eco obfuscate` | `afterObfuscate` |
| `eco release` | `afterPackage`, `afterObfuscate`, `afterRelease` |
| `eco release --skip-obfuscate` | `afterPackage`, `afterRelease` |

## Tipos

```ts
import type { EcoConfig, EcoHook, Platform } from '@dixavier27/eco'

const hook: EcoHook = async (config) => {
  // config: EcoConfig com todos os campos resolvidos
  // config.platforms: Platform[]
  // config.entry: string
  // ...
}
```
