package lib

import (
	"bytes"
	"compress/zlib"
	"fmt"
	"io"
)

func Compress(b []byte) ([]byte, error) {
	var buffer bytes.Buffer
	w, err := zlib.NewWriterLevel(&buffer, zlib.DefaultCompression)
	if err != nil {
		return nil, err
	}
	n, err := w.Write(b)
	if err != nil {
		return nil, err
	}
	if n != len(b) {
		return nil, fmt.Errorf("readed %d, but need %d", n, len(b))
	}
	w.Close()

	return buffer.Bytes(), nil
}

func Decompress(b []byte) ([]byte, error) {
	buffer := bytes.NewBuffer(b)
	r, err := zlib.NewReader(buffer)
	if err != nil {
		return nil, err
	}

	return io.ReadAll(r)
}
