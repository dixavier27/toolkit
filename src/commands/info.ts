import { readFileSync } from "node:fs";
import { dirname, resolve } from "node:path";
import { fileURLToPath } from "node:url";
import { $ } from "bun";
import { getConfigPath } from "../config.ts";
import type { CommandMeta } from "../utils/command-meta.ts";
import { log, pc } from "../utils/logger.ts";

export const meta: CommandMeta = {
  name: "info",
  description:
    "Mostra versão do eco, Bun, javascript-obfuscator e config detectado",
  examples: ["eco info"],
};

function readEcoVersion(): string {
  try {
    const here = dirname(fileURLToPath(import.meta.url));
    const pkgPath = resolve(here, "..", "package.json");
    const pkg = JSON.parse(readFileSync(pkgPath, "utf8"));
    return pkg.version ?? "desconhecida";
  } catch {
    return "desconhecida";
  }
}

async function checkCommand(cmd: string): Promise<string> {
  try {
    const result = await $`${{ raw: cmd }} --version`.quiet().text();
    return result.trim().split("\n")[0] ?? "instalado";
  } catch {
    return pc.red("ausente");
  }
}

export async function runInfo() {
  const ecoVersion = readEcoVersion();
  const bunVersion = await checkCommand("bun");
  const obfuscatorVersion = await checkCommand("javascript-obfuscator");
  const configPath = getConfigPath();

  log.info(pc.bold("eco — informações do ambiente"));
  log.info("");
  log.info(`  ${pc.cyan("eco")}                  ${pc.bold(ecoVersion)}`);
  log.info(`  ${pc.cyan("Bun")}                  ${bunVersion}`);
  log.info(`  ${pc.cyan("javascript-obfuscator")} ${obfuscatorVersion}`);
  log.info(
    `  ${pc.cyan("Plataforma host")}      ${process.platform} (${process.arch})`,
  );
  log.info(
    `  ${pc.cyan("Config detectado")}     ${configPath ?? pc.dim("nenhum (rode 'eco init')")}`,
  );
}
