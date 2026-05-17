import { createHash } from "node:crypto";
import { readFile, writeFile } from "node:fs/promises";
import { basename } from "node:path";

export async function sha256File(path: string): Promise<string> {
  const data = await readFile(path);
  return createHash("sha256").update(data).digest("hex");
}

export async function writeChecksumsFile(
  files: string[],
  outputPath: string,
): Promise<void> {
  const lines: string[] = [];
  for (const file of files) {
    const hash = await sha256File(file);
    lines.push(`${hash}  ${basename(file)}`);
  }
  await writeFile(outputPath, `${lines.join("\n")}\n`, "utf8");
}
