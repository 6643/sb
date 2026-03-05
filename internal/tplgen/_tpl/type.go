package {{.Package}}

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
{{range .Types}}
{{- if and (ne .Name "Bin") (ne .Name "Text") -}}
type {{.Name}} {{.Go}}
func Get{{.Name}}(buf *bytes.Buffer) ({{.Go}}, error) {
	{{- if eq .Go "uint8" "int8"}}
	b, err := buf.ReadByte()
	return {{.Go}}(b), err
	{{- else}}
	var b [{{if eq .Go "uint16" "int16"}}2{{else if eq .Go "uint32" "int32" "float32"}}4{{else}}8{{end}}]byte
	if _, err := buf.Read(b[:]); err != nil { return 0, err }
	{{- if eq .Go "uint16" "int16"}}
	return {{.Go}}(binary.LittleEndian.Uint16(b[:])), nil
	{{- else if eq .Go "uint32" "int32"}}
	return {{.Go}}(binary.LittleEndian.Uint32(b[:])), nil
	{{- else if eq .Go "uint64" "int64"}}
	return {{.Go}}(binary.LittleEndian.Uint64(b[:])), nil
	{{- else if eq .Go "float32"}}
	return math.Float32frombits(binary.LittleEndian.Uint32(b[:])), nil
	{{- else if eq .Go "float64"}}
	return math.Float64frombits(binary.LittleEndian.Uint64(b[:])), nil
	{{- end}}
	{{- end}}
}
func Set{{.Name}}(buf *bytes.Buffer, v {{.Go}}) error {
	{{- if eq .Go "uint8" "int8"}}
	return buf.WriteByte(uint8(v))
	{{- else}}
	var b [{{if eq .Go "uint16" "int16"}}2{{else if eq .Go "uint32" "int32" "float32"}}4{{else}}8{{end}}]byte
	{{- if eq .Go "uint16" "int16"}}
	binary.LittleEndian.PutUint16(b[:], uint16(v))
	{{- else if eq .Go "uint32" "int32"}}
	binary.LittleEndian.PutUint32(b[:], uint32(v))
	{{- else if eq .Go "uint64" "int64"}}
	binary.LittleEndian.PutUint64(b[:], uint64(v))
	{{- else if eq .Go "float32"}}
	binary.LittleEndian.PutUint32(b[:], math.Float32bits(v))
	{{- else if eq .Go "float64"}}
	binary.LittleEndian.PutUint64(b[:], math.Float64bits(v))
	{{- end}}
	_, err := buf.Write(b[:])
	return err
	{{- end}}
}
func Eq{{.Name}}(a, b {{.Go}}) bool { return {{if .IsFloat}}math.Abs(float64(a-b)) < {{.Eps}}{{else}}a == b{{end}} }

type {{.Name}}List []{{.Go}}
func Get{{.Name}}List(buf *bytes.Buffer) ([]{{.Go}}, error) { return getList[{{.Go}}, []{{.Go}}](buf, Get{{.Name}}) }
func Set{{.Name}}List(buf *bytes.Buffer, v []{{.Go}}) error { return setList(buf, v, Set{{.Name}}) }
func Eq{{.Name}}List(a, b []{{.Go}}) bool { return slices.Equal(a, b) }
{{end}}
{{- end}}

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
