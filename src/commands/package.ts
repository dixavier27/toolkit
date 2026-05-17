import { $ } from 'bun'
import { mkdirSync } from 'fs'
import type { ToolkitConfig } from '../config.ts'

export async function runPackage(config: ToolkitConfig, _flags: string[] = []) {
  mkdirSync(config.outDir, { recursive: true })

  for (const platform of config.platforms) {
    const outfile = `${config.outDir}/bundle-${platform}.js`

    console.log(`📦 Bundling ${platform} → ${outfile}`)
    await $`bun build ${config.entry} --outfile ${outfile} --target=bun --minify`
  }
}
