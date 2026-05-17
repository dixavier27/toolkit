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
  // Exemplo de cópia declarativa para apps com assets estáticos:
  // assets: [
  //   { from: 'node_modules/@fastify/swagger-ui/static', to: 'dist/static' },
  // ],
}
