#!/usr/bin/env bun

import { readFileSync } from "node:fs";
import { dirname, resolve } from "node:path";
import { fileURLToPath } from "node:url";
import { meta as checkMeta, runCheck } from "./commands/check.ts";
import { meta as ciMeta, runCiGenerate } from "./commands/ci.ts";
import { meta as configMeta, runConfigShow } from "./commands/config-show.ts";
import { meta as doctorMeta, runDoctor } from "./commands/doctor.ts";
import { meta as infoMeta, runInfo } from "./commands/info.ts";
import { meta as initMeta, runInit } from "./commands/init.ts";
import { meta as newMeta, runNew } from "./commands/new.ts";
import { meta as obfuscateMeta, runObfuscate } from "./commands/obfuscate.ts";
import { meta as packageMeta, runPackage } from "./commands/package.ts";
import { meta as releaseMeta, runRelease } from "./commands/release.ts";
import { runScriptsInject, meta as scriptsMeta } from "./commands/scripts.ts";
import { loadConfig } from "./config.ts";
import { renderCommandHelp } from "./utils/command-meta.ts";
import { parseArgs } from "./utils/flags.ts";
import { log, pc, setLogLevel } from "./utils/logger.ts";

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

const COMMAND_META = {
  new: newMeta,
  init: initMeta,
  check: checkMeta,
  doctor: doctorMeta,
  info: infoMeta,
  config: configMeta,
  scripts: scriptsMeta,
  ci: ciMeta,
  package: packageMeta,
  obfuscate: obfuscateMeta,
  release: releaseMeta,
} as const;

function rootHelp(): string {
  const lines: string[] = [];
  lines.push(
    `${pc.bold("eco")} — toolkit de build, package, obfuscate e release para apps Bun`,
  );
  lines.push("");
  lines.push(`Uso:  ${pc.cyan("eco <comando> [flags]")}`);
  lines.push("");
  lines.push(pc.bold("Comandos:"));
  for (const meta of Object.values(COMMAND_META)) {
    lines.push(`  ${pc.cyan(meta.name.padEnd(12))} ${meta.description}`);
  }
  lines.push("");
  lines.push(pc.bold("Flags globais:"));
  lines.push(
    `  ${"-h, --help".padEnd(22)} Mostra esta mensagem ou ajuda do comando`,
  );
  lines.push(`  ${"-v, --version".padEnd(22)} Mostra a versão do eco`);
  lines.push(`  ${"--config <path>".padEnd(22)} Caminho customizado do config`);
  lines.push(
    `  ${"--platforms <list>".padEnd(22)} Override de plataformas (comma-separated)`,
  );
  lines.push(`  ${"--verbose".padEnd(22)} Saída detalhada`);
  lines.push(`  ${"--quiet".padEnd(22)} Silencia logs`);
  lines.push(
    `  ${"--dry-run".padEnd(22)} Mostra o que seria executado sem rodar`,
  );
  lines.push("");
  lines.push(pc.dim("Para detalhes de um comando: eco <comando> --help"));
  return lines.join("\n");
}

const args = parseArgs(process.argv.slice(2));

if (args.flags.verbose) setLogLevel("verbose");
if (args.flags.quiet) setLogLevel("silent");

if (args.flags.version) {
  console.log(readVersion());
  process.exit(0);
}

if (args.flags.help && !args.command) {
  console.log(rootHelp());
  process.exit(0);
}

const { command, flags, rest, positional } = args;

if (!command) {
  console.error(rootHelp());
  process.exit(1);
}

if (flags.help && command in COMMAND_META) {
  console.log(
    renderCommandHelp(COMMAND_META[command as keyof typeof COMMAND_META]),
  );
  process.exit(0);
}

try {
  if (command === "new") {
    const projectName = positional[0];
    if (!projectName) {
      console.error("Uso: eco new <nome-do-projeto> [--template=<tipo>]\n");
      console.error(renderCommandHelp(newMeta));
      process.exit(1);
    }
    const templateFlag = rest.find((f) => f.startsWith("--template="));
    const template = templateFlag?.slice("--template=".length) as
      | "cli-tool"
      | "library"
      | undefined;
    const force = rest.includes("--force");
    runNew({ projectName, template, force });
    process.exit(0);
  }

  if (command === "init") {
    const force = rest.includes("--force");
    await runInit({ force });
    process.exit(0);
  }

  if (command === "info") {
    await runInfo();
    process.exit(0);
  }

  if (command === "config") {
    const sub = positional[0];
    if (sub !== "show") {
      console.error(`Subcomando inválido: ${sub ?? "(vazio)"}\n`);
      console.error(renderCommandHelp(configMeta));
      process.exit(1);
    }
    const config = await loadConfig(flags.config);
    if (flags.platforms) {
      config.platforms = flags.platforms as typeof config.platforms;
    }
    runConfigShow(config);
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

  if (command === "doctor") {
    const config = await loadConfig(flags.config);
    if (flags.platforms) {
      config.platforms = flags.platforms as typeof config.platforms;
    }
    const fix = rest.includes("--fix");
    await runDoctor(config, { fix });
    process.exit(0);
  }

  if (command === "scripts") {
    const sub = positional[0];
    if (sub !== "inject") {
      console.error(`Subcomando inválido: ${sub ?? "(vazio)"}\n`);
      console.error(renderCommandHelp(scriptsMeta));
      process.exit(1);
    }
    const force = rest.includes("--force");
    runScriptsInject({ force });
    process.exit(0);
  }

  if (command === "ci") {
    const sub = positional[0];
    if (sub !== "generate") {
      console.error(`Subcomando inválido: ${sub ?? "(vazio)"}\n`);
      console.error(renderCommandHelp(ciMeta));
      process.exit(1);
    }
    const config = await loadConfig(flags.config);
    if (flags.platforms) {
      config.platforms = flags.platforms as typeof config.platforms;
    }
    const onlyFlag = rest.find((f) => f.startsWith("--only="));
    const only = onlyFlag?.slice("--only=".length) as
      | "ci"
      | "release"
      | undefined;
    const update = rest.includes("--update");
    runCiGenerate(config, { only, update });
    process.exit(0);
  }

  const config = await loadConfig(flags.config);
  if (flags.platforms) {
    config.platforms = flags.platforms as typeof config.platforms;
  }

  const opts = { dryRun: flags.dryRun };

  if (command === "package") {
    const watch = rest.includes("--watch");
    await runPackage(config, { ...opts, watch });
  } else if (command === "obfuscate") {
    await runObfuscate(config, opts);
  } else if (command === "release") {
    const skipObfuscate = rest.includes("--skip-obfuscate");
    const keepGoing = rest.includes("--keep-going");
    const parallel = !rest.includes("--no-parallel");
    await runRelease(config, { ...opts, skipObfuscate, keepGoing, parallel });
  } else {
    console.error(`Comando desconhecido: ${command}\n`);
    console.error(rootHelp());
    process.exit(1);
  }
} catch (err) {
  const message = err instanceof Error ? err.message : String(err);
  log.error(`\n❌ ${message}`);
  process.exit(1);
}
