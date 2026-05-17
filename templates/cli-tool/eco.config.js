/** @type {import('@dixavier27/eco').EcoConfig} */
export default {
  entry: 'src/main.ts',
  outDir: 'dist',
  bundleName: 'bundle.js',
  releaseName: '{{name}}',
  obfuscatorConfig: 'obfuscator.config.cjs',
  platforms: ['linux', 'win'],
  embedVersion: true,
  checksums: true,
}
