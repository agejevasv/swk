package encode

import (
	"github.com/spf13/cobra"

	encLib "github.com/agejevasv/swk/internal/encode"
	"github.com/agejevasv/swk/internal/ioutil"
)

var gzipDecode bool
var gzipLevel int

var gzipCmd = &cobra.Command{
	Use:     "gzip [input]",
	Aliases: []string{"gz"},
	Short:   "Gzip compress or decompress",
	RunE: func(cmd *cobra.Command, args []string) error {
		input, err := ioutil.ReadInput(args, cmd.InOrStdin())
		if err != nil {
			return err
		}

		if gzipDecode {
			result, err := encLib.GzipDecompress(input)
			if err != nil {
				return err
			}
			_, err = cmd.OutOrStdout().Write(result)
			return err
		}

		result, err := encLib.GzipCompress(input, gzipLevel)
		if err != nil {
			return err
		}
		_, err = cmd.OutOrStdout().Write(result)
		return err
	},
}

func init() {
	gzipCmd.Flags().BoolVarP(&gzipDecode, "decode", "d", false, "decompress gzip input")
	gzipCmd.Flags().IntVarP(&gzipLevel, "level", "l", 6, "compression level (1-9)")

	Cmd.AddCommand(gzipCmd)
}
