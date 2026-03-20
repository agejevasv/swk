package encode

import (
	"fmt"

	"github.com/spf13/cobra"

	encLib "github.com/agejevasv/swk/internal/encode"
	"github.com/agejevasv/swk/internal/ioutil"
)

var qrCmd = &cobra.Command{
	Use:   "qr [input]",
	Short: "Generate QR code",
	RunE: func(cmd *cobra.Command, args []string) error {
		input, err := ioutil.ReadInputString(args, cmd.InOrStdin())
		if err != nil {
			return err
		}

		qrOutput, _ := cmd.Flags().GetString("output")
		qrSize, _ := cmd.Flags().GetInt("size")
		qrLevel, _ := cmd.Flags().GetString("level")

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
	qrCmd.Flags().StringP("output", "o", "terminal", "output format: terminal, png")
	qrCmd.Flags().IntP("size", "s", 256, "image size in pixels (for PNG output)")
	qrCmd.Flags().StringP("level", "l", "M", "error correction level: L, M, Q, H")
	Cmd.AddCommand(qrCmd)
}
