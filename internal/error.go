package internal

import "io"

func must[T any](v T, err error) T {
	if err != nil {
		panic(err)
	}
	return v
}

func MustClose(c io.Closer) {
	if err := c.Close(); err != nil {
		panic(err)
	}
}
