package sb

import (
	"bytes"
	"fmt"
	"math"
	"slices"
)

type {{.Name | PascalCase}} struct {
	{{- range .Fields}}
	{{.Name | PascalCase}} {{GoLogicType .Type}} {{GoTag .}} {{if .Note}}// {{.Note}}{{end}}
	{{- end}}
}

func Get{{.Name | PascalCase}}(buf *bytes.Buffer, s *{{.Name | PascalCase}}) error {
	if s == nil { return nil }
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
		val, err := Get{{.Type.Name | PascalCase}}List(buf)
		if err != nil { return fmt.Errorf("Get{{$.Name | PascalCase}} {{.Name | PascalCase}}: %w", err) }
		s.{{$field.Name | PascalCase}} = val
		{{- else}}
		val, err := Get{{.Type.Name | PascalCase}}(buf)
		if err != nil { return fmt.Errorf("Get{{$.Name | PascalCase}} {{.Name | PascalCase}}: %w", err) }
		s.{{$field.Name | PascalCase}} = val
		if !Is{{.Type.Name | PascalCase}}(s.{{$field.Name | PascalCase}}) { return fmt.Errorf("Get{{$.Name | PascalCase}} {{.Name | PascalCase}}: 非法枚举值: %d", s.{{$field.Name | PascalCase}}) }
		{{- end}}
		{{- else}}
		{{- if .Type.IsList}}
		val, err := Get{{.Type.Name | PascalCase}}List(buf)
		if err != nil { return fmt.Errorf("Get{{$.Name | PascalCase}} {{.Name | PascalCase}}: %w", err) }
		s.{{$field.Name | PascalCase}} = val
		{{- else}}
		val, err := Read{{.Type.Name | PascalCase}}(buf)
		if err != nil { return fmt.Errorf("Get{{$.Name | PascalCase}} {{.Name | PascalCase}}: %w", err) }
		s.{{$field.Name | PascalCase}} = val
		{{- end}}
		{{- end}}
		{{- end}}
	}
	{{- end}}
	{{- end}}
	if err := Validate{{$.Name | PascalCase}}(s); err != nil { return fmt.Errorf("Validate{{$.Name | PascalCase}}: %w", err) }
	return nil
}

