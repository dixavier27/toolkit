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
		tupa   bool
	)

	cmd := &cobra.Command{
		Use:   "new <nome>",
		Short: "Cria um novo projeto de API REST em Go",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			out := cmd.OutOrStdout()
			name := args[0]

			// Pré-check do Go.
			res := gocheck.CheckGo()
			if !res.OK() {
				fmt.Fprintf(out, "%s\n", res.Message)
				if res.Suggestion != "" {
					fmt.Fprintf(out, "  %s\n", res.Suggestion)
				}
				return fmt.Errorf("ambiente Go inválido — rode `eco doctor` para detalhes")
			}

			path, err := scaffold.Generate(scaffold.Options{
				Name:   name,
				Module: module,
				Force:  force,
				Tupa:   tupa,
			})
			if err != nil {
				return err
			}

			fmt.Fprintf(out, "\n✓ Projeto criado em %s\n\n", path)
			fmt.Fprintln(out, "Próximos passos:")
			fmt.Fprintf(out, "  cd %s\n", name)
			fmt.Fprintln(out, "  go run ./cmd/api")
			fmt.Fprintln(out, "  # em outro terminal:")
			if tupa {
				fmt.Fprintln(out, "  curl -X POST http://localhost:8080/tarefas -H 'Content-Type: application/json' -d '{\"id\":\"1\",\"titulo\":\"comprar pão\"}'")
				fmt.Fprintln(out, "  curl http://localhost:8080/tarefas")
			} else {
				fmt.Fprintln(out, "  curl http://localhost:8080/healthz")
			}
			return nil
		},
	}

	cmd.Flags().StringVar(&module, "module", "", "go module path (default: nome do projeto)")
	cmd.Flags().BoolVar(&force, "force", false, "sobrescrever diretório não vazio")
	cmd.Flags().BoolVar(&tupa, "tupa", false, "scaffolda usando tupa-go (vendored em internal/tupa) com Recurso[T] de exemplo")

	return cmd
}
