import pc from "picocolors";

export interface Column {
  header: string;
  align?: "left" | "right";
}

export function renderTable(columns: Column[], rows: string[][]): string {
  const widths = columns.map((col, i) => {
    const cellWidths = rows.map((row) => (row[i] ?? "").length);
    return Math.max(col.header.length, ...cellWidths);
  });

  const pad = (text: string, width: number, align: "left" | "right") =>
    align === "right" ? text.padStart(width) : text.padEnd(width);

  const renderRow = (cells: string[], styler?: (s: string) => string) =>
    cells
      .map((cell, i) => {
        const w = widths[i] ?? 0;
        const align = columns[i]?.align ?? "left";
        const padded = pad(cell, w, align);
        return styler ? styler(padded) : padded;
      })
      .join("  ");

  const header = renderRow(
    columns.map((c) => c.header),
    pc.bold,
  );
  const separator = widths.map((w) => "─".repeat(w)).join("  ");
  const body = rows.map((r) => renderRow(r));

  return [header, pc.dim(separator), ...body].join("\n");
}

export function fileSize(bytes: number): string {
  if (bytes < 1024) return `${bytes} B`;
  if (bytes < 1024 * 1024) return `${(bytes / 1024).toFixed(1)} KB`;
  return `${(bytes / 1024 / 1024).toFixed(2)} MB`;
}
