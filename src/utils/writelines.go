package utils

import (
	"os"
)

func WriteLines(path string, lines []string) error {
	fout, err := os.Create(path)
	if err != nil {
		return err
	}
	defer fout.Close()

	for _, line := range lines {
		_, err = fout.WriteString(line + "\n")
		if err != nil {
			return err
		}
	}

	return nil
}

func WriteBuffer(path string, content []byte) error {
	fout, err := os.Create(path)
	if err != nil {
		return err
	}
	defer fout.Close()

	_, err = fout.Write(content)

	return err
}
