package sb

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"math"
	"slices"
)

type Serializable interface { Set(*bytes.Buffer) error }
type Deserializable interface { Get(*bytes.Buffer) error }

func SetAll(buf *bytes.Buffer, args ...Serializable) error {
	for _, arg := range args { if err := arg.Set(buf); err != nil { return err } }
	return nil
}

func GetAll(buf *bytes.Buffer, args ...Deserializable) error {
	for _, arg := range args { if err := arg.Get(buf); err != nil { return err } }
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
func (v Bool) Set(buf *bytes.Buffer) error { return SetBool(buf, bool(v)) }
func (v *Bool) Get(buf *bytes.Buffer) error { val, err := GetBool(buf); if err == nil { *v = Bool(val) }; return err }
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
func (v BoolList) Set(buf *bytes.Buffer) error { return SetBoolList(buf, v) }
func (v *BoolList) Get(buf *bytes.Buffer) error { val, err := GetBoolList(buf); if err == nil { *v = val }; return err }
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
func (v I8) Set(buf *bytes.Buffer) error { return SetI8(buf, int8(v)) }
func (v *I8) Get(buf *bytes.Buffer) error { val, err := GetI8(buf); if err == nil { *v = I8(val) }; return err }
func GetI8(buf *bytes.Buffer) (int8, error) {
	b, err := buf.ReadByte()
	return int8(b), err
}
func SetI8(buf *bytes.Buffer, v int8) error {
	return buf.WriteByte(uint8(v))
}
func EqI8(a, b int8) bool { return a == b }

type I8List []int8
func (v I8List) Set(buf *bytes.Buffer) error { return SetI8List(buf, v) }
func (v *I8List) Get(buf *bytes.Buffer) error { val, err := GetI8List(buf); if err == nil { *v = val }; return err }
func GetI8List(buf *bytes.Buffer) ([]int8, error) { return getList[int8, []int8](buf, GetI8) }
func SetI8List(buf *bytes.Buffer, v []int8) error { return setList(buf, v, SetI8) }
func EqI8List(a, b []int8) bool { return slices.Equal(a, b) }
type U8 uint8
func (v U8) Set(buf *bytes.Buffer) error { return SetU8(buf, uint8(v)) }
func (v *U8) Get(buf *bytes.Buffer) error { val, err := GetU8(buf); if err == nil { *v = U8(val) }; return err }
func GetU8(buf *bytes.Buffer) (uint8, error) {
	b, err := buf.ReadByte()
	return uint8(b), err
}
func SetU8(buf *bytes.Buffer, v uint8) error {
	return buf.WriteByte(uint8(v))
}
func EqU8(a, b uint8) bool { return a == b }

type U8List []uint8
func (v U8List) Set(buf *bytes.Buffer) error { return SetU8List(buf, v) }
func (v *U8List) Get(buf *bytes.Buffer) error { val, err := GetU8List(buf); if err == nil { *v = val }; return err }
func GetU8List(buf *bytes.Buffer) ([]uint8, error) { return getList[uint8, []uint8](buf, GetU8) }
func SetU8List(buf *bytes.Buffer, v []uint8) error { return setList(buf, v, SetU8) }
func EqU8List(a, b []uint8) bool { return slices.Equal(a, b) }
type I16 int16
func (v I16) Set(buf *bytes.Buffer) error { return SetI16(buf, int16(v)) }
func (v *I16) Get(buf *bytes.Buffer) error { val, err := GetI16(buf); if err == nil { *v = I16(val) }; return err }
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
func (v I16List) Set(buf *bytes.Buffer) error { return SetI16List(buf, v) }
func (v *I16List) Get(buf *bytes.Buffer) error { val, err := GetI16List(buf); if err == nil { *v = val }; return err }
func GetI16List(buf *bytes.Buffer) ([]int16, error) { return getList[int16, []int16](buf, GetI16) }
func SetI16List(buf *bytes.Buffer, v []int16) error { return setList(buf, v, SetI16) }
func EqI16List(a, b []int16) bool { return slices.Equal(a, b) }
type U16 uint16
func (v U16) Set(buf *bytes.Buffer) error { return SetU16(buf, uint16(v)) }
func (v *U16) Get(buf *bytes.Buffer) error { val, err := GetU16(buf); if err == nil { *v = U16(val) }; return err }
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
func (v U16List) Set(buf *bytes.Buffer) error { return SetU16List(buf, v) }
func (v *U16List) Get(buf *bytes.Buffer) error { val, err := GetU16List(buf); if err == nil { *v = val }; return err }
func GetU16List(buf *bytes.Buffer) ([]uint16, error) { return getList[uint16, []uint16](buf, GetU16) }
func SetU16List(buf *bytes.Buffer, v []uint16) error { return setList(buf, v, SetU16) }
func EqU16List(a, b []uint16) bool { return slices.Equal(a, b) }
type I32 int32
func (v I32) Set(buf *bytes.Buffer) error { return SetI32(buf, int32(v)) }
func (v *I32) Get(buf *bytes.Buffer) error { val, err := GetI32(buf); if err == nil { *v = I32(val) }; return err }
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
func (v I32List) Set(buf *bytes.Buffer) error { return SetI32List(buf, v) }
func (v *I32List) Get(buf *bytes.Buffer) error { val, err := GetI32List(buf); if err == nil { *v = val }; return err }
func GetI32List(buf *bytes.Buffer) ([]int32, error) { return getList[int32, []int32](buf, GetI32) }
func SetI32List(buf *bytes.Buffer, v []int32) error { return setList(buf, v, SetI32) }
func EqI32List(a, b []int32) bool { return slices.Equal(a, b) }
type U32 uint32
func (v U32) Set(buf *bytes.Buffer) error { return SetU32(buf, uint32(v)) }
func (v *U32) Get(buf *bytes.Buffer) error { val, err := GetU32(buf); if err == nil { *v = U32(val) }; return err }
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
func (v U32List) Set(buf *bytes.Buffer) error { return SetU32List(buf, v) }
func (v *U32List) Get(buf *bytes.Buffer) error { val, err := GetU32List(buf); if err == nil { *v = val }; return err }
func GetU32List(buf *bytes.Buffer) ([]uint32, error) { return getList[uint32, []uint32](buf, GetU32) }
func SetU32List(buf *bytes.Buffer, v []uint32) error { return setList(buf, v, SetU32) }
func EqU32List(a, b []uint32) bool { return slices.Equal(a, b) }
type I64 int64
func (v I64) Set(buf *bytes.Buffer) error { return SetI64(buf, int64(v)) }
func (v *I64) Get(buf *bytes.Buffer) error { val, err := GetI64(buf); if err == nil { *v = I64(val) }; return err }
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
func (v I64List) Set(buf *bytes.Buffer) error { return SetI64List(buf, v) }
func (v *I64List) Get(buf *bytes.Buffer) error { val, err := GetI64List(buf); if err == nil { *v = val }; return err }
func GetI64List(buf *bytes.Buffer) ([]int64, error) { return getList[int64, []int64](buf, GetI64) }
func SetI64List(buf *bytes.Buffer, v []int64) error { return setList(buf, v, SetI64) }
func EqI64List(a, b []int64) bool { return slices.Equal(a, b) }
type U64 uint64
func (v U64) Set(buf *bytes.Buffer) error { return SetU64(buf, uint64(v)) }
func (v *U64) Get(buf *bytes.Buffer) error { val, err := GetU64(buf); if err == nil { *v = U64(val) }; return err }
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
func (v U64List) Set(buf *bytes.Buffer) error { return SetU64List(buf, v) }
func (v *U64List) Get(buf *bytes.Buffer) error { val, err := GetU64List(buf); if err == nil { *v = val }; return err }
func GetU64List(buf *bytes.Buffer) ([]uint64, error) { return getList[uint64, []uint64](buf, GetU64) }
func SetU64List(buf *bytes.Buffer, v []uint64) error { return setList(buf, v, SetU64) }
func EqU64List(a, b []uint64) bool { return slices.Equal(a, b) }
type F32 float32
func (v F32) Set(buf *bytes.Buffer) error { return SetF32(buf, float32(v)) }
func (v *F32) Get(buf *bytes.Buffer) error { val, err := GetF32(buf); if err == nil { *v = F32(val) }; return err }
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
func (v F32List) Set(buf *bytes.Buffer) error { return SetF32List(buf, v) }
func (v *F32List) Get(buf *bytes.Buffer) error { val, err := GetF32List(buf); if err == nil { *v = val }; return err }
func GetF32List(buf *bytes.Buffer) ([]float32, error) { return getList[float32, []float32](buf, GetF32) }
func SetF32List(buf *bytes.Buffer, v []float32) error { return setList(buf, v, SetF32) }
func EqF32List(a, b []float32) bool { return slices.Equal(a, b) }
type F64 float64
func (v F64) Set(buf *bytes.Buffer) error { return SetF64(buf, float64(v)) }
func (v *F64) Get(buf *bytes.Buffer) error { val, err := GetF64(buf); if err == nil { *v = F64(val) }; return err }
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
func (v F64List) Set(buf *bytes.Buffer) error { return SetF64List(buf, v) }
func (v *F64List) Get(buf *bytes.Buffer) error { val, err := GetF64List(buf); if err == nil { *v = val }; return err }
func GetF64List(buf *bytes.Buffer) ([]float64, error) { return getList[float64, []float64](buf, GetF64) }
func SetF64List(buf *bytes.Buffer, v []float64) error { return setList(buf, v, SetF64) }
func EqF64List(a, b []float64) bool { return slices.Equal(a, b) }


// Bin
type Bin []byte
func (v Bin) Set(buf *bytes.Buffer) error { return SetBin(buf, []byte(v)) }
func (v *Bin) Get(buf *bytes.Buffer) error { val, err := GetBin(buf); if err == nil { *v = Bin(val) }; return err }
func GetBin(buf *bytes.Buffer) ([]byte, error) {
	l, err := GetU16(buf)
	if err != nil { return nil, err }
	if uint16(buf.Len()) < l { return nil, fmt.Errorf("not enough data") }
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
func (v BinList) Set(buf *bytes.Buffer) error { return SetBinList(buf, v) }
func (v *BinList) Get(buf *bytes.Buffer) error { val, err := GetBinList(buf); if err == nil { *v = val }; return err }
func GetBinList(buf *bytes.Buffer) ([][]byte, error) { return getList[[]byte, [][]byte](buf, GetBin) }
func SetBinList(buf *bytes.Buffer, v [][]byte) error { return setList(buf, v, SetBin) }
func EqBinList(a, b [][]byte) bool { return slices.EqualFunc(a, b, bytes.Equal) }

// Text
type Text string
func (v Text) Set(buf *bytes.Buffer) error { return SetText(buf, string(v)) }
func (v *Text) Get(buf *bytes.Buffer) error { val, err := GetText(buf); if err == nil { *v = Text(val) }; return err }
func GetText(buf *bytes.Buffer) (string, error) { b, err := GetBin(buf); return string(b), err }
func SetText(buf *bytes.Buffer, v string) error { return SetBin(buf, []byte(v)) }
func EqText(a, b string) bool { return a == b }

type TextList []string
func (v TextList) Set(buf *bytes.Buffer) error { return SetTextList(buf, v) }
func (v *TextList) Get(buf *bytes.Buffer) error { val, err := GetTextList(buf); if err == nil { *v = val }; return err }
func GetTextList(buf *bytes.Buffer) ([]string, error) { return getList[string, []string](buf, GetText) }
func SetTextList(buf *bytes.Buffer, v []string) error { return setList(buf, v, SetText) }
func EqTextList(a, b []string) bool { return slices.Equal(a, b) }
