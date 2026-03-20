package generate

import (
	"fmt"

	"github.com/spf13/cobra"

	genLib "github.com/agejevasv/swk/internal/gen"
)

var uuidCmd = &cobra.Command{
	Use:     "uuid",
	Aliases: []string{"uid"},
	Short:   "Generate UUIDs",
	Long:    "Generate UUIDs of various versions (1, 4, 5, 7).",
	RunE: func(cmd *cobra.Command, args []string) error {
		uuidVersion, _ := cmd.Flags().GetInt("version")
		uuidCount, _ := cmd.Flags().GetInt("count")
		uuidNamespace, _ := cmd.Flags().GetString("namespace")
		uuidName, _ := cmd.Flags().GetString("name")

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
	uuidCmd.Flags().IntP("version", "v", 4, "UUID version (1, 4, 5, 7)")
	uuidCmd.Flags().IntP("count", "n", 1, "number of UUIDs to generate")
	uuidCmd.Flags().String("namespace", "", "namespace for v5 (dns, url, oid, x500, or UUID)")
	uuidCmd.Flags().String("name", "", "name for v5")
	Cmd.AddCommand(uuidCmd)
}
