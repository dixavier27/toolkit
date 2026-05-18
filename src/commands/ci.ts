import { existsSync, mkdirSync, writeFileSync } from "node:fs";
import { resolve } from "node:path";
import type { EcoConfig, Platform } from "../config.ts";
import type { CommandMeta } from "../utils/command-meta.ts";
import { log, pc } from "../utils/logger.ts";

export const meta: CommandMeta = {
  name: "ci",
  description: "Gera workflows GitHub Actions (ci.yml e release.yml)",
  flags: [
    { name: "generate", description: "Subcomando: gera os workflows" },
    {
      name: "--only=<ci|release>",
      description: "Gera apenas um dos workflows",
    },
    { name: "--update", description: "Sobrescreve arquivos existentes" },
  ],
  examples: [
    "eco ci generate",
    "eco ci generate --only=release",
    "eco ci generate --update",
  ],
};

const platformToRunner: Record<Platform, string> = {
  linux: "ubuntu-latest",
  win: "windows-latest",
  macos: "macos-latest",
  "macos-arm64": "macos-latest",
};

const ciTemplate = () =>
  `name: CI

on:
  push:
    branches: [main, master]
  pull_request:
    types: [opened, synchronize, reopened]
  workflow_dispatch:

concurrency:
  group: \${{ github.workflow }}-\${{ github.ref }}
  cancel-in-progress: true

permissions:
  contents: read

jobs:
  ci:
    name: Typecheck + lint + build
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - uses: oven-sh/setup-bun@v2
        with:
          bun-version: '1.2.0'
      - run: bun install --frozen-lockfile
      - run: bun run typecheck
      - run: bun run lint
      - run: bun run build
`;

const releaseTemplate = (config: EcoConfig) => {
  const matrix = config.platforms
    .map((p) => `          - { os: ${platformToRunner[p]}, platform: ${p} }`)
    .join("\n");

  return `name: Release

on:
  push:
    tags: ['v*']
  workflow_dispatch:

concurrency:
  group: \${{ github.workflow }}-\${{ github.ref }}
  cancel-in-progress: false

jobs:
  release:
    name: Release \${{ matrix.platform }}
    runs-on: \${{ matrix.os }}
    permissions:
      contents: write
    strategy:
      fail-fast: false
      matrix:
        include:
${matrix}
    steps:
      - uses: actions/checkout@v4
      - uses: dixavier27/toolkit/composite-action@v2.7.0
        with:
          command: release --keep-going
          platforms: \${{ matrix.platform }}
      - uses: softprops/action-gh-release@v2
        with:
          files: |
            release/${config.releaseName}-\${{ matrix.platform }}*
            release/checksums.txt
`;
};

export interface CiOptions {
  only?: "ci" | "release";
  update?: boolean;
}

export function runCiGenerate(config: EcoConfig, opts: CiOptions = {}) {
  const cwd = process.cwd();
  const workflowsDir = resolve(cwd, ".github", "workflows");
  mkdirSync(workflowsDir, { recursive: true });

  const tasks: Array<{ file: string; content: string }> = [];

  if (opts.only !== "release") {
    tasks.push({
      file: resolve(workflowsDir, "ci.yml"),
      content: ciTemplate(),
    });
  }

  if (opts.only !== "ci") {
    tasks.push({
      file: resolve(workflowsDir, "release.yml"),
      content: releaseTemplate(config),
    });
  }

  const created: string[] = [];
  const skipped: string[] = [];
  const updated: string[] = [];

  for (const task of tasks) {
    const relativePath = task.file
      .replace(`${cwd}\\`, "")
      .replace(`${cwd}/`, "");
    if (existsSync(task.file) && !opts.update) {
      skipped.push(relativePath);
      continue;
    }
    const existed = existsSync(task.file);
    writeFileSync(task.file, task.content, "utf8");
    (existed ? updated : created).push(relativePath);
  }

  if (created.length > 0) {
    log.success(
      `✅ ${created.length} workflow${created.length !== 1 ? "s" : ""} criado${created.length !== 1 ? "s" : ""}:`,
    );
    for (const f of created) log.info(`   ${pc.cyan(f)}`);
  }

  if (updated.length > 0) {
    log.warn(
      `⚠️  ${updated.length} workflow${updated.length !== 1 ? "s" : ""} sobrescrito${updated.length !== 1 ? "s" : ""}:`,
    );
    for (const f of updated) log.info(`   ${pc.cyan(f)}`);
  }

  if (skipped.length > 0) {
    log.info(
      pc.dim(`Skipped (já existem, use --update): ${skipped.join(", ")}`),
    );
  }

  if (created.length + updated.length > 0) {
    log.info("");
    log.info(pc.bold("Próximos passos:"));
    log.info(`  ${pc.dim("1.")} Revise os workflows gerados`);
    log.info(`  ${pc.dim("2.")} Commit e push`);
    log.info(
      `  ${pc.dim("3.")} Crie uma tag v* para disparar o release: ${pc.cyan("git tag v1.0.0 && git push --tags")}`,
    );
  }
}
