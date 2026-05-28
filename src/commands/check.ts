import { existsSync } from "node:fs";
import { resolve } from "node:path";
import { $ } from "bun";
import { type EcoConfig, getConfigPath } from "../config.ts";
import type { CommandMeta } from "../utils/command-meta.ts";
import { detectLang } from "../utils/detect-lang.ts";
import { log, pc } from "../utils/logger.ts";

export const meta: CommandMeta = {
  name: "check",
  description: "Valida configuração e ambiente para uso do eco",
  examples: ["eco check"],
};

type Status = "ok" | "warn" | "fail";

interface CheckResult {
  status: Status;
  label: string;
  detail?: string;
}

function symbol(s: Status): string {
  if (s === "ok") return pc.green("✅");
  if (s === "warn") return pc.yellow("⚠️ ");
  return pc.red("❌");
}

async function checkBunVersion(): Promise<CheckResult> {
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
      detail: "não encontrado no PATH",
    };
  }
}

async function checkObfuscator(): Promise<CheckResult> {
  try {
    await $`javascript-obfuscator --version`.quiet();
    return { status: "ok", label: "javascript-obfuscator no PATH" };
  } catch {
    return {
      status: "warn",
      label: "javascript-obfuscator no PATH",
      detail: "necessário para 'eco obfuscate' e 'eco release'",
    };
  }
}

function checkConfigFile(): CheckResult {
  const path = getConfigPath();
  if (!path) {
    return {
      status: "fail",
      label: "Arquivo de config",
      detail: "eco.config.js / eco.config.ts ausente — rode 'eco init'",
    };
  }
  return { status: "ok", label: "Arquivo de config", detail: path };
}

function checkEntry(config: EcoConfig): CheckResult {
  const path = resolve(process.cwd(), config.entry);
  return existsSync(path)
    ? { status: "ok", label: `entry: ${config.entry}` }
    : {
        status: "fail",
        label: `entry: ${config.entry}`,
        detail: "arquivo não encontrado",
      };
}

function checkObfuscatorConfig(config: EcoConfig): CheckResult {
  const path = resolve(process.cwd(), config.obfuscatorConfig);
  return existsSync(path)
    ? { status: "ok", label: `obfuscatorConfig: ${config.obfuscatorConfig}` }
    : {
        status: "warn",
        label: `obfuscatorConfig: ${config.obfuscatorConfig}`,
        detail: "arquivo não encontrado — obfuscate falhará",
      };
}

function checkPlatformsHost(config: EcoConfig): CheckResult {
  const host = process.platform;
  const macosPlatforms = config.platforms.filter((p: string) =>
    p.startsWith("macos"),
  );
  if (macosPlatforms.length > 0 && host !== "darwin") {
    return {
      status: "warn",
      label: `platforms: ${config.platforms.join(", ")}`,
      detail: `cross-compile de ${macosPlatforms.join(", ")} a partir de ${host} pode falhar`,
    };
  }
  return { status: "ok", label: `platforms: ${config.platforms.join(", ")}` };
}

async function checkGoVersion(): Promise<CheckResult> {
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
      detail: "não encontrado no PATH",
    };
  }
}

function checkGoMod(): CheckResult {
  const path = resolve(process.cwd(), "go.mod");
  return existsSync(path)
    ? { status: "ok", label: "go.mod", detail: path }
    : { status: "fail", label: "go.mod", detail: "arquivo não encontrado" };
}

export async function runCheck(config: EcoConfig) {
  const lang = detectLang();
  log.info(
    pc.bold(
      `🔍 Verificando ambiente eco ${lang === "go" ? pc.dim("(projeto Go)") : ""}`,
    ),
  );
  log.info("");

  const results: CheckResult[] =
    lang === "go"
      ? [checkGoMod(), checkPlatformsHost(config), await checkGoVersion()]
      : [
          checkConfigFile(),
          checkEntry(config),
          checkObfuscatorConfig(config),
          checkPlatformsHost(config),
          await checkBunVersion(),
          await checkObfuscator(),
        ];

  for (const r of results) {
    const detail = r.detail ? pc.dim(` — ${r.detail}`) : "";
    log.info(`${symbol(r.status)} ${r.label}${detail}`);
  }

  const failed = results.filter((r) => r.status === "fail").length;
  const warned = results.filter((r) => r.status === "warn").length;
  const ok = results.length - failed - warned;

  log.info("");
  log.info(
    `Resumo: ${pc.green(`${ok} ok`)}, ${pc.yellow(`${warned} warning${warned !== 1 ? "s" : ""}`)}, ${pc.red(`${failed} falha${failed !== 1 ? "s" : ""}`)}`,
  );

  if (failed > 0) process.exit(1);
}
