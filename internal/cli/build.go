package cli

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/dixavier27/eco/internal/builder"
)

func newBuildCmd() *cobra.Command {
	var (
		entry   string
		outPath string
		version string
		strip   bool
		verbose bool
	)

	cmd := &cobra.Command{
		Use:   "build",
		Short: "Compila o projeto Go atual",
		Long: "Wrapper de `go build` com injeção automática de versão (via git describe)\n" +
			"e detecção do entrypoint (cmd/api, cmd/<único>, main.go).",
		RunE: func(cmd *cobra.Command, args []string) error {
			stdout := cmd.OutOrStdout()

			res, err := builder.Build(cmd.Context(), builder.Options{
				Entry:   entry,
				Out:     outPath,
				Version: version,
				Strip:   strip,
				Verbose: verbose,
			})
			if err != nil {
				return err
			}

			fmt.Fprintf(stdout, "✓ %s (%s/%s, %s)\n", res.Out, res.GOOS, res.GOARCH, humanSize(res.Size))
			return nil
		},
	}

	cmd.Flags().StringVar(&entry, "entry", "", "pacote a compilar (default: detectado)")
	cmd.Flags().StringVar(&outPath, "out", "", "caminho do binário (default: bin/<nome>[.exe])")
	cmd.Flags().StringVar(&version, "version", "", "versão injetada em main.version (default: git describe)")
	cmd.Flags().BoolVar(&strip, "strip", false, "aplica -ldflags '-s -w' (binário menor, sem símbolos)")
	cmd.Flags().BoolVar(&verbose, "verbose", false, "ecoa o comando `go build` executado")

	return cmd
}

func humanSize(n int64) string {
	const (
		KB = 1024
		MB = KB * 1024
	)
	switch {
	case n >= MB:
		return fmt.Sprintf("%.1f MB", float64(n)/float64(MB))
	case n >= KB:
		return fmt.Sprintf("%.1f KB", float64(n)/float64(KB))
	default:
		return fmt.Sprintf("%d B", n)
	}
}
