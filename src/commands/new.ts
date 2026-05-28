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
    "Cria um projeto novo a partir de um template curado (cli-tool, library, backend-fastify, frontend-angular-tauri, go)",
  flags: [
    {
      name: "--template=<tipo>",
      description:
        "Template Bun: cli-tool | library | backend-fastify | frontend-angular-tauri (default: cli-tool)",
    },
    {
      name: "--go",
      description: "Cria projeto Go (combine com --cli e/ou --api)",
    },
    {
      name: "--cli",
      description:
        "Inclui interface de linha de comando (Cobra). Default quando --go é usado",
    },
    {
      name: "--api",
      description: "Inclui API REST (net/http stdlib + validator)",
    },
    {
      name: "--module=<path>",
      description: "Path do módulo Go (default: example.com/<nome>)",
    },
    { name: "--force", description: "Sobrescreve diretório existente" },
  ],
  examples: [
    "eco new minha-cli",
    "eco new minha-lib --template=library",
    "eco new minha-api --template=backend-fastify",
    "eco new meu-app --template=frontend-angular-tauri",
    "eco new minha-cli --go",
    "eco new minha-api --go --api",
    "eco new meu-app --go --cli --api --module=github.com/user/meu-app",
  ],
};

const AVAILABLE_TEMPLATES = [
  "cli-tool",
  "library",
  "backend-fastify",
  "frontend-angular-tauri",
] as const;
type TemplateName = (typeof AVAILABLE_TEMPLATES)[number];

type GoVariant = "cli" | "api" | "cli-api";

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

function applyReplacements(
  content: string,
  replacements: Record<string, string>,
): string {
  return Object.entries(replacements).reduce(
    (acc, [key, value]) => acc.replaceAll(`{{${key}}}`, value),
    content,
  );
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
    writeFileSync(destPath, applyReplacements(content, replacements), "utf8");
  }
}

function copyFileWithReplacements(
  srcFile: string,
  destFile: string,
  replacements: Record<string, string>,
) {
  mkdirSync(dirname(destFile), { recursive: true });
  const content = readFileSync(srcFile, "utf8");
  writeFileSync(destFile, applyReplacements(content, replacements), "utf8");
}

function copyDirSelective(
  srcDir: string,
  destDir: string,
  replacements: Record<string, string>,
  skipNames: Set<string>,
) {
  for (const entry of readdirSync(srcDir)) {
    if (skipNames.has(entry)) continue;
    const srcPath = join(srcDir, entry);
    const destPath = join(destDir, entry);
    const stat = statSync(srcPath);

    if (stat.isDirectory()) {
      mkdirSync(destPath, { recursive: true });
      copyDirSelective(srcPath, destPath, replacements, skipNames);
      continue;
    }

    const content = readFileSync(srcPath, "utf8");
    writeFileSync(destPath, applyReplacements(content, replacements), "utf8");
  }
}

function copyGoTemplate(
  srcDir: string,
  destDir: string,
  projectName: string,
  moduleName: string,
  variant: GoVariant,
) {
  const replacements = { name: projectName, module: moduleName };
  const includesCli = variant === "cli" || variant === "cli-api";
  const includesApi = variant === "api" || variant === "cli-api";

  const skip = new Set<string>(["_variants"]);
  if (!includesCli) skip.add("cli");
  if (!includesApi) skip.add("api");

  // copia tudo exceto _variants e os subdirs internal/* não usados
  for (const entry of readdirSync(srcDir)) {
    if (entry === "_variants") continue;
    const srcPath = join(srcDir, entry);
    const destPath = join(destDir, entry);
    const stat = statSync(srcPath);

    if (stat.isDirectory()) {
      if (entry === "internal") {
        // filtra cli/ ou api/ conforme variant
        const internalSkip = new Set<string>();
        if (!includesCli) internalSkip.add("cli");
        if (!includesApi) internalSkip.add("api");
        mkdirSync(destPath, { recursive: true });
        copyDirSelective(srcPath, destPath, replacements, internalSkip);
      } else {
        mkdirSync(destPath, { recursive: true });
        copyTemplate(srcPath, destPath, replacements);
      }
      continue;
    }

    const content = readFileSync(srcPath, "utf8");
    writeFileSync(destPath, applyReplacements(content, replacements), "utf8");
  }

  const variants = join(srcDir, "_variants");

  // main.go → cmd/<name>/main.go
  copyFileWithReplacements(
    join(variants, `main.${variant}.go`),
    join(destDir, "cmd", projectName, "main.go"),
    replacements,
  );

  // go.mod
  copyFileWithReplacements(
    join(variants, `go.mod.${variant}`),
    join(destDir, "go.mod"),
    replacements,
  );

  // Makefile
  copyFileWithReplacements(
    join(variants, `Makefile.${variant}`),
    join(destDir, "Makefile"),
    replacements,
  );

  // README.md
  copyFileWithReplacements(
    join(variants, `README.${variant}.md`),
    join(destDir, "README.md"),
    replacements,
  );

  // .air.toml apenas se há API (server precisa de hot-reload)
  if (includesApi) {
    copyFileWithReplacements(
      join(variants, ".air.toml"),
      join(destDir, ".air.toml"),
      replacements,
    );
  }

  // serve_cmd.go apenas em cli-api (api pura usa api.Run direto)
  if (variant === "cli-api") {
    copyFileWithReplacements(
      join(variants, "serve_cmd.go"),
      join(destDir, "internal", "api", "serve_cmd.go"),
      replacements,
    );
  }
}

