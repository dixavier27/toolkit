import { $ } from 'bun'
import { mkdirSync } from 'fs'
import type { ToolkitConfig } from '../config.ts'
import { runObfuscate } from './obfuscate.ts'

const bunTargets: Record<string, string> = {
  linux:       'bun-linux-x64',
  win:         'bun-windows-x64',
  macos:       'bun-darwin-x64',
  'macos-arm64': 'bun-darwin-arm64',
}

const ext: Record<string, string> = {
  win: '.exe',
}

export async function runRelease(config: ToolkitConfig) {
  await runObfuscate(config)

  mkdirSync('release', { recursive: true })

  for (const platform of config.platforms) {
    const bundle  = `${config.outDir}/bundle-${platform}.js`
    const target  = bunTargets[platform]
    const outfile = `release/${config.releaseName}-${platform}${ext[platform] ?? ''}`

    console.log(`🚀 Compiling ${platform} → ${outfile}`)
    await $`bun build ${bundle} --compile --target=${target} --outfile ${outfile}`
  }
}
