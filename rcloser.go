package main

import (
	"bufio"
	"io"
	"os"
)

type RCloser struct {
	io.Reader
}

func NewRCloser(file *os.File) RCloser {
	reader := bufio.NewReader(file)
	return RCloser{reader}
}

func (RCloser) Close() error {
	return nil
}
