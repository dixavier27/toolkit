# toolkit

CLI toolkit para build, package, obfuscate e release de apps Biglaw.

## Instalação

```bash
pnpm add -D github:dixavier27/toolkit
```

## Uso

Crie um `biglaw.config.js` na raiz do projeto:

```js
/** @type {import('toolkit/dist/index.js').ToolkitConfig} */
export default {
  entry: 'src/main.ts',
  outDir: 'dist',
  releaseName: 'meu-app',
  obfuscatorConfig: 'obfuscator.config.cjs',
  platforms: ['linux', 'win'],
}
```

Adicione os scripts no `package.json`:

```json
"package":   "biglaw-scripts package",
"obfuscate": "biglaw-scripts obfuscate",
"release":   "biglaw-scripts release"
```

## Comandos

| Comando | Descrição |
|---------|-----------|
| `biglaw-scripts package` | Gera bundles JS por plataforma via `bun build` |
| `biglaw-scripts obfuscate` | Roda package + ofusca com `javascript-obfuscator` |
| `biglaw-scripts release` | Roda obfuscate + compila binários nativos com `bun --compile` |

## Requisitos

- [Bun](https://bun.sh) >= 1.2.0
