import { existsSync } from 'fs'
import { resolve } from 'path'

export interface ToolkitConfig {
  entry: string
  outDir: string
  releaseName: string
  obfuscatorConfig: string
  platforms: Array<'linux' | 'win' | 'macos' | 'macos-arm64'>
}

const defaults: ToolkitConfig = {
  entry: 'src/main.ts',
  outDir: 'dist',
  releaseName: 'app',
  obfuscatorConfig: 'obfuscator.config.cjs',
  platforms: ['linux', 'win'],
}

export async function loadConfig(): Promise<ToolkitConfig> {
  const configPath = resolve(process.cwd(), 'biglaw.config.js')

  if (!existsSync(configPath)) {
    return defaults
  }

  const userConfig = (await import(configPath)).default
  return { ...defaults, ...userConfig }
}
