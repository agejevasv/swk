package convert

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/spf13/cobra"

	graphicLib "github.com/agejevasv/swk/internal/graphic"
	"github.com/agejevasv/swk/internal/ioutil"
)

var imageCmd = &cobra.Command{
	Use:     "image [file]",
	Args:    cobra.MaximumNArgs(1),
	Aliases: []string{"img"},
	Short:   "Convert image formats",
	RunE: func(cmd *cobra.Command, args []string) error {
		input, err := ioutil.ReadFileInput(args, cmd.InOrStdin())
		if err != nil {
			return err
		}

		imageToFormat := ioutil.MustGetString(cmd, "to")
		imageQuality, err := ioutil.IntInRange(cmd, "quality", 1, 100)
		if err != nil {
			return err
		}

		imageResize := ioutil.MustGetString(cmd, "resize")
		imageOutput := ioutil.MustGetString(cmd, "output")

		var width, height int
		if imageResize != "" {
			width, height, err = parseResize(imageResize)
			if err != nil {
				return err
			}
		}

		output, err := graphicLib.ConvertImage(input, imageToFormat, imageQuality, width, height)
		if err != nil {
			return err
		}

		if imageOutput != "" {
			return os.WriteFile(imageOutput, output, 0o666)
		}
		_, err = cmd.OutOrStdout().Write(output)
		return err
	},
}

func parseResize(s string) (int, int, error) {
	parts := strings.SplitN(strings.ToLower(s), "x", 2)
	if len(parts) != 2 {
		return 0, 0, fmt.Errorf("invalid resize format %q: use WxH (e.g., 100x100)", s)
	}
	w, err := strconv.Atoi(parts[0])
	if err != nil {
		return 0, 0, fmt.Errorf("invalid width %q: %w", parts[0], err)
	}
	h, err := strconv.Atoi(parts[1])
	if err != nil {
		return 0, 0, fmt.Errorf("invalid height %q: %w", parts[1], err)
	}
	if w <= 0 || h <= 0 {
		return 0, 0, fmt.Errorf("width and height must be positive")
	}
	if w > graphicLib.MaxImageDim || h > graphicLib.MaxImageDim {
		return 0, 0, fmt.Errorf("width and height must be at most %d, got %dx%d",
			graphicLib.MaxImageDim, w, h)
	}
	return w, h, nil
}

func init() {
	imageCmd.Flags().StringP("to", "t", "", "output format: png, jpeg, gif")
	imageCmd.Flags().IntP("quality", "q", 85, "JPEG quality (1-100)")
	imageCmd.Flags().StringP("resize", "r", "", "resize to WxH (e.g., 100x100)")
	imageCmd.Flags().StringP("output", "o", "", "output file (default: stdout)")
	imageCmd.MarkFlagRequired("to")
	Cmd.AddCommand(imageCmd)
}
