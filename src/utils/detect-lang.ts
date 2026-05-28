import { existsSync } from "node:fs";
import { resolve } from "node:path";

export type Lang = "bun" | "go";
export type DetectResult = Lang | "ambiguous" | "none";

export function detectLang(cwd: string = process.cwd()): DetectResult {
  const hasGo = existsSync(resolve(cwd, "go.mod"));
  const hasBun = existsSync(resolve(cwd, "package.json"));
  if (hasGo && hasBun) return "ambiguous";
  if (hasGo) return "go";
  if (hasBun) return "bun";
  return "none";
}

export function resolveLang(
  cwd: string = process.cwd(),
  override?: Lang,
): Lang {
  if (override) return override;
  const detected = detectLang(cwd);
  if (detected === "ambiguous") {
    throw new Error(
      "Diretório contém go.mod e package.json. Use --lang=go ou --lang=bun para escolher.",
    );
  }
  if (detected === "none") {
    throw new Error(
      "Não detectei projeto Bun (package.json) nem Go (go.mod) no diretório atual.\nCrie um com 'eco new <nome>' ou 'eco new <nome> --go'.",
    );
  }
  return detected;
}
