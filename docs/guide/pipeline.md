# Pipeline

O eco organiza o build em **3 etapas independentes**, que podem ser executadas em sequência ou isoladamente.

```
package  →  obfuscate  →  release
   ↓            ↓             ↓
bundle JS    ofusca       compila binários nativos
(1 arquivo)  inplace      (1 por plataforma)
```

## Etapa 1 — `package`

Gera **um único** bundle JavaScript:

```bash
eco package
# 📦 Bundled → dist/bundle.js (1.2 MB, 80ms)
```

Internamente: `bun build src/main.ts --outfile dist/bundle.js --target=bun --minify`.

Configura via:
- `entry` — arquivo de entrada
- `outDir` — diretório de saída
- `bundleName` — nome do bundle
- `sourcemap` — `false | 'inline' | 'external'`
- `define` — env vars / constantes injetadas no bundle
- `embedVersion` — injeta `__VERSION__` do `package.json` (default `true`)
- `assets` — cópia declarativa pós-bundle (`[{ from, to }]`)

## Etapa 2 — `obfuscate`

Ofusca o bundle gerado pela etapa anterior, **inplace**:

```bash
eco obfuscate
# 🔒 Obfuscated → dist/bundle.js (2.4 MB, 1.8s)
```

Internamente: `javascript-obfuscator dist/bundle.js --output dist/bundle.js --config obfuscator.config.cjs`.

::: warning Requer bundle prévio
`eco obfuscate` falha se `dist/bundle.js` não existir. Rode `eco package` antes, ou use `eco release` que orquestra tudo.
:::

## Etapa 3 — `release`

Compila o bundle ofuscado em **binários nativos** para cada plataforma:

```bash
eco release
# 📦 Bundled → dist/bundle.js (1.2 MB, 80ms)
# 🔒 Obfuscated → dist/bundle.js (2.4 MB, 1.8s)
#
# 🔧 Compilando 2 plataformas em paralelo…
#
# 🚀 linux    → release/app-linux (3.2s)
# 🚀 win      → release/app-win.exe (3.5s)
# 🔐 Checksums → release/checksums.txt
#
# Release pronto:
# Plataforma  Tamanho  Tempo  Status  Path
# ──────────  ───────  ─────  ──────  ────
# linux       52 MB    3.2s   ok      release/app-linux
# win         54 MB    3.5s   ok      release/app-win.exe
```

Internamente, para cada plataforma: `bun build dist/bundle.js --compile --target=<bun-target> --outfile release/<name>-<platform>`.

Flags úteis:
- `--skip-obfuscate` — pipeline sem ofuscação (debug)
- `--keep-going` — continua mesmo se uma plataforma falhar
- `--no-parallel` — execução sequencial
- `--platforms=linux,win` — override de plataformas

## Etapas independentes

Cada comando pode rodar isolado:

```bash
eco package                    # só bundle (rebuild rápido em dev)
eco package && eco obfuscate   # bundle + ofusca (sem compilar)
eco release --skip-obfuscate   # bundle + compila (debug, sem ofuscação)
```

## Customização via hooks

Veja [Hooks](/guide/hooks) para customizar cada etapa.
