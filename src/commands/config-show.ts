import type { EcoConfig } from "../config.ts";
import type { CommandMeta } from "../utils/command-meta.ts";
import { log, pc } from "../utils/logger.ts";

export const meta: CommandMeta = {
  name: "config",
  description: "Inspeciona a configuração resolvida (com defaults aplicados)",
  flags: [{ name: "show", description: "Imprime o config em JSON" }],
  examples: ["eco config show", "eco config show --config=custom.config.js"],
};

export function runConfigShow(config: EcoConfig) {
  log.info(pc.bold("Config resolvido:"));
  log.info("");

  const serializable = {
    entry: config.entry,
    outDir: config.outDir,
    bundleName: config.bundleName,
    releaseName: config.releaseName,
    obfuscatorConfig: config.obfuscatorConfig,
    platforms: config.platforms,
    afterPackage: config.afterPackage ? "<function>" : undefined,
    afterObfuscate: config.afterObfuscate ? "<function>" : undefined,
    afterRelease: config.afterRelease ? "<function>" : undefined,
  };

  const json = JSON.stringify(serializable, null, 2);
  log.info(colorize(json));
}

function colorize(json: string): string {
  return json
    .replace(/"([^"]+)":/g, (_, key) => `${pc.cyan(`"${key}"`)}:`)
    .replace(/: "([^"]*)"/g, (_, val) => `: ${pc.green(`"${val}"`)}`)
    .replace(/: (true|false|null)/g, (_, val) => `: ${pc.yellow(val)}`)
    .replace(/: (\d+)/g, (_, val) => `: ${pc.yellow(val)}`);
}
