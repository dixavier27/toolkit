package api

import (
	"github.com/spf13/cobra"
)

// ServeCmd retorna o subcomando "serve" que sobe o servidor HTTP.
func ServeCmd(version string) *cobra.Command {
	return &cobra.Command{
		Use:   "serve",
		Short: "Sobe o servidor HTTP",
		RunE: func(_ *cobra.Command, _ []string) error {
			return Run(version)
		},
	}
}
