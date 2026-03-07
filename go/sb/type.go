package sb

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"math"
	"slices"
)

type Setter func(*bytes.Buffer) error
type Getter func(*bytes.Buffer) error

func SetAll(buf *bytes.Buffer, args ...Setter) error {
	for i, arg := range args {
		if arg == nil { return fmt.Errorf("SetAll arg[%d] is nil", i) }
		if err := arg(buf); err != nil { return err }
	}
	return nil
}

func GetAll(buf *bytes.Buffer, args ...Getter) error {
	for i, arg := range args {
		if arg == nil { return fmt.Errorf("GetAll arg[%d] is nil", i) }
		if err := arg(buf); err != nil { return err }
	}
	return nil
}

// Helpers
func getList[T any, L ~[]T](buf *bytes.Buffer, getItem func(*bytes.Buffer) (T, error)) (L, error) {
	count, err := GetU16(buf)
	if err != nil { return nil, err }
	list := make([]T, count)
	for i := range list {
		if list[i], err = getItem(buf); err != nil { return nil, err }
	}
	return L(list), nil
}

func setList[T any](buf *bytes.Buffer, list []T, setItem func(*bytes.Buffer, T) error) error {
	if len(list) > 65535 { return fmt.Errorf("list length exceeds uint16 max") }
	if err := SetU16(buf, uint16(len(list))); err != nil { return err }
	for _, item := range list {
		if err := setItem(buf, item); err != nil { return err }
	}
	return nil
}

func sizeFixedList(length int, elemSize int) (int, error) {
	if length > 65535 { return 0, fmt.Errorf("list length exceeds uint16 max") }
	return 2 + length*elemSize, nil
}

func sizeBoolList(v []bool) (int, error) {
	if len(v) > 65535 { return 0, fmt.Errorf("list length exceeds uint16 max") }
	return 2 + (len(v)+7)/8, nil
}

func sizeBin(v []byte) (int, error) {
	if uint64(len(v)) > uint64(^uint32(0)) { return 0, fmt.Errorf("bin length exceeds uint32 max") }
	return 4 + len(v), nil
}

func sizeBinList(v [][]byte) (int, error) {
	if len(v) > 65535 { return 0, fmt.Errorf("list length exceeds uint16 max") }
	size := 2
	for i, item := range v {
		itemSize, err := sizeBin(item)
		if err != nil { return 0, fmt.Errorf("bin list[%d]: %w", i, err) }
		size += itemSize
	}
	return size, nil
}

func sizeText(v string) (int, error) {
	if len(v) > 65535 { return 0, fmt.Errorf("text length exceeds uint16 max") }
	return 2 + len(v), nil
}

func sizeTextList(v []string) (int, error) {
	if len(v) > 65535 { return 0, fmt.Errorf("list length exceeds uint16 max") }
	size := 2
	for i, item := range v {
		itemSize, err := sizeText(item)
		if err != nil { return 0, fmt.Errorf("text list[%d]: %w", i, err) }
		size += itemSize
	}
	return size, nil
}

func GetBit(bits []byte, i uint8) bool {
	if int(i/8) >= len(bits) { return false }
	return (bits[i/8] & (1 << (i % 8))) != 0
}

func SetBit(bits []byte, i uint8, v bool) {
	if int(i/8) >= len(bits) { return }
	if v { bits[i/8] |= (1 << (i % 8)) } else { bits[i/8] &= ^(1 << (i % 8)) }
}

// Bool
type Bool bool
func GetBool(buf *bytes.Buffer) (bool, error) {
	b, err := buf.ReadByte()
	return b == 1, err
}
func SetBool(buf *bytes.Buffer, v bool) error {
	if v { return buf.WriteByte(1) }
	return buf.WriteByte(0)
}
func EqBool(a, b bool) bool { return a == b }

