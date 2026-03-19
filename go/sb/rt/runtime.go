package rt

import (
	"bytes"
	"fmt"
	"math"
	"unicode/utf8"
)

const (
	StateZero uint8 = 0
	StateU8   uint8 = 1
	StateU16  uint8 = 2
	StateU24  uint8 = 3
)

func HeaderSize(bitCount int) int {
	if bitCount <= 0 {
		return 0
	}
	return (bitCount + 7) / 8
}

func headerBitCount(widths []uint8) (int, error) {
	bitCount := 0
	for i, width := range widths {
		if width != 1 && width != 2 {
			return 0, fmt.Errorf("header field[%d] invalid width: %d", i, width)
		}
		bitCount += int(width)
	}
	return bitCount, nil
}

func ReadHeader(data []byte, widths []uint8, out []uint8, kind string) error {
	if len(widths) != len(out) {
		return fmt.Errorf("%s widths/out length mismatch: %d != %d", kind, len(widths), len(out))
	}
	bitCount, err := headerBitCount(widths)
	if err != nil {
		return fmt.Errorf("%s: %w", kind, err)
	}
	headerSize := HeaderSize(bitCount)
	if len(data) != headerSize {
		return fmt.Errorf("%s invalid header size: %d != %d", kind, len(data), headerSize)
	}
	if err := ValidatePaddingZero(data, bitCount, kind); err != nil {
		return err
	}
	reader := NewBitReader(data, bitCount)
	for i, width := range widths {
		value, err := reader.ReadBits(int(width))
		if err != nil {
			return fmt.Errorf("%s field[%d]: %w", kind, i, err)
		}
		out[i] = value
	}
	return nil
}

func WriteHeader(dst []byte, widths []uint8, values []uint8) error {
	if len(widths) != len(values) {
		return fmt.Errorf("header widths/values length mismatch: %d != %d", len(widths), len(values))
	}
	bitCount, err := headerBitCount(widths)
	if err != nil {
		return err
	}
	headerSize := HeaderSize(bitCount)
	if len(dst) != headerSize {
		return fmt.Errorf("header dst invalid size: %d != %d", len(dst), headerSize)
	}
	for i := range dst {
		dst[i] = 0
	}
	bitOffset := 0
	for i, width := range widths {
		value := values[i]
		maxValue := uint8(0xFF)
		if width < 8 {
			maxValue = (1 << width) - 1
		}
		if value > maxValue {
			return fmt.Errorf("header field[%d] value %d exceeds %d-bit max %d", i, value, width, maxValue)
		}
		for shift := int(width) - 1; shift >= 0; shift-- {
			if (value>>shift)&1 == 1 {
				byteIndex := bitOffset / 8
				bitIndex := 7 - (bitOffset % 8)
				dst[byteIndex] |= 1 << bitIndex
			}
			bitOffset++
		}
	}
	return nil
}

func bitmapSize(count int) int {
	if count <= 0 {
		return 0
	}
	return (count + 7) / 8
}

func itemHeaderSize(count int) int {
	if count <= 0 {
		return 0
	}
	return (count*2 + 7) / 8
}

type BitReader struct {
	data      []byte
	bitLimit  int
	bitOffset int
}

func NewBitReader(data []byte, bitLimit int) *BitReader {
	maxBits := len(data) * 8
	if bitLimit < 0 || bitLimit > maxBits {
		bitLimit = maxBits
	}
	return &BitReader{data: data, bitLimit: bitLimit}
}

func (r *BitReader) ReadBit() (bool, error) {
	if r.bitOffset >= r.bitLimit {
		return false, fmt.Errorf("not enough bits")
	}
	byteIndex := r.bitOffset / 8
	bitIndex := 7 - (r.bitOffset % 8)
	r.bitOffset++
	return (r.data[byteIndex] & (1 << bitIndex)) != 0, nil
}

