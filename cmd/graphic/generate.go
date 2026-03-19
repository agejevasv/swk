package graphic

import (
	"os"

	"github.com/spf13/cobra"

	graphicLib "github.com/agejevasv/swk/internal/graphic"
)

var (
	genWidth  int
	genHeight int
	genStyle  string
	genOutput string
)

var generateCmd = &cobra.Command{
	Use:     "generate",
	Aliases: []string{"gen"},
	Short:   "Generate abstract placeholder images",
	Example: `  swk graphic generate -o art.png
  swk graphic generate --width 800 --height 600 --style circles -o out.png
  swk graphic generate --style mixed > art.png`,
	RunE: func(cmd *cobra.Command, args []string) error {
		data, err := graphicLib.GenerateImage(genWidth, genHeight, genStyle)
		if err != nil {
			return err
		}

		if genOutput != "" {
			return os.WriteFile(genOutput, data, 0644)
		}
		_, err = cmd.OutOrStdout().Write(data)
		return err
	},
}

func init() {
	generateCmd.Flags().IntVar(&genWidth, "width", 800, "image width in pixels")
	generateCmd.Flags().IntVar(&genHeight, "height", 600, "image height in pixels")
	generateCmd.Flags().StringVar(&genStyle, "style", "mixed", "art style (circles, squares, lines, mixed)")
	generateCmd.Flags().StringVarP(&genOutput, "output", "o", "", "output file (default: stdout)")

	Cmd.AddCommand(generateCmd)
}
