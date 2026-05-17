import { existsSync, readFileSync, writeFileSync } from "node:fs";
import { resolve } from "node:path";
import { log } from "../utils/logger.ts";

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

  log.info(`✅ Criado eco.config.js`);
  log.info(`   entry:       ${entry}`);
  log.info(`   releaseName: ${releaseName}`);
  log.info("");
  log.info("Próximos passos:");
  log.info("  1. Ajuste eco.config.js conforme necessário");
  log.info("  2. Rode 'eco check' para validar");
  log.info("  3. Rode 'eco package' para gerar o bundle");
}