type BoolList []bool
func GetBoolList(buf *bytes.Buffer) ([]bool, error) {
	count, err := GetU16(buf)
	if err != nil { return nil, err }
	bitSize := (int(count) + 7) / 8
	if buf.Len() < bitSize { return nil, fmt.Errorf("not enough data") }
	bits := buf.Next(bitSize)
	bools := make([]bool, int(count))
	for i := 0; i < int(count); i++ { bools[i] = GetBit(bits, uint8(i)) }
	return bools, nil
}
func SetBoolList(buf *bytes.Buffer, v []bool) error {
	if len(v) > 65535 { return fmt.Errorf("list length exceeds uint16 max") }
	bitSize := (len(v) + 7) / 8
	buf.Grow(2 + bitSize)
	b := buf.AvailableBuffer()
	b = binary.LittleEndian.AppendUint16(b, uint16(len(v)))
	start := len(b)
	for i := 0; i < bitSize; i++ { b = append(b, 0) }
	for i, val := range v {
		if val { b[start+(i/8)] |= 1 << (uint(i) % 8) }
	}
	_, err := buf.Write(b)
	return err
}
func EqBoolList(a, b []bool) bool { return slices.Equal(a, b) }

// Primitives Macro
type I8 int8
func GetI8(buf *bytes.Buffer) (int8, error) {
	b, err := buf.ReadByte()
	return int8(b), err
}
func SetI8(buf *bytes.Buffer, v int8) error {
	return buf.WriteByte(uint8(v))
}
func EqI8(a, b int8) bool { return a == b }

type I8List []int8
func GetI8List(buf *bytes.Buffer) ([]int8, error) {
	count, err := GetU16(buf)
	if err != nil { return nil, err }
	if buf.Len() < int(count) { return nil, fmt.Errorf("not enough data") }
	data := buf.Next(int(count))
	list := make([]int8, count)
	for i := range list { list[i] = int8(data[i]) }
	return list, nil
}
func SetI8List(buf *bytes.Buffer, v []int8) error {
	if len(v) > 65535 { return fmt.Errorf("list length exceeds uint16 max") }
	buf.Grow(2 + len(v))
	if err := SetU16(buf, uint16(len(v))); err != nil { return err }
	if len(v) == 0 { return nil }
	b := buf.AvailableBuffer()
	for _, item := range v { b = append(b, byte(item)) }
	_, err := buf.Write(b)
	return err
}
func EqI8List(a, b []int8) bool { return slices.Equal(a, b) }
type U8 uint8
func GetU8(buf *bytes.Buffer) (uint8, error) {
	b, err := buf.ReadByte()
	return uint8(b), err
}
func SetU8(buf *bytes.Buffer, v uint8) error {
	return buf.WriteByte(uint8(v))
}
func EqU8(a, b uint8) bool { return a == b }

type U8List []uint8
func GetU8List(buf *bytes.Buffer) ([]uint8, error) {
	count, err := GetU16(buf)
	if err != nil { return nil, err }
	if buf.Len() < int(count) { return nil, fmt.Errorf("not enough data") }
	return bytes.Clone(buf.Next(int(count))), nil
}
func SetU8List(buf *bytes.Buffer, v []uint8) error {
	if len(v) > 65535 { return fmt.Errorf("list length exceeds uint16 max") }
	buf.Grow(2 + len(v))
	if err := SetU16(buf, uint16(len(v))); err != nil { return err }
	if len(v) == 0 { return nil }
	_, err := buf.Write(v)
	return err
}
func EqU8List(a, b []uint8) bool { return slices.Equal(a, b) }
type I16 int16
func GetI16(buf *bytes.Buffer) (int16, error) {
	var b [2]byte
	if _, err := buf.Read(b[:]); err != nil { return 0, err }
	return int16(binary.LittleEndian.Uint16(b[:])), nil
}
func SetI16(buf *bytes.Buffer, v int16) error {
	buf.Grow(2)
	b := buf.AvailableBuffer()
	b = binary.LittleEndian.AppendUint16(b, uint16(v))
	_, err := buf.Write(b)
	return err
}
func EqI16(a, b int16) bool { return a == b }

