import { existsSync, statSync } from "node:fs";
import { $ } from "bun";
import type { EcoConfig } from "../config.ts";
import type { CommandMeta } from "../utils/command-meta.ts";
import { log, pc } from "../utils/logger.ts";
import { spinner } from "../utils/spinner.ts";
import { fileSize } from "../utils/table.ts";
import { formatDuration, timed } from "../utils/timing.ts";
import type { RunOptions } from "./package.ts";

export const meta: CommandMeta = {
  name: "obfuscate",
  description: "Ofusca o bundle JS gerado por 'eco package'",
  flags: [{ name: "--dry-run", description: "Mostra o comando sem executar" }],
  examples: ["eco package && eco obfuscate", "eco obfuscate --dry-run"],
};

export async function runObfuscate(config: EcoConfig, opts: RunOptions = {}) {
  const bundle = `${config.outDir}/${config.bundleName}`;

  if (opts.dryRun) {
    log.info(pc.cyan(`🔒 [dry-run] Obfuscate ${bundle}`));
    log.dim(
      `   javascript-obfuscator ${bundle} --output ${bundle} --config ${config.obfuscatorConfig}`,
    );
    return;
  }

  if (!existsSync(bundle)) {
    throw new Error(
      `Bundle não encontrado em ${bundle}. Rode 'eco package' antes (ou 'eco release' que orquestra tudo).`,
    );
  }

  const sp = spinner(`Obfuscating ${bundle}`);
  sp.start();

  const { durationMs } = await timed(async () => {
    await $`javascript-obfuscator ${bundle} --output ${bundle} --config ${config.obfuscatorConfig}`.quiet();
  });

  const size = fileSize(statSync(bundle).size);
  sp.stop(
    `${pc.green("🔒")} ${pc.bold("Obfuscated")} ${pc.dim(`→ ${bundle}`)} ${pc.dim(`(${size}, ${formatDuration(durationMs)})`)}`,
  );

  if (config.afterObfuscate) {
    log.verbose("⚙️  Executando hook afterObfuscate");
    await config.afterObfuscate(config);
  }
}
