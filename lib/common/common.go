package common

import (
	"bytes"
	"io"
	"math/rand"
)

const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890"

func GenerateNonce(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

func CompareBytesWithStream(ref []byte, input io.Reader) (bool, error) {
	buf := make([]byte, len(ref))
	_, err := io.ReadFull(input, buf)
	if err != nil {
		return false, err
	}

	return bytes.Equal(ref, buf), nil
}
