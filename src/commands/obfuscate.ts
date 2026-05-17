import { $ } from 'bun'
import type { ToolkitConfig } from '../config.ts'
import { runPackage } from './package.ts'

export async function runObfuscate(config: ToolkitConfig) {
  await runPackage(config)

  for (const platform of config.platforms) {
    const bundle = `${config.outDir}/bundle-${platform}.js`

    console.log(`🔒 Obfuscating ${bundle}`)
    await $`javascript-obfuscator ${bundle} --output ${bundle} --config ${config.obfuscatorConfig}`
  }
}
