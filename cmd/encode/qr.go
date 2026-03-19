package encode

import (
	"fmt"

	"github.com/spf13/cobra"

	encLib "github.com/agejevasv/swk/internal/encode"
	"github.com/agejevasv/swk/internal/ioutil"
)

var (
	qrOutput string
	qrSize   int
	qrLevel  string
)

var qrCmd = &cobra.Command{
	Use:   "qr [input]",
	Short: "Generate QR code",
	RunE: func(cmd *cobra.Command, args []string) error {
		input, err := ioutil.ReadInputString(args, cmd.InOrStdin())
		if err != nil {
			return err
		}

		switch qrOutput {
		case "terminal":
			result, err := encLib.QRTerminal(input, qrLevel)
			if err != nil {
				return err
			}
			fmt.Fprint(cmd.OutOrStdout(), result)
		case "png":
			data, err := encLib.QRGenerate(input, qrSize, qrLevel)
			if err != nil {
				return err
			}
			_, err = cmd.OutOrStdout().Write(data)
			return err
		default:
			return fmt.Errorf("invalid output format %q: must be terminal or png", qrOutput)
		}

		return nil
	},
}

func init() {
	qrCmd.Flags().StringVarP(&qrOutput, "output", "o", "terminal", "output format: terminal, png")
	qrCmd.Flags().IntVarP(&qrSize, "size", "s", 256, "image size in pixels (for PNG output)")
	qrCmd.Flags().StringVarP(&qrLevel, "level", "l", "M", "error correction level: L, M, Q, H")

	Cmd.AddCommand(qrCmd)
}
