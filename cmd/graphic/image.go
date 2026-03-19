package graphic

import (
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"

	"github.com/spf13/cobra"

	graphicLib "github.com/agejevasv/swk/internal/graphic"
)

var (
	imageToFormat string
	imageQuality  int
	imageResize   string
	imageInput    string
	imageOutput   string
)

var imageCmd = &cobra.Command{
	Use:     "image",
	Aliases: []string{"img"},
	Short:   "Convert image formats",
	RunE: func(cmd *cobra.Command, args []string) error {
		var input []byte
		var err error
		if imageInput != "" {
			input, err = os.ReadFile(imageInput)
			if err != nil {
				return fmt.Errorf("reading input file: %w", err)
			}
		} else {
			input, err = io.ReadAll(cmd.InOrStdin())
			if err != nil {
				return fmt.Errorf("reading stdin: %w", err)
			}
			if len(input) == 0 {
				return fmt.Errorf("no input provided")
			}
		}

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
			return os.WriteFile(imageOutput, output, 0644)
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
	return w, h, nil
}

func init() {
	imageCmd.Flags().StringVarP(&imageToFormat, "to", "t", "", "output format: png, jpeg, gif")
	imageCmd.Flags().IntVarP(&imageQuality, "quality", "q", 85, "JPEG quality (1-100)")
	imageCmd.Flags().StringVarP(&imageResize, "resize", "r", "", "resize to WxH (e.g., 100x100)")
	imageCmd.Flags().StringVarP(&imageInput, "input", "i", "", "input file (default: stdin)")
	imageCmd.Flags().StringVarP(&imageOutput, "output", "o", "", "output file (default: stdout)")
	imageCmd.MarkFlagRequired("to")

	Cmd.AddCommand(imageCmd)
}
