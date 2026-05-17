import { $ } from 'bun'
import type { ToolkitConfig } from '../config.ts'

const bunTargets: Record<string, string> = {
  linux:       'bun-linux-x64',
  win:         'bun-windows-x64',
  macos:       'bun-darwin-x64',
  'macos-arm64': 'bun-darwin-arm64',
}

export async function runPackage(config: ToolkitConfig) {
  for (const platform of config.platforms) {
    const target = bunTargets[platform]
    const outfile = `${config.outDir}/bundle-${platform}.js`

    console.log(`📦 Bundling ${platform} → ${outfile}`)
    await $`bun build ${config.entry} --outfile ${outfile} --target=${target}`
  }
}
