#!/usr/bin/env bun

declare const __VERSION__: string;

const args = process.argv.slice(2);

if (args.includes("--version") || args.includes("-v")) {
  console.log(__VERSION__);
  process.exit(0);
}

if (args.includes("--help") || args.includes("-h") || args.length === 0) {
  console.log(`{{name}} v${__VERSION__}

Uso: {{name}} <comando>

Comandos:
  hello   Mensagem de boas-vindas
`);
  process.exit(0);
}

if (args[0] === "hello") {
  console.log("Olá, mundo! 🌱");
  process.exit(0);
}

console.error(`Comando desconhecido: ${args[0]}`);
process.exit(1);
