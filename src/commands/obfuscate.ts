import { existsSync } from "node:fs";
import { $ } from "bun";
import type { EcoConfig } from "../config.ts";
import { log } from "../utils/logger.ts";
import type { RunOptions } from "./package.ts";

export async function runObfuscate(config: EcoConfig, opts: RunOptions = {}) {
  const bundle = `${config.outDir}/${config.bundleName}`;

  if (!opts.dryRun && !existsSync(bundle)) {
    throw new Error(
      `Bundle não encontrado em ${bundle}. Rode 'eco package' antes (ou 'eco release' que orquestra tudo).`,
    );
  }

  log.info(`🔒 Obfuscating ${bundle}`);

  if (opts.dryRun) {
    log.info(
      `   [dry-run] javascript-obfuscator ${bundle} --output ${bundle} --config ${config.obfuscatorConfig}`,
    );
  } else {
    await $`javascript-obfuscator ${bundle} --output ${bundle} --config ${config.obfuscatorConfig}`;
  }

  if (config.afterObfuscate) {
    log.verbose("⚙️  Executando hook afterObfuscate");
    if (!opts.dryRun) await config.afterObfuscate(config);
  }
}
