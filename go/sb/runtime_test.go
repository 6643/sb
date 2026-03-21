package sb

import (
	"bytes"
	"slices"
	"testing"
)

type gameVectorForTest struct {
	ID   uint32
	Name string
}

func isZeroGameVectorForTest(v gameVectorForTest) bool {
	return v.ID == 0 && v.Name == ""
}

func encodeGameWireForTest(id uint32, name string) ([]byte, error) {
	widths := []uint8{1, 2}
	nameState, err := TextState(len(name))
	if err != nil {
		return nil, err
	}
	header := make([]byte, HeaderSize(3))
	if err := WriteHeader(header, widths, []uint8{boolToState1ForTest(id != 0), nameState}); err != nil {
		return nil, err
	}
	buf := bytes.NewBuffer(header)
	if id != 0 {
		if err := SetU32(buf, id); err != nil {
			return nil, err
		}
	}
	if err := SetText(buf, nameState, name); err != nil {
		return nil, err
	}
	return bytes.Clone(buf.Bytes()), nil
}

func decodeGameWireForTest(buf *bytes.Buffer) (uint32, string, error) {
	header := buf.Next(HeaderSize(3))
	if len(header) != HeaderSize(3) {
		return 0, "", bytes.ErrTooLarge
	}
	states := make([]uint8, 2)
	if err := ReadHeader(header, []uint8{1, 2}, states, "Game header"); err != nil {
		return 0, "", err
	}
	var id uint32
	if states[0] == 1 {
		value, err := GetU32(buf)
		if err != nil {
			return 0, "", err
		}
		id = value
	}
	name, err := GetText(buf, states[1])
	if err != nil {
		return 0, "", err
	}
	return id, name, nil
}

func boolToState1ForTest(v bool) uint8 {
	if v {
		return 1
	}
	return 0
}

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

func TestProtocolGoldenVectors(t *testing.T) {
	t.Run("game vectors", func(t *testing.T) {
		cases := []struct {
			name string
			id   uint32
			text string
			want []byte
		}{
			{name: "id zero name lol", id: 0, text: "lol", want: []byte{0x20, 0x03, 0x6C, 0x6F, 0x6C}},
			{name: "id seven name empty", id: 7, text: "", want: []byte{0x80, 0x07, 0x00, 0x00, 0x00}},
			{name: "id seven name lol", id: 7, text: "lol", want: []byte{0xA0, 0x07, 0x00, 0x00, 0x00, 0x03, 0x6C, 0x6F, 0x6C}},
		}

		for _, tc := range cases {
			t.Run(tc.name, func(t *testing.T) {
				wire, err := encodeGameWireForTest(tc.id, tc.text)
				if err != nil {
					t.Fatalf("encodeGameWireForTest failed: %v", err)
				}
				if !bytes.Equal(wire, tc.want) {
					t.Fatalf("wire = % X want % X", wire, tc.want)
				}
				id, text, err := decodeGameWireForTest(bytes.NewBuffer(bytes.Clone(wire)))
				if err != nil {
					t.Fatalf("decodeGameWireForTest failed: %v", err)
				}
				if id != tc.id || text != tc.text {
					t.Fatalf("decoded = (%d,%q) want (%d,%q)", id, text, tc.id, tc.text)
				}
			})
		}
	})

	t.Run("demo cross-byte header", func(t *testing.T) {
		widths := []uint8{1, 2, 2, 2, 2}
		values := []uint8{1, 1, 2, 0, 1}
		header := make([]byte, HeaderSize(9))

		if err := WriteHeader(header, widths, values); err != nil {
			t.Fatalf("WriteHeader failed: %v", err)
		}
		if got, want := header, []byte{0xB0, 0x80}; !bytes.Equal(got, want) {
			t.Fatalf("header = % X want % X", got, want)
		}

		decoded := make([]uint8, len(values))
		if err := ReadHeader(header, widths, decoded, "Demo header"); err != nil {
			t.Fatalf("ReadHeader failed: %v", err)
		}
		if !slices.Equal(decoded, values) {
			t.Fatalf("decoded = %v want %v", decoded, values)
		}
	})

	t.Run("text list item header block", func(t *testing.T) {
		header, err := writeStateBlock([]uint8{StateU8, StateZero, StateU16, StateU8, StateZero})
		if err != nil {
			t.Fatalf("writeStateBlock failed: %v", err)
		}
		if got, want := header, []byte{0x49, 0x00}; !bytes.Equal(got, want) {
			t.Fatalf("header = % X want % X", got, want)
		}

		buf := bytes.NewBuffer(bytes.Clone(header))
		states, err := readStateBlock(buf, 5, "text list state block")
		if err != nil {
			t.Fatalf("readStateBlock failed: %v", err)
		}
		if got, want := states, []uint8{StateU8, StateZero, StateU16, StateU8, StateZero}; !slices.Equal(got, want) {
			t.Fatalf("states = %v want %v", got, want)
		}
	})

	t.Run("u32 bitmap list body", func(t *testing.T) {
		bitmap, count := writeBitmapFromDefaults([]uint32{1, 2, 0, 0, 3}, func(v uint32) bool { return v == 0 })
		if got, want := count, 3; got != want {
			t.Fatalf("non-default count = %d want %d", got, want)
		}
		if got, want := bitmap, []byte{0xC8}; !bytes.Equal(got, want) {
			t.Fatalf("bitmap = % X want % X", got, want)
		}
	})

	t.Run("i8 bitmap list body", func(t *testing.T) {
		bitmap, count := writeBitmapFromDefaults([]int8{1, 2, 0, 0, 3}, func(v int8) bool { return v == 0 })
		if got, want := count, 3; got != want {
			t.Fatalf("non-default count = %d want %d", got, want)
		}
		if got, want := bitmap, []byte{0xC8}; !bytes.Equal(got, want) {
			t.Fatalf("bitmap = % X want % X", got, want)
		}
	})

	t.Run("struct list bitmap semantics", func(t *testing.T) {
		games := []gameVectorForTest{{}, {ID: 7, Name: "lol"}, {}, {ID: 9, Name: "go"}}
		bitmap, count := writeBitmapFromDefaults(games, isZeroGameVectorForTest)
		if got, want := count, 2; got != want {
			t.Fatalf("non-default count = %d want %d", got, want)
		}
		if got, want := bitmap, []byte{0x50}; !bytes.Equal(got, want) {
			t.Fatalf("bitmap = % X want % X", got, want)
		}
	})
}

