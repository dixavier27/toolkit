import type { EcoConfig } from "../config.ts";
import type { CommandMeta } from "../utils/command-meta.ts";
import { detectLang } from "../utils/detect-lang.ts";
import { runBuildGo } from "./build-go.ts";
import { type RunOptions, runPackage } from "./package.ts";

export const meta: CommandMeta = {
  name: "build",
  description:
    "Compila o projeto detectando linguagem (Bun → bundle JS; Go → binário nativo)",
  flags: [
    { name: "--dry-run", description: "Mostra o comando sem executar" },
    {
      name: "--watch",
      description: "Bun: rebundla automaticamente em mudanças",
    },
  ],
  examples: ["eco build", "eco build --dry-run"],
};

export async function runBuild(config: EcoConfig, opts: RunOptions = {}) {
  const lang = detectLang();
  if (lang === "ambiguous") {
    throw new Error(
      "Diretório contém go.mod e package.json. eco build não sabe qual usar — separe os projetos ou rode 'eco package' (Bun) / build manual.",
    );
  }
  if (lang === "none") {
    throw new Error(
      "Não detectei projeto Bun nem Go no diretório atual. Crie um com 'eco new <nome>' (Bun) ou 'eco new <nome> --go'.",
    );
  }
  if (lang === "go") {
    await runBuildGo(config, { dryRun: opts.dryRun });
    return;
  }
  await runPackage(config, opts);
}
