package sb

import (
	"bytes"
	"fmt"
	"math"
	"slices"
)

type Child struct {
	Id uint32  
}

func GetChild(buf *bytes.Buffer, s *Child) error {
	if s == nil { return nil }
	bitSize := int(math.Ceil(float64(1) / 8.0))
	if buf.Len() < bitSize { return fmt.Errorf("GetChild bitmask: %d - %d", buf.Len(), bitSize) }
	bits := buf.Next(bitSize)
	if GetBit(bits, uint8(0)) {
		val, err := GetU32(buf)
		if err != nil { return fmt.Errorf("GetChild Id: %w", err) }
		s.Id = val
	}
	if err := ValidateChild(s); err != nil { return fmt.Errorf("ValidateChild: %w", err) }
	return nil
}

func SetChild(buf *bytes.Buffer, s *Child) error {
	if s == nil { return fmt.Errorf("SetChild: nil value") }
	if err := ValidateChild(s); err != nil { return fmt.Errorf("ValidateChild: %w", err) }
	bits := make([]byte, uint8(math.Ceil(float64(1)/8.0)))
	body := bytes.NewBuffer(nil)
	if s.Id != 0 {
		if err := SetU32(body, s.Id); err != nil { return fmt.Errorf("SetChild Id: %w", err) }
		SetBit(bits, uint8(0), true)
	}

	if _, err := buf.Write(bits); err != nil { return fmt.Errorf("SetChild write bitmask: %w", err) }
	_, err := body.WriteTo(buf); return err
}

func ValidateChild(s *Child) error {
	if s == nil { return nil }
	return nil
}

func EqChild(a, b *Child) bool {
	if a == b { return true }
	if a == nil || b == nil { return false }
	if !EqU32(a.Id, b.Id) { return false }
	return true
}

// Standalone functions
func ReadChild(buf *bytes.Buffer) (*Child, error) {
	s := new(Child)
	return s, GetChild(buf, s)
}

type ChildList []*Child
func GetChildList(buf *bytes.Buffer) (ChildList, error) { return getList[*Child, ChildList](buf, ReadChild) }
func SetChildList(buf *bytes.Buffer, v ChildList) error { return setList(buf, v, SetChild) }
func ValidateChildList(v ChildList) error {
	for i, item := range v {
		if item == nil { return fmt.Errorf("ChildList[%d]: nil item", i) }
		if err := ValidateChild(item); err != nil { return fmt.Errorf("ChildList[%d]: %w", i, err) }
	}
	return nil
}
func EqChildList(a, b ChildList) bool { return slices.EqualFunc(a, b, EqChild) }