func (r *BitReader) ReadBits(width int) (uint8, error) {
	if width <= 0 || width > 8 {
		return 0, fmt.Errorf("invalid bit width: %d", width)
	}
	var value uint8
	for i := 0; i < width; i++ {
		bit, err := r.ReadBit()
		if err != nil {
			return 0, err
		}
		value <<= 1
		if bit {
			value |= 1
		}
	}
	return value, nil
}

func ValidatePaddingZero(data []byte, usedBits int, kind string) error {
	maxBits := len(data) * 8
	if usedBits < 0 || usedBits > maxBits {
		return fmt.Errorf("%s invalid used bits: %d not in [0,%d]", kind, usedBits, maxBits)
	}
	for bitOffset := usedBits; bitOffset < len(data)*8; bitOffset++ {
		byteIndex := bitOffset / 8
		bitIndex := 7 - (bitOffset % 8)
		if (data[byteIndex] & (1 << bitIndex)) != 0 {
			return fmt.Errorf("%s padding bit %d is not zero", kind, bitOffset-usedBits)
		}
	}
	return nil
}

func GetU8(buf *bytes.Buffer) (uint8, error) {
	if buf.Len() < 1 {
		return 0, fmt.Errorf("not enough data")
	}
	return buf.Next(1)[0], nil
}

func SetU8(buf *bytes.Buffer, value uint8) error {
	return buf.WriteByte(value)
}

func GetI8(buf *bytes.Buffer) (int8, error) {
	value, err := GetU8(buf)
	return int8(value), err
}

func SetI8(buf *bytes.Buffer, value int8) error {
	return SetU8(buf, uint8(value))
}

func GetU16(buf *bytes.Buffer) (uint16, error) {
	if buf.Len() < 2 {
		return 0, fmt.Errorf("not enough data")
	}
	data := buf.Next(2)
	return uint16(data[0]) | uint16(data[1])<<8, nil
}

func SetU16(buf *bytes.Buffer, value uint16) error {
	_, err := buf.Write([]byte{byte(value), byte(value >> 8)})
	return err
}

func GetI16(buf *bytes.Buffer) (int16, error) {
	value, err := GetU16(buf)
	return int16(value), err
}

func SetI16(buf *bytes.Buffer, value int16) error {
	return SetU16(buf, uint16(value))
}

func GetU24(buf *bytes.Buffer) (uint32, error) {
	if buf.Len() < 3 {
		return 0, fmt.Errorf("not enough data")
	}
	data := buf.Next(3)
	return uint32(data[0]) | uint32(data[1])<<8 | uint32(data[2])<<16, nil
}

func SetU24(buf *bytes.Buffer, value uint32) error {
	if value > 0xFFFFFF {
		return fmt.Errorf("u24 out of range: %d", value)
	}
	_, err := buf.Write([]byte{byte(value), byte(value >> 8), byte(value >> 16)})
	return err
}

func GetU32(buf *bytes.Buffer) (uint32, error) {
	if buf.Len() < 4 {
		return 0, fmt.Errorf("not enough data")
	}
	data := buf.Next(4)
	return uint32(data[0]) |
		uint32(data[1])<<8 |
		uint32(data[2])<<16 |
		uint32(data[3])<<24, nil
}

func SetU32(buf *bytes.Buffer, value uint32) error {
	_, err := buf.Write([]byte{byte(value), byte(value >> 8), byte(value >> 16), byte(value >> 24)})
	return err
}

func GetI32(buf *bytes.Buffer) (int32, error) {
	value, err := GetU32(buf)
	return int32(value), err
}

func SetI32(buf *bytes.Buffer, value int32) error {
	return SetU32(buf, uint32(value))
}

func GetU64(buf *bytes.Buffer) (uint64, error) {
	if buf.Len() < 8 {
		return 0, fmt.Errorf("not enough data")
	}
	data := buf.Next(8)
	return uint64(data[0]) |
		uint64(data[1])<<8 |
		uint64(data[2])<<16 |
		uint64(data[3])<<24 |
		uint64(data[4])<<32 |
		uint64(data[5])<<40 |
		uint64(data[6])<<48 |
		uint64(data[7])<<56, nil
}

