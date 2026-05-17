import { mkdirSync, statSync } from "node:fs";
import { $ } from "bun";
import type { EcoConfig } from "../config.ts";
import { writeChecksumsFile } from "../utils/checksum.ts";
import type { CommandMeta } from "../utils/command-meta.ts";
import { log, pc } from "../utils/logger.ts";
import { type Column, fileSize, renderTable } from "../utils/table.ts";
import { formatDuration, timed } from "../utils/timing.ts";
import { runObfuscate } from "./obfuscate.ts";
import { type RunOptions, runPackage } from "./package.ts";

export const meta: CommandMeta = {
  name: "release",
  description: "Pipeline completo: package → obfuscate → binários nativos",
  flags: [
    { name: "--skip-obfuscate", description: "Pula a etapa de ofuscação" },
    {
      name: "--keep-going",
      description: "Continua mesmo se uma plataforma falhar",
    },
    { name: "--no-parallel", description: "Compila plataformas em sequência" },
    { name: "--dry-run", description: "Mostra os comandos sem executar" },
  ],
  examples: [
    "eco release",
    "eco release --skip-obfuscate",
    "eco release --platforms=linux,win --verbose",
    "eco release --keep-going",
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
  keepGoing?: boolean;
  parallel?: boolean;
}

interface ReleaseArtifact {
  platform: string;
  outfile: string;
  durationMs: number;
  status: "ok" | "failed";
  error?: string;
}

async function compileOne(
  platform: string,
  bundle: string,
  releaseName: string,
  dryRun: boolean,
): Promise<ReleaseArtifact> {
  const target = bunTargets[platform];
  const outfile = `release/${releaseName}-${platform}${ext[platform] ?? ""}`;

  if (dryRun) {
    log.info(pc.cyan(`🚀 [dry-run] Compile ${platform} → ${outfile}`));
    log.dim(
      `   bun build ${bundle} --compile --target=${target} --outfile ${outfile}`,
    );
    return { platform, outfile, durationMs: 0, status: "ok" };
  }

  try {
    const { durationMs } = await timed(async () => {
      await $`bun build ${bundle} --compile --target=${target} --outfile ${outfile}`.quiet();
    });
    log.info(
      `${pc.green("🚀")} ${pc.bold(platform.padEnd(11))} ${pc.dim(`→ ${outfile}`)} ${pc.dim(`(${formatDuration(durationMs)})`)}`,
    );
    return { platform, outfile, durationMs, status: "ok" };
  } catch (err) {
    const message = err instanceof Error ? err.message : String(err);
    log.error(`${pc.red("✗")} ${platform}: ${message}`);
    return {
      platform,
      outfile,
      durationMs: 0,
      status: "failed",
      error: message,
    };
  }
}

export async function runRelease(config: EcoConfig, opts: ReleaseOptions = {}) {
  await runPackage(config, opts);
  if (!opts.skipObfuscate) await runObfuscate(config, opts);

  const bundle = `${config.outDir}/${config.bundleName}`;
  if (!opts.dryRun) mkdirSync("release", { recursive: true });

  const useParallel = opts.parallel !== false && config.parallel;
  const dryRun = opts.dryRun === true;

  log.info(
    pc.dim(
      `\n🔧 Compilando ${config.platforms.length} plataforma${config.platforms.length !== 1 ? "s" : ""} ${useParallel ? "em paralelo" : "em sequência"}…\n`,
    ),
  );

  let artifacts: ReleaseArtifact[];
  if (useParallel) {
    artifacts = await Promise.all(
      config.platforms.map((p) =>
        compileOne(p, bundle, config.releaseName, dryRun),
      ),
    );
  } else {
    artifacts = [];
    for (const p of config.platforms) {
      const r = await compileOne(p, bundle, config.releaseName, dryRun);
      artifacts.push(r);
      if (r.status === "failed" && !opts.keepGoing) {
        throw new Error(
          `Compilação de ${p} falhou. Use --keep-going para continuar mesmo assim.`,
        );
      }
    }
  }

  const failures = artifacts.filter((a) => a.status === "failed");
  if (failures.length > 0 && !opts.keepGoing) {
    throw new Error(
      `${failures.length} plataforma(s) falharam: ${failures.map((f) => f.platform).join(", ")}. Use --keep-going para gerar os artefatos restantes.`,
    );
  }

  const successes = artifacts.filter((a) => a.status === "ok");

  if (!dryRun && config.checksums && successes.length > 0) {
    const checksumPath = "release/checksums.txt";
    await writeChecksumsFile(
      successes.map((a) => a.outfile),
      checksumPath,
    );
    log.info(
      `${pc.green("🔐")} ${pc.bold("Checksums")} ${pc.dim(`→ ${checksumPath}`)}`,
    );
  }

  if (config.afterRelease) {
    log.verbose("⚙️  Executando hook afterRelease");
    if (!dryRun) await config.afterRelease(config);
  }

  if (!dryRun && artifacts.length > 0) {
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
    { header: "Status" },
    { header: "Path" },
  ];
  const rows = artifacts.map((a) => {
    const status =
      a.status === "ok" ? pc.green("ok") : pc.red(`falhou: ${a.error ?? "?"}`);
    const size =
      a.status === "ok"
        ? (() => {
            try {
              return fileSize(statSync(a.outfile).size);
            } catch {
              return "-";
            }
          })()
        : "-";
    return [a.platform, size, formatDuration(a.durationMs), status, a.outfile];
  });
  log.info(renderTable(columns, rows));
}
