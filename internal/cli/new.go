package cli

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/dixavier27/eco/internal/gocheck"
	"github.com/dixavier27/eco/internal/scaffold"
)

func newNewCmd() *cobra.Command {
	var (
		module string
		force  bool
	)

	cmd := &cobra.Command{
		Use:   "new <nome>",
		Short: "Cria um novo projeto de API REST em Go",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			out := cmd.OutOrStdout()
			name := args[0]

			// Pré-check do Go.
			res := gocheck.Check()
			if !res.Found {
				fmt.Fprintln(out, gocheck.InstallInstructions())
				return fmt.Errorf("go não encontrado — rode `eco doctor` para detalhes")
			}
			if !res.MeetsMin {
				return fmt.Errorf("versão do Go (%s) abaixo do mínimo %d.%d",
					res.Version, gocheck.MinMajor, gocheck.MinMinor)
			}

			path, err := scaffold.Generate(scaffold.Options{
				Name:   name,
				Module: module,
				Force:  force,
			})
			if err != nil {
				return err
			}

			fmt.Fprintf(out, "\n✓ Projeto criado em %s\n\n", path)
			fmt.Fprintln(out, "Próximos passos:")
			fmt.Fprintf(out, "  cd %s\n", name)
			fmt.Fprintln(out, "  go run ./cmd/api")
			fmt.Fprintln(out, "  # em outro terminal:")
			fmt.Fprintln(out, "  curl http://localhost:8080/healthz")
			return nil
		},
	}

	cmd.Flags().StringVar(&module, "module", "", "go module path (default: nome do projeto)")
	cmd.Flags().BoolVar(&force, "force", false, "sobrescrever diretório não vazio")

	return cmd
}
