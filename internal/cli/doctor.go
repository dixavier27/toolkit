package cli

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"text/tabwriter"

	"github.com/spf13/cobra"

	"github.com/dixavier27/eco/internal/gocheck"
)

func newDoctorCmd() *cobra.Command {
	var fix bool

	cmd := &cobra.Command{
		Use:   "doctor",
		Short: "Inspeciona o ambiente Go e ferramentas auxiliares",
		Long: "Verifica a presença e versão da toolchain Go, do golangci-lint, do air\n" +
			"e detecta projeto Go local (go.mod). Reporta status em formato tabular.",
		RunE: func(cmd *cobra.Command, args []string) error {
			out := cmd.OutOrStdout()
			cwd, _ := os.Getwd()

			checks := gocheck.CheckAll(cwd)
			printChecks(out, checks)

			if fix {
				if err := applyFixes(out, cwd, checks); err != nil {
					return err
				}
			}

			// Falha se houver qualquer check em estado crítico (missing/error)
			// para a toolchain Go base. Warns e missing de ferramentas opcionais
			// (air, golangci-lint) NÃO causam exit != 0 — só informam.
			for _, c := range checks {
				if c.Name == "go" && !c.OK() {
					return fmt.Errorf("toolchain Go indisponível ou abaixo do mínimo")
				}
			}
			return nil
		},
	}

	cmd.Flags().BoolVar(&fix, "fix", false, "cria arquivos de configuração ausentes (ex.: .golangci.yml)")
	return cmd
}

func printChecks(w io.Writer, checks []gocheck.Check) {
	tw := tabwriter.NewWriter(w, 0, 0, 2, ' ', 0)
	for _, c := range checks {
		line := fmt.Sprintf("%s  %s\t%s", c.Status.Symbol(), c.Name, c.Message)
		fmt.Fprintln(tw, line)
		if !c.OK() && c.Suggestion != "" {
			fmt.Fprintf(tw, "   \t  → %s\n", c.Suggestion)
		}
	}
	tw.Flush()
}

// applyFixes cria arquivos de configuração faltantes que o eco sabe gerar
// sem ambiguidade. Não instala binários nem altera o sistema.
func applyFixes(w io.Writer, dir string, checks []gocheck.Check) error {
	// Cria .golangci.yml se faltar e houver projeto Go local.
	hasGoMod := false
	for _, c := range checks {
		if c.Name == "go.mod" && c.OK() {
			hasGoMod = true
		}
	}
	if !hasGoMod {
		return nil
	}

	cfgPath := filepath.Join(dir, ".golangci.yml")
	if _, err := os.Stat(cfgPath); err == nil {
		return nil
	}
	if err := os.WriteFile(cfgPath, []byte(defaultGolangciYAML), 0o644); err != nil {
		return fmt.Errorf("criar .golangci.yml: %w", err)
	}
	fmt.Fprintln(w, strings.Repeat("─", 40))
	fmt.Fprintf(w, "✓ criado %s\n", cfgPath)
	return nil
}

const defaultGolangciYAML = `# Configuração mínima do golangci-lint para projetos eco.
# Veja https://golangci-lint.run/usage/configuration/ para opções avançadas.
run:
  timeout: 3m

linters:
  enable:
    - errcheck
    - govet
    - ineffassign
    - staticcheck
    - unused

issues:
  exclude-use-default: false
`