type I16List []int16
func GetI16List(buf *bytes.Buffer) ([]int16, error) {
	count, err := GetU16(buf)
	if err != nil { return nil, err }
	totalSize := int(count) * 2
	if buf.Len() < totalSize { return nil, fmt.Errorf("not enough data") }
	data := buf.Next(totalSize)
	list := make([]int16, count)
	for i := range list { list[i] = int16(binary.LittleEndian.Uint16(data[i*2:])) }
	return list, nil
}
func SetI16List(buf *bytes.Buffer, v []int16) error {
	if len(v) > 65535 { return fmt.Errorf("list length exceeds uint16 max") }
	buf.Grow(2 + len(v)*2)
	if err := SetU16(buf, uint16(len(v))); err != nil { return err }
	if len(v) == 0 { return nil }
	b := buf.AvailableBuffer()
	for _, item := range v { b = binary.LittleEndian.AppendUint16(b, uint16(item)) }
	_, err := buf.Write(b)
	return err
}
func EqI16List(a, b []int16) bool { return slices.Equal(a, b) }
type U16 uint16
func GetU16(buf *bytes.Buffer) (uint16, error) {
	var b [2]byte
	if _, err := buf.Read(b[:]); err != nil { return 0, err }
	return uint16(binary.LittleEndian.Uint16(b[:])), nil
}
func SetU16(buf *bytes.Buffer, v uint16) error {
	buf.Grow(2)
	b := buf.AvailableBuffer()
	b = binary.LittleEndian.AppendUint16(b, uint16(v))
	_, err := buf.Write(b)
	return err
}
func EqU16(a, b uint16) bool { return a == b }

type U16List []uint16
func GetU16List(buf *bytes.Buffer) ([]uint16, error) {
	count, err := GetU16(buf)
	if err != nil { return nil, err }
	totalSize := int(count) * 2
	if buf.Len() < totalSize { return nil, fmt.Errorf("not enough data") }
	data := buf.Next(totalSize)
	list := make([]uint16, count)
	for i := range list { list[i] = binary.LittleEndian.Uint16(data[i*2:]) }
	return list, nil
}
func SetU16List(buf *bytes.Buffer, v []uint16) error {
	if len(v) > 65535 { return fmt.Errorf("list length exceeds uint16 max") }
	buf.Grow(2 + len(v)*2)
	if err := SetU16(buf, uint16(len(v))); err != nil { return err }
	if len(v) == 0 { return nil }
	b := buf.AvailableBuffer()
	for _, item := range v { b = binary.LittleEndian.AppendUint16(b, item) }
	_, err := buf.Write(b)
	return err
}
func EqU16List(a, b []uint16) bool { return slices.Equal(a, b) }
type I32 int32
func GetI32(buf *bytes.Buffer) (int32, error) {
	var b [4]byte
	if _, err := buf.Read(b[:]); err != nil { return 0, err }
	return int32(binary.LittleEndian.Uint32(b[:])), nil
}
func SetI32(buf *bytes.Buffer, v int32) error {
	buf.Grow(4)
	b := buf.AvailableBuffer()
	b = binary.LittleEndian.AppendUint32(b, uint32(v))
	_, err := buf.Write(b)
	return err
}
func EqI32(a, b int32) bool { return a == b }

