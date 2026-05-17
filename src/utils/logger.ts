import pc from "picocolors";

export type LogLevel = "silent" | "normal" | "verbose";

let level: LogLevel = "normal";

export function setLogLevel(newLevel: LogLevel) {
  level = newLevel;
}

export function getLogLevel(): LogLevel {
  return level;
}

export const log = {
  info(...args: unknown[]) {
    if (level !== "silent") console.log(...args);
  },
  verbose(...args: unknown[]) {
    if (level === "verbose") console.log(pc.dim(args.join(" ")));
  },
  warn(...args: unknown[]) {
    if (level !== "silent") console.warn(pc.yellow(args.join(" ")));
  },
  error(...args: unknown[]) {
    console.error(pc.red(args.join(" ")));
  },
  success(...args: unknown[]) {
    if (level !== "silent") console.log(pc.green(args.join(" ")));
  },
  step(...args: unknown[]) {
    if (level !== "silent") console.log(pc.cyan(args.join(" ")));
  },
  dim(...args: unknown[]) {
    if (level !== "silent") console.log(pc.dim(args.join(" ")));
  },
};

export { pc };
