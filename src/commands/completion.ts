import type { CommandMeta } from "../utils/command-meta.ts";
import { log } from "../utils/logger.ts";

export const meta: CommandMeta = {
  name: "completion",
  description: "Emite script de autocomplete para shell (bash | zsh | fish)",
  examples: [
    "eco completion bash > ~/.eco-completion.bash",
    "eco completion zsh  > ~/.zsh/completions/_eco",
    "eco completion fish > ~/.config/fish/completions/eco.fish",
  ],
};

const COMMANDS = [
  "init",
  "new",
  "check",
  "doctor",
  "info",
  "config",
  "scripts",
  "ci",
  "package",
  "obfuscate",
  "release",
  "completion",
] as const;

const SUBCOMMANDS: Record<string, string[]> = {
  config: ["show"],
  scripts: ["inject"],
  ci: ["generate"],
  completion: ["bash", "zsh", "fish"],
  new: [
    "--template=cli-tool",
    "--template=library",
    "--template=backend-fastify",
    "--template=frontend-angular-tauri",
  ],
};

const GLOBAL_FLAGS = [
  "--help",
  "--version",
  "--config",
  "--platforms",
  "--verbose",
  "--quiet",
  "--dry-run",
];

function bashScript(): string {
  const cmds = COMMANDS.join(" ");
  const subEntries = Object.entries(SUBCOMMANDS)
    .map(
      ([cmd, subs]) =>
        `    ${cmd}) COMPREPLY=( $(compgen -W "${subs.join(" ")}" -- "$cur") ); return 0 ;;`,
    )
    .join("\n");

  return `# eco bash completion
# Source this file: source ~/.eco-completion.bash
_eco_completion() {
  local cur prev words cword
  _init_completion || return

  if [ "$cword" -eq 1 ]; then
    COMPREPLY=( $(compgen -W "${cmds}" -- "$cur") )
    return 0
  fi

  local cmd="\${words[1]}"
  case "$cmd" in
${subEntries}
  esac

  if [[ "$cur" == --* ]]; then
    COMPREPLY=( $(compgen -W "${GLOBAL_FLAGS.join(" ")}" -- "$cur") )
  fi
}
complete -F _eco_completion eco
`;
}

function zshScript(): string {
  const cmdLines = COMMANDS.map((c) => `    '${c}:eco ${c} command'`).join(
    "\n",
  );
  const subCases = Object.entries(SUBCOMMANDS)
    .map(
      ([cmd, subs]) =>
        `      ${cmd}) _values 'subcommand' ${subs.map((s) => `'${s}'`).join(" ")} ;;`,
    )
    .join("\n");

  return `#compdef eco
# Place in a directory of your $fpath, e.g. ~/.zsh/completions/
_eco() {
  local context curcontext="$curcontext" state line
  typeset -A opt_args

  _arguments -C \\
    '1: :->command' \\
    '*::arg:->args'

  case $state in
    command)
      local -a commands
      commands=(
${cmdLines}
      )
      _describe 'command' commands
      ;;
    args)
      case $words[1] in
${subCases}
      esac
      _values 'global flag' ${GLOBAL_FLAGS.map((f) => `'${f}'`).join(" ")}
      ;;
  esac
}
_eco "$@"
`;
}

function fishScript(): string {
  const cmdLines = COMMANDS.map(
    (c) => `complete -c eco -n "__fish_use_subcommand" -a "${c}"`,
  ).join("\n");
  const subLines = Object.entries(SUBCOMMANDS)
    .flatMap(([cmd, subs]) =>
      subs.map(
        (sub) =>
          `complete -c eco -n "__fish_seen_subcommand_from ${cmd}" -a "${sub}"`,
      ),
    )
    .join("\n");
  const flagLines = GLOBAL_FLAGS.map(
    (f) => `complete -c eco -l "${f.replace(/^--/, "")}"`,
  ).join("\n");

  return `# eco fish completion
${cmdLines}
${subLines}
${flagLines}
`;
}

export function runCompletion(shell: string | undefined) {
  if (!shell || !["bash", "zsh", "fish"].includes(shell)) {
    log.error("Uso: eco completion <bash|zsh|fish>");
    process.exit(1);
  }

  if (shell === "bash") process.stdout.write(bashScript());
  else if (shell === "zsh") process.stdout.write(zshScript());
  else process.stdout.write(fishScript());
}