type I32List []int32
func GetI32List(buf *bytes.Buffer) ([]int32, error) {
	count, err := GetU16(buf)
	if err != nil { return nil, err }
	totalSize := int(count) * 4
	if buf.Len() < totalSize { return nil, fmt.Errorf("not enough data") }
	data := buf.Next(totalSize)
	list := make([]int32, count)
	for i := range list { list[i] = int32(binary.LittleEndian.Uint32(data[i*4:])) }
	return list, nil
}
func SetI32List(buf *bytes.Buffer, v []int32) error {
	if len(v) > 65535 { return fmt.Errorf("list length exceeds uint16 max") }
	buf.Grow(2 + len(v)*4)
	if err := SetU16(buf, uint16(len(v))); err != nil { return err }
	if len(v) == 0 { return nil }
	b := buf.AvailableBuffer()
	for _, item := range v { b = binary.LittleEndian.AppendUint32(b, uint32(item)) }
	_, err := buf.Write(b)
	return err
}
func EqI32List(a, b []int32) bool { return slices.Equal(a, b) }
type U32 uint32
func GetU32(buf *bytes.Buffer) (uint32, error) {
	var b [4]byte
	if _, err := buf.Read(b[:]); err != nil { return 0, err }
	return uint32(binary.LittleEndian.Uint32(b[:])), nil
}
func SetU32(buf *bytes.Buffer, v uint32) error {
	buf.Grow(4)
	b := buf.AvailableBuffer()
	b = binary.LittleEndian.AppendUint32(b, uint32(v))
	_, err := buf.Write(b)
	return err
}
func EqU32(a, b uint32) bool { return a == b }

type U32List []uint32
func GetU32List(buf *bytes.Buffer) ([]uint32, error) {
	count, err := GetU16(buf)
	if err != nil { return nil, err }
	totalSize := int(count) * 4
	if buf.Len() < totalSize { return nil, fmt.Errorf("not enough data") }
	data := buf.Next(totalSize)
	list := make([]uint32, count)
	for i := range list { list[i] = binary.LittleEndian.Uint32(data[i*4:]) }
	return list, nil
}
func SetU32List(buf *bytes.Buffer, v []uint32) error {
	if len(v) > 65535 { return fmt.Errorf("list length exceeds uint16 max") }
	buf.Grow(2 + len(v)*4)
	if err := SetU16(buf, uint16(len(v))); err != nil { return err }
	if len(v) == 0 { return nil }
	b := buf.AvailableBuffer()
	for _, item := range v { b = binary.LittleEndian.AppendUint32(b, item) }
	_, err := buf.Write(b)
	return err
}
func EqU32List(a, b []uint32) bool { return slices.Equal(a, b) }
type I64 int64
func GetI64(buf *bytes.Buffer) (int64, error) {
	var b [8]byte
	if _, err := buf.Read(b[:]); err != nil { return 0, err }
	return int64(binary.LittleEndian.Uint64(b[:])), nil
}
func SetI64(buf *bytes.Buffer, v int64) error {
	buf.Grow(8)
	b := buf.AvailableBuffer()
	b = binary.LittleEndian.AppendUint64(b, uint64(v))
	_, err := buf.Write(b)
	return err
}
func EqI64(a, b int64) bool { return a == b }

type I64List []int64
func GetI64List(buf *bytes.Buffer) ([]int64, error) {
	count, err := GetU16(buf)
	if err != nil { return nil, err }
	totalSize := int(count) * 8
	if buf.Len() < totalSize { return nil, fmt.Errorf("not enough data") }
	data := buf.Next(totalSize)
	list := make([]int64, count)
	for i := range list { list[i] = int64(binary.LittleEndian.Uint64(data[i*8:])) }
	return list, nil
}
func SetI64List(buf *bytes.Buffer, v []int64) error {
	if len(v) > 65535 { return fmt.Errorf("list length exceeds uint16 max") }
	buf.Grow(2 + len(v)*8)
	if err := SetU16(buf, uint16(len(v))); err != nil { return err }
	if len(v) == 0 { return nil }
	b := buf.AvailableBuffer()
	for _, item := range v { b = binary.LittleEndian.AppendUint64(b, uint64(item)) }
	_, err := buf.Write(b)
	return err
}
func EqI64List(a, b []int64) bool { return slices.Equal(a, b) }
type U64 uint64
func GetU64(buf *bytes.Buffer) (uint64, error) {
	var b [8]byte
	if _, err := buf.Read(b[:]); err != nil { return 0, err }
	return uint64(binary.LittleEndian.Uint64(b[:])), nil
}
func SetU64(buf *bytes.Buffer, v uint64) error {
	buf.Grow(8)
	b := buf.AvailableBuffer()
	b = binary.LittleEndian.AppendUint64(b, uint64(v))
	_, err := buf.Write(b)
	return err
}
func EqU64(a, b uint64) bool { return a == b }

