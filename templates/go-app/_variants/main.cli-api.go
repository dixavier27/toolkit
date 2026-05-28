package main

import (
	"fmt"
	"os"

	"{{module}}/internal/api"
	"{{module}}/internal/cli"

	"github.com/spf13/cobra"
)

var version = "dev"

func main() {
	root := &cobra.Command{
		Use:     "{{name}}",
		Short:   "{{name}} — aplicação Go criada com eco",
		Version: version,
	}
	root.AddCommand(cli.HelloCmd())
	root.AddCommand(api.ServeCmd(version))
	if err := root.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