func SetU64(buf *bytes.Buffer, value uint64) error {
	_, err := buf.Write([]byte{
		byte(value),
		byte(value >> 8),
		byte(value >> 16),
		byte(value >> 24),
		byte(value >> 32),
		byte(value >> 40),
		byte(value >> 48),
		byte(value >> 56),
	})
	return err
}

func GetI64(buf *bytes.Buffer) (int64, error) {
	value, err := GetU64(buf)
	return int64(value), err
}

func SetI64(buf *bytes.Buffer, value int64) error {
	return SetU64(buf, uint64(value))
}

func GetF32(buf *bytes.Buffer) (float32, error) {
	value, err := GetU32(buf)
	if err != nil {
		return 0, err
	}
	return math.Float32frombits(value), nil
}

func SetF32(buf *bytes.Buffer, value float32) error {
	return SetU32(buf, math.Float32bits(value))
}

func GetF64(buf *bytes.Buffer) (float64, error) {
	value, err := GetU64(buf)
	if err != nil {
		return 0, err
	}
	return math.Float64frombits(value), nil
}

func SetF64(buf *bytes.Buffer, value float64) error {
	return SetU64(buf, math.Float64bits(value))
}

func textLengthState(length int) (uint8, error) {
	if length < 0 {
		return 0, fmt.Errorf("negative text length: %d", length)
	}
	if length == 0 {
		return StateZero, nil
	}
	if length <= 0xFF {
		return StateU8, nil
	}
	if length <= 0xFFFF {
		return StateU16, nil
	}
	return 0, fmt.Errorf("text length exceeds u16 max: %d", length)
}

func TextState(length int) (uint8, error) {
	return textLengthState(length)
}

func BinState(length int) (uint8, error) {
	if length < 0 {
		return 0, fmt.Errorf("negative bin length: %d", length)
	}
	if length == 0 {
		return StateZero, nil
	}
	if length <= 0xFF {
		return StateU8, nil
	}
	if length <= 0xFFFF {
		return StateU16, nil
	}
	if length <= 0xFFFFFF {
		return StateU24, nil
	}
	return 0, fmt.Errorf("bin length exceeds u24 max: %d", length)
}

func ListCountState(count int) (uint8, error) {
	if count < 0 {
		return 0, fmt.Errorf("negative list count: %d", count)
	}
	if count == 0 {
		return StateZero, nil
	}
	if count <= 0xFF {
		return StateU8, nil
	}
	return 0, fmt.Errorf("list count exceeds u8 max: %d", count)
}

func readStateLength(buf *bytes.Buffer, state uint8, max int, kind string) (int, error) {
	switch state {
	case StateZero:
		return 0, nil
	case StateU8:
		value, err := GetU8(buf)
		if err != nil {
			return 0, err
		}
		if value == 0 {
			return 0, fmt.Errorf("%s state %d encoded zero length", kind, state)
		}
		return int(value), nil
	case StateU16:
		value, err := GetU16(buf)
		if err != nil {
			return 0, err
		}
		if value == 0 {
			return 0, fmt.Errorf("%s state %d encoded zero length", kind, state)
		}
		if value <= 0xFF {
			return 0, fmt.Errorf("%s state %d is not canonical for length %d", kind, state, value)
		}
		return int(value), nil
	case StateU24:
		if max < 0x10000 {
			return 0, fmt.Errorf("%s state %d is illegal", kind, state)
		}
		value, err := GetU24(buf)
		if err != nil {
			return 0, err
		}
		if value == 0 {
			return 0, fmt.Errorf("%s state %d encoded zero length", kind, state)
		}
		if value <= 0xFFFF {
			return 0, fmt.Errorf("%s state %d is not canonical for length %d", kind, state, value)
		}
		return int(value), nil
	default:
		return 0, fmt.Errorf("invalid %s state: %d", kind, state)
	}
}