type U64List []uint64
func GetU64List(buf *bytes.Buffer) ([]uint64, error) {
	count, err := GetU16(buf)
	if err != nil { return nil, err }
	totalSize := int(count) * 8
	if buf.Len() < totalSize { return nil, fmt.Errorf("not enough data") }
	data := buf.Next(totalSize)
	list := make([]uint64, count)
	for i := range list { list[i] = binary.LittleEndian.Uint64(data[i*8:]) }
	return list, nil
}
func SetU64List(buf *bytes.Buffer, v []uint64) error {
	if len(v) > 65535 { return fmt.Errorf("list length exceeds uint16 max") }
	buf.Grow(2 + len(v)*8)
	if err := SetU16(buf, uint16(len(v))); err != nil { return err }
	if len(v) == 0 { return nil }
	b := buf.AvailableBuffer()
	for _, item := range v { b = binary.LittleEndian.AppendUint64(b, item) }
	_, err := buf.Write(b)
	return err
}
func EqU64List(a, b []uint64) bool { return slices.Equal(a, b) }
type F32 float32
func GetF32(buf *bytes.Buffer) (float32, error) {
	var b [4]byte
	if _, err := buf.Read(b[:]); err != nil { return 0, err }
	return math.Float32frombits(binary.LittleEndian.Uint32(b[:])), nil
}
func SetF32(buf *bytes.Buffer, v float32) error {
	buf.Grow(4)
	b := buf.AvailableBuffer()
	b = binary.LittleEndian.AppendUint32(b, math.Float32bits(v))
	_, err := buf.Write(b)
	return err
}
func EqF32(a, b float32) bool { return math.Abs(float64(a-b)) < 1e-6 }

type F32List []float32
func GetF32List(buf *bytes.Buffer) ([]float32, error) {
	count, err := GetU16(buf)
	if err != nil { return nil, err }
	totalSize := int(count) * 4
	if buf.Len() < totalSize { return nil, fmt.Errorf("not enough data") }
	data := buf.Next(totalSize)
	list := make([]float32, count)
	for i := range list { list[i] = math.Float32frombits(binary.LittleEndian.Uint32(data[i*4:])) }
	return list, nil
}
func SetF32List(buf *bytes.Buffer, v []float32) error {
	if len(v) > 65535 { return fmt.Errorf("list length exceeds uint16 max") }
	buf.Grow(2 + len(v)*4)
	if err := SetU16(buf, uint16(len(v))); err != nil { return err }
	if len(v) == 0 { return nil }
	b := buf.AvailableBuffer()
	for _, item := range v { b = binary.LittleEndian.AppendUint32(b, math.Float32bits(item)) }
	_, err := buf.Write(b)
	return err
}
func EqF32List(a, b []float32) bool { return slices.Equal(a, b) }
type F64 float64
func GetF64(buf *bytes.Buffer) (float64, error) {
	var b [8]byte
	if _, err := buf.Read(b[:]); err != nil { return 0, err }
	return math.Float64frombits(binary.LittleEndian.Uint64(b[:])), nil
}
func SetF64(buf *bytes.Buffer, v float64) error {
	buf.Grow(8)
	b := buf.AvailableBuffer()
	b = binary.LittleEndian.AppendUint64(b, math.Float64bits(v))
	_, err := buf.Write(b)
	return err
}
func EqF64(a, b float64) bool { return math.Abs(float64(a-b)) < 1e-9 }

