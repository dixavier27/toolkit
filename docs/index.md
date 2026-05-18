---
layout: home

hero:
  name: "eco"
  text: "Toolkit Bun para múltiplas plataformas"
  tagline: Ecossistema de ferramentas para empacotar, ofuscar e distribuir aplicações Bun como binários nativos.
  image:
    src: /logo.svg
    alt: eco
  actions:
    - theme: brand
      text: Começar
      link: /guide/getting-started
    - theme: alt
      text: Ver no GitHub
      link: https://github.com/dixavier27/toolkit

features:
  - icon: 🌱
    title: Setup em 30 segundos
    details: "<code>eco new minha-app</code> e está pronto. Templates curados para CLI, library, backend Fastify e desktop Angular+Tauri."
  - icon: 📦
    title: Pipeline declarativo
    details: "Bundle → ofuscação → binários nativos. Tudo configurado em <code>eco.config.js</code> com validação Zod e hooks customizáveis."
  - icon: 🚀
    title: Multi-plataforma
    details: "Compila Linux, Windows, macOS (x64 + arm64) em paralelo via <code>bun --compile</code>. Checksums SHA256 automáticos."
  - icon: 🔍
    title: Diagnóstico inteligente
    details: "<code>eco doctor --fix</code> verifica e corrige problemas comuns. <code>eco check</code> valida config + ambiente antes de buildar."
  - icon: 🎨
    title: DX cuidada
    details: "Cores semânticas, spinner, tabelas, timing por etapa. Help contextual por comando. Autocomplete bash/zsh/fish."
  - icon: 🤝
    title: GitHub Actions nativo
    details: "Composite Action publicada — substitui 30 linhas de YAML por 3. <code>eco ci generate</code> cria workflows do zero."
---

## Início rápido

::: code-group

```bash [Novo projeto]
bunx @dixavier27/eco new minha-app
cd minha-app
bun install
bun run dev
```

```bash [Projeto existente]
bun add -D github:dixavier27/toolkit#v2.8.0
eco init           # cria eco.config.js
eco scripts inject # adiciona scripts no package.json
eco ci generate    # gera workflows GitHub Actions
eco doctor --fix   # corrige .gitignore + obfuscator config
```

:::

## Comandos principais

```bash
eco new <nome>     # scaffolda projeto a partir de template
eco init           # cria eco.config.js
eco check          # valida config + ambiente
eco doctor --fix   # diagnóstico com autofix
eco package        # gera bundle JS
eco obfuscate      # ofusca o bundle
eco release        # pipeline completo → binários nativos
```

Ver todos em [comandos](/commands).

## Por que eco?

| Sem eco | Com eco |
|---------|---------|
| Cada projeto reescreve build scripts no package.json | `eco scripts inject` |
| Cada repositório copia-cola CI/CD YAML de outro | `eco ci generate` ou `dixavier27/toolkit/composite-action@v2.8.0` |
| Tooling inconsistente entre projetos | Templates curados garantem padronização |
| Build de produção espalhado em flags de bun build | `eco.config.js` declarativo |