func GetText(buf *bytes.Buffer, state uint8) (string, error) {
	length, err := readStateLength(buf, state, 0xFFFF, "text")
	if err != nil {
		return "", err
	}
	if length == 0 {
		return "", nil
	}
	if buf.Len() < length {
		return "", fmt.Errorf("text body not enough data: %d - %d", buf.Len(), length)
	}
	data := buf.Next(length)
	if !utf8.Valid(data) {
		return "", fmt.Errorf("text body invalid utf-8")
	}
	return string(data), nil
}

func SetText(buf *bytes.Buffer, state uint8, value string) error {
	if !utf8.ValidString(value) {
		return fmt.Errorf("text value invalid utf-8")
	}
	canonical, err := TextState(len(value))
	if err != nil {
		return err
	}
	if state != canonical {
		return fmt.Errorf("text state %d is not canonical for length %d", state, len(value))
	}
	switch state {
	case StateZero:
		return nil
	case StateU8:
		if err := SetU8(buf, uint8(len(value))); err != nil {
			return err
		}
	case StateU16:
		if err := SetU16(buf, uint16(len(value))); err != nil {
			return err
		}
	default:
		return fmt.Errorf("invalid text state: %d", state)
	}
	_, err = buf.WriteString(value)
	return err
}

func SizeText(value string) (int, error) {
	if !utf8.ValidString(value) {
		return 0, fmt.Errorf("text value invalid utf-8")
	}
	state, err := TextState(len(value))
	if err != nil {
		return 0, err
	}
	switch state {
	case StateZero:
		return 0, nil
	case StateU8:
		return 1 + len(value), nil
	case StateU16:
		return 2 + len(value), nil
	default:
		return 0, fmt.Errorf("invalid text state: %d", state)
	}
}

func GetBoolListInto(buf *bytes.Buffer, state uint8, dst []bool) ([]bool, error) {
	count, err := getListCount(buf, state)
	if err != nil {
		return nil, err
	}
	if count == 0 {
		if dst != nil {
			return dst[:0], nil
		}
		return nil, nil
	}
	bitmap, err := readBitmap(buf, count, "bool list bitmap")
	if err != nil {
		return nil, err
	}
	if cap(dst) >= count {
		dst = dst[:count]
	} else {
		dst = make([]bool, count)
	}
	for i := 0; i < count; i++ {
		dst[i] = bitmapBit(bitmap, i)
	}
	return dst, nil
}

func SetBoolList(buf *bytes.Buffer, state uint8, values []bool) error {
	canonical, err := ListCountState(len(values))
	if err != nil {
		return err
	}
	if state != canonical {
		return fmt.Errorf("bool list state %d is not canonical for count %d", state, len(values))
	}
	if err := setListCount(buf, state, len(values)); err != nil {
		return err
	}
	if len(values) == 0 {
		return nil
	}
	bitmap := make([]byte, bitmapSize(len(values)))
	for i, value := range values {
		if !value {
			continue
		}
		byteIndex := i / 8
		bitIndex := 7 - (i % 8)
		bitmap[byteIndex] |= 1 << bitIndex
	}
	_, err = buf.Write(bitmap)
	return err
}

func SizeBoolList(values []bool) (int, error) {
	state, err := ListCountState(len(values))
	if err != nil {
		return 0, err
	}
	if state == StateZero {
		return 0, nil
	}
	return 1 + bitmapSize(len(values)), nil
}

func GetBinInto(buf *bytes.Buffer, state uint8, dst []byte) ([]byte, error) {
	length, err := readStateLength(buf, state, 0xFFFFFF, "bin")
	if err != nil {
		return nil, err
	}
	if length == 0 {
		if dst != nil {
			return dst[:0], nil
		}
		return nil, nil
	}
	if buf.Len() < length {
		return nil, fmt.Errorf("bin body not enough data: %d - %d", buf.Len(), length)
	}
	if cap(dst) >= length {
		dst = dst[:length]
	} else {
		dst = make([]byte, length)
	}
	copy(dst, buf.Next(length))
	return dst, nil
}

