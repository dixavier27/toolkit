# Comandos

Referência completa de todos os 12 comandos do eco. Cada comando aceita `--help` para detalhes contextuais.

## Setup

### `eco new <nome>`

Cria projeto novo a partir de template curado.

```bash
eco new minha-cli                                  # default: cli-tool
eco new minha-lib  --template=library
eco new minha-api  --template=backend-fastify
eco new meu-app    --template=frontend-angular-tauri
eco new existing   --force                         # sobrescreve diretório
```

### `eco init`

Cria `eco.config.js` no diretório atual com defaults inferidos do `package.json`.

```bash
eco init
eco init --force   # sobrescreve config existente
```

### `eco scripts inject`

Adiciona scripts (`package`, `obfuscate`, `release`, `check`) no `package.json` sem sobrescrever os existentes.

```bash
eco scripts inject
eco scripts inject --force   # sobrescreve scripts existentes
```

### `eco ci generate`

Gera `.github/workflows/ci.yml` e `release.yml` baseado no `eco.config.js`.

```bash
eco ci generate
eco ci generate --only=ci        # só ci.yml
eco ci generate --only=release   # só release.yml
eco ci generate --update         # sobrescreve workflows existentes
```

Os workflows gerados usam o [Composite Action](/reference/action).

## Pipeline

### `eco package`

Gera **um único** bundle JS em `dist/bundle.js`.

```bash
eco package
eco package --watch              # rebundla em mudanças
eco package --dry-run            # preview sem executar
eco package --platforms=linux    # ignorado nesta etapa (single bundle)
```

### `eco obfuscate`

Ofusca o bundle gerado por `eco package`, **inplace**.

```bash
eco obfuscate
```

::: warning Requer `eco package` antes
Falha se `dist/bundle.js` não existir.
:::

### `eco release`

Pipeline completo: `package` → `obfuscate` → compila binários nativos.

```bash
eco release
eco release --skip-obfuscate                # pipeline sem ofuscação (debug)
eco release --keep-going                    # continua mesmo se uma plataforma falhar
eco release --no-parallel                   # execução sequencial
eco release --platforms=linux,win,macos     # override de plataformas
eco release --dry-run                       # preview sem executar
```

## Diagnóstico

### `eco check`

Valida config + ambiente. Sai com código 1 se houver falhas críticas.

```bash
eco check
```

Verifica:
- `eco.config.js` presente e válido (schema Zod)
- `entry` aponta para arquivo existente
- `obfuscatorConfig` aponta para arquivo existente
- Bun >= 1.2.0 no PATH
- `javascript-obfuscator` no PATH
- Plataformas compatíveis com o host (alerta de cross-compile macOS)

### `eco doctor`

Superset do `check` com **autofix**:

```bash
eco doctor          # diagnóstico
eco doctor --fix    # aplica correções automáticas
```

Pode corrigir:
- `.gitignore` ausente ou sem `dist/`/`release/`
- `obfuscator.config.cjs` ausente (cria com preset `medium`)
- Scripts faltando no `package.json`

### `eco info`

Mostra versões instaladas, plataforma do host e path do config.

```bash
eco info
```

Saída:

```
eco — informações do ambiente

  eco                  2.8.0
  Bun                  1.3.14
  javascript-obfuscator 5.4.2
  Plataforma host      linux (x64)
  Config detectado     /home/user/projeto/eco.config.js
```

### `eco config show`

Imprime o config resolvido (com defaults aplicados) em JSON colorizado.

```bash
eco config show
eco config show --config=custom.config.js
```

## Outros

### `eco completion <shell>`

Emite script de autocomplete para `bash`, `zsh` ou `fish`.

```bash
eco completion bash > ~/.eco-completion.bash
eco completion zsh  > ~/.zsh/completions/_eco
eco completion fish > ~/.config/fish/completions/eco.fish
```

## Flags globais

| Flag | Descrição |
|------|-----------|
| `-h, --help` | Ajuda raiz ou específica do comando |
| `-v, --version` | Mostra versão do eco |
| `--config <path>` | Caminho customizado do config |
| `--platforms <list>` | Override de plataformas (comma-separated) |
| `--verbose` | Saída detalhada |
| `--quiet` | Silencia logs |
| `--dry-run` | Mostra o que faria sem executar |
