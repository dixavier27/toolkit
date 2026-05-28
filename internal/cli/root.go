package cli

import "github.com/spf13/cobra"

// Execute is the entrypoint called by cmd/eco/main.go.
// The version string comes from the binary's build-time ldflags.
func Execute(version string) error {
	root := &cobra.Command{
		Use:           "eco",
		Short:         "Gerenciador de ambientes de desenvolvimento de APIs REST em Go",
		Long:          "eco — CLI enxuta para criar e gerenciar projetos de API REST em Go.",
		SilenceUsage:  true,
		SilenceErrors: true,
	}

	root.AddCommand(newVersionCmd(version))
	root.AddCommand(newDoctorCmd())
	root.AddCommand(newNewCmd())
	root.AddCommand(newBuildCmd())

	return root.Execute()
}
