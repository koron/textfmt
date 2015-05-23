package fmtinfo

import (
	"io"
	"os"
)

// Detect detects format of text file: encoding and end of line.
func Detect(path string) (*Info, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	return detect(f)
}

func detect(r io.Reader) (*Info, error) {
	// TODO:
	return &Info{Encode: UTF8, EOL: LF}, nil
}
