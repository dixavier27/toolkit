package main

import (
	"fmt"
	"os"

	"{{module}}/internal/cli"

	"github.com/spf13/cobra"
)

var version = "dev"

func main() {
	root := &cobra.Command{
		Use:     "{{name}}",
		Short:   "{{name}} — aplicação CLI criada com eco",
		Version: version,
	}
	root.AddCommand(cli.HelloCmd())
	if err := root.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
