package generate

import (
	"os"

	"github.com/spf13/cobra"

	graphicLib "github.com/agejevasv/swk/internal/graphic"
	"github.com/agejevasv/swk/internal/ioutil"
)

var imageCmd = &cobra.Command{
	Use:   "image",
	Short: "Generate abstract placeholder images",
	Example: `  swk generate image -o art.png
  swk generate image --width 800 --height 600 --style circles -o out.png
  swk generate image --style mixed > art.png`,
	RunE: func(cmd *cobra.Command, args []string) error {
		genWidth := ioutil.MustGetInt(cmd, "width")
		genHeight := ioutil.MustGetInt(cmd, "height")
		genStyle := ioutil.MustGetString(cmd, "style")
		genOutput := ioutil.MustGetString(cmd, "output")

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
	imageCmd.Flags().Int("width", 800, "image width in pixels")
	imageCmd.Flags().Int("height", 600, "image height in pixels")
	imageCmd.Flags().String("style", "mixed", "art style (circles, squares, lines, mixed)")
	imageCmd.Flags().StringP("output", "o", "", "output file (default: stdout)")
	Cmd.AddCommand(imageCmd)
}
