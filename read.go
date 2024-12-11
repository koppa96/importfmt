package main

import (
	"fmt"
	"github.com/dave/dst"
	"github.com/dave/dst/decorator"
	"os"
)

func readFile(path ...string) (*dst.File, error) {
	reader := os.Stdin
	if len(path) > 0 {
		file, err := os.Open(path[0])
		if err != nil {
			return nil, fmt.Errorf("failed to open file '%s': %w", path[0], err)
		}
		defer file.Close()

		reader = file
	}

	return decorator.Parse(reader)
}
