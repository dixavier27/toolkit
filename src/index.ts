#!/usr/bin/env bun

import { loadConfig } from './config.ts'
import { runPackage } from './commands/package.ts'
import { runObfuscate } from './commands/obfuscate.ts'
import { runRelease } from './commands/release.ts'

export type { ToolkitConfig } from './config.ts'

const commands: Record<string, (cfg: Awaited<ReturnType<typeof loadConfig>>) => Promise<void>> = {
  package:   runPackage,
  obfuscate: runObfuscate,
  release:   runRelease,
}

const command = process.argv[2]

if (!command || !commands[command]) {
  console.error(`Uso: biglaw-scripts <comando>`)
  console.error(`Comandos disponíveis: ${Object.keys(commands).join(', ')}`)
  process.exit(1)
}

const config = await loadConfig()
await commands[command](config)
