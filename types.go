package main

import "strings"

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
