import { existsSync, readFileSync, writeFileSync } from "node:fs";
import { resolve } from "node:path";
import type { CommandMeta } from "../utils/command-meta.ts";
import { log, pc } from "../utils/logger.ts";

export const meta: CommandMeta = {
  name: "init",
  description: "Cria eco.config.js com defaults inferidos do package.json",
  flags: [{ name: "--force", description: "Sobrescreve config existente" }],
  examples: ["eco init", "eco init --force"],
};

const TEMPLATE = (entry: string, releaseName: string) =>
  `/** @type {import('@dixavier27/eco').EcoConfig} */
export default {
  entry: '${entry}',
  outDir: 'dist',
  bundleName: 'bundle.js',
  releaseName: '${releaseName}',
  obfuscatorConfig: 'obfuscator.config.cjs',
  platforms: ['linux', 'win'],
}
`;

function inferEntry(cwd: string): string {
  const pkgPath = resolve(cwd, "package.json");
  if (existsSync(pkgPath)) {
    try {
      const pkg = JSON.parse(readFileSync(pkgPath, "utf8"));
      const candidate = pkg.module ?? pkg.main;
      if (typeof candidate === "string" && candidate.endsWith(".ts")) {
        return candidate;
      }
    } catch {}
  }
  return "src/main.ts";
}

function inferReleaseName(cwd: string): string {
  const pkgPath = resolve(cwd, "package.json");
  if (existsSync(pkgPath)) {
    try {
      const pkg = JSON.parse(readFileSync(pkgPath, "utf8"));
      if (typeof pkg.name === "string") {
        return pkg.name.replace(/^@[^/]+\//, "").replace(/[^a-z0-9-]/gi, "-");
      }
    } catch {}
  }
  return "app";
}

export interface InitOptions {
  force?: boolean;
}

export async function runInit(opts: InitOptions = {}) {
  const cwd = process.cwd();
  const target = resolve(cwd, "eco.config.js");

  if (existsSync(target) && !opts.force) {
    log.error(`❌ ${target} já existe. Use --force para sobrescrever.`);
    process.exit(1);
  }

  const entry = inferEntry(cwd);
  const releaseName = inferReleaseName(cwd);

  writeFileSync(target, TEMPLATE(entry, releaseName), "utf8");

  log.success("✅ Criado eco.config.js");
  log.info(`   ${pc.dim("entry:")}       ${entry}`);
  log.info(`   ${pc.dim("releaseName:")} ${releaseName}`);
  log.info("");
  log.info(pc.bold("Próximos passos:"));
  log.info(
    `  ${pc.dim("1.")} Ajuste ${pc.cyan("eco.config.js")} conforme necessário`,
  );
  log.info(`  ${pc.dim("2.")} Rode ${pc.cyan("eco check")} para validar`);
  log.info(
    `  ${pc.dim("3.")} Rode ${pc.cyan("eco package")} para gerar o bundle`,
  );
}
