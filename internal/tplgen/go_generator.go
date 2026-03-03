package tplgen

import (
	"math"
	"sb/internal/gen"
	"sb/internal/util"
	"text/template"
)

type GoGenerator struct {
	Config   gen.Config
	FuncMap  template.FuncMap
	tplCache map[string]*template.Template
}

func NewGoGenerator(cfg gen.Config) *GoGenerator {
	g := &GoGenerator{Config: cfg, tplCache: make(map[string]*template.Template)}
	g.FuncMap = template.FuncMap{
		"PascalCase":  util.PascalCase,
		"SnakeCase":   util.SnakeCase,
		"CamelCase":   util.CamelCase,
		"GoType":      g.getGoType,
		"GoLogicType": g.getGoLogicType,
		"GoValue":     g.getGoValue,
		"GoTag":       g.getGoTag,
		"GoRpcType":   g.getGoRpcType,
		"IsBaseType":  func(t Type) bool { return t.Kind == KindBase },
		"IsEnum":      func(t Type) bool { return t.Kind == KindEnum },
		"IsStruct":    func(t Type) bool { return t.Kind == KindStruct },
		"IsList":      func(t Type) bool { return t.IsList },
		"Ceil":        func(n int) int { return int(math.Ceil(float64(n) / 8.0)) },
	}
	return g
}
