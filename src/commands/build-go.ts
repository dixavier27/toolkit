import {
  existsSync,
  mkdirSync,
  readdirSync,
  readFileSync,
  statSync,
} from "node:fs";
import { basename, resolve } from "node:path";
import { $ } from "bun";
import type { EcoConfig } from "../config.ts";
import { log, pc } from "../utils/logger.ts";
import { spinner } from "../utils/spinner.ts";
import { fileSize } from "../utils/table.ts";
import { formatDuration, timed } from "../utils/timing.ts";

export interface BuildGoOptions {
  dryRun?: boolean;
}

function readVersion(cwd: string): string {
  try {
    return readFileSync(resolve(cwd, "VERSION"), "utf8").trim() || "dev";
  } catch {
    return "dev";
  }
}

function detectEntryPackage(cwd: string, releaseName: string): string {
  const cmdDir = resolve(cwd, "cmd");
  if (existsSync(cmdDir)) {
    const candidates = readdirSync(cmdDir).filter((entry) => {
      const path = resolve(cmdDir, entry);
      return (
        statSync(path).isDirectory() && existsSync(resolve(path, "main.go"))
      );
    });
    if (candidates.includes(releaseName)) return `./cmd/${releaseName}`;
    if (candidates.length === 1) return `./cmd/${candidates[0]}`;
    if (candidates.length > 1) {
      throw new Error(
        `Múltiplos pacotes em cmd/: ${candidates.join(", ")}. Renomeie ou ajuste 'releaseName' no eco.config.js.`,
      );
    }
  }
  if (existsSync(resolve(cwd, "main.go"))) return ".";
  throw new Error(
    "Não encontrei entry Go. Esperado cmd/<nome>/main.go ou main.go na raiz.",
  );
}

export async function runBuildGo(config: EcoConfig, opts: BuildGoOptions = {}) {
  const cwd = process.cwd();
  const name =
    config.releaseName !== "app" ? config.releaseName : basename(cwd);
  const entry = detectEntryPackage(cwd, name);
  const outDir = config.outDir;
  const outfile = `${outDir}/${name}`;
  const version = readVersion(cwd);
  const ldflags = `-s -w -X main.version=${version}`;

  if (opts.dryRun) {
    log.info(pc.cyan(`📦 [dry-run] Build Go → ${outfile}`));
    log.dim(`   go build -ldflags "${ldflags}" -o ${outfile} ${entry}`);
    return;
  }

  mkdirSync(outDir, { recursive: true });

  const sp = spinner(`Compilando ${entry} → ${outfile}`);
  sp.start();

  const { durationMs } = await timed(async () => {
    await $`go build -ldflags ${ldflags} -o ${outfile} ${entry}`.quiet();
  });

  const size = fileSize(statSync(outfile).size);
  sp.stop(
    `${pc.green("📦")} ${pc.bold("Built")} ${pc.dim(`→ ${outfile}`)} ${pc.dim(`(${size}, ${formatDuration(durationMs)})`)}`,
  );
}
