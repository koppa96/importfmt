package main

import (
	"fmt"
	"io"
	"os"

	"github.com/spf13/cobra"
)

var cmd = &cobra.Command{
	Use:   "importfmt [filename]",
	Short: "Groups imports in go files like goimports, but also cleans up unnecessary empty lines.",
	Args:  cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		write, _ := cmd.Flags().GetBool("write")
		if write && len(args) == 0 {
			return fmt.Errorf("the argument 'filename' must be specified when the flag 'write' is set")
		}

		f, err := readFile(args...)
		if err != nil {
			return err
		}

		var writer io.Writer
		if len(args) > 0 && write {
			file, err := os.Create(args[0])
			if err != nil {
				return fmt.Errorf("failed to create file '%s': %w", args[0], err)
			}
			defer file.Close()

			writer = file
		} else {
			writer = os.Stdout
		}

		return format(f, writer)
	},
}

func main() {
	cmd.Flags().BoolP("write", "w", false, "Set to true to write the result to the file")

	if err := cmd.Execute(); err != nil {
		println(err)
	}
}
