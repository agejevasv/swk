package generate

import (
	"fmt"

	"github.com/spf13/cobra"

	genLib "github.com/agejevasv/swk/internal/gen"
	"github.com/agejevasv/swk/internal/ioutil"
)

var uuidCmd = &cobra.Command{
	Use:   "uuid",
	Args:  cobra.NoArgs,
	Short: "Generate UUIDs",
	Long:  "Generate UUIDs of various versions (1, 4, 5, 7).",
	RunE: func(cmd *cobra.Command, args []string) error {
		uuidVersion := ioutil.MustGetInt(cmd, "version")

		uuidCount, err := ioutil.IntInRange(cmd, "count", 1, 10000)
		if err != nil {
			return err
		}
		uuidNamespace := ioutil.MustGetString(cmd, "namespace")
		uuidName := ioutil.MustGetString(cmd, "name")

		// Without this the request for a deterministic UUID is dropped and a
		// random one is printed instead.
		if uuidVersion != 5 && (uuidNamespace != "" || uuidName != "") {
			return fmt.Errorf("--namespace and --name only apply to --version 5")
		}

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
