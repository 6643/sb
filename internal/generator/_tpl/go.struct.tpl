{{- $needUnsafe := false -}}
{{- range .Fields -}}
	{{- if and (IsEnum .Type) .Type.IsList -}}
		{{- $needUnsafe = true -}}
	{{- end -}}
{{- end -}}
package sb

import (
	"bytes"
	"fmt"
	"math"
	"slices"
	{{if $needUnsafe}}"unsafe"{{end}}
)

type {{.Name | PascalCase}} struct {
	{{- range .Fields}}
	{{.Name | PascalCase}} {{GoLogicType .Type}} {{GoTag .}} {{if .Note}}// {{.Note}}{{end}}
	{{- end}}
}

func (s *{{.Name | PascalCase}}) Get(buf *bytes.Buffer) error {
	if buf.Len() == 0 { return nil }
	bitSize := int(math.Ceil(float64({{len .Fields}}) / 8.0))
	if buf.Len() < bitSize { return fmt.Errorf("Get{{$.Name | PascalCase}} bitmask: %d - %d", buf.Len(), bitSize) }
	bits := buf.Next(bitSize)

	{{- range $i, $field := .Fields}}
	{{- if eq .Type.Name "bool"}}
	s.{{$field.Name | PascalCase}} = GetBit(bits, uint8({{$i}}))
	{{- else}}
	if GetBit(bits, uint8({{$i}})) {
		{{- if IsBaseType .Type}}
		val, err := Get{{.Type.Name | PascalCase}}{{if .Type.IsList}}List{{end}}(buf)
		if err != nil { return fmt.Errorf("Get{{$.Name | PascalCase}} {{.Name | PascalCase}}: %w", err) }
		s.{{$field.Name | PascalCase}} = val
		{{- else}}
		{{- if IsEnum .Type}}
		{{- if .Type.IsList}}
		val, err := GetU8List(buf)
		if err != nil { return fmt.Errorf("Get{{$.Name | PascalCase}} {{.Name | PascalCase}}: %w", err) }
		s.{{$field.Name | PascalCase}} = *(*[]{{.Type.Name | PascalCase}})(unsafe.Pointer(&val))
		{{- else}}
		val, err := GetU8(buf)
		if err != nil { return fmt.Errorf("Get{{$.Name | PascalCase}} {{.Name | PascalCase}}: %w", err) }
		s.{{$field.Name | PascalCase}} = {{.Type.Name | PascalCase}}(val)
		{{- end}}
		{{- else}}
		{{- if .Type.IsList}}
		var val {{.Type.Name | PascalCase}}List
		if err := val.Get(buf); err != nil { return fmt.Errorf("Get{{$.Name | PascalCase}} {{.Name | PascalCase}}: %w", err) }
		s.{{$field.Name | PascalCase}} = val
		{{- else}}
		if s.{{$field.Name | PascalCase}} == nil { s.{{$field.Name | PascalCase}} = new({{.Type.Name | PascalCase}}) }
		if err := s.{{$field.Name | PascalCase}}.Get(buf); err != nil { return fmt.Errorf("Get{{$.Name | PascalCase}} {{.Name | PascalCase}}: %w", err) }
		{{- end}}
		{{- end}}
		{{- end}}
	}
	{{- end}}
	{{- end}}
	return nil
}

