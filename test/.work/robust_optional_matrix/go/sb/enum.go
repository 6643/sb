package sb
import (
	"bytes"
	"slices"
	"unsafe"
)



type Kind uint8

const (
	KindA Kind = 1 
	KindB Kind = 2 
)

func IsKind(v Kind) bool {
	switch v {
	case KindA, KindB:
		return true
	default:
		return false
	}
}

func GetKind(buf *bytes.Buffer) (Kind, error) {
	val, err := GetU8(buf)
	return Kind(val), err
}

func SetKind(buf *bytes.Buffer, v Kind) error { return SetU8(buf, uint8(v)) }

type KindList []Kind
func GetKindList(buf *bytes.Buffer) (KindList, error) {
	val, err := GetU8List(buf)
	if err != nil { return nil, err }
	return *(*KindList)(unsafe.Pointer(&val)), nil
}
func SetKindList(buf *bytes.Buffer, v KindList) error { return SetU8List(buf, *(*[]uint8)(unsafe.Pointer(&v))) }
func IsKindList(v KindList) bool {
	for _, item := range v {
		if !IsKind(item) { return false }
	}
	return true
}
func EqKindList(a, b KindList) bool { return slices.Equal(a, b) }

