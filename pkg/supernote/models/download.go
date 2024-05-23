package models

import (
	"fmt"
	"io"
)

type ProgressReader struct {
	Reader io.Reader
	Size   int64
	Pos    int64
}

func (pr *ProgressReader) Read(p []byte) (int, error) {
	n, err := pr.Reader.Read(p)
	if err == nil {
		pr.Pos += int64(n)
		fmt.Printf("\rdownloading... %.2f%%", float64(pr.Pos)/float64(pr.Size)*100)
	}

	return n, err
}