func Set{{.Name | PascalCase}}(buf *bytes.Buffer, s *{{.Name | PascalCase}}) error {
	if s == nil { return fmt.Errorf("Set{{$.Name | PascalCase}}: nil value") }
	if err := Validate{{$.Name | PascalCase}}(s); err != nil { return fmt.Errorf("Validate{{$.Name | PascalCase}}: %w", err) }
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
		if err := Set{{.Type.Name | PascalCase}}List(body, s.{{$field.Name | PascalCase}}); err != nil { return fmt.Errorf("Set{{$.Name | PascalCase}} {{.Name | PascalCase}}: %w", err) }
		SetBit(bits, uint8({{$i}}), true)
	}
	{{- else}}
	if s.{{$field.Name | PascalCase}} != 0 {
		if err := Set{{.Type.Name | PascalCase}}(body, s.{{$field.Name | PascalCase}}); err != nil { return fmt.Errorf("Set{{$.Name | PascalCase}} {{.Name | PascalCase}}: %w", err) }
		SetBit(bits, uint8({{$i}}), true)
	}
	{{- end}}
	{{- else}}
	{{- if .Type.IsList}}
	if len(s.{{$field.Name | PascalCase}}) > 0 {
		if err := Set{{.Type.Name | PascalCase}}List(body, ({{.Type.Name | PascalCase}}List)(s.{{$field.Name | PascalCase}})); err != nil { return fmt.Errorf("Set{{$.Name | PascalCase}} {{.Name | PascalCase}}: %w", err) }
		SetBit(bits, uint8({{$i}}), true)
	}
	{{- else}}
	if s.{{$field.Name | PascalCase}} != nil {
		if err := Set{{.Type.Name | PascalCase}}(body, s.{{$field.Name | PascalCase}}); err != nil { return fmt.Errorf("Set{{$.Name | PascalCase}} {{.Name | PascalCase}}: %w", err) }
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

func Validate{{.Name | PascalCase}}(s *{{.Name | PascalCase}}) error {
	if s == nil { return nil }
	{{- range .Fields}}
	{{- if IsEnum .Type}}
	{{- if .Type.IsList}}
	for i, item := range s.{{.Name | PascalCase}} {
		if !Is{{.Type.Name | PascalCase}}(item) { return fmt.Errorf("{{.Name | PascalCase}}[%d] 非法枚举值: %d", i, item) }
	}
	{{- else}}
	if s.{{.Name | PascalCase}} != 0 && !Is{{.Type.Name | PascalCase}}(s.{{.Name | PascalCase}}) { return fmt.Errorf("{{.Name | PascalCase}} 非法枚举值: %d", s.{{.Name | PascalCase}}) }
	{{- end}}
	{{- else if IsStruct .Type}}
	{{- if .Type.IsList}}
	if err := Validate{{.Type.Name | PascalCase}}List(s.{{.Name | PascalCase}}); err != nil { return fmt.Errorf("{{.Name | PascalCase}}: %w", err) }
	{{- else}}
	if s.{{.Name | PascalCase}} != nil {
		if err := Validate{{.Type.Name | PascalCase}}(s.{{.Name | PascalCase}}); err != nil { return fmt.Errorf("{{.Name | PascalCase}}: %w", err) }
	}
	{{- end}}
	{{- end}}
	{{- end}}
	return nil
}

func Eq{{.Name | PascalCase}}(a, b *{{.Name | PascalCase}}) bool {
	if a == b { return true }
	if a == nil || b == nil { return false }
	{{- range .Fields}}
	{{- if IsBaseType .Type}}
	if !Eq{{.Type.Name | PascalCase}}{{if .Type.IsList}}List{{end}}(a.{{.Name | PascalCase}}, b.{{.Name | PascalCase}}) { return false }
	{{- else}}
	{{- if IsEnum .Type}}
	{{- if .Type.IsList}}
	if !Eq{{.Type.Name | PascalCase}}List(a.{{.Name | PascalCase}}, b.{{.Name | PascalCase}}) { return false }
	{{- else}}
	if a.{{.Name | PascalCase}} != b.{{.Name | PascalCase}} { return false }
	{{- end}}
	{{- else}}
	{{- if .Type.IsList}}
	if !Eq{{.Type.Name | PascalCase}}List(({{.Type.Name | PascalCase}}List)(a.{{.Name | PascalCase}}), ({{.Type.Name | PascalCase}}List)(b.{{.Name | PascalCase}})) { return false }
	{{- else}}
	if !Eq{{.Type.Name | PascalCase}}(a.{{.Name | PascalCase}}, b.{{.Name | PascalCase}}) { return false }
	{{- end}}
	{{- end}}
	{{- end}}
	{{- end}}
	return true
}

// Standalone functions
func Read{{.Name | PascalCase}}(buf *bytes.Buffer) (*{{.Name | PascalCase}}, error) {
	s := new({{.Name | PascalCase}})
	return s, Get{{.Name | PascalCase}}(buf, s)
}

type {{.Name | PascalCase}}List []*{{.Name | PascalCase}}
func Get{{.Name | PascalCase}}List(buf *bytes.Buffer) ({{.Name | PascalCase}}List, error) { return getList[*{{.Name | PascalCase}}, {{.Name | PascalCase}}List](buf, Read{{.Name | PascalCase}}) }
func Set{{.Name | PascalCase}}List(buf *bytes.Buffer, v {{.Name | PascalCase}}List) error { return setList(buf, v, Set{{.Name | PascalCase}}) }
func Validate{{.Name | PascalCase}}List(v {{.Name | PascalCase}}List) error {
	for i, item := range v {
		if item == nil { return fmt.Errorf("{{.Name | PascalCase}}List[%d]: nil item", i) }
		if err := Validate{{.Name | PascalCase}}(item); err != nil { return fmt.Errorf("{{.Name | PascalCase}}List[%d]: %w", i, err) }
	}
	return nil
}
func Eq{{.Name | PascalCase}}List(a, b {{.Name | PascalCase}}List) bool { return slices.EqualFunc(a, b, Eq{{.Name | PascalCase}}) }
