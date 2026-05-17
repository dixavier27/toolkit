import { mkdirSync, statSync } from "node:fs";
import { $ } from "bun";
import type { EcoConfig } from "../config.ts";
import type { CommandMeta } from "../utils/command-meta.ts";
import { log, pc } from "../utils/logger.ts";
import { spinner } from "../utils/spinner.ts";
import { type Column, fileSize, renderTable } from "../utils/table.ts";
import { formatDuration, timed } from "../utils/timing.ts";
import { runObfuscate } from "./obfuscate.ts";
import { type RunOptions, runPackage } from "./package.ts";

export const meta: CommandMeta = {
  name: "release",
  description: "Pipeline completo: package → obfuscate → binários nativos",
  flags: [
    { name: "--skip-obfuscate", description: "Pula a etapa de ofuscação" },
    { name: "--dry-run", description: "Mostra os comandos sem executar" },
  ],
  examples: [
    "eco release",
    "eco release --skip-obfuscate",
    "eco release --platforms=linux,win --verbose",
  ],
};

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

interface ReleaseArtifact {
  platform: string;
  outfile: string;
  durationMs: number;
}

export async function runRelease(config: EcoConfig, opts: ReleaseOptions = {}) {
  await runPackage(config, opts);
  if (!opts.skipObfuscate) await runObfuscate(config, opts);

  const bundle = `${config.outDir}/${config.bundleName}`;

  if (!opts.dryRun) mkdirSync("release", { recursive: true });

  const artifacts: ReleaseArtifact[] = [];

  for (const platform of config.platforms) {
    const target = bunTargets[platform];
    const outfile = `release/${config.releaseName}-${platform}${ext[platform] ?? ""}`;

    if (opts.dryRun) {
      log.info(pc.cyan(`🚀 [dry-run] Compile ${platform} → ${outfile}`));
      log.dim(
        `   bun build ${bundle} --compile --target=${target} --outfile ${outfile}`,
      );
      continue;
    }

    const sp = spinner(`Compiling ${platform} → ${outfile}`);
    sp.start();

    const { durationMs } = await timed(async () => {
      await $`bun build ${bundle} --compile --target=${target} --outfile ${outfile}`.quiet();
    });

    sp.stop(
      `${pc.green("🚀")} ${pc.bold(platform.padEnd(11))} ${pc.dim(`→ ${outfile}`)} ${pc.dim(`(${formatDuration(durationMs)})`)}`,
    );

    artifacts.push({ platform, outfile, durationMs });
  }

  if (config.afterRelease) {
    log.verbose("⚙️  Executando hook afterRelease");
    if (!opts.dryRun) await config.afterRelease(config);
  }

  if (!opts.dryRun && artifacts.length > 0) {
    printSummary(artifacts);
  }
}

function printSummary(artifacts: ReleaseArtifact[]) {
  log.info("");
  log.info(pc.bold("Release pronto:"));
  const columns: Column[] = [
    { header: "Plataforma" },
    { header: "Tamanho", align: "right" },
    { header: "Tempo", align: "right" },
    { header: "Path" },
  ];
  const rows = artifacts.map((a) => [
    a.platform,
    fileSize(statSync(a.outfile).size),
    formatDuration(a.durationMs),
    a.outfile,
  ]);
  log.info(renderTable(columns, rows));
}
