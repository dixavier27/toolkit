export interface ParsedArgs {
  command: string | undefined;
  positional: string[];
  flags: {
    help: boolean;
    version: boolean;
    verbose: boolean;
    quiet: boolean;
    dryRun: boolean;
    config: string | undefined;
    platforms: string[] | undefined;
  };
  rest: string[];
}

const BOOLEAN_FLAGS = new Set([
  "--help",
  "-h",
  "--version",
  "-v",
  "--verbose",
  "--quiet",
  "--dry-run",
]);

const VALUE_FLAGS = new Set(["--config", "--platforms"]);

export function parseArgs(argv: string[]): ParsedArgs {
  const result: ParsedArgs = {
    command: undefined,
    positional: [],
    flags: {
      help: false,
      version: false,
      verbose: false,
      quiet: false,
      dryRun: false,
      config: undefined,
      platforms: undefined,
    },
    rest: [],
  };

  for (let i = 0; i < argv.length; i++) {
    const arg = argv[i];
    if (arg === undefined) continue;

    if (arg === "--help" || arg === "-h") {
      result.flags.help = true;
      continue;
    }
    if (arg === "--version" || arg === "-v") {
      result.flags.version = true;
      continue;
    }
    if (arg === "--verbose") {
      result.flags.verbose = true;
      continue;
    }
    if (arg === "--quiet") {
      result.flags.quiet = true;
      continue;
    }
    if (arg === "--dry-run") {
      result.flags.dryRun = true;
      continue;
    }

    if (arg.startsWith("--config=")) {
      result.flags.config = arg.slice("--config=".length);
      continue;
    }
    if (arg === "--config") {
      const next = argv[++i];
      if (next) result.flags.config = next;
      continue;
    }

    if (arg.startsWith("--platforms=")) {
      result.flags.platforms = arg
        .slice("--platforms=".length)
        .split(",")
        .map((p) => p.trim())
        .filter(Boolean);
      continue;
    }
    if (arg === "--platforms") {
      const next = argv[++i];
      if (next)
        result.flags.platforms = next
          .split(",")
          .map((p) => p.trim())
          .filter(Boolean);
      continue;
    }

    if (arg.startsWith("-")) {
      result.rest.push(arg);
      continue;
    }

    if (result.command === undefined) {
      result.command = arg;
    } else {
      result.positional.push(arg);
    }
  }

  return result;
}

export function _isKnownFlag(arg: string): boolean {
  return BOOLEAN_FLAGS.has(arg) || VALUE_FLAGS.has(arg);
}