type F64List []float64
func GetF64List(buf *bytes.Buffer) ([]float64, error) {
	count, err := GetU16(buf)
	if err != nil { return nil, err }
	totalSize := int(count) * 8
	if buf.Len() < totalSize { return nil, fmt.Errorf("not enough data") }
	data := buf.Next(totalSize)
	list := make([]float64, count)
	for i := range list { list[i] = math.Float64frombits(binary.LittleEndian.Uint64(data[i*8:])) }
	return list, nil
}
func SetF64List(buf *bytes.Buffer, v []float64) error {
	if len(v) > 65535 { return fmt.Errorf("list length exceeds uint16 max") }
	buf.Grow(2 + len(v)*8)
	if err := SetU16(buf, uint16(len(v))); err != nil { return err }
	if len(v) == 0 { return nil }
	b := buf.AvailableBuffer()
	for _, item := range v { b = binary.LittleEndian.AppendUint64(b, math.Float64bits(item)) }
	_, err := buf.Write(b)
	return err
}
func EqF64List(a, b []float64) bool { return slices.Equal(a, b) }


// Bin
type Bin []byte
func GetBinView(buf *bytes.Buffer) ([]byte, error) {
	l, err := GetU32(buf)
	if err != nil { return nil, err }
	if uint64(buf.Len()) < uint64(l) { return nil, fmt.Errorf("not enough data") }
	return buf.Next(int(l)), nil
}
func GetBin(buf *bytes.Buffer) ([]byte, error) {
	res, err := GetBinView(buf)
	if err != nil { return nil, err }
	return bytes.Clone(res), nil
}
func SetBin(buf *bytes.Buffer, v []byte) error {
	if uint64(len(v)) > uint64(^uint32(0)) { return fmt.Errorf("bin length exceeds uint32 max") }
	buf.Grow(4 + len(v))
	if err := SetU32(buf, uint32(len(v))); err != nil { return err }
	_, err := buf.Write(v)
	return err
}
func EqBin(a, b []byte) bool { return bytes.Equal(a, b) }

type BinList [][]byte
func GetBinList(buf *bytes.Buffer) ([][]byte, error) {
	count, err := GetU16(buf)
	if err != nil { return nil, err }
	list := make([][]byte, count)
	for i := range list {
		item, err := GetBinView(buf)
		if err != nil { return nil, err }
		list[i] = bytes.Clone(item)
	}
	return list, nil
}
func SetBinList(buf *bytes.Buffer, v [][]byte) error {
	totalSize, err := sizeBinList(v)
	if err != nil { return err }
	buf.Grow(totalSize)
	if err := SetU16(buf, uint16(len(v))); err != nil { return err }
	for _, item := range v {
		if err := SetU32(buf, uint32(len(item))); err != nil { return err }
		if len(item) == 0 { continue }
		if _, err := buf.Write(item); err != nil { return err }
	}
	return nil
}
func EqBinList(a, b [][]byte) bool { return slices.EqualFunc(a, b, bytes.Equal) }

// Text
type Text string
func GetText(buf *bytes.Buffer) (string, error) {
	l, err := GetU16(buf)
	if err != nil { return "", err }
	if buf.Len() < int(l) { return "", fmt.Errorf("not enough data") }
	return string(buf.Next(int(l))), nil
}
func SetText(buf *bytes.Buffer, v string) error {
	if len(v) > 65535 { return fmt.Errorf("text length exceeds uint16 max") }
	buf.Grow(2 + len(v))
	if err := SetU16(buf, uint16(len(v))); err != nil { return err }
	_, err := buf.WriteString(v)
	return err
}
func EqText(a, b string) bool { return a == b }

type TextList []string
func GetTextList(buf *bytes.Buffer) ([]string, error) {
	count, err := GetU16(buf)
	if err != nil { return nil, err }
	list := make([]string, count)
	for i := range list {
		item, err := GetText(buf)
		if err != nil { return nil, err }
		list[i] = item
	}
	return list, nil
}
func SetTextList(buf *bytes.Buffer, v []string) error {
	totalSize, err := sizeTextList(v)
	if err != nil { return err }
	buf.Grow(totalSize)
	if err := SetU16(buf, uint16(len(v))); err != nil { return err }
	for _, item := range v {
		if err := SetU16(buf, uint16(len(item))); err != nil { return err }
		if len(item) == 0 { continue }
		if _, err := buf.WriteString(item); err != nil { return err }
	}
	return nil
}
func EqTextList(a, b []string) bool { return slices.Equal(a, b) }
