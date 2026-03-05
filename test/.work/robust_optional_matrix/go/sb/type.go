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
	count, err := GetU8(buf)
	if err != nil { return nil, err }
	list := make([]T, count)
	for i := range list {
		if list[i], err = getItem(buf); err != nil { return nil, err }
	}
	return L(list), nil
}

func setList[T any](buf *bytes.Buffer, list []T, setItem func(*bytes.Buffer, T) error) error {
	if len(list) > 255 { return fmt.Errorf("list length exceeds uint8 max") }
	if err := SetU8(buf, uint8(len(list))); err != nil { return err }
	for _, item := range list {
		if err := setItem(buf, item); err != nil { return err }
	}
	return nil
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
	count, err := GetU8(buf)
	if err != nil { return nil, err }
	bitSize := (int(count) + 7) / 8
	if buf.Len() < bitSize { return nil, fmt.Errorf("not enough data") }
	bits := buf.Next(bitSize)
	bools := make([]bool, count)
	for i := 0; i < int(count); i++ { bools[i] = GetBit(bits, uint8(i)) }
	return bools, nil
}
func SetBoolList(buf *bytes.Buffer, v []bool) error {
	if len(v) > 255 { return fmt.Errorf("list length exceeds uint8 max") }
	if err := SetU8(buf, uint8(len(v))); err != nil { return err }
	bits := make([]byte, (len(v)+7)/8)
	for i, val := range v { SetBit(bits, uint8(i), val) }
	_, err := buf.Write(bits)
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
func GetI8List(buf *bytes.Buffer) ([]int8, error) { return getList[int8, []int8](buf, GetI8) }
func SetI8List(buf *bytes.Buffer, v []int8) error { return setList(buf, v, SetI8) }
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
func GetU8List(buf *bytes.Buffer) ([]uint8, error) { return getList[uint8, []uint8](buf, GetU8) }
func SetU8List(buf *bytes.Buffer, v []uint8) error { return setList(buf, v, SetU8) }
func EqU8List(a, b []uint8) bool { return slices.Equal(a, b) }
type I16 int16
func GetI16(buf *bytes.Buffer) (int16, error) {
	var b [2]byte
	if _, err := buf.Read(b[:]); err != nil { return 0, err }
	return int16(binary.LittleEndian.Uint16(b[:])), nil
}
func SetI16(buf *bytes.Buffer, v int16) error {
	var b [2]byte
	binary.LittleEndian.PutUint16(b[:], uint16(v))
	_, err := buf.Write(b[:])
	return err
}
func EqI16(a, b int16) bool { return a == b }

type I16List []int16
func GetI16List(buf *bytes.Buffer) ([]int16, error) { return getList[int16, []int16](buf, GetI16) }
func SetI16List(buf *bytes.Buffer, v []int16) error { return setList(buf, v, SetI16) }
func EqI16List(a, b []int16) bool { return slices.Equal(a, b) }
type U16 uint16
func GetU16(buf *bytes.Buffer) (uint16, error) {
	var b [2]byte
	if _, err := buf.Read(b[:]); err != nil { return 0, err }
	return uint16(binary.LittleEndian.Uint16(b[:])), nil
}
func SetU16(buf *bytes.Buffer, v uint16) error {
	var b [2]byte
	binary.LittleEndian.PutUint16(b[:], uint16(v))
	_, err := buf.Write(b[:])
	return err
}
func EqU16(a, b uint16) bool { return a == b }

type U16List []uint16
func GetU16List(buf *bytes.Buffer) ([]uint16, error) { return getList[uint16, []uint16](buf, GetU16) }
func SetU16List(buf *bytes.Buffer, v []uint16) error { return setList(buf, v, SetU16) }
func EqU16List(a, b []uint16) bool { return slices.Equal(a, b) }
type I32 int32
func GetI32(buf *bytes.Buffer) (int32, error) {
	var b [4]byte
	if _, err := buf.Read(b[:]); err != nil { return 0, err }
	return int32(binary.LittleEndian.Uint32(b[:])), nil
}
func SetI32(buf *bytes.Buffer, v int32) error {
	var b [4]byte
	binary.LittleEndian.PutUint32(b[:], uint32(v))
	_, err := buf.Write(b[:])
	return err
}
func EqI32(a, b int32) bool { return a == b }

type I32List []int32
func GetI32List(buf *bytes.Buffer) ([]int32, error) { return getList[int32, []int32](buf, GetI32) }
func SetI32List(buf *bytes.Buffer, v []int32) error { return setList(buf, v, SetI32) }
func EqI32List(a, b []int32) bool { return slices.Equal(a, b) }
type U32 uint32
func GetU32(buf *bytes.Buffer) (uint32, error) {
	var b [4]byte
	if _, err := buf.Read(b[:]); err != nil { return 0, err }
	return uint32(binary.LittleEndian.Uint32(b[:])), nil
}
func SetU32(buf *bytes.Buffer, v uint32) error {
	var b [4]byte
	binary.LittleEndian.PutUint32(b[:], uint32(v))
	_, err := buf.Write(b[:])
	return err
}
func EqU32(a, b uint32) bool { return a == b }

type U32List []uint32
func GetU32List(buf *bytes.Buffer) ([]uint32, error) { return getList[uint32, []uint32](buf, GetU32) }
func SetU32List(buf *bytes.Buffer, v []uint32) error { return setList(buf, v, SetU32) }
func EqU32List(a, b []uint32) bool { return slices.Equal(a, b) }
type I64 int64
func GetI64(buf *bytes.Buffer) (int64, error) {
	var b [8]byte
	if _, err := buf.Read(b[:]); err != nil { return 0, err }
	return int64(binary.LittleEndian.Uint64(b[:])), nil
}
func SetI64(buf *bytes.Buffer, v int64) error {
	var b [8]byte
	binary.LittleEndian.PutUint64(b[:], uint64(v))
	_, err := buf.Write(b[:])
	return err
}
func EqI64(a, b int64) bool { return a == b }

type I64List []int64
func GetI64List(buf *bytes.Buffer) ([]int64, error) { return getList[int64, []int64](buf, GetI64) }
func SetI64List(buf *bytes.Buffer, v []int64) error { return setList(buf, v, SetI64) }
func EqI64List(a, b []int64) bool { return slices.Equal(a, b) }
type U64 uint64
func GetU64(buf *bytes.Buffer) (uint64, error) {
	var b [8]byte
	if _, err := buf.Read(b[:]); err != nil { return 0, err }
	return uint64(binary.LittleEndian.Uint64(b[:])), nil
}
func SetU64(buf *bytes.Buffer, v uint64) error {
	var b [8]byte
	binary.LittleEndian.PutUint64(b[:], uint64(v))
	_, err := buf.Write(b[:])
	return err
}
func EqU64(a, b uint64) bool { return a == b }

type U64List []uint64
func GetU64List(buf *bytes.Buffer) ([]uint64, error) { return getList[uint64, []uint64](buf, GetU64) }
func SetU64List(buf *bytes.Buffer, v []uint64) error { return setList(buf, v, SetU64) }
func EqU64List(a, b []uint64) bool { return slices.Equal(a, b) }
type F32 float32
func GetF32(buf *bytes.Buffer) (float32, error) {
	var b [4]byte
	if _, err := buf.Read(b[:]); err != nil { return 0, err }
	return math.Float32frombits(binary.LittleEndian.Uint32(b[:])), nil
}
func SetF32(buf *bytes.Buffer, v float32) error {
	var b [4]byte
	binary.LittleEndian.PutUint32(b[:], math.Float32bits(v))
	_, err := buf.Write(b[:])
	return err
}
func EqF32(a, b float32) bool { return math.Abs(float64(a-b)) < 1e-6 }

type F32List []float32
func GetF32List(buf *bytes.Buffer) ([]float32, error) { return getList[float32, []float32](buf, GetF32) }
func SetF32List(buf *bytes.Buffer, v []float32) error { return setList(buf, v, SetF32) }
func EqF32List(a, b []float32) bool { return slices.Equal(a, b) }
type F64 float64
func GetF64(buf *bytes.Buffer) (float64, error) {
	var b [8]byte
	if _, err := buf.Read(b[:]); err != nil { return 0, err }
	return math.Float64frombits(binary.LittleEndian.Uint64(b[:])), nil
}
func SetF64(buf *bytes.Buffer, v float64) error {
	var b [8]byte
	binary.LittleEndian.PutUint64(b[:], math.Float64bits(v))
	_, err := buf.Write(b[:])
	return err
}
func EqF64(a, b float64) bool { return math.Abs(float64(a-b)) < 1e-9 }

type F64List []float64
func GetF64List(buf *bytes.Buffer) ([]float64, error) { return getList[float64, []float64](buf, GetF64) }
func SetF64List(buf *bytes.Buffer, v []float64) error { return setList(buf, v, SetF64) }
func EqF64List(a, b []float64) bool { return slices.Equal(a, b) }


// Bin
type Bin []byte
func GetBin(buf *bytes.Buffer) ([]byte, error) {
	l, err := GetU16(buf)
	if err != nil { return nil, err }
	if buf.Len() < int(l) { return nil, fmt.Errorf("not enough data") }
	res := make([]byte, l)
	copy(res, buf.Next(int(l)))
	return res, nil
}
func SetBin(buf *bytes.Buffer, v []byte) error {
	if len(v) > 65535 { return fmt.Errorf("bin length exceeds uint16 max") }
	if err := SetU16(buf, uint16(len(v))); err != nil { return err }
	_, err := buf.Write(v)
	return err
}
func EqBin(a, b []byte) bool { return bytes.Equal(a, b) }

type BinList [][]byte
func GetBinList(buf *bytes.Buffer) ([][]byte, error) { return getList[[]byte, [][]byte](buf, GetBin) }
func SetBinList(buf *bytes.Buffer, v [][]byte) error { return setList(buf, v, SetBin) }
func EqBinList(a, b [][]byte) bool { return slices.EqualFunc(a, b, bytes.Equal) }

// Text
type Text string
func GetText(buf *bytes.Buffer) (string, error) { b, err := GetBin(buf); return string(b), err }
func SetText(buf *bytes.Buffer, v string) error { return SetBin(buf, []byte(v)) }
func EqText(a, b string) bool { return a == b }

type TextList []string
func GetTextList(buf *bytes.Buffer) ([]string, error) { return getList[string, []string](buf, GetText) }
func SetTextList(buf *bytes.Buffer, v []string) error { return setList(buf, v, SetText) }
func EqTextList(a, b []string) bool { return slices.Equal(a, b) }