func SetBin(buf *bytes.Buffer, state uint8, value []byte) error {
	canonical, err := BinState(len(value))
	if err != nil {
		return err
	}
	if state != canonical {
		return fmt.Errorf("bin state %d is not canonical for length %d", state, len(value))
	}
	switch state {
	case StateZero:
		return nil
	case StateU8:
		if err := SetU8(buf, uint8(len(value))); err != nil {
			return err
		}
	case StateU16:
		if err := SetU16(buf, uint16(len(value))); err != nil {
			return err
		}
	case StateU24:
		if err := SetU24(buf, uint32(len(value))); err != nil {
			return err
		}
	default:
		return fmt.Errorf("invalid bin state: %d", state)
	}
	_, err = buf.Write(value)
	return err
}

func SizeBin(value []byte) (int, error) {
	state, err := BinState(len(value))
	if err != nil {
		return 0, err
	}
	switch state {
	case StateZero:
		return 0, nil
	case StateU8:
		return 1 + len(value), nil
	case StateU16:
		return 2 + len(value), nil
	case StateU24:
		return 3 + len(value), nil
	default:
		return 0, fmt.Errorf("invalid bin state: %d", state)
	}
}

func GetBinListInto(buf *bytes.Buffer, state uint8, dst [][]byte) ([][]byte, error) {
	count, err := getListCount(buf, state)
	if err != nil {
		return nil, err
	}
	if count == 0 {
		if dst != nil {
			return dst[:0], nil
		}
		return nil, nil
	}
	states, err := readStateBlock(buf, count, "bin list state block")
	if err != nil {
		return nil, err
	}
	if cap(dst) >= count {
		dst = dst[:count]
	} else {
		dst = make([][]byte, count)
	}
	for i, itemState := range states {
		item, itemErr := GetBinInto(buf, itemState, dst[i])
		if itemErr != nil {
			return nil, fmt.Errorf("bin list[%d]: %w", i, itemErr)
		}
		dst[i] = item
	}
	return dst, nil
}

func SetBinList(buf *bytes.Buffer, state uint8, values [][]byte) error {
	canonical, err := ListCountState(len(values))
	if err != nil {
		return err
	}
	if state != canonical {
		return fmt.Errorf("bin list state %d is not canonical for count %d", state, len(values))
	}
	if err := setListCount(buf, state, len(values)); err != nil {
		return err
	}
	if len(values) == 0 {
		return nil
	}
	states := make([]uint8, len(values))
	for i, item := range values {
		itemState, itemErr := BinState(len(item))
		if itemErr != nil {
			return fmt.Errorf("bin list[%d]: %w", i, itemErr)
		}
		states[i] = itemState
	}
	header, err := writeStateBlock(states)
	if err != nil {
		return err
	}
	if _, err := buf.Write(header); err != nil {
		return err
	}
	for i, item := range values {
		if err := SetBin(buf, states[i], item); err != nil {
			return fmt.Errorf("bin list[%d]: %w", i, err)
		}
	}
	return nil
}

func SizeBinList(values [][]byte) (int, error) {
	state, err := ListCountState(len(values))
	if err != nil {
		return 0, err
	}
	if state == StateZero {
		return 0, nil
	}
	size := 1 + itemHeaderSize(len(values))
	for i, item := range values {
		itemSize, itemErr := SizeBin(item)
		if itemErr != nil {
			return 0, fmt.Errorf("bin list[%d]: %w", i, itemErr)
		}
		size += itemSize
	}
	return size, nil
}

func getListCount(buf *bytes.Buffer, state uint8) (int, error) {
	switch state {
	case StateZero:
		return 0, nil
	case StateU8:
		value, err := GetU8(buf)
		if err != nil {
			return 0, err
		}
		if value == 0 {
			return 0, fmt.Errorf("list count state %d encoded zero count", state)
		}
		return int(value), nil
	case StateU16, StateU24:
		return 0, fmt.Errorf("list count state %d is illegal", state)
	default:
		return 0, fmt.Errorf("invalid list count state: %d", state)
	}
}