func TestReadHeaderRejectsOversizedBuffer(t *testing.T) {
	widths := []uint8{1, 2, 1}
	out := make([]uint8, len(widths))
	data := make([]byte, HeaderSize(4)+1)

	if err := ReadHeader(data, widths, out, "test header"); err == nil {
		t.Fatalf("ReadHeader should fail on oversized buffer")
	}
}

func TestWriteHeaderRejectsOversizedBuffer(t *testing.T) {
	widths := []uint8{1, 2, 1}
	values := []uint8{1, 2, 1}
	dst := make([]byte, HeaderSize(4)+1)

	if err := WriteHeader(dst, widths, values); err == nil {
		t.Fatalf("WriteHeader should fail on oversized buffer")
	}
}

func TestHeaderRejectsInvalidWidthThree(t *testing.T) {
	widths := []uint8{1, 3}
	data := make([]byte, HeaderSize(4))
	out := make([]uint8, len(widths))

	if err := ReadHeader(data, widths, out, "test header"); err == nil {
		t.Fatalf("ReadHeader should fail on invalid width")
	}
	if err := WriteHeader(data, widths, []uint8{1, 0}); err == nil {
		t.Fatalf("WriteHeader should fail on invalid width")
	}
}

func TestValidatePaddingZeroRejectsInvalidUsedBits(t *testing.T) {
	t.Run("negative", func(t *testing.T) {
		if err := ValidatePaddingZero([]byte{0}, -1, "padding"); err == nil {
			t.Fatalf("ValidatePaddingZero should fail on negative usedBits")
		}
	})

	t.Run("beyond buffer", func(t *testing.T) {
		if err := ValidatePaddingZero([]byte{0}, 9, "padding"); err == nil {
			t.Fatalf("ValidatePaddingZero should fail when usedBits exceed buffer")
		}
	})
}

func TestSetTextRejectsInvalidUTF8(t *testing.T) {
	buf := &bytes.Buffer{}
	value := string([]byte{0xff})

	err := SetText(buf, StateU8, value)
	if err == nil {
		t.Fatalf("SetText should reject invalid UTF-8")
	}
}

func TestGetTextRejectsInvalidUTF8(t *testing.T) {
	buf := bytes.NewBuffer([]byte{0x01, 0xff})

	_, err := GetText(buf, StateU8)
	if err == nil {
		t.Fatalf("GetText should reject invalid UTF-8")
	}
}

func TestTextRoundTripPreservesValidUTF8(t *testing.T) {
	buf := &bytes.Buffer{}
	value := "你好, runtime"

	if err := SetText(buf, StateU8, value); err != nil {
		t.Fatalf("SetText failed: %v", err)
	}

	got, err := GetText(bytes.NewBuffer(buf.Bytes()), StateU8)
	if err != nil {
		t.Fatalf("GetText failed: %v", err)
	}
	if got != value {
		t.Fatalf("text mismatch: got %q want %q", got, value)
	}
}

func TestSizeTextRejectsInvalidUTF8(t *testing.T) {
	value := string([]byte{0xff})

	_, err := SizeText(value)
	if err == nil {
		t.Fatalf("SizeText should reject invalid UTF-8")
	}
}

func TestGetListCountReportsZeroCount(t *testing.T) {
	buf := bytes.NewBuffer([]byte{0x00})

	_, err := getListCount(buf, StateU8)
	if err == nil {
		t.Fatalf("getListCount should reject encoded zero count")
	}
	if got, want := err.Error(), "list count state 1 encoded zero count"; got != want {
		t.Fatalf("error = %q want %q", got, want)
	}
}

func TestRuntimeRejectsNonCanonicalStates(t *testing.T) {
	t.Run("text u16 for short payload", func(t *testing.T) {
		buf := bytes.NewBuffer([]byte{0x03, 0x00, 'l', 'o', 'l'})

		_, err := GetText(buf, StateU16)
		if err == nil {
			t.Fatalf("GetText should reject non-canonical state")
		}
		if got, want := err.Error(), "text state 2 is not canonical for length 3"; got != want {
			t.Fatalf("error = %q want %q", got, want)
		}
	})

	t.Run("bin u16 for short payload", func(t *testing.T) {
		buf := bytes.NewBuffer([]byte{0x03, 0x00, 0x09, 0x08, 0x07})

		_, err := GetBinInto(buf, StateU16, nil)
		if err == nil {
			t.Fatalf("GetBinInto should reject non-canonical state")
		}
		if got, want := err.Error(), "bin state 2 is not canonical for length 3"; got != want {
			t.Fatalf("error = %q want %q", got, want)
		}
	})

	t.Run("list illegal state", func(t *testing.T) {
		_, err := getListCount(bytes.NewBuffer(nil), StateU16)
		if err == nil {
			t.Fatalf("getListCount should reject illegal state")
		}
		if got, want := err.Error(), "list count state 2 is illegal"; got != want {
			t.Fatalf("error = %q want %q", got, want)
		}
	})
}
