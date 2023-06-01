package util

import "io"

type zero struct{}

func (z zero) Read(d []byte) (int, error) {
	return len(d), nil
}

var Zero io.Reader = zero{}
