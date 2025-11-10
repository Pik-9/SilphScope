package utils

import (
	"fmt"
	"os"
)

func WriteLines(path string, lines []string) error {
	fout, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	defer fout.Close()

	for _, ln := range lines {
		fmt.Fprintln(fout, ln)
	}

	return nil
}
