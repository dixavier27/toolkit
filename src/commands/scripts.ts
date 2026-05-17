import { existsSync, readFileSync, writeFileSync } from "node:fs";
import { resolve } from "node:path";
import type { CommandMeta } from "../utils/command-meta.ts";
import { log, pc } from "../utils/logger.ts";

export const meta: CommandMeta = {
  name: "scripts",
  description:
    "Injeta scripts do eco no package.json (sem sobrescrever existentes)",
  flags: [
    {
      name: "inject",
      description: "Subcomando: adiciona scripts ao package.json",
    },
    { name: "--force", description: "Sobrescreve scripts existentes" },
  ],
  examples: ["eco scripts inject", "eco scripts inject --force"],
};

const ECO_SCRIPTS: Record<string, string> = {
  package: "eco package",
  obfuscate: "eco obfuscate",
  release: "eco release",
  check: "eco check",
};

export interface ScriptsOptions {
  force?: boolean;
}

export function runScriptsInject(opts: ScriptsOptions = {}) {
  const cwd = process.cwd();
  const pkgPath = resolve(cwd, "package.json");

  if (!existsSync(pkgPath)) {
    log.error("❌ package.json não encontrado no diretório atual.");
    process.exit(1);
  }

  const original = readFileSync(pkgPath, "utf8");
  const pkg = JSON.parse(original);
  pkg.scripts = pkg.scripts ?? {};

  const added: string[] = [];
  const skipped: string[] = [];
  const overwritten: string[] = [];

  for (const [name, command] of Object.entries(ECO_SCRIPTS)) {
    if (pkg.scripts[name]) {
      if (opts.force) {
        overwritten.push(name);
        pkg.scripts[name] = command;
      } else {
        skipped.push(name);
      }
    } else {
      added.push(name);
      pkg.scripts[name] = command;
    }
  }

  if (added.length === 0 && overwritten.length === 0) {
    log.info(pc.dim("Nenhum script adicionado — todos já existem."));
    log.info(
      pc.dim(`Skipped: ${skipped.join(", ")}. Use --force para sobrescrever.`),
    );
    return;
  }

  const indentMatch = original.match(/\n(\s+)"/);
  const indent = indentMatch?.[1]?.length ?? 2;
  writeFileSync(pkgPath, `${JSON.stringify(pkg, null, indent)}\n`, "utf8");

  if (added.length > 0) {
    log.success(
      `✅ ${added.length} script${added.length !== 1 ? "s" : ""} adicionado${added.length !== 1 ? "s" : ""}:`,
    );
    for (const name of added) {
      log.info(`   ${pc.cyan(name)}: ${pc.dim(ECO_SCRIPTS[name] ?? "")}`);
    }
  }

  if (overwritten.length > 0) {
    log.warn(
      `⚠️  ${overwritten.length} script${overwritten.length !== 1 ? "s" : ""} sobrescrito${overwritten.length !== 1 ? "s" : ""}: ${overwritten.join(", ")}`,
    );
  }

  if (skipped.length > 0) {
    log.info(pc.dim(`Skipped (já existiam): ${skipped.join(", ")}`));
  }
}
