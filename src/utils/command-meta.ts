import pc from "picocolors";

export interface FlagMeta {
  name: string;
  description: string;
}

export interface CommandMeta {
  name: string;
  description: string;
  flags?: FlagMeta[];
  examples?: string[];
}

export function renderCommandHelp(meta: CommandMeta): string {
  const lines: string[] = [];
  lines.push(`${pc.bold(`eco ${meta.name}`)} — ${meta.description}`);
  lines.push("");
  lines.push(`Uso:  eco ${meta.name} [flags]`);

  if (meta.flags && meta.flags.length > 0) {
    lines.push("");
    lines.push(pc.bold("Flags:"));
    const maxLen = Math.max(...meta.flags.map((f) => f.name.length));
    for (const f of meta.flags) {
      lines.push(`  ${f.name.padEnd(maxLen + 2)}${f.description}`);
    }
  }

  if (meta.examples && meta.examples.length > 0) {
    lines.push("");
    lines.push(pc.bold("Exemplos:"));
    for (const ex of meta.examples) {
      lines.push(`  ${pc.dim("$")} ${ex}`);
    }
  }

  return lines.join("\n");
}
