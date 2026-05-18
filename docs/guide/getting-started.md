# Início rápido

O eco gera projetos prontos pra rodar em **menos de 60 segundos**.

## Pré-requisitos

- [Bun](https://bun.sh) >= 1.2.0
- Git
- (Opcional, para Tauri) [Rust](https://rustup.rs/)

## Setup de um projeto novo

Use um template curado:

```bash
bunx @dixavier27/eco new minha-app --template=cli-tool
cd minha-app
bun install
bun run dev
```

Templates disponíveis:

| Template | Para que serve |
|----------|----------------|
| `cli-tool` | CLI distribuída como binário (default) |
| `library` | Biblioteca TypeScript publicável |
| `backend-fastify` | API Fastify + Zod |
| `frontend-angular-tauri` | App desktop Angular + Tauri |

## Setup em projeto existente

```bash
bun add -D github:dixavier27/toolkit#v2.8.0
eco init           # cria eco.config.js
eco scripts inject # adiciona "package", "obfuscate", "release", "check"
eco ci generate    # gera .github/workflows/ci.yml e release.yml
eco doctor --fix   # corrige problemas comuns automaticamente
```

## Primeiro release

```bash
bun run release    # gera bundle → ofusca → compila binários nativos
```

Saída em `release/`:

```
release/
├── meu-app-linux
├── meu-app-win.exe
└── checksums.txt
```

## Próximos passos

- [Pipeline](/guide/pipeline) — entenda as etapas package → obfuscate → release
- [Hooks](/guide/hooks) — customize cada etapa
- [Composite Action](/reference/action) — use no GitHub Actions
