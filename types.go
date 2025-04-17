package main

import (
	"fmt"
	"strings"
)

type HeaderList []string

func (h *HeaderList) String() string {
	out := make([]string, len(*h))
	for _, i := range *h {
		out = append(out, i)
	}

	return strings.Join(out, ",")
}

func (h *HeaderList) Set(s string) error {
	*h = append(*h, s)
	return nil
}

type Bytes struct {
	Size float64
}

func (b Bytes) String() string {
	var rt float64
	var suffix string
	const (
		Byte  = 1
		KByte = Byte * 1024
		MByte = KByte * 1024
		GByte = MByte * 1024
	)

	if b.Size > GByte {
		rt = b.Size / GByte
		suffix = "GB"
	} else if b.Size > MByte {
		rt = b.Size / MByte
		suffix = "MB"
	} else if b.Size > KByte {
		rt = b.Size / KByte
		suffix = "KB"
	} else {
		rt = b.Size
		suffix = "bytes"
	}

	srt := fmt.Sprintf("%.2f%v", rt, suffix)

	return srt
}