func (s *{{.Name | PascalCase}}) Set(buf *bytes.Buffer) error {
	if s == nil { return nil }
	bits := make([]byte, uint8(math.Ceil(float64({{len .Fields}})/8.0)))
	body := bytes.NewBuffer(nil)

	{{- range $i, $field := .Fields}}
	{{- if eq .Type.Name "bool"}}
	SetBit(bits, uint8({{$i}}), s.{{$field.Name | PascalCase}})
	{{- else}}
	{{- if IsBaseType .Type}}
	{{- if .Type.IsList}}
	if len(s.{{$field.Name | PascalCase}}) > 0 {
		if err := Set{{.Type.Name | PascalCase}}List(body, s.{{$field.Name | PascalCase}}); err != nil { return fmt.Errorf("Set{{$.Name | PascalCase}} {{.Name | PascalCase}}: %w", err) }
		SetBit(bits, uint8({{$i}}), true)
	}
	{{- else}}
	if s.{{$field.Name | PascalCase}} != {{GoValue .Type.Name}} {
		if err := Set{{.Type.Name | PascalCase}}(body, s.{{$field.Name | PascalCase}}); err != nil { return fmt.Errorf("Set{{$.Name | PascalCase}} {{.Name | PascalCase}}: %w", err) }
		SetBit(bits, uint8({{$i}}), true)
	}
	{{- end}}
	{{- else}}
	{{- if IsEnum .Type}}
	{{- if .Type.IsList}}
	if len(s.{{$field.Name | PascalCase}}) > 0 {
		if err := SetU8List(body, *(*[]uint8)(unsafe.Pointer(&s.{{$field.Name | PascalCase}}))); err != nil { return fmt.Errorf("Set{{$.Name | PascalCase}} {{.Name | PascalCase}}: %w", err) }
		SetBit(bits, uint8({{$i}}), true)
	}
	{{- else}}
	if s.{{$field.Name | PascalCase}} != 0 {
		if err := SetU8(body, uint8(s.{{$field.Name | PascalCase}})); err != nil { return fmt.Errorf("Set{{$.Name | PascalCase}} {{.Name | PascalCase}}: %w", err) }
		SetBit(bits, uint8({{$i}}), true)
	}
	{{- end}}
	{{- else}}
	{{- if .Type.IsList}}
	if len(s.{{$field.Name | PascalCase}}) > 0 {
		if err := ({{.Type.Name | PascalCase}}List)(s.{{$field.Name | PascalCase}}).Set(body); err != nil { return fmt.Errorf("Set{{$.Name | PascalCase}} {{.Name | PascalCase}}: %w", err) }
		SetBit(bits, uint8({{$i}}), true)
	}
	{{- else}}
	if s.{{$field.Name | PascalCase}} != nil {
		if err := s.{{$field.Name | PascalCase}}.Set(body); err != nil { return fmt.Errorf("Set{{$.Name | PascalCase}} {{.Name | PascalCase}}: %w", err) }
		SetBit(bits, uint8({{$i}}), true)
	}
	{{- end}}
	{{- end}}
	{{- end}}
	{{- end}}
	{{- end}}

	if _, err := buf.Write(bits); err != nil { return fmt.Errorf("Set{{$.Name | PascalCase}} write bitmask: %w", err) }
	_, err := body.WriteTo(buf); return err
}

func (s *{{.Name | PascalCase}}) Eq(other *{{.Name | PascalCase}}) bool {
	if s == other { return true }
	if s == nil || other == nil { return false }
	{{- range .Fields}}
	{{- if IsBaseType .Type}}
	if !Eq{{.Type.Name | PascalCase}}{{if .Type.IsList}}List{{end}}(s.{{.Name | PascalCase}}, other.{{.Name | PascalCase}}) { return false }
	{{- else}}
	{{- if IsEnum .Type}}
	{{- if .Type.IsList}}
	if !slices.Equal(s.{{.Name | PascalCase}}, other.{{.Name | PascalCase}}) { return false }
	{{- else}}
	if s.{{.Name | PascalCase}} != other.{{.Name | PascalCase}} { return false }
	{{- end}}
	{{- else}}
	{{- if .Type.IsList}}
	if !({{.Type.Name | PascalCase}}List)(s.{{.Name | PascalCase}}).Eq(({{.Type.Name | PascalCase}}List)(other.{{.Name | PascalCase}})) { return false }
	{{- else}}
	if !s.{{.Name | PascalCase}}.Eq(other.{{.Name | PascalCase}}) { return false }
	{{- end}}
	{{- end}}
	{{- end}}
	{{- end}}
	return true
}

// Standalone functions for compatibility
func Get{{.Name | PascalCase}}(buf *bytes.Buffer) (*{{.Name | PascalCase}}, error) {
	s := new({{.Name | PascalCase}}); return s, s.Get(buf)
}
func Set{{.Name | PascalCase}}(buf *bytes.Buffer, s *{{.Name | PascalCase}}) error { return s.Set(buf) }
func Eq{{.Name | PascalCase}}(a, b *{{.Name | PascalCase}}) bool { return a.Eq(b) }

type {{.Name | PascalCase}}List []*{{.Name | PascalCase}}
func (v {{.Name | PascalCase}}List) Set(buf *bytes.Buffer) error { return setList(buf, v, Set{{.Name | PascalCase}}) }
func (v *{{.Name | PascalCase}}List) Get(buf *bytes.Buffer) error {
	val, err := getList[*{{.Name | PascalCase}}, {{.Name | PascalCase}}List](buf, Get{{.Name | PascalCase}})
	if err == nil { *v = val }; return err
}
func (v {{.Name | PascalCase}}List) Eq(other {{.Name | PascalCase}}List) bool { return slices.EqualFunc(v, other, Eq{{.Name | PascalCase}}) }