func setListCount(buf *bytes.Buffer, state uint8, count int) error {
	canonical, err := ListCountState(count)
	if err != nil {
		return err
	}
	if state != canonical {
		return fmt.Errorf("list count state %d is not canonical for count %d", state, count)
	}
	switch state {
	case StateZero:
		return nil
	case StateU8:
		return SetU8(buf, uint8(count))
	default:
		return fmt.Errorf("invalid list count state: %d", state)
	}
}

func readBitmap(buf *bytes.Buffer, count int, kind string) ([]byte, error) {
	size := bitmapSize(count)
	if buf.Len() < size {
		return nil, fmt.Errorf("%s not enough data: %d - %d", kind, buf.Len(), size)
	}
	data := append([]byte(nil), buf.Next(size)...)
	if err := ValidatePaddingZero(data, count, "bitmap"); err != nil {
		return nil, err
	}
	return data, nil
}

func bitmapBit(bitmap []byte, index int) bool {
	byteIndex := index / 8
	bitIndex := 7 - (index % 8)
	return (bitmap[byteIndex] & (1 << bitIndex)) != 0
}

func writeBitmapFromDefaults[T any](values []T, isDefault func(T) bool) ([]byte, int) {
	bitmap := make([]byte, bitmapSize(len(values)))
	nonDefaultCount := 0
	for i, item := range values {
		if isDefault(item) {
			continue
		}
		byteIndex := i / 8
		bitIndex := 7 - (i % 8)
		bitmap[byteIndex] |= 1 << bitIndex
		nonDefaultCount++
	}
	return bitmap, nonDefaultCount
}

func GetDefaultListInto[T any](buf *bytes.Buffer, state uint8, dst []T, defaultItem func() T, getItem func(*bytes.Buffer) (T, error)) ([]T, error) {
	count, err := getListCount(buf, state)
	if err != nil {
		return nil, err
	}
	if count == 0 {
		if dst != nil {
			return dst[:0], nil
		}
		return nil, nil
	}
	bitmap, err := readBitmap(buf, count, "bitmap list")
	if err != nil {
		return nil, err
	}
	if cap(dst) >= count {
		dst = dst[:count]
	} else {
		dst = make([]T, count)
	}
	for i := 0; i < count; i++ {
		if !bitmapBit(bitmap, i) {
			dst[i] = defaultItem()
			continue
		}
		item, itemErr := getItem(buf)
		if itemErr != nil {
			return nil, fmt.Errorf("bitmap list[%d]: %w", i, itemErr)
		}
		dst[i] = item
	}
	return dst, nil
}

func SetDefaultList[T any](buf *bytes.Buffer, state uint8, values []T, isDefault func(T) bool, sizeItem func(T) (int, error), setItem func(*bytes.Buffer, T) error) error {
	canonical, err := ListCountState(len(values))
	if err != nil {
		return err
	}
	if state != canonical {
		return fmt.Errorf("bitmap list state %d is not canonical for count %d", state, len(values))
	}
	if err := setListCount(buf, state, len(values)); err != nil {
		return err
	}
	if len(values) == 0 {
		return nil
	}
	bitmap, _ := writeBitmapFromDefaults(values, isDefault)
	if _, err := buf.Write(bitmap); err != nil {
		return err
	}
	for i, item := range values {
		if isDefault(item) {
			continue
		}
		if _, err := sizeItem(item); err != nil {
			return fmt.Errorf("bitmap list[%d]: %w", i, err)
		}
		if err := setItem(buf, item); err != nil {
			return fmt.Errorf("bitmap list[%d]: %w", i, err)
		}
	}
	return nil
}

