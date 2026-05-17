import { existsSync } from "node:fs";
import { resolve } from "node:path";
import { z } from "zod";

const PlatformSchema = z.enum(["linux", "win", "macos", "macos-arm64"]);
export type Platform = z.infer<typeof PlatformSchema>;

const SourcemapSchema = z
  .union([z.literal(false), z.literal("inline"), z.literal("external")])
  .default(false);

const AssetSchema = z.object({
  from: z.string(),
  to: z.string(),
});
export type Asset = z.infer<typeof AssetSchema>;

const BaseConfigSchema = z.object({
  entry: z.string().default("src/main.ts"),
  outDir: z.string().default("dist"),
  bundleName: z.string().default("bundle.js"),
  releaseName: z.string().default("app"),
  obfuscatorConfig: z.string().default("obfuscator.config.cjs"),
  platforms: z.array(PlatformSchema).default(["linux", "win"]),
  sourcemap: SourcemapSchema,
  assets: z.array(AssetSchema).default([]),
  define: z.record(z.string()).default({}),
  embedVersion: z.boolean().default(true),
  parallel: z.boolean().default(true),
  checksums: z.boolean().default(false),
});

type BaseConfig = z.infer<typeof BaseConfigSchema>;

export type EcoHook = (config: EcoConfig) => void | Promise<void>;

export interface EcoConfig extends BaseConfig {
  afterPackage?: EcoHook;
  afterObfuscate?: EcoHook;
  afterRelease?: EcoHook;
}

export const EcoConfigSchema = BaseConfigSchema.extend({
  afterPackage: z.function().optional() as z.ZodOptional<z.ZodType<EcoHook>>,
  afterObfuscate: z.function().optional() as z.ZodOptional<z.ZodType<EcoHook>>,
  afterRelease: z.function().optional() as z.ZodOptional<z.ZodType<EcoHook>>,
});

const DEFAULT_CONFIG_NAMES = [
  "eco.config.ts",
  "eco.config.js",
  "eco.config.mjs",
];

function findConfigFile(cwd: string, override?: string): string | undefined {
  if (override) {
    const path = resolve(cwd, override);
    return existsSync(path) ? path : undefined;
  }
  for (const name of DEFAULT_CONFIG_NAMES) {
    const path = resolve(cwd, name);
    if (existsSync(path)) return path;
  }
  return undefined;
}

export async function loadConfig(configPath?: string): Promise<EcoConfig> {
  const cwd = process.cwd();
  const resolved = findConfigFile(cwd, configPath);

  if (!resolved) {
    return EcoConfigSchema.parse({}) as EcoConfig;
  }

  const mod = await import(resolved);
  const userConfig = mod.default ?? mod;

  const parsed = EcoConfigSchema.safeParse(userConfig);
  if (!parsed.success) {
    const issues = parsed.error.issues
      .map(
        (i: z.ZodIssue) => `  - ${i.path.join(".") || "<root>"}: ${i.message}`,
      )
      .join("\n");
    throw new Error(
      `Configuração inválida em ${resolved}:\n${issues}\n\nConsulte 'eco --help' ou rode 'eco check'.`,
    );
  }

  return parsed.data as EcoConfig;
}

export function getConfigPath(configPath?: string): string | undefined {
  return findConfigFile(process.cwd(), configPath);
}
