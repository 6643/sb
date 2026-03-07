package rt

import (
	"slices"
	"testing"
)

func TestHeaderReadWriteRoundTrip(t *testing.T) {
	widths := []uint8{1, 2, 1, 2, 1, 2}
	values := []uint8{1, 3, 0, 2, 1, 1}
	header := make([]byte, HeaderSize(9))

	if err := WriteHeader(header, widths, values); err != nil {
		t.Fatalf("WriteHeader failed: %v", err)
	}

	got := make([]uint8, len(values))
	if err := ReadHeader(header, widths, got, "test header"); err != nil {
		t.Fatalf("ReadHeader failed: %v", err)
	}
	if !slices.Equal(values, got) {
		t.Fatalf("header values mismatch: got %v want %v", got, values)
	}
}

func TestHeaderPaddingAndRangeValidation(t *testing.T) {
	widths := []uint8{1, 2, 1}
	header := make([]byte, HeaderSize(4))

	if err := WriteHeader(header, widths, []uint8{1, 2, 1}); err != nil {
		t.Fatalf("WriteHeader failed: %v", err)
	}
	header[0] |= 0x08
	got := make([]uint8, len(widths))
	if err := ReadHeader(header, widths, got, "bad header"); err == nil {
		t.Fatalf("ReadHeader should fail on non-zero padding")
	}

	if err := WriteHeader(make([]byte, HeaderSize(4)), widths, []uint8{2, 0, 0}); err == nil {
		t.Fatalf("WriteHeader should fail on out-of-range value")
	}
}
