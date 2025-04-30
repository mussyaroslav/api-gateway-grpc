package lib

import (
	"bytes"
	"testing"
)

func TestCompressAndDecompress(t *testing.T) {
	input := []byte{'h', 'e', 'l', 'l', 'o'}

	compressed, err := Compress(input)
	if err != nil {
		t.Errorf("Compression failed with error: %v", err)
	}

	decompressed, err := Decompress(compressed)
	if err != nil {
		t.Errorf("Decompression failed with error: %v", err)
	}

	if !bytes.Equal(decompressed, input) {
		t.Errorf("Round-trip compression and decompression failed. Expected: %v, got: %v", input, decompressed)
	}
}
