# Referência: Composite GitHub Action

Action reusable do eco para usar em workflows GitHub Actions sem reescrever 30 linhas de YAML.

## Localização

```yaml
uses: dixavier27/eco/composite-action@v2.8.0
```

## Inputs

| Input | Default | Descrição |
|-------|---------|-----------|
| `command` | (required) | Comando eco a rodar (ex: `release`, `package`, `check`, `release --keep-going`) |
| `platforms` | — | Override de plataformas (comma-separated). Passado como `--platforms=` |
| `bun-version` | `1.2.0` | Versão do Bun a instalar |
| `working-directory` | `.` | Diretório onde rodar (útil em monorepos) |
| `install` | `true` | Rodar `bun install --frozen-lockfile` antes |
| `upload-artifacts` | `false` | Upload de `release/` como artefato do workflow |
| `artifact-name` | `eco-release` | Nome do artefato (se `upload-artifacts=true`) |

## Outputs

| Output | Descrição |
|--------|-----------|
| `release-dir` | Caminho absoluto para `release/` |

## Exemplo: workflow básico

```yaml
name: Release

on:
  push:
    tags: ['v*']

jobs:
  release:
    runs-on: ubuntu-latest
    permissions:
      contents: write
    steps:
      - uses: actions/checkout@v4
      - uses: dixavier27/eco/composite-action@v2.8.0
        with:
          command: release
          platforms: linux,win
      - uses: softprops/action-gh-release@v2
        with:
          files: release/*
```

## Exemplo: matrix por plataforma (recomendado)

Permite cross-compile correto (macOS precisa rodar em `macos-latest`):

```yaml
jobs:
  release:
    strategy:
      fail-fast: false
      matrix:
        include:
          - { os: ubuntu-latest,  platform: linux }
          - { os: windows-latest, platform: win }
          - { os: macos-latest,   platform: macos }
          - { os: macos-latest,   platform: macos-arm64 }
    runs-on: ${{ matrix.os }}
    permissions:
      contents: write
    steps:
      - uses: actions/checkout@v4
      - uses: dixavier27/eco/composite-action@v2.8.0
        with:
          command: release --keep-going
          platforms: ${{ matrix.platform }}
          upload-artifacts: true
          artifact-name: release-${{ matrix.platform }}
```

::: tip Geração automática
`eco ci generate` produz exatamente esse padrão baseado no seu `eco.config.js`.
:::

## Exemplo: monorepo

```yaml
- uses: dixavier27/eco/composite-action@v2.8.0
  with:
    command: release
    working-directory: apps/backend
    platforms: linux,win
```

## Equivalência manual

A action é equivalente a:

```yaml
- uses: oven-sh/setup-bun@v2
  with: { bun-version: '1.2.0' }
- run: bun install --frozen-lockfile
  working-directory: ${{ inputs.working-directory }}
- run: bunx eco ${{ inputs.command }} --platforms=${{ inputs.platforms }}
  working-directory: ${{ inputs.working-directory }}
- uses: actions/upload-artifact@v4   # se upload-artifacts=true
  with:
    name: ${{ inputs.artifact-name }}
    path: ${{ inputs.working-directory }}/release/
```

## Versionamento

Use sempre tag exata (`@v2.8.0`) ou major fluida (`@v2`) para receber patches automaticamente:

```yaml
uses: dixavier27/eco/composite-action@v2.8.0  # tag exata (preferido em prod)
uses: dixavier27/eco/composite-action@v2      # major fluida (semver patch+minor)
uses: dixavier27/eco/composite-action@main    # ❌ NÃO recomendado
```
