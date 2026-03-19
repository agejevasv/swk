package generate

import (
	"fmt"

	"github.com/spf13/cobra"

	genLib "github.com/agejevasv/swk/internal/gen"
)

var (
	uuidVersion   int
	uuidCount     int
	uuidNamespace string
	uuidName      string
)

var uuidCmd = &cobra.Command{
	Use:     "uuid",
	Aliases: []string{"uid"},
	Short:   "Generate UUIDs",
	Long:    "Generate UUIDs of various versions (1, 4, 5, 7).",
	RunE: func(cmd *cobra.Command, args []string) error {
		for i := 0; i < uuidCount; i++ {
			id, err := genLib.GenerateUUID(uuidVersion, uuidNamespace, uuidName)
			if err != nil {
				return err
			}
			fmt.Fprintln(cmd.OutOrStdout(), id)
		}
		return nil
	},
}

func init() {
	uuidCmd.Flags().IntVarP(&uuidVersion, "version", "v", 4, "UUID version (1, 4, 5, 7)")
	uuidCmd.Flags().IntVarP(&uuidCount, "count", "n", 1, "number of UUIDs to generate")
	uuidCmd.Flags().StringVar(&uuidNamespace, "namespace", "", "namespace for v5 (dns, url, oid, x500, or UUID)")
	uuidCmd.Flags().StringVar(&uuidName, "name", "", "name for v5")
	Cmd.AddCommand(uuidCmd)
}
