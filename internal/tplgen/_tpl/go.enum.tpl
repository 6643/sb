package {{.Package}}

import (
	"bytes"
	"slices"
	"unsafe"
)

{{range .Enums}}
{{$enumName := .Name | PascalCase}}
{{- if .Note}}// {{$enumName}} {{.Note}}{{end}}
type {{$enumName}} uint8

const (
{{- range .Children}}
	{{$enumName}}{{.Name | PascalCase}} {{$enumName}} = {{.ID}} {{if .Note}}// {{.Note}}{{end}}
{{- end}}
)

type {{$enumName}}List []{{$enumName}}
func (v {{$enumName}}List) Set(buf *bytes.Buffer) error { return SetU8List(buf, *(*[]uint8)(unsafe.Pointer(&v))) }
func (v *{{$enumName}}List) Get(buf *bytes.Buffer) error {
	val, err := GetU8List(buf)
	if err == nil { *v = *(*{{$enumName}}List)(unsafe.Pointer(&val)) }
	return err
}
func (v {{$enumName}}List) Eq(other {{$enumName}}List) bool { return slices.Equal(v, other) }
{{end}}