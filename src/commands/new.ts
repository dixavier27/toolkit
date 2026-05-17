import {
  existsSync,
  mkdirSync,
  readdirSync,
  readFileSync,
  statSync,
  writeFileSync,
} from "node:fs";
import { dirname, join, resolve } from "node:path";
import { fileURLToPath } from "node:url";
import type { CommandMeta } from "../utils/command-meta.ts";
import { log, pc } from "../utils/logger.ts";

export const meta: CommandMeta = {
  name: "new",
  description:
    "Cria um projeto novo a partir de um template curado (cli-tool, library, backend-fastify, frontend-angular-tauri)",
  flags: [
    {
      name: "--template=<tipo>",
      description:
        "Template: cli-tool | library | backend-fastify | frontend-angular-tauri (default: cli-tool)",
    },
    { name: "--force", description: "Sobrescreve diretório existente" },
  ],
  examples: [
    "eco new minha-cli",
    "eco new minha-lib --template=library",
    "eco new minha-api --template=backend-fastify",
    "eco new meu-app --template=frontend-angular-tauri",
  ],
};

const AVAILABLE_TEMPLATES = [
  "cli-tool",
  "library",
  "backend-fastify",
  "frontend-angular-tauri",
] as const;
type TemplateName = (typeof AVAILABLE_TEMPLATES)[number];

function findTemplatesDir(): string {
  const here = dirname(fileURLToPath(import.meta.url));
  const candidates = [
    resolve(here, "..", "..", "templates"),
    resolve(here, "..", "templates"),
  ];
  for (const c of candidates) {
    if (existsSync(c)) return c;
  }
  throw new Error(
    `Diretório de templates não encontrado. Procurado em:\n${candidates.map((c) => `  - ${c}`).join("\n")}`,
  );
}

function isValidProjectName(name: string): boolean {
  return /^[a-z0-9][a-z0-9._-]*$/i.test(name) && !name.startsWith(".");
}

function copyTemplate(
  srcDir: string,
  destDir: string,
  replacements: Record<string, string>,
) {
  for (const entry of readdirSync(srcDir)) {
    const srcPath = join(srcDir, entry);
    const destPath = join(destDir, entry);
    const stat = statSync(srcPath);

    if (stat.isDirectory()) {
      mkdirSync(destPath, { recursive: true });
      copyTemplate(srcPath, destPath, replacements);
      continue;
    }

    const content = readFileSync(srcPath, "utf8");
    const replaced = Object.entries(replacements).reduce(
      (acc, [key, value]) => acc.replaceAll(`{{${key}}}`, value),
      content,
    );
    writeFileSync(destPath, replaced, "utf8");
  }
}

export interface NewOptions {
  projectName: string;
  template?: TemplateName;
  force?: boolean;
}

export function runNew(opts: NewOptions) {
  const { projectName, force } = opts;
  const template = opts.template ?? "cli-tool";

  if (!isValidProjectName(projectName)) {
    log.error(
      `❌ Nome inválido: "${projectName}". Use letras, números, "-", "_" ou "." (sem começar com ".")`,
    );
    process.exit(1);
  }

  if (!AVAILABLE_TEMPLATES.includes(template)) {
    log.error(
      `❌ Template desconhecido: "${template}". Disponíveis: ${AVAILABLE_TEMPLATES.join(", ")}`,
    );
    process.exit(1);
  }

  const templatesDir = findTemplatesDir();
  const srcDir = join(templatesDir, template);
  const destDir = resolve(process.cwd(), projectName);

  if (existsSync(destDir) && !force) {
    log.error(
      `❌ Diretório ${destDir} já existe. Use --force para sobrescrever.`,
    );
    process.exit(1);
  }

  mkdirSync(destDir, { recursive: true });

  log.info(
    pc.cyan(`📋 Copiando template ${pc.bold(template)} → ${projectName}/`),
  );
  copyTemplate(srcDir, destDir, { name: projectName });

  log.success(
    `\n✅ Projeto ${pc.bold(projectName)} criado a partir de "${template}"`,
  );
  log.info("");
  log.info(pc.bold("Próximos passos:"));
  log.info(`  ${pc.dim("1.")} ${pc.cyan(`cd ${projectName}`)}`);
  log.info(`  ${pc.dim("2.")} ${pc.cyan("bun install")}`);
  if (template === "cli-tool") {
    log.info(`  ${pc.dim("3.")} ${pc.cyan("bun run dev hello")}`);
    log.info(
      `  ${pc.dim("4.")} ${pc.cyan("bunx eco ci generate")} ${pc.dim("(opcional, para CI/CD)")}`,
    );
  } else if (template === "backend-fastify") {
    log.info(`  ${pc.dim("3.")} ${pc.cyan("cp .env.example .env")}`);
    log.info(`  ${pc.dim("4.")} ${pc.cyan("bun run dev")}`);
    log.info(
      `  ${pc.dim("5.")} ${pc.cyan("curl http://localhost:3000/hello/mundo")}`,
    );
  } else if (template === "frontend-angular-tauri") {
    log.info(`  ${pc.dim("3.")} ${pc.cyan("bun run dev")} ${pc.dim("(web)")}`);
    log.info(
      `  ${pc.dim("4.")} ${pc.cyan("bun run tauri:dev")} ${pc.dim("(desktop, requer Rust)")}`,
    );
    log.info(
      `  ${pc.dim("5.")} ${pc.cyan("bunx @tauri-apps/cli icon logo.png")} ${pc.dim("(gera ícones antes do build)")}`,
    );
  } else {
    log.info(`  ${pc.dim("3.")} ${pc.cyan("bun test")}`);
    log.info(`  ${pc.dim("4.")} ${pc.cyan("bun run build")}`);
  }
}
