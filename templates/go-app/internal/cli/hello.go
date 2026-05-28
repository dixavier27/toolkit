package cli

import (
	"fmt"

	"github.com/spf13/cobra"
)

// HelloCmd retorna o subcomando "hello" que imprime uma saudação.
func HelloCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "hello [nome]",
		Short: "Imprime saudação",
		Args:  cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			nome := "mundo"
			if len(args) > 0 {
				nome = args[0]
			}
			_, err := fmt.Fprintf(cmd.OutOrStdout(), "Olá, %s!\n", nome)
			return err
		},
	}
}
