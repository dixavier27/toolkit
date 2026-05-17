import { mkdirSync, statSync } from "node:fs";
import { $ } from "bun";
import type { EcoConfig } from "../config.ts";
import type { CommandMeta } from "../utils/command-meta.ts";
import { log, pc } from "../utils/logger.ts";
import { spinner } from "../utils/spinner.ts";
import { fileSize } from "../utils/table.ts";
import { formatDuration, timed } from "../utils/timing.ts";

export const meta: CommandMeta = {
  name: "package",
  description: "Gera o bundle JS único em outDir/bundleName",
  flags: [
    { name: "--watch", description: "Rebundla automaticamente em mudanças" },
    { name: "--dry-run", description: "Mostra o comando sem executar" },
  ],
  examples: [
    "eco package",
    "eco package --watch",
    "eco package --platforms=linux",
  ],
};

export interface RunOptions {
  dryRun?: boolean;
  watch?: boolean;
}

export async function runPackage(config: EcoConfig, opts: RunOptions = {}) {
  const bundlePath = `${config.outDir}/${config.bundleName}`;

  if (opts.dryRun) {
    log.info(pc.cyan(`📦 [dry-run] Bundle → ${bundlePath}`));
    log.dim(
      `   bun build ${config.entry} --outfile ${bundlePath} --target=bun --minify`,
    );
    return;
  }

  mkdirSync(config.outDir, { recursive: true });

  if (opts.watch) {
    log.step(`📦 Bundling em watch mode → ${bundlePath}`);
    await $`bun build ${config.entry} --outfile ${bundlePath} --target=bun --watch`;
    return;
  }

  const sp = spinner(`Bundling → ${bundlePath}`);
  sp.start();

  const { durationMs } = await timed(async () => {
    await $`bun build ${config.entry} --outfile ${bundlePath} --target=bun --minify`.quiet();
  });

  const size = fileSize(statSync(bundlePath).size);
  sp.stop(
    `${pc.green("📦")} ${pc.bold("Bundled")} ${pc.dim(`→ ${bundlePath}`)} ${pc.dim(`(${size}, ${formatDuration(durationMs)})`)}`,
  );

  if (config.afterPackage) {
    log.verbose("⚙️  Executando hook afterPackage");
    await config.afterPackage(config);
  }
}
