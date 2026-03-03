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

func Is{{$enumName}}(v {{$enumName}}) bool {
	switch v {
	case {{range $i, $child := .Children}}{{if gt $i 0}}, {{end}}{{$enumName}}{{.Name | PascalCase}}{{end}}:
		return true
	default:
		return false
	}
}

func Get{{$enumName}}(buf *bytes.Buffer) ({{$enumName}}, error) {
	val, err := GetU8(buf)
	return {{$enumName}}(val), err
}

func Set{{$enumName}}(buf *bytes.Buffer, v {{$enumName}}) error { return SetU8(buf, uint8(v)) }

type {{$enumName}}List []{{$enumName}}
func Get{{$enumName}}List(buf *bytes.Buffer) ({{$enumName}}List, error) {
	val, err := GetU8List(buf)
	if err != nil { return nil, err }
	return *(*{{$enumName}}List)(unsafe.Pointer(&val)), nil
}
func Set{{$enumName}}List(buf *bytes.Buffer, v {{$enumName}}List) error { return SetU8List(buf, *(*[]uint8)(unsafe.Pointer(&v))) }
func Is{{$enumName}}List(v {{$enumName}}List) bool {
	for _, item := range v {
		if !Is{{$enumName}}(item) { return false }
	}
	return true
}
func Eq{{$enumName}}List(a, b {{$enumName}}List) bool { return slices.Equal(a, b) }
{{end}}
