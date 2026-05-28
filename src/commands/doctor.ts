import {
  appendFileSync,
  existsSync,
  readFileSync,
  writeFileSync,
} from "node:fs";
import { resolve } from "node:path";
import { $ } from "bun";
import { type EcoConfig, getConfigPath } from "../config.ts";
import type { CommandMeta } from "../utils/command-meta.ts";
import { detectLang } from "../utils/detect-lang.ts";
import { log, pc } from "../utils/logger.ts";
import { runScriptsInject } from "./scripts.ts";

export const meta: CommandMeta = {
  name: "doctor",
  description:
    "Diagnóstico do ambiente e do projeto. Use --fix para autocorrigir problemas",
  flags: [
    {
      name: "--fix",
      description: "Aplica correções automáticas quando possível",
    },
  ],
  examples: ["eco doctor", "eco doctor --fix"],
};

type Status = "ok" | "warn" | "fail";

interface DiagnosticResult {
  status: Status;
  label: string;
  detail?: string;
  fix?: () => void | Promise<void>;
  fixLabel?: string;
}

function symbol(s: Status): string {
  if (s === "ok") return pc.green("✅");
  if (s === "warn") return pc.yellow("⚠️ ");
  return pc.red("❌");
}

async function checkBun(): Promise<DiagnosticResult> {
  try {
    const result = await $`bun --version`.text();
    const version = result.trim();
    const [major, minor] = version.split(".").map(Number);
    const ok = (major ?? 0) > 1 || ((major ?? 0) === 1 && (minor ?? 0) >= 2);
    return {
      status: ok ? "ok" : "fail",
      label: "Bun >= 1.2.0",
      detail: `instalado: ${version}`,
    };
  } catch {
    return {
      status: "fail",
      label: "Bun >= 1.2.0",
      detail: "não encontrado no PATH — instale em https://bun.sh",
    };
  }
}

async function checkObfuscator(): Promise<DiagnosticResult> {
  try {
    await $`javascript-obfuscator --version`.quiet();
    return { status: "ok", label: "javascript-obfuscator no PATH" };
  } catch {
    return {
      status: "warn",
      label: "javascript-obfuscator no PATH",
      detail:
        "instale com 'bun add -D javascript-obfuscator' para usar obfuscate/release",
    };
  }
}

function checkConfigFile(): DiagnosticResult {
  const path = getConfigPath();
  if (!path) {
    return {
      status: "fail",
      label: "eco.config.{js,ts,mjs}",
      detail: "rode 'eco init'",
    };
  }
  return { status: "ok", label: "eco.config", detail: path };
}

function checkEntry(config: EcoConfig): DiagnosticResult {
  const path = resolve(process.cwd(), config.entry);
  return existsSync(path)
    ? { status: "ok", label: `entry: ${config.entry}` }
    : {
        status: "fail",
        label: `entry: ${config.entry}`,
        detail: "arquivo não encontrado",
      };
}

function checkObfuscatorConfig(config: EcoConfig): DiagnosticResult {
  const path = resolve(process.cwd(), config.obfuscatorConfig);
  if (existsSync(path)) {
    return {
      status: "ok",
      label: `obfuscatorConfig: ${config.obfuscatorConfig}`,
    };
  }
  return {
    status: "warn",
    label: `obfuscatorConfig: ${config.obfuscatorConfig}`,
    detail: "arquivo não encontrado",
    fixLabel: "criar com preset 'medium'",
    fix: () => {
      const template = `// Preset 'medium' — equilíbrio entre proteção e performance.
// Documentação: https://github.com/javascript-obfuscator/javascript-obfuscator
module.exports = {
  compact: true,
  controlFlowFlattening: true,
  controlFlowFlatteningThreshold: 0.5,
  deadCodeInjection: false,
  stringArray: true,
  stringArrayEncoding: ['base64'],
  stringArrayThreshold: 0.75,
  identifierNamesGenerator: 'mangled',
  numbersToExpressions: true,
  simplify: true,
  transformObjectKeys: false,
  unicodeEscapeSequence: false,
};
`;
      writeFileSync(path, template, "utf8");
    },
  };
}

function checkGitignore(): DiagnosticResult {
  const path = resolve(process.cwd(), ".gitignore");
  if (!existsSync(path)) {
    return {
      status: "warn",
      label: ".gitignore",
      detail: "arquivo não existe",
      fixLabel: "criar com entries do eco",
      fix: () => {
        writeFileSync(
          path,
          "node_modules/\ndist/\nrelease/\n.eco-cache/\n",
          "utf8",
        );
      },
    };
  }

  const content = readFileSync(path, "utf8");
  const missing: string[] = [];
  for (const entry of ["dist/", "release/"]) {
    if (!new RegExp(`(^|\\n)${entry.replace(/[/]/g, "\\/")}`).test(content)) {
      missing.push(entry);
    }
  }

  if (missing.length === 0) {
    return { status: "ok", label: ".gitignore inclui dist/ e release/" };
  }

  return {
    status: "warn",
    label: ".gitignore",
    detail: `faltando: ${missing.join(", ")}`,
    fixLabel: "adicionar entries faltantes",
    fix: () => {
      const additions = missing.join("\n");
      appendFileSync(path, `\n# eco artifacts\n${additions}\n`, "utf8");
    },
  };
}

