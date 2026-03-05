package sb
import (
	"bytes"
	"slices"
	"unsafe"
)



type Level uint8

const (
	LevelLow Level = 1 
	LevelMid Level = 2 
	LevelHigh Level = 3 
)

func IsLevel(v Level) bool {
	switch v {
	case LevelLow, LevelMid, LevelHigh:
		return true
	default:
		return false
	}
}

func GetLevel(buf *bytes.Buffer) (Level, error) {
	val, err := GetU8(buf)
	return Level(val), err
}

func SetLevel(buf *bytes.Buffer, v Level) error { return SetU8(buf, uint8(v)) }

type LevelList []Level
func GetLevelList(buf *bytes.Buffer) (LevelList, error) {
	val, err := GetU8List(buf)
	if err != nil { return nil, err }
	return *(*LevelList)(unsafe.Pointer(&val)), nil
}
func SetLevelList(buf *bytes.Buffer, v LevelList) error { return SetU8List(buf, *(*[]uint8)(unsafe.Pointer(&v))) }
func IsLevelList(v LevelList) bool {
	for _, item := range v {
		if !IsLevel(item) { return false }
	}
	return true
}
func EqLevelList(a, b LevelList) bool { return slices.Equal(a, b) }

