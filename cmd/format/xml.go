package format

import (
	"github.com/spf13/cobra"

	fmtLib "github.com/agejevasv/swk/internal/fmt"
	"github.com/agejevasv/swk/internal/ioutil"
)

var xmlCmd = &cobra.Command{
	Use:   "xml [input]",
	Short: "Prettify or minify XML",
	Example: `  # Prettify XML
  echo '<root><a>1</a></root>' | swk format xml

  # Minify XML
  swk format xml --minify file.xml`,
	RunE: func(cmd *cobra.Command, args []string) error {
		input, err := ioutil.ReadFileInputString(args, cmd.InOrStdin())
		if err != nil {
			return err
		}

		opts := fmtLib.XMLOptions{
			Indent: ioutil.MustGetInt(cmd, "indent"),
			Minify: ioutil.MustGetBool(cmd, "minify"),
		}
		result, err := fmtLib.FormatXML([]byte(input), opts)
		if err != nil {
			return err
		}

		if len(result) > 0 && result[len(result)-1] != '\n' {
			result = append(result, '\n')
		}
		_, err = cmd.OutOrStdout().Write(result)
		return err
	},
}

func init() {
	xmlCmd.Flags().BoolP("minify", "m", false, "minify XML")
	xmlCmd.Flags().IntP("indent", "i", 2, "indentation spaces")
	Cmd.AddCommand(xmlCmd)
}
