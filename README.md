# eco

🌱 Toolkit de build, package, obfuscate e release para apps Bun — um ecossistema de ferramentas para empacotar e distribuir aplicações em múltiplas plataformas.

## Instalação

```bash
bun add -D github:dixavier27/toolkit#v2.0.0
```

> O pacote é publicado como `@dixavier27/eco` (o repositório continua sendo `toolkit`).

## Início rápido

```bash
eco init        # cria eco.config.js inferindo do package.json
eco check       # valida configuração e ambiente
eco info        # exibe versões instaladas e config detectado
eco package     # gera o bundle
eco release     # pipeline completo: package → obfuscate → binários
```

## Comandos

| Comando | Descrição |
|---------|-----------|
| `eco new <nome>` | Cria projeto novo a partir de template (`--template=cli-tool\|library`) |
| `eco init` | Cria `eco.config.js` com defaults inferidos do `package.json` |
| `eco check` | Valida config e ambiente (entry, obfuscator, Bun, cross-compile) |
| `eco doctor` | Diagnóstico extendido com `--fix` para autocorrigir problemas |
| `eco info` | Mostra versão do eco, Bun, javascript-obfuscator e config path |
| `eco config show` | Imprime o config resolvido (com defaults) em JSON |
| `eco scripts inject` | Adiciona scripts do eco no `package.json` |
| `eco ci generate` | Gera `.github/workflows/ci.yml` e `release.yml` |
| `eco package` | Gera o bundle JS único |
| `eco obfuscate` | Ofusca o bundle (requer `package` antes) |
| `eco release` | Pipeline completo: package → obfuscate → binários nativos |
| `eco completion <shell>` | Emite script de autocomplete para `bash`, `zsh` ou `fish` |

Para detalhes de cada comando: `eco <comando> --help`.

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

Todos os campos têm defaults — comece com `export default {}` se quiser.

### Build features

```js
export default {
  entry: 'src/main.ts',
  // sourcemaps:
  sourcemap: 'external',  // false | 'inline' | 'external'

  // cópia declarativa de arquivos:
  assets: [
    { from: 'node_modules/@fastify/swagger-ui/static', to: 'dist/static' },
    { from: 'public', to: 'dist/public' },
  ],

  // env-var injection no bundle (usado como bun build --define):
  define: {
    'process.env.NODE_ENV': '"production"',
  },

  // versão do package.json é injetada automaticamente como __VERSION__:
  embedVersion: true,

  // paralelização das compilações no release (default true):
  parallel: true,

  // SHA256 checksums em release/checksums.txt:
  checksums: true,
}
```

### Hooks de pipeline

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

Hooks: `afterPackage`, `afterObfuscate`, `afterRelease`.

## Flags globais

| Flag | Descrição |
|------|-----------|
| `-h, --help` | Ajuda raiz ou específica de um comando |
| `-v, --version` | Mostra versão do eco |
| `--config <path>` | Caminho customizado do config |
| `--platforms <list>` | Override de plataformas (`linux,win,macos,macos-arm64`) |
| `--verbose` | Saída detalhada |
| `--quiet` | Silencia logs (apenas erros) |
| `--dry-run` | Mostra o que faria sem executar |

## Flags específicas

| Comando | Flag | Descrição |
|---------|------|-----------|
| `package` | `--watch` | Rebundla em mudanças (delega para `bun build --watch`) |
| `release` | `--skip-obfuscate` | Pula a etapa de ofuscação |
| `release` | `--keep-going` | Continua mesmo se uma plataforma falhar |
| `release` | `--no-parallel` | Compila em sequência (em vez de paralelo) |
| `init` | `--force` | Sobrescreve `eco.config.js` existente |

## Plataformas suportadas

- `linux` (x64)
- `win` (Windows x64, `.exe`)
- `macos` (x64)
- `macos-arm64` (Apple Silicon)

> ⚠️ Cross-compile de `macos-*` em runners não-Darwin pode falhar. Use `macos-latest` no CI ou rode `eco check` para validar.

## Integração no `package.json` do seu projeto

```json
"scripts": {
  "package":   "eco package",
  "obfuscate": "eco obfuscate",
  "release":   "eco release",
  "check":     "eco check"
}
```

## Autocomplete no shell

```bash
# Bash
eco completion bash > ~/.eco-completion.bash
echo 'source ~/.eco-completion.bash' >> ~/.bashrc

# Zsh
eco completion zsh > ~/.zsh/completions/_eco
# Garanta que ~/.zsh/completions/ está em $fpath

# Fish
eco completion fish > ~/.config/fish/completions/eco.fish
```

## GitHub Action (Composite)

Para usar o eco em workflows de outros projetos, sem reescrever 30 linhas de YAML:

```yaml
- uses: dixavier27/toolkit/composite-action@v2.7.0
  with:
    command: release
    platforms: linux,win
    upload-artifacts: true
```

Detalhes em [`composite-action/README.md`](composite-action/README.md).

## Requisitos

- [Bun](https://bun.sh) >= 1.2.0

## Setup de um projeto novo em 30 segundos

A partir do zero, com template curado:

```bash
bunx @dixavier27/eco new minha-app
cd minha-app
bun install
bun run dev hello
```

Templates disponíveis:

- **`cli-tool`** (default) — CLI Bun com `bin/`, build via `eco package`/`release`
- **`library`** — Biblioteca TypeScript com `bun test`, types `.d.ts` via `tsc`
- **`backend-fastify`** — API Fastify + Zod (helmet, cors, error handler, validação de env, exemplo `/hello/:name` com teste via `app.inject()`)
- **`frontend-angular-tauri`** — App desktop Angular 21 standalone + Tauri 2 + Tailwind + Playwright (component com signal, command Rust de exemplo `greet`, e2e do contador)

Em um projeto **já existente**:

```bash
cd meu-projeto-existente
bun add -D github:dixavier27/toolkit#v2.4.0
eco init                 # cria eco.config.js
eco scripts inject       # adiciona scripts no package.json
eco ci generate          # gera workflows ci.yml + release.yml
eco doctor --fix         # cria .gitignore e obfuscator.config.cjs
```

## Roadmap

- ✅ **v2.0** — Fundação: rename, Zod, hooks, comandos novos, flags globais, DX polish
- ✅ **v2.2** — Build features: sourcemaps, assets declarativos, define, embed de versão, paralelização, checksums
- ✅ **v2.3** — Platform engineering: `eco scripts inject`, `eco ci generate`, `eco doctor`
- ✅ **v2.4** — Templates: `eco new <nome> --template=cli-tool|library`
- ✅ **v2.5** — Template `backend-fastify` (Fastify + Zod + Vitest)
- ✅ **v2.6** — Template `frontend-angular-tauri` (Angular + Tauri + Tailwind)
- ✅ **v2.7** — Ecossistema: Composite GitHub Action + autocompletion (bash/zsh/fish)
- ⏳ **v2.8** — Code signing (Windows signtool, macOS codesign + notarization)
- ⏳ **v2.9** — Docs site (Astro/VitePress) + guias de migração
