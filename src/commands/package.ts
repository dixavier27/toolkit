import { existsSync, mkdirSync, readFileSync, statSync } from "node:fs";
import { resolve } from "node:path";
import { $ } from "bun";
import type { EcoConfig } from "../config.ts";
import { copyAssets } from "../utils/assets.ts";
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

function resolveProjectVersion(): string | undefined {
  try {
    const pkg = JSON.parse(
      readFileSync(resolve(process.cwd(), "package.json"), "utf8"),
    );
    return typeof pkg.version === "string" ? pkg.version : undefined;
  } catch {
    return undefined;
  }
}

function buildDefineArgs(config: EcoConfig): string[] {
  const define: Record<string, string> = { ...config.define };

  if (config.embedVersion) {
    const version = resolveProjectVersion();
    if (version && !("__VERSION__" in define)) {
      define.__VERSION__ = JSON.stringify(version);
    }
  }

  const args: string[] = [];
  for (const [key, value] of Object.entries(define)) {
    args.push("--define", `${key}=${value}`);
  }
  return args;
}

function buildSourcemapArgs(config: EcoConfig): string[] {
  if (config.sourcemap === false) return [];
  return ["--sourcemap", config.sourcemap];
}

export async function runPackage(config: EcoConfig, opts: RunOptions = {}) {
  const bundlePath = `${config.outDir}/${config.bundleName}`;
  const defineArgs = buildDefineArgs(config);
  const sourcemapArgs = buildSourcemapArgs(config);

  if (opts.dryRun) {
    log.info(pc.cyan(`📦 [dry-run] Bundle → ${bundlePath}`));
    const extras = [...defineArgs, ...sourcemapArgs].join(" ");
    log.dim(
      `   bun build ${config.entry} --outfile ${bundlePath} --target=bun --minify ${extras}`.trim(),
    );
    if (config.assets.length > 0) {
      log.dim(`   [dry-run] copiar ${config.assets.length} asset(s)`);
    }
    return;
  }

  mkdirSync(config.outDir, { recursive: true });

  if (opts.watch) {
    log.step(`📦 Bundling em watch mode → ${bundlePath}`);
    await $`bun build ${config.entry} --outfile ${bundlePath} --target=bun --watch ${defineArgs} ${sourcemapArgs}`;
    return;
  }

  const sp = spinner(`Bundling → ${bundlePath}`);
  sp.start();

  const { durationMs } = await timed(async () => {
    await $`bun build ${config.entry} --outfile ${bundlePath} --target=bun --minify ${defineArgs} ${sourcemapArgs}`.quiet();
  });

  const size = fileSize(statSync(bundlePath).size);
  sp.stop(
    `${pc.green("📦")} ${pc.bold("Bundled")} ${pc.dim(`→ ${bundlePath}`)} ${pc.dim(`(${size}, ${formatDuration(durationMs)})`)}`,
  );

  if (config.sourcemap === "external" && existsSync(`${bundlePath}.map`)) {
    log.verbose(`🗺️  sourcemap: ${bundlePath}.map`);
  }

  if (config.assets.length > 0) {
    await copyAssets(config.assets);
  }

  if (config.afterPackage) {
    log.verbose("⚙️  Executando hook afterPackage");
    await config.afterPackage(config);
  }
}