func SizeDefaultList[T any](values []T, isDefault func(T) bool, sizeItem func(T) (int, error)) (int, error) {
	state, err := ListCountState(len(values))
	if err != nil {
		return 0, err
	}
	if state == StateZero {
		return 0, nil
	}
	size := 1 + bitmapSize(len(values))
	for i, item := range values {
		if isDefault(item) {
			continue
		}
		itemSize, itemErr := sizeItem(item)
		if itemErr != nil {
			return 0, fmt.Errorf("bitmap list[%d]: %w", i, itemErr)
		}
		size += itemSize
	}
	return size, nil
}

func GetZeroListInto[T comparable](buf *bytes.Buffer, state uint8, dst []T, getItem func(*bytes.Buffer) (T, error)) ([]T, error) {
	var zero T
	return GetDefaultListInto(buf, state, dst, func() T { return zero }, getItem)
}

func SetZeroList[T comparable](buf *bytes.Buffer, state uint8, values []T, itemSize int, setItem func(*bytes.Buffer, T) error) error {
	var zero T
	return SetDefaultList(
		buf,
		state,
		values,
		func(item T) bool { return item == zero },
		func(T) (int, error) { return itemSize, nil },
		setItem,
	)
}

func SizeZeroList[T comparable](values []T, itemSize int) (int, error) {
	var zero T
	return SizeDefaultList(
		values,
		func(item T) bool { return item == zero },
		func(T) (int, error) { return itemSize, nil },
	)
}

func GetDefaultPtrListInto[T any](buf *bytes.Buffer, state uint8, dst []*T, defaultItem func() *T, getItem func(*bytes.Buffer, *T) error) ([]*T, error) {
	count, err := getListCount(buf, state)
	if err != nil {
		return nil, err
	}
	if count == 0 {
		if dst != nil {
			return dst[:0], nil
		}
		return nil, nil
	}
	bitmap, err := readBitmap(buf, count, "bitmap ptr list")
	if err != nil {
		return nil, err
	}
	if cap(dst) >= count {
		dst = dst[:count]
	} else {
		dst = make([]*T, count)
	}
	for i := 0; i < count; i++ {
		if !bitmapBit(bitmap, i) {
			dst[i] = nil
			continue
		}
		item := dst[i]
		if item == nil {
			item = defaultItem()
		}
		if err := getItem(buf, item); err != nil {
			return nil, fmt.Errorf("bitmap list[%d]: %w", i, err)
		}
		dst[i] = item
	}
	return dst, nil
}

func SetDefaultPtrList[T any](buf *bytes.Buffer, state uint8, values []*T, typeName string, validate func(*T) error, isZero func(*T) bool, sizeItem func(*T) (int, error), setItem func(*bytes.Buffer, *T) error) error {
	canonical, err := ListCountState(len(values))
	if err != nil {
		return err
	}
	if state != canonical {
		return fmt.Errorf("%s list state %d is not canonical for count %d", typeName, state, len(values))
	}
	if err := setListCount(buf, state, len(values)); err != nil {
		return err
	}
	if len(values) == 0 {
		return nil
	}
	bitmap, _ := writeBitmapFromDefaults(values, isZero)
	if _, err := buf.Write(bitmap); err != nil {
		return err
	}
	for i, item := range values {
		if isZero(item) {
			continue
		}
		if err := validate(item); err != nil {
			return fmt.Errorf("%s list[%d]: %w", typeName, i, err)
		}
		if _, err := sizeItem(item); err != nil {
			return fmt.Errorf("%s list[%d]: %w", typeName, i, err)
		}
		if err := setItem(buf, item); err != nil {
			return fmt.Errorf("%s list[%d]: %w", typeName, i, err)
		}
	}
	return nil
}

