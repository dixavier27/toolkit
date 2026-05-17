import { cp, mkdir } from "node:fs/promises";
import { dirname, resolve } from "node:path";
import type { Asset } from "../config.ts";
import { log, pc } from "./logger.ts";

export async function copyAssets(assets: Asset[], cwd = process.cwd()) {
  if (assets.length === 0) return;

  for (const asset of assets) {
    const from = resolve(cwd, asset.from);
    const to = resolve(cwd, asset.to);

    log.verbose(`📁 Copiando ${pc.dim(asset.from)} → ${pc.dim(asset.to)}`);

    await mkdir(dirname(to), { recursive: true });
    await cp(from, to, { recursive: true, force: true });
  }

  log.info(
    `${pc.green("📁")} ${pc.bold(`${assets.length}`)} ${pc.dim(`asset${assets.length !== 1 ? "s" : ""} copiado${assets.length !== 1 ? "s" : ""}`)}`,
  );
}
