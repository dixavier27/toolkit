import { mkdirSync } from "node:fs";
import { $ } from "bun";
import type { EcoConfig } from "../config.ts";
import { log } from "../utils/logger.ts";
import { runObfuscate } from "./obfuscate.ts";
import { type RunOptions, runPackage } from "./package.ts";

const bunTargets: Record<string, string> = {
  linux: "bun-linux-x64",
  win: "bun-windows-x64",
  macos: "bun-darwin-x64",
  "macos-arm64": "bun-darwin-arm64",
};

const ext: Record<string, string> = {
  win: ".exe",
};

export interface ReleaseOptions extends RunOptions {
  skipObfuscate?: boolean;
}

export async function runRelease(config: EcoConfig, opts: ReleaseOptions = {}) {
  await runPackage(config, opts);
  if (!opts.skipObfuscate) await runObfuscate(config, opts);

  const bundle = `${config.outDir}/${config.bundleName}`;

  if (!opts.dryRun) mkdirSync("release", { recursive: true });

  for (const platform of config.platforms) {
    const target = bunTargets[platform];
    const outfile = `release/${config.releaseName}-${platform}${ext[platform] ?? ""}`;

    log.info(`🚀 Compiling ${platform} → ${outfile}`);

    if (opts.dryRun) {
      log.info(
        `   [dry-run] bun build ${bundle} --compile --target=${target} --outfile ${outfile}`,
      );
    } else {
      await $`bun build ${bundle} --compile --target=${target} --outfile ${outfile}`;
    }
  }

  if (config.afterRelease) {
    log.verbose("⚙️  Executando hook afterRelease");
    if (!opts.dryRun) await config.afterRelease(config);
  }
}