function checkPackageScripts(): DiagnosticResult {
  const pkgPath = resolve(process.cwd(), "package.json");
  if (!existsSync(pkgPath)) {
    return {
      status: "warn",
      label: "package.json",
      detail: "arquivo não encontrado",
    };
  }
  const pkg = JSON.parse(readFileSync(pkgPath, "utf8"));
  const scripts = pkg.scripts ?? {};
  const missing = ["package", "obfuscate", "release", "check"].filter(
    (s) => !scripts[s],
  );

  if (missing.length === 0) {
    return { status: "ok", label: "scripts do eco no package.json" };
  }

  return {
    status: "warn",
    label: "scripts do eco no package.json",
    detail: `faltando: ${missing.join(", ")}`,
    fixLabel: "rodar 'eco scripts inject'",
    fix: () => runScriptsInject(),
  };
}

function checkPlatformsHost(config: EcoConfig): DiagnosticResult {
  const host = process.platform;
  const macosPlatforms = config.platforms.filter((p: string) =>
    p.startsWith("macos"),
  );
  if (macosPlatforms.length > 0 && host !== "darwin") {
    return {
      status: "warn",
      label: `platforms: ${config.platforms.join(", ")}`,
      detail: `cross-compile de ${macosPlatforms.join(", ")} em ${host} pode falhar — use macos-latest no CI`,
    };
  }
  return {
    status: "ok",
    label: `platforms: ${config.platforms.join(", ")}`,
  };
}

async function checkGo(): Promise<DiagnosticResult> {
  try {
    const result = await $`go version`.text();
    const match = result.match(/go(\d+)\.(\d+)/);
    if (!match) {
      return {
        status: "warn",
        label: "Go >= 1.22",
        detail: `versão não reconhecida: ${result.trim()}`,
      };
    }
    const major = Number(match[1]);
    const minor = Number(match[2]);
    const ok = major > 1 || (major === 1 && minor >= 22);
    return {
      status: ok ? "ok" : "fail",
      label: "Go >= 1.22",
      detail: `instalado: ${result.trim()}`,
    };
  } catch {
    return {
      status: "fail",
      label: "Go >= 1.22",
      detail: "não encontrado no PATH — instale em https://go.dev/dl/",
    };
  }
}

async function checkAir(): Promise<DiagnosticResult> {
  try {
    await $`air -v`.quiet();
    return { status: "ok", label: "air (hot-reload) no PATH" };
  } catch {
    return {
      status: "warn",
      label: "air (hot-reload) no PATH",
      detail: "instale com 'go install github.com/air-verse/air@latest'",
    };
  }
}

async function checkGolangciLint(): Promise<DiagnosticResult> {
  try {
    await $`golangci-lint --version`.quiet();
    return { status: "ok", label: "golangci-lint no PATH" };
  } catch {
    return {
      status: "warn",
      label: "golangci-lint no PATH",
      detail:
        "instale com 'go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest'",
    };
  }
}

function checkGoMod(): DiagnosticResult {
  const path = resolve(process.cwd(), "go.mod");
  if (!existsSync(path)) {
    return {
      status: "fail",
      label: "go.mod",
      detail: "arquivo não encontrado",
    };
  }
  return { status: "ok", label: "go.mod", detail: path };
}

export interface DoctorOptions {
  fix?: boolean;
}

export async function runDoctor(config: EcoConfig, opts: DoctorOptions = {}) {
  const lang = detectLang();
  log.info(
    pc.bold(
      `🩺 Diagnóstico eco ${lang === "go" ? pc.dim("(projeto Go)") : ""}`,
    ),
  );
  log.info("");

  const diagnostics: DiagnosticResult[] =
    lang === "go"
      ? [
          checkGoMod(),
          checkGitignore(),
          await checkGo(),
          await checkAir(),
          await checkGolangciLint(),
        ]
      : [
          checkConfigFile(),
          checkEntry(config),
          checkObfuscatorConfig(config),
          checkPlatformsHost(config),
          checkGitignore(),
          checkPackageScripts(),
          await checkBun(),
          await checkObfuscator(),
        ];

  for (const d of diagnostics) {
    const detail = d.detail ? pc.dim(` — ${d.detail}`) : "";
    log.info(`${symbol(d.status)} ${d.label}${detail}`);
  }

  const fixable = diagnostics.filter((d) => d.fix);
  const failed = diagnostics.filter((d) => d.status === "fail").length;
  const warned = diagnostics.filter((d) => d.status === "warn").length;
  const ok = diagnostics.length - failed - warned;

  log.info("");
  log.info(
    `Resumo: ${pc.green(`${ok} ok`)}, ${pc.yellow(`${warned} warning${warned !== 1 ? "s" : ""}`)}, ${pc.red(`${failed} falha${failed !== 1 ? "s" : ""}`)}`,
  );

  if (!opts.fix) {
    if (fixable.length > 0) {
      log.info("");
      log.info(
        pc.dim(
          `${fixable.length} problema${fixable.length !== 1 ? "s" : ""} pode${fixable.length !== 1 ? "m" : ""} ser corrigido${fixable.length !== 1 ? "s" : ""} automaticamente — rode ${pc.cyan("eco doctor --fix")}`,
        ),
      );
    }
    if (failed > 0) process.exit(1);
    return;
  }

  if (fixable.length === 0) {
    log.info(pc.dim("\nNenhum problema fixável."));
    if (failed > 0) process.exit(1);
    return;
  }

  log.info("");
  log.info(pc.bold("🔧 Aplicando correções:"));
  for (const d of fixable) {
    log.info(`   ${pc.cyan(d.label)} → ${pc.dim(d.fixLabel ?? "fix")}`);
    try {
      await d.fix?.();
    } catch (err) {
      const msg = err instanceof Error ? err.message : String(err);
      log.error(`   ❌ falhou: ${msg}`);
    }
  }
  log.success(
    "\n✅ Correções aplicadas. Rode 'eco doctor' novamente para verificar.",
  );
}
