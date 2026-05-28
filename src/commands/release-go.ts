import {
  existsSync,
  mkdirSync,
  readdirSync,
  readFileSync,
  statSync,
} from "node:fs";
import { basename, resolve } from "node:path";
import { $ } from "bun";
import type { EcoConfig, Platform } from "../config.ts";
import { writeChecksumsFile } from "../utils/checksum.ts";
import { log, pc } from "../utils/logger.ts";
import { type Column, fileSize, renderTable } from "../utils/table.ts";
import { formatDuration, timed } from "../utils/timing.ts";

export interface ReleaseGoOptions {
  dryRun?: boolean;
  keepGoing?: boolean;
  parallel?: boolean;
}

interface GoTarget {
  goos: string;
  goarch: string;
  ext: string;
}

const GO_TARGETS: Record<Platform, GoTarget> = {
  linux: { goos: "linux", goarch: "amd64", ext: "" },
  win: { goos: "windows", goarch: "amd64", ext: ".exe" },
  macos: { goos: "darwin", goarch: "amd64", ext: "" },
  "macos-arm64": { goos: "darwin", goarch: "arm64", ext: "" },
};

interface GoArtifact {
  platform: Platform;
  outfile: string;
  durationMs: number;
  status: "ok" | "failed";
  error?: string;
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

async function compileOnePlatform(
  platform: Platform,
  entry: string,
  releaseName: string,
  ldflags: string,
  dryRun: boolean,
): Promise<GoArtifact> {
  const target = GO_TARGETS[platform];
  const outfile = `release/${releaseName}-${platform}${target.ext}`;

  if (dryRun) {
    log.info(pc.cyan(`🚀 [dry-run] Compile ${platform} → ${outfile}`));
    log.dim(
      `   GOOS=${target.goos} GOARCH=${target.goarch} go build -ldflags "${ldflags}" -o ${outfile} ${entry}`,
    );
    return { platform, outfile, durationMs: 0, status: "ok" };
  }

  try {
    const env = {
      ...process.env,
      GOOS: target.goos,
      GOARCH: target.goarch,
      CGO_ENABLED: "0",
    };
    const { durationMs } = await timed(async () => {
      await $`go build -ldflags ${ldflags} -o ${outfile} ${entry}`
        .env(env)
        .quiet();
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

export async function runReleaseGo(
  config: EcoConfig,
  opts: ReleaseGoOptions = {},
) {
  const cwd = process.cwd();
  const name =
    config.releaseName !== "app" ? config.releaseName : basename(cwd);
  const entry = detectEntryPackage(cwd, name);
  const version = readVersion(cwd);
  const ldflags = `-s -w -X main.version=${version}`;
  const dryRun = opts.dryRun === true;

  if (!dryRun) mkdirSync("release", { recursive: true });

  const useParallel = opts.parallel !== false && config.parallel;
  log.info(
    pc.dim(
      `\n🔧 Compilando ${config.platforms.length} plataforma${config.platforms.length !== 1 ? "s" : ""} ${useParallel ? "em paralelo" : "em sequência"}…\n`,
    ),
  );

  let artifacts: GoArtifact[];
  if (useParallel) {
    artifacts = await Promise.all(
      config.platforms.map((p) =>
        compileOnePlatform(p, entry, name, ldflags, dryRun),
      ),
    );
  } else {
    artifacts = [];
    for (const p of config.platforms) {
      const r = await compileOnePlatform(p, entry, name, ldflags, dryRun);
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

function printSummary(artifacts: GoArtifact[]) {
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
