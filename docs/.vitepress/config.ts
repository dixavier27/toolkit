import { defineConfig } from "vitepress";

export default defineConfig({
  title: "eco",
  description: "🌱 Toolkit de build, package, obfuscate e release para apps Bun",
  lang: "pt-BR",
  cleanUrls: true,
  base: "/toolkit/",
  lastUpdated: true,

  head: [
    ["meta", { name: "theme-color", content: "#22c55e" }],
    ["meta", { property: "og:title", content: "eco — toolkit Bun" }],
    [
      "meta",
      {
        property: "og:description",
        content:
          "Ecossistema de ferramentas para empacotar e distribuir aplicações Bun em múltiplas plataformas.",
      },
    ],
  ],

  themeConfig: {
    nav: [
      { text: "Guia", link: "/guide/getting-started" },
      { text: "Comandos", link: "/commands" },
      { text: "Templates", link: "/templates" },
      { text: "Referência", link: "/reference/config" },
      {
        text: "v2.8",
        items: [
          { text: "Changelog", link: "/changelog" },
          { text: "Migração v1→v2", link: "/guide/migration-v1-to-v2" },
          {
            text: "GitHub",
            link: "https://github.com/dixavier27/toolkit",
          },
        ],
      },
    ],

    sidebar: {
      "/guide/": [
        {
          text: "Introdução",
          items: [
            { text: "Início rápido", link: "/guide/getting-started" },
            { text: "Instalação", link: "/guide/installation" },
            { text: "Migração v1 → v2", link: "/guide/migration-v1-to-v2" },
          ],
        },
        {
          text: "Conceitos",
          items: [
            { text: "Pipeline", link: "/guide/pipeline" },
            { text: "Hooks", link: "/guide/hooks" },
          ],
        },
      ],
      "/commands": [
        {
          text: "Comandos",
          items: [{ text: "Todos os comandos", link: "/commands" }],
        },
      ],
      "/templates": [
        {
          text: "Templates",
          items: [{ text: "Todos os templates", link: "/templates" }],
        },
      ],
      "/reference/": [
        {
          text: "Referência",
          items: [
            { text: "Config (eco.config.js)", link: "/reference/config" },
            { text: "Hooks", link: "/reference/hooks" },
            { text: "Composite Action", link: "/reference/action" },
          ],
        },
      ],
    },

    socialLinks: [
      { icon: "github", link: "https://github.com/dixavier27/toolkit" },
    ],

    footer: {
      message:
        "Released under the <a href='https://github.com/dixavier27/toolkit/blob/master/LICENSE'>MIT License</a>",
      copyright: "Copyright © 2026 Affonso Xavier",
    },

    search: {
      provider: "local",
    },

    editLink: {
      pattern:
        "https://github.com/dixavier27/toolkit/edit/master/docs/:path",
      text: "Editar esta página no GitHub",
    },

    outline: {
      level: [2, 3],
      label: "Nesta página",
    },

    docFooter: {
      prev: "Anterior",
      next: "Próximo",
    },
  },
});
