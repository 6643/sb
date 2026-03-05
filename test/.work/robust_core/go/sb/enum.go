package sb
import (
	"bytes"
	"slices"
	"unsafe"
)



type Color uint8

const (
	ColorRed Color = 1 
	ColorGreen Color = 2 
	ColorBlue Color = 3 
)

func IsColor(v Color) bool {
	switch v {
	case ColorRed, ColorGreen, ColorBlue:
		return true
	default:
		return false
	}
}

func GetColor(buf *bytes.Buffer) (Color, error) {
	val, err := GetU8(buf)
	return Color(val), err
}

func SetColor(buf *bytes.Buffer, v Color) error { return SetU8(buf, uint8(v)) }

type ColorList []Color
func GetColorList(buf *bytes.Buffer) (ColorList, error) {
	val, err := GetU8List(buf)
	if err != nil { return nil, err }
	return *(*ColorList)(unsafe.Pointer(&val)), nil
}
func SetColorList(buf *bytes.Buffer, v ColorList) error { return SetU8List(buf, *(*[]uint8)(unsafe.Pointer(&v))) }
func IsColorList(v ColorList) bool {
	for _, item := range v {
		if !IsColor(item) { return false }
	}
	return true
}
func EqColorList(a, b ColorList) bool { return slices.Equal(a, b) }

