# Eco GitHub Action

Composite action que roda comandos do [eco](../README.md) em workflows GitHub Actions com Bun pré-instalado.

## Uso

```yaml
- uses: dixavier27/toolkit/composite-action@v2.7.0
  with:
    command: release
    platforms: linux,win
    upload-artifacts: true
```

Substitui ~30 linhas de YAML (checkout, setup-bun, bun install, bunx eco, upload-artifact) por 3 linhas.

## Inputs

| Input | Default | Descrição |
|-------|---------|-----------|
| `command` | — (required) | Comando eco a rodar (ex: `release`, `package`, `check`, `release --keep-going`) |
| `platforms` | — | Override de plataformas (comma-separated). Passado como `--platforms=`. |
| `bun-version` | `1.2.0` | Versão do Bun a instalar |
| `working-directory` | `.` | Diretório onde rodar (útil em monorepos) |
| `install` | `true` | Rodar `bun install --frozen-lockfile` antes |
| `upload-artifacts` | `false` | Upload de `release/` como artefato do workflow |
| `artifact-name` | `eco-release` | Nome do artefato (se `upload-artifacts=true`) |

## Outputs

| Output | Descrição |
|--------|-----------|
| `release-dir` | Caminho absoluto para `release/` |

## Exemplo: workflow completo de release

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
      - uses: dixavier27/toolkit/composite-action@v2.7.0
        with:
          command: release
          platforms: linux,win
      - uses: softprops/action-gh-release@v2
        with:
          files: release/*
```

## Exemplo: matrix por plataforma

```yaml
jobs:
  release:
    strategy:
      matrix:
        include:
          - { os: ubuntu-latest,  platform: linux }
          - { os: windows-latest, platform: win }
          - { os: macos-latest,   platform: macos }
    runs-on: ${{ matrix.os }}
    steps:
      - uses: actions/checkout@v4
      - uses: dixavier27/toolkit/composite-action@v2.7.0
        with:
          command: release --keep-going
          platforms: ${{ matrix.platform }}
          upload-artifacts: true
          artifact-name: release-${{ matrix.platform }}
```

## Equivalência manual

A action acima é equivalente a:

```yaml
- uses: actions/checkout@v4
- uses: oven-sh/setup-bun@v2
  with: { bun-version: '1.2.0' }
- run: bun install --frozen-lockfile
- run: bunx eco release --platforms=linux,win
- uses: actions/upload-artifact@v4
  with:
    name: eco-release
    path: release/
```
