package query

import "github.com/spf13/cobra"

var Cmd = &cobra.Command{
	Use:     "query",
	Aliases: []string{"q"},
	Short:   "Query and search data",
}
