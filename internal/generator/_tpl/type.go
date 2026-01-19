package {{.Package}}

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
func get[T any](buf *bytes.Buffer) (T, error) {
	var v T; err := binary.Read(buf, binary.LittleEndian, &v); return v, err
}
func set[T any](buf *bytes.Buffer, v T) error {
	return binary.Write(buf, binary.LittleEndian, v)
}
func getList[T any, L ~[]T](buf *bytes.Buffer, getItem func(*bytes.Buffer) (T, error)) (L, error) {
	count, err := GetU8(buf); if err != nil { return nil, err }
	list := make([]T, count)
	for i := range list { if list[i], err = getItem(buf); err != nil { return nil, err } }
	return L(list), nil
}
func setList[T any](buf *bytes.Buffer, list []T, setItem func(*bytes.Buffer, T) error) error {
	if len(list) > 255 { return fmt.Errorf("list length exceeds uint8 max") }
	if err := SetU8(buf, uint8(len(list))); err != nil { return err }
	for _, item := range list { if err := setItem(buf, item); err != nil { return err } }
	return nil
}

func GetBit(bits []byte, i uint8) bool {
	if int(i/8) >= len(bits) { return false }; return (bits[i/8] & (1 << (i % 8))) != 0
}
func SetBit(bits []byte, i uint8, v bool) {
	if int(i/8) >= len(bits) { return }
	if v { bits[i/8] |= (1 << (i % 8)) } else { bits[i/8] &= ^(1 << (i % 8)) }
}

// Bool
type Bool bool
func (v Bool) Set(buf *bytes.Buffer) error { return SetBool(buf, bool(v)) }
func (v *Bool) Get(buf *bytes.Buffer) error { val, err := GetBool(buf); if err == nil { *v = Bool(val) }; return err }
func GetBool(buf *bytes.Buffer) (bool, error) { b, err := get[uint8](buf); return b == 1, err }
func SetBool(buf *bytes.Buffer, v bool) error { val := uint8(0); if v { val = 1 }; return set(buf, val) }
func EqBool(a, b bool) bool { return a == b }

type BoolList []bool
func (v BoolList) Set(buf *bytes.Buffer) error { return SetBoolList(buf, v) }
func (v *BoolList) Get(buf *bytes.Buffer) error { val, err := GetBoolList(buf); if err == nil { *v = val }; return err }
func GetBoolList(buf *bytes.Buffer) ([]bool, error) {
	count, err := GetU8(buf); if err != nil { return nil, err }
	bits := buf.Next((int(count) + 7) / 8)
	bools := make([]bool, count)
	for i := 0; i < int(count); i++ { bools[i] = GetBit(bits, uint8(i)) }
	return bools, nil
}
func SetBoolList(buf *bytes.Buffer, v []bool) error {
	if err := SetU8(buf, uint8(len(v))); err != nil { return err }
	bits := make([]byte, (len(v)+7)/8)
	for i, val := range v { SetBit(bits, uint8(i), val) }
	_, err := buf.Write(bits); return err
}
func EqBoolList(a, b []bool) bool { return slices.Equal(a, b) }

// Primitives Macro
{{range .Types}}
{{- if and (ne .Name "Bin") (ne .Name "Text") -}}
type {{.Name}} {{.Go}}
func (v {{.Name}}) Set(buf *bytes.Buffer) error { return Set{{.Name}}(buf, {{.Go}}(v)) }
func (v *{{.Name}}) Get(buf *bytes.Buffer) error { val, err := Get{{.Name}}(buf); if err == nil { *v = {{.Name}}(val) }; return err }
func Get{{.Name}}(buf *bytes.Buffer) ({{.Go}}, error) { return get[{{.Go}}](buf) }
func Set{{.Name}}(buf *bytes.Buffer, v {{.Go}}) error { return set(buf, v) }
func Eq{{.Name}}(a, b {{.Go}}) bool { return {{if .IsFloat}}math.Abs(float64(a-b)) < {{.Eps}}{{else}}a == b{{end}} }

type {{.Name}}List []{{.Go}}
func (v {{.Name}}List) Set(buf *bytes.Buffer) error { return Set{{.Name}}List(buf, v) }
func (v *{{.Name}}List) Get(buf *bytes.Buffer) error { val, err := Get{{.Name}}List(buf); if err == nil { *v = val }; return err }
func Get{{.Name}}List(buf *bytes.Buffer) ([]{{.Go}}, error) { return getList[{{.Go}}, []{{.Go}}](buf, Get{{.Name}}) }
func Set{{.Name}}List(buf *bytes.Buffer, v []{{.Go}}) error { return setList(buf, v, Set{{.Name}}) }
func Eq{{.Name}}List(a, b []{{.Go}}) bool { return slices.Equal(a, b) }
{{end}}
{{- end}}

// Bin
type Bin []byte
func (v Bin) Set(buf *bytes.Buffer) error { return SetBin(buf, []byte(v)) }
func (v *Bin) Get(buf *bytes.Buffer) error { val, err := GetBin(buf); if err == nil { *v = Bin(val) }; return err }
func GetBin(buf *bytes.Buffer) ([]byte, error) {
	l, err := GetU16(buf); if err != nil { return nil, err }
	if uint16(buf.Len()) < l { return nil, fmt.Errorf("not enough data") }; return buf.Next(int(l)), nil
}
func SetBin(buf *bytes.Buffer, v []byte) error {
	if err := SetU16(buf, uint16(len(v))); err != nil { return err }; _, err := buf.Write(v); return err
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