func SizeDefaultPtrList[T any](values []*T, typeName string, validate func(*T) error, isZero func(*T) bool, sizeItem func(*T) (int, error)) (int, error) {
	state, err := ListCountState(len(values))
	if err != nil {
		return 0, err
	}
	if state == StateZero {
		return 0, nil
	}
	size := 1 + bitmapSize(len(values))
	for i, item := range values {
		if isZero(item) {
			continue
		}
		if err := validate(item); err != nil {
			return 0, fmt.Errorf("%s list[%d]: %w", typeName, i, err)
		}
		itemSize, itemErr := sizeItem(item)
		if itemErr != nil {
			return 0, fmt.Errorf("%s list[%d]: %w", typeName, i, itemErr)
		}
		size += itemSize
	}
	return size, nil
}

func readStateBlock(buf *bytes.Buffer, count int, kind string) ([]uint8, error) {
	size := itemHeaderSize(count)
	if buf.Len() < size {
		return nil, fmt.Errorf("%s not enough data: %d - %d", kind, buf.Len(), size)
	}
	data := buf.Next(size)
	states := make([]uint8, count)
	reader := NewBitReader(data, count*2)
	for i := 0; i < count; i++ {
		state, err := reader.ReadBits(2)
		if err != nil {
			return nil, err
		}
		states[i] = state
	}
	if err := ValidatePaddingZero(data, count*2, "state block"); err != nil {
		return nil, err
	}
	return states, nil
}

func writeStateBlock(states []uint8) ([]byte, error) {
	data := make([]byte, itemHeaderSize(len(states)))
	bitOffset := 0
	for _, state := range states {
		if state > 0x03 {
			return nil, fmt.Errorf("invalid state: %d", state)
		}
		for shift := 1; shift >= 0; shift-- {
			bit := (state >> shift) & 1
			byteIndex := bitOffset / 8
			bitIndex := 7 - (bitOffset % 8)
			data[byteIndex] |= bit << bitIndex
			bitOffset++
		}
	}
	return data, nil
}

func GetTextListInto(buf *bytes.Buffer, state uint8, dst []string) ([]string, error) {
	count, err := getListCount(buf, state)
	if err != nil {
		return nil, err
	}
	if count == 0 {
		if dst != nil {
			return dst[:0], nil
		}
		return nil, nil
	}
	states, err := readStateBlock(buf, count, "text list state block")
	if err != nil {
		return nil, err
	}
	if cap(dst) >= count {
		dst = dst[:count]
	} else {
		dst = make([]string, count)
	}
	for i, itemState := range states {
		item, itemErr := GetText(buf, itemState)
		if itemErr != nil {
			return nil, fmt.Errorf("text list[%d]: %w", i, itemErr)
		}
		dst[i] = item
	}
	return dst, nil
}

func SetTextList(buf *bytes.Buffer, state uint8, values []string) error {
	canonical, err := ListCountState(len(values))
	if err != nil {
		return err
	}
	if state != canonical {
		return fmt.Errorf("text list state %d is not canonical for count %d", state, len(values))
	}
	if err := setListCount(buf, state, len(values)); err != nil {
		return err
	}
	if len(values) == 0 {
		return nil
	}
	states := make([]uint8, len(values))
	for i, item := range values {
		itemState, itemErr := TextState(len(item))
		if itemErr != nil {
			return fmt.Errorf("text list[%d]: %w", i, itemErr)
		}
		states[i] = itemState
	}
	header, err := writeStateBlock(states)
	if err != nil {
		return err
	}
	if _, err := buf.Write(header); err != nil {
		return err
	}
	for i, item := range values {
		if err := SetText(buf, states[i], item); err != nil {
			return fmt.Errorf("text list[%d]: %w", i, err)
		}
	}
	return nil
}

func SizeTextList(values []string) (int, error) {
	state, err := ListCountState(len(values))
	if err != nil {
		return 0, err
	}
	if state == StateZero {
		return 0, nil
	}
	size := 1 + itemHeaderSize(len(values))
	for i, item := range values {
		itemSize, itemErr := SizeText(item)
		if itemErr != nil {
			return 0, fmt.Errorf("text list[%d]: %w", i, itemErr)
		}
		size += itemSize
	}
	return size, nil
}
