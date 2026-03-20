package inspect

import (
	"encoding/json"
	"fmt"
	"text/tabwriter"

	"github.com/spf13/cobra"

	"github.com/agejevasv/swk/internal/ioutil"
	textLib "github.com/agejevasv/swk/internal/text"
)

var textCmd = &cobra.Command{
	Use:     "text [input]",
	Aliases: []string{"txt"},
	Short:   "Analyze text and show statistics",
	RunE: func(cmd *cobra.Command, args []string) error {
		input, err := ioutil.ReadFileInputString(args, cmd.InOrStdin())
		if err != nil {
			return err
		}

		info := textLib.Inspect(input)

		inspectJSON, _ := cmd.Flags().GetBool("json")

		if inspectJSON {
			enc := json.NewEncoder(cmd.OutOrStdout())
			enc.SetIndent("", "  ")
			return enc.Encode(info)
		}

		w := tabwriter.NewWriter(cmd.OutOrStdout(), 0, 0, 2, ' ', 0)
		fmt.Fprintf(w, "Characters:\t%d\n", info.Characters)
		fmt.Fprintf(w, "Words:\t%d\n", info.Words)
		fmt.Fprintf(w, "Lines:\t%d\n", info.Lines)
		fmt.Fprintf(w, "Sentences:\t%d\n", info.Sentences)
		fmt.Fprintf(w, "Bytes:\t%d (%s)\n", info.Bytes, info.BytesHuman)
		fmt.Fprintf(w, "Is ASCII:\t%v\n", info.IsASCII)
		fmt.Fprintf(w, "Has Unicode:\t%v\n", info.HasUnicode)
		w.Flush()

		return nil
	},
}

func init() {
	textCmd.Flags().Bool("json", false, "Output as JSON")
	Cmd.AddCommand(textCmd)
}
