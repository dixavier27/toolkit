import { $ } from "bun";
import type { ToolkitConfig } from "../config.ts";

export async function runObfuscate(
  config: ToolkitConfig,
  _flags: string[] = [],
) {
  for (const platform of config.platforms) {
    const bundle = `${config.outDir}/bundle-${platform}.js`;

    console.log(`🔒 Obfuscating ${bundle}`);
    await $`javascript-obfuscator ${bundle} --output ${bundle} --config ${config.obfuscatorConfig}`;
  }
}
