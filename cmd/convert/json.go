package convert

import (
	"fmt"

	"github.com/spf13/cobra"

	convLib "github.com/agejevasv/swk/internal/convert"
	fmtLib "github.com/agejevasv/swk/internal/fmt"
	"github.com/agejevasv/swk/internal/ioutil"
)

var jsonCmd = &cobra.Command{
	Use:   "json [input]",
	Short:   "Convert and format JSON",
	Long: `Convert between JSON, YAML, and CSV. Also prettify or minify JSON.

When --from and --to are both json (the default), it formats the input.`,
	Example: `  # Prettify JSON
  echo '{"a":1}' | swk convert json

  # Minify JSON
  echo '{"a": 1}' | swk convert json --minify

  # JSON to YAML
  echo '{"a":1}' | swk convert json --to yaml

  # YAML to JSON
  echo 'a: 1' | swk convert json --from yaml

  # JSON to CSV
  echo '[{"name":"alice"}]' | swk convert json --to csv

  # CSV to JSON
  echo 'name,age\nalice,30' | swk convert json --from csv`,
	RunE: func(cmd *cobra.Command, args []string) error {
		input, err := ioutil.ReadFileInputString(args, cmd.InOrStdin())
		if err != nil {
			return err
		}

		jsonFrom := ioutil.MustGetString(cmd, "from")
		jsonTo := ioutil.MustGetString(cmd, "to")
		jsonMinify := ioutil.MustGetBool(cmd, "minify")
		jsonIndent := ioutil.MustGetInt(cmd, "indent")
		jsonDelimiter := ioutil.MustGetString(cmd, "delimiter")

		validFormats := map[string]bool{"json": true, "yaml": true, "csv": true}
		if !validFormats[jsonFrom] {
			return fmt.Errorf("unsupported --from format: %s", jsonFrom)
		}
		if !validFormats[jsonTo] {
			return fmt.Errorf("unsupported --to format: %s", jsonTo)
		}

		switch {
		case jsonFrom == "yaml":
			output, err := convLib.YAMLToJSON([]byte(input), jsonIndent)
			if err != nil {
				return err
			}
			fmt.Fprint(cmd.OutOrStdout(), string(output))

		case jsonFrom == "csv":
			delimiter := ','
			if len(jsonDelimiter) > 0 {
				delimiter = rune(jsonDelimiter[0])
			}
			output, err := convLib.CSVToJSON([]byte(input), delimiter)
			if err != nil {
				return err
			}
			fmt.Fprint(cmd.OutOrStdout(), string(output))

		case jsonTo == "yaml":
			output, err := convLib.JSONToYAML([]byte(input))
			if err != nil {
				return err
			}
			fmt.Fprint(cmd.OutOrStdout(), string(output))

		case jsonTo == "csv":
			delimiter := ','
			if len(jsonDelimiter) > 0 {
				delimiter = rune(jsonDelimiter[0])
			}
			output, err := convLib.JSONToCSV([]byte(input), delimiter)
			if err != nil {
				return err
			}
			fmt.Fprint(cmd.OutOrStdout(), string(output))

		default:
			// Format JSON (prettify or minify)
			opts := fmtLib.JSONOptions{
				Indent: jsonIndent,
				Minify: jsonMinify,
			}
			result, err := fmtLib.FormatJSON([]byte(input), opts)
			if err != nil {
				return err
			}
			fmt.Fprint(cmd.OutOrStdout(), string(result))
		}

		return nil
	},
}

func init() {
	jsonCmd.Flags().String("from", "json", "input format (json, yaml, csv)")
	jsonCmd.Flags().String("to", "json", "output format (json, yaml, csv)")
	jsonCmd.Flags().BoolP("minify", "m", false, "minify JSON output")
	jsonCmd.Flags().IntP("indent", "i", 2, "indentation spaces")
	jsonCmd.Flags().StringP("delimiter", "d", ",", "CSV delimiter character")
	Cmd.AddCommand(jsonCmd)
}
