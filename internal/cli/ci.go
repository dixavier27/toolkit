package cli

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/dixavier27/eco/internal/cigen"
)

func newCICmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "ci",
		Short: "Gerencia workflows de CI/CD do projeto",
	}
	cmd.AddCommand(newCIGenerateCmd())
	return cmd
}

func newCIGenerateCmd() *cobra.Command {
	var (
		only      string
		force     bool
		name      string
		entry     string
		goVersion string
	)

	cmd := &cobra.Command{
		Use:   "generate",
		Short: "Gera workflows GitHub Actions (.github/workflows/{ci,release}.yml)",
		Long: "Cria ou atualiza os arquivos ci.yml e release.yml em .github/workflows/.\n" +
			"ci.yml: vet, lint (golangci-lint), test (com race), build em push/PR.\n" +
			"release.yml: cross-compile multi-OS + GitHub Release em tag v*.",
		RunE: func(cmd *cobra.Command, args []string) error {
			out := cmd.OutOrStdout()

			res, err := cigen.Generate(cigen.Options{
				Only:      only,
				Force:     force,
				Name:      name,
				Entry:     entry,
				GoVersion: goVersion,
			})
			if err != nil {
				return err
			}

			for _, w := range res.Written {
				fmt.Fprintf(out, "✓ %s\n", w)
			}
			for _, s := range res.Skipped {
				fmt.Fprintf(out, "- %s (já existe — use --force para sobrescrever)\n", s)
			}
			if len(res.Written) == 0 && len(res.Skipped) > 0 {
				return fmt.Errorf("nada gerado")
			}
			return nil
		},
	}

	cmd.Flags().StringVar(&only, "only", "", "gera só um workflow: 'ci' ou 'release' (default: ambos)")
	cmd.Flags().BoolVar(&force, "force", false, "sobrescreve arquivos existentes")
	cmd.Flags().StringVar(&name, "name", "", "nome do binário (default: nome do diretório)")
	cmd.Flags().StringVar(&entry, "entry", "./...", "pacote a compilar no release")
	cmd.Flags().StringVar(&goVersion, "go-version", "1.26", "versão do Go no workflow")

	return cmd
}
