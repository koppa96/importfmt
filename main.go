package main

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"

	"github.com/dave/dst/decorator"
	"github.com/spf13/cobra"
)

var cmd = &cobra.Command{
	Use:   "importfmt [filename]",
	Short: "Groups imports in go files like goimports, but also cleans up unnecessary empty lines.",
	Args:  cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		write, _ := cmd.Flags().GetBool("write")
		if write && len(args) == 0 {
			return fmt.Errorf(
				"the argument 'filename' must be specified when the flag 'write' is set",
			)
		}

		var reader io.ReadCloser
		if len(args) > 0 {
			file, err := os.Open(args[0])
			if err != nil {
				return fmt.Errorf("failed to open file '%s': %w", args[0], err)
			}
			defer file.Close()

			reader = file
		} else {
			reader = os.Stdin
		}

		var preformattedCode io.Reader
		runGolines, _ := cmd.Flags().GetBool("invoke-golines")
		if runGolines {
			var golinesBuf bytes.Buffer
			err := invokeGolines(reader, &golinesBuf)
			if err != nil {
				return err
			}

			preformattedCode = &golinesBuf
		}

		f, err := decorator.Parse(preformattedCode)
		if err != nil {
			return err
		}

		var buf bytes.Buffer
		err = format(f, &buf)
		if err != nil {
			return fmt.Errorf("failed to format file: %w", err)
		}

		if write {
			reader.Close()
			file, err := os.Create(args[0])
			if err != nil {
				return fmt.Errorf("failed to create file '%s': %w", args[0], err)
			}
			defer file.Close()

			_, err = buf.WriteTo(file)
			if err != nil {
				return fmt.Errorf("failed to write formatted file to '%s': %w", args[0], err)
			}
		} else {
			_, err = buf.WriteTo(os.Stdout)
			if err != nil {
				return fmt.Errorf("failed to write formatted file to stdout: %w", err)
			}
		}

		return nil
	},
}

func invokeGolines(reader io.Reader, writer io.Writer) error {
	cmd := exec.Command("golines")
	cmd.Stdin = reader
	cmd.Stdout = writer

	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("failed to invoke golines on file: %w", err)
	}

	return nil
}

func main() {
	cmd.Flags().BoolP("write", "w", false, "Set to true to write the result to the file")
	cmd.Flags().
		Bool("invoke-golines", false, "If set to true, this tool invokes golines to format, then formats the imports")

	if err := cmd.Execute(); err != nil {
		println(err)
	}
}
