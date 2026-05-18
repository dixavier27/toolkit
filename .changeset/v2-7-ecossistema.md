---
"@dixavier27/eco": minor
---

## v2.7 — Ecossistema

### Composite GitHub Action

Nova action em `composite-action/`, invocada como `dixavier27/toolkit/composite-action@v2.7.0`. Substitui ~30 linhas de YAML (checkout + setup-bun + install + bunx eco + upload-artifact) por 3 linhas.

```yaml
- uses: dixavier27/toolkit/composite-action@v2.7.0
  with:
    command: release
    platforms: linux,win
    upload-artifacts: true
```

Inputs: `command`, `platforms`, `bun-version`, `working-directory`, `install`, `upload-artifacts`, `artifact-name`. Output: `release-dir`.

### Comando `eco completion`

Emite script de autocomplete para `bash`, `zsh` ou `fish`:

```bash
eco completion bash > ~/.eco-completion.bash
eco completion zsh  > ~/.zsh/completions/_eco
eco completion fish > ~/.config/fish/completions/eco.fish
```

Completa comandos, subcomandos (config show, scripts inject, ci generate) e flags globais.

### Dogfooding

`eco ci generate` agora produz workflows que usam o próprio Composite Action — projetos novos consomem o pattern correto desde o primeiro dia.

### Próximas releases

- **v2.8** — code signing (Windows signtool, macOS codesign + notarization)
- **v2.9** — docs site (Astro/VitePress)
