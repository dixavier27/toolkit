import { mkdirSync } from "node:fs";
import { $ } from "bun";
import type { EcoConfig } from "../config.ts";
import { log } from "../utils/logger.ts";

export interface RunOptions {
  dryRun?: boolean;
}

export async function runPackage(config: EcoConfig, opts: RunOptions = {}) {
  const bundlePath = `${config.outDir}/${config.bundleName}`;

  log.info(`📦 Bundling → ${bundlePath}`);

  if (opts.dryRun) {
    log.info(
      `   [dry-run] bun build ${config.entry} --outfile ${bundlePath} --target=bun --minify`,
    );
  } else {
    mkdirSync(config.outDir, { recursive: true });
    await $`bun build ${config.entry} --outfile ${bundlePath} --target=bun --minify`;
  }

  if (config.afterPackage) {
    log.verbose("⚙️  Executando hook afterPackage");
    if (!opts.dryRun) await config.afterPackage(config);
  }
}
