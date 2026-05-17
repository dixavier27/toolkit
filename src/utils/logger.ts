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
    if (level === "verbose") console.log(...args);
  },
  warn(...args: unknown[]) {
    if (level !== "silent") console.warn(...args);
  },
  error(...args: unknown[]) {
    console.error(...args);
  },
};
