package utils

import (
	"os"
	"runtime"
)

var lineEnd []byte

func init() {
	if runtime.GOOS == "windows" {
		lineEnd = []byte{'\r', '\n'}
	} else {
		lineEnd = []byte{'\n'}
	}
}

func LinesToBuffer(lines []string) []byte {
	ret := make([]byte, 0)
	for _, line := range lines {
		ret = append(ret, []byte(line)...)
		ret = append(ret, lineEnd...)
	}
	return ret
}

func WriteLines(path string, lines []string) error {
	return WriteBuffer(path, LinesToBuffer(lines))
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
