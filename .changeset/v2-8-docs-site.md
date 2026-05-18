---
"@dixavier27/eco": minor
---

## v2.8 — Docs site

Site de documentação estático com VitePress, deployado automaticamente em GitHub Pages a cada push em `master`.

### URL

[dixavier27.github.io/toolkit](https://dixavier27.github.io/toolkit/)

### Conteúdo

- **Landing page** com hero, features, comparativo "sem eco / com eco"
- **Guia** (5 páginas): início rápido, instalação, migração v1→v2, pipeline, hooks
- **Comandos** — referência consolidada dos 12 comandos
- **Templates** — referência dos 4 templates curados
- **Referência técnica**: config (EcoConfig schema), hooks, Composite Action
- **Changelog** estruturado por versão (v2.0 → v2.8)
- **Search local** (Lunr) embutida no theme
- **Edit on GitHub** link em cada página
- **i18n** pt-BR

### Build e deploy

- `bun run docs:dev` — dev server (vitepress)
- `bun run docs:build` — gera `docs/.vitepress/dist/` (1.7 MB, ~11s)
- `bun run docs:preview` — serve build local
- `.github/workflows/docs.yml` — deploy automático para GitHub Pages

### Decisão de escopo

Code signing (originalmente planejado como v2.8) foi movido para "deferred" — sem consumidor real exigindo, a implementação seria especulativa. Será adicionado em release futura quando houver requisito concreto.
