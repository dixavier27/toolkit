#!/usr/bin/env bun

import { runObfuscate } from "./commands/obfuscate.ts";
import { runPackage } from "./commands/package.ts";
import { runRelease } from "./commands/release.ts";
import { loadConfig } from "./config.ts";

export type { ToolkitConfig } from "./config.ts";

type CommandFn = (
  cfg: Awaited<ReturnType<typeof loadConfig>>,
  flags: string[],
) => Promise<void>;

const commands: Record<string, CommandFn> = {
  package: runPackage,
  obfuscate: runObfuscate,
  release: runRelease,
};

const command = process.argv[2];
const flags = process.argv.slice(3);

if (!command || !commands[command]) {
  console.error("Uso: biglaw-scripts <comando> [flags]");
  console.error(`Comandos disponíveis: ${Object.keys(commands).join(", ")}`);
  console.error("Flags:  --skip-obfuscate  (somente release)");
  process.exit(1);
}

const config = await loadConfig();
await commands[command](config, flags);
