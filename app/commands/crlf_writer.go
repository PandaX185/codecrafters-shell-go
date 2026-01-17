package commands

import "io"

type crlfWriter struct {
	w io.Writer
}

func (c *crlfWriter) Write(p []byte) (int, error) {
	out := make([]byte, 0, len(p)*2)
	for _, b := range p {
		if b == '\n' {
			out = append(out, '\r', '\n')
		} else {
			out = append(out, b)
		}
	}
	return c.w.Write(out)
}
