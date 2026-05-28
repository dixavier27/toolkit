package cli

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	"github.com/dixavier27/eco/internal/builder"
)

func newReleaseCmd() *cobra.Command {
	var (
		outDir    string
		name      string
		version   string
		strip     bool
		verbose   bool
		targets   []string
		keepGoing bool
	)

	cmd := &cobra.Command{
		Use:   "release",
		Short: "Cross-compila o projeto para múltiplos SO/arquiteturas",
		Long: "Gera binários para a matriz padrão (linux/amd64, linux/arm64,\n" +
			"darwin/amd64, darwin/arm64, windows/amd64) e um checksums.txt.",
		RunE: func(cmd *cobra.Command, args []string) error {
			stdout := cmd.OutOrStdout()

			ts, err := parseTargets(targets)
			if err != nil {
				return err
			}

			res, err := builder.Release(cmd.Context(), builder.ReleaseOptions{
				OutDir:    outDir,
				Name:      name,
				Version:   version,
				Strip:     strip,
				Verbose:   verbose,
				Targets:   ts,
				KeepGoing: keepGoing,
			})
			if err != nil {
				return err
			}

			fmt.Fprintf(stdout, "\n%d artefato(s) gerados:\n", len(res.Artifacts))
			for _, a := range res.Artifacts {
				fmt.Fprintf(stdout, "  ✓ %s/%s  %s  (%s)\n", a.GOOS, a.GOARCH, a.Out, humanSize(a.Size))
			}
			if res.ChecksumsFile != "" {
				fmt.Fprintf(stdout, "\n  checksums: %s\n", res.ChecksumsFile)
			}
			if len(res.Failed) > 0 {
				fmt.Fprintf(stdout, "\n%d falha(s):\n", len(res.Failed))
				for _, f := range res.Failed {
					fmt.Fprintf(stdout, "  ✗ %s/%s: %v\n", f.Target.GOOS, f.Target.GOARCH, f.Err)
				}
				return fmt.Errorf("release com falhas parciais")
			}
			return nil
		},
	}

	cmd.Flags().StringVar(&outDir, "out", "release", "diretório de saída")
	cmd.Flags().StringVar(&name, "name", "", "nome base do binário (default: nome do projeto)")
	cmd.Flags().StringVar(&version, "version", "", "versão injetada (default: git describe)")
	cmd.Flags().BoolVar(&strip, "strip", true, "binário stripped (-s -w)")
	cmd.Flags().BoolVar(&verbose, "verbose", false, "ecoa cada comando go build")
	cmd.Flags().StringSliceVar(&targets, "targets", nil, "lista os/arch (ex.: linux/amd64,darwin/arm64); default: matriz completa")
	cmd.Flags().BoolVar(&keepGoing, "keep-going", false, "continua quando um target falha")

	return cmd
}

// parseTargets converte ["linux/amd64", "darwin/arm64"] em []Target.
// Lista vazia devolve nil (builder usa DefaultTargets).
func parseTargets(specs []string) ([]builder.Target, error) {
	if len(specs) == 0 {
		return nil, nil
	}
	out := make([]builder.Target, 0, len(specs))
	for _, s := range specs {
		parts := strings.SplitN(strings.TrimSpace(s), "/", 2)
		if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
			return nil, fmt.Errorf("target inválido %q (formato: os/arch)", s)
		}
		out = append(out, builder.Target{GOOS: parts[0], GOARCH: parts[1]})
	}
	return out, nil
}
