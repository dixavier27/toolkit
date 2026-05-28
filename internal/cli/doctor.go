package cli

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/dixavier27/eco/internal/gocheck"
)

func newDoctorCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "doctor",
		Short: "Verifica se a toolchain Go está pronta para uso",
		RunE: func(cmd *cobra.Command, args []string) error {
			out := cmd.OutOrStdout()
			res := gocheck.Check()

			if !res.Found {
				fmt.Fprintln(out, gocheck.InstallInstructions())
				return fmt.Errorf("go não encontrado")
			}

			fmt.Fprintf(out, "✓ Go encontrado: %s\n", res.RawOutput)
			fmt.Fprintf(out, "  caminho: %s\n", res.Path)

			if !res.MeetsMin {
				fmt.Fprintf(out, "✗ Versão mínima exigida: %d.%d (atual: %s)\n",
					gocheck.MinMajor, gocheck.MinMinor, res.Version)
				fmt.Fprintln(out, "  Atualize em https://go.dev/dl/")
				return fmt.Errorf("versão do Go abaixo do mínimo")
			}

			fmt.Fprintln(out, "✓ Toolchain Go pronta.")
			return nil
		},
	}
}