function resolveGoVariant(cli: boolean, api: boolean): GoVariant {
  if (cli && api) return "cli-api";
  if (api) return "api";
  return "cli";
}

export interface NewOptions {
  projectName: string;
  template?: TemplateName;
  force?: boolean;
  go?: boolean;
  cli?: boolean;
  api?: boolean;
  module?: string;
}

export function runNew(opts: NewOptions) {
  const { projectName, force } = opts;

  if (!isValidProjectName(projectName)) {
    log.error(
      `❌ Nome inválido: "${projectName}". Use letras, números, "-", "_" ou "." (sem começar com ".")`,
    );
    process.exit(1);
  }

  if (opts.go && opts.template) {
    log.error(
      `❌ --go é incompatível com --template. Use --cli e/ou --api para projetos Go.`,
    );
    process.exit(1);
  }

  if (!opts.go && (opts.cli || opts.api)) {
    log.error(`❌ --cli e --api só são válidas em conjunto com --go.`);
    process.exit(1);
  }

  const templatesDir = findTemplatesDir();
  const destDir = resolve(process.cwd(), projectName);

  if (existsSync(destDir) && !force) {
    log.error(
      `❌ Diretório ${destDir} já existe. Use --force para sobrescrever.`,
    );
    process.exit(1);
  }

  mkdirSync(destDir, { recursive: true });

  if (opts.go) {
    const variant = resolveGoVariant(opts.cli ?? false, opts.api ?? false);
    const moduleName = opts.module ?? `example.com/${projectName}`;
    const srcDir = join(templatesDir, "go-app");

    log.info(
      pc.cyan(
        `📋 Copiando template ${pc.bold(`go (${variant})`)} → ${projectName}/`,
      ),
    );
    copyGoTemplate(srcDir, destDir, projectName, moduleName, variant);

    log.success(
      `\n✅ Projeto Go ${pc.bold(projectName)} criado (variante "${variant}", módulo ${moduleName})`,
    );
    log.info("");
    log.info(pc.bold("Próximos passos:"));
    log.info(`  ${pc.dim("1.")} ${pc.cyan(`cd ${projectName}`)}`);
    log.info(`  ${pc.dim("2.")} ${pc.cyan("go mod tidy")}`);
    if (variant === "cli") {
      log.info(
        `  ${pc.dim("3.")} ${pc.cyan(`go run ./cmd/${projectName} hello mundo`)}`,
      );
    } else if (variant === "api") {
      log.info(
        `  ${pc.dim("3.")} ${pc.cyan("make dev")} ${pc.dim("(requer air)")}`,
      );
      log.info(
        `  ${pc.dim("4.")} ${pc.cyan("curl http://localhost:3000/hello/mundo")}`,
      );
    } else {
      log.info(
        `  ${pc.dim("3.")} ${pc.cyan(`go run ./cmd/${projectName} hello mundo`)}`,
      );
      log.info(
        `  ${pc.dim("4.")} ${pc.cyan(`go run ./cmd/${projectName} serve`)} ${pc.dim("(sobe API)")}`,
      );
    }
    log.info(
      `  ${pc.dim("•")} ${pc.cyan("eco release")} ${pc.dim("para cross-compile")}`,
    );
    return;
  }

  const template = opts.template ?? "cli-tool";

  if (!AVAILABLE_TEMPLATES.includes(template)) {
    log.error(
      `❌ Template desconhecido: "${template}". Disponíveis: ${AVAILABLE_TEMPLATES.join(", ")}`,
    );
    process.exit(1);
  }

  const srcDir = join(templatesDir, template);

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
