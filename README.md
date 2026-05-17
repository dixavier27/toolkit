# eco

🌱 Toolkit de build, package, obfuscate e release para apps Bun — um ecossistema de ferramentas para empacotar e distribuir aplicações em múltiplas plataformas.

## Instalação

```bash
bun add -D github:dixavier27/toolkit#v2.0.0
```

> O pacote é publicado como `@dixavier27/eco` (o repositório continua sendo `toolkit`).

## Início rápido

```bash
eco init        # cria eco.config.js
eco check       # valida configuração e ambiente
eco package     # gera o bundle
eco release     # pipeline completo: package → obfuscate → binários
```

## Configuração

Crie um `eco.config.js` (ou `.ts`, `.mjs`) na raiz do projeto:

```js
/** @type {import('@dixavier27/eco').EcoConfig} */
export default {
  entry:            'src/main.ts',
  outDir:           'dist',
  bundleName:       'bundle.js',
  releaseName:      'meu-app',
  obfuscatorConfig: 'obfuscator.config.cjs',
  platforms:        ['linux', 'win'],
}
```

Todos os campos têm defaults sensatos — você pode começar com `export default {}` e ajustar conforme necessário.

### Hooks de pipeline

Use hooks para customizar cada etapa do pipeline:

```js
import { cp } from 'node:fs/promises'

export default {
  entry: 'src/main.ts',
  // copia arquivos estáticos depois do bundle:
  afterPackage: async (cfg) => {
    await cp('node_modules/@fastify/swagger-ui/static',
             `${cfg.outDir}/static`,
             { recursive: true })
  },
}
```

Hooks disponíveis: `afterPackage`, `afterObfuscate`, `afterRelease`.

## Comandos

| Comando | Descrição |
|---------|-----------|
| `eco init` | Cria `eco.config.js` com defaults inferidos do `package.json` |
| `eco check` | Valida config e ambiente (entry existe, Bun no PATH, etc.) |
| `eco package` | Gera o bundle JS único |
| `eco obfuscate` | Ofusca o bundle (requer `package` antes) |
| `eco release` | Pipeline completo: package → obfuscate → binários nativos |

## Flags globais

| Flag | Descrição |
|------|-----------|
| `-h, --help` | Mostra ajuda |
| `-v, --version` | Mostra versão |
| `--config <path>` | Caminho customizado do config |
| `--platforms <list>` | Override de plataformas (`linux,win,macos,macos-arm64`) |
| `--verbose` | Saída detalhada |
| `--quiet` | Silencia logs |
| `--dry-run` | Mostra o que faria sem executar |

## Plataformas suportadas

- `linux` (x64)
- `win` (Windows x64, `.exe`)
- `macos` (x64)
- `macos-arm64` (Apple Silicon)

> ⚠️ Cross-compile de `macos-*` em runners Linux/Windows pode falhar. Use `macos-latest` no CI ou rode `eco check` para verificar.

## Integração no `package.json` do seu projeto

```json
"scripts": {
  "package":   "eco package",
  "obfuscate": "eco obfuscate",
  "release":   "eco release",
  "check":     "eco check"
}
```

## Requisitos

- [Bun](https://bun.sh) >= 1.2.0
