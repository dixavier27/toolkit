#!/usr/bin/env bun

import { readFileSync } from "node:fs";
import { dirname, resolve } from "node:path";
import { fileURLToPath } from "node:url";
import { runCheck } from "./commands/check.ts";
import { runInit } from "./commands/init.ts";
import { runObfuscate } from "./commands/obfuscate.ts";
import { runPackage } from "./commands/package.ts";
import { runRelease } from "./commands/release.ts";
import { loadConfig } from "./config.ts";
import { parseArgs } from "./utils/flags.ts";
import { setLogLevel } from "./utils/logger.ts";

export type { EcoConfig, Platform } from "./config.ts";

function readVersion(): string {
  try {
    const here = dirname(fileURLToPath(import.meta.url));
    const pkgPath = resolve(here, "..", "package.json");
    const pkg = JSON.parse(readFileSync(pkgPath, "utf8"));
    return pkg.version ?? "0.0.0";
  } catch {
    return "0.0.0";
  }
}

const HELP = `eco — toolkit de build, package, obfuscate e release para apps Bun

Uso:  eco <comando> [flags]

Comandos:
  init                Cria eco.config.js no diretório atual
  check               Valida configuração e ambiente
  package             Gera o bundle JS único
  obfuscate           Ofusca o bundle (requer 'package' antes)
  release             Pipeline completo: package → obfuscate → binários nativos

Flags globais:
  -h, --help          Mostra esta mensagem
  -v, --version       Mostra a versão do eco
      --config <path> Caminho customizado do arquivo de config
      --platforms <l> Override de plataformas (comma-separated: linux,win,macos,macos-arm64)
      --verbose       Saída detalhada
      --quiet         Silencia logs (apenas erros)
      --dry-run       Mostra o que seria executado sem rodar

Flags específicas:
  release --skip-obfuscate    Pula a etapa de ofuscação

Exemplos:
  eco init
  eco check
  eco package --platforms=linux
  eco release --skip-obfuscate --verbose
  eco release --config=custom.config.js
`;

const args = parseArgs(process.argv.slice(2));

if (args.flags.version) {
  console.log(readVersion());
  process.exit(0);
}

if (args.flags.help && !args.command) {
  console.log(HELP);
  process.exit(0);
}

if (args.flags.verbose) setLogLevel("verbose");
if (args.flags.quiet) setLogLevel("silent");

const { command, flags, rest } = args;

if (!command) {
  console.error(HELP);
  process.exit(1);
}

try {
  if (command === "init") {
    const force = rest.includes("--force");
    await runInit({ force });
    process.exit(0);
  }

  if (command === "check") {
    const config = await loadConfig(flags.config);
    if (flags.platforms) {
      config.platforms = flags.platforms as typeof config.platforms;
    }
    await runCheck(config);
    process.exit(0);
  }

  const config = await loadConfig(flags.config);
  if (flags.platforms) {
    config.platforms = flags.platforms as typeof config.platforms;
  }

  const opts = { dryRun: flags.dryRun };

  if (command === "package") {
    await runPackage(config, opts);
  } else if (command === "obfuscate") {
    await runObfuscate(config, opts);
  } else if (command === "release") {
    const skipObfuscate = rest.includes("--skip-obfuscate");
    await runRelease(config, { ...opts, skipObfuscate });
  } else {
    console.error(`Comando desconhecido: ${command}\n`);
    console.error(HELP);
    process.exit(1);
  }
} catch (err) {
  const message = err instanceof Error ? err.message : String(err);
  console.error(`\n❌ ${message}`);
  process.exit(1);
}
