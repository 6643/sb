package sb

import (
	"bytes"
	"fmt"
	"math"
	"slices"
)

type Matrix struct {
	Single *Child  
	Many []*Child  
	Kind Kind  
	Kinds []Kind  
}

func GetMatrix(buf *bytes.Buffer, s *Matrix) error {
	if s == nil { return nil }
	bitSize := int(math.Ceil(float64(4) / 8.0))
	if buf.Len() < bitSize { return fmt.Errorf("GetMatrix bitmask: %d - %d", buf.Len(), bitSize) }
	bits := buf.Next(bitSize)
	if GetBit(bits, uint8(0)) {
		val, err := ReadChild(buf)
		if err != nil { return fmt.Errorf("GetMatrix Single: %w", err) }
		s.Single = val
	}
	if GetBit(bits, uint8(1)) {
		val, err := GetChildList(buf)
		if err != nil { return fmt.Errorf("GetMatrix Many: %w", err) }
		s.Many = val
	}
	if GetBit(bits, uint8(2)) {
		val, err := GetKind(buf)
		if err != nil { return fmt.Errorf("GetMatrix Kind: %w", err) }
		s.Kind = val
		if !IsKind(s.Kind) { return fmt.Errorf("GetMatrix Kind: 非法枚举值: %d", s.Kind) }
	}
	if GetBit(bits, uint8(3)) {
		val, err := GetKindList(buf)
		if err != nil { return fmt.Errorf("GetMatrix Kinds: %w", err) }
		s.Kinds = val
	}
	if err := ValidateMatrix(s); err != nil { return fmt.Errorf("ValidateMatrix: %w", err) }
	return nil
}

func SetMatrix(buf *bytes.Buffer, s *Matrix) error {
	if s == nil { return fmt.Errorf("SetMatrix: nil value") }
	if err := ValidateMatrix(s); err != nil { return fmt.Errorf("ValidateMatrix: %w", err) }
	bits := make([]byte, uint8(math.Ceil(float64(4)/8.0)))
	body := bytes.NewBuffer(nil)
	if s.Single != nil {
		if err := SetChild(body, s.Single); err != nil { return fmt.Errorf("SetMatrix Single: %w", err) }
		SetBit(bits, uint8(0), true)
	}
	if len(s.Many) > 0 {
		if err := SetChildList(body, (ChildList)(s.Many)); err != nil { return fmt.Errorf("SetMatrix Many: %w", err) }
		SetBit(bits, uint8(1), true)
	}
	if s.Kind != 0 {
		if err := SetKind(body, s.Kind); err != nil { return fmt.Errorf("SetMatrix Kind: %w", err) }
		SetBit(bits, uint8(2), true)
	}
	if len(s.Kinds) > 0 {
		if err := SetKindList(body, s.Kinds); err != nil { return fmt.Errorf("SetMatrix Kinds: %w", err) }
		SetBit(bits, uint8(3), true)
	}

	if _, err := buf.Write(bits); err != nil { return fmt.Errorf("SetMatrix write bitmask: %w", err) }
	_, err := body.WriteTo(buf); return err
}

func ValidateMatrix(s *Matrix) error {
	if s == nil { return nil }
	if s.Single != nil {
		if err := ValidateChild(s.Single); err != nil { return fmt.Errorf("Single: %w", err) }
	}
	if err := ValidateChildList(s.Many); err != nil { return fmt.Errorf("Many: %w", err) }
	if s.Kind != 0 && !IsKind(s.Kind) { return fmt.Errorf("Kind 非法枚举值: %d", s.Kind) }
	for i, item := range s.Kinds {
		if !IsKind(item) { return fmt.Errorf("Kinds[%d] 非法枚举值: %d", i, item) }
	}
	return nil
}

func EqMatrix(a, b *Matrix) bool {
	if a == b { return true }
	if a == nil || b == nil { return false }
	if !EqChild(a.Single, b.Single) { return false }
	if !EqChildList((ChildList)(a.Many), (ChildList)(b.Many)) { return false }
	if a.Kind != b.Kind { return false }
	if !EqKindList(a.Kinds, b.Kinds) { return false }
	return true
}

// Standalone functions
func ReadMatrix(buf *bytes.Buffer) (*Matrix, error) {
	s := new(Matrix)
	return s, GetMatrix(buf, s)
}

type MatrixList []*Matrix
func GetMatrixList(buf *bytes.Buffer) (MatrixList, error) { return getList[*Matrix, MatrixList](buf, ReadMatrix) }
func SetMatrixList(buf *bytes.Buffer, v MatrixList) error { return setList(buf, v, SetMatrix) }
func ValidateMatrixList(v MatrixList) error {
	for i, item := range v {
		if item == nil { return fmt.Errorf("MatrixList[%d]: nil item", i) }
		if err := ValidateMatrix(item); err != nil { return fmt.Errorf("MatrixList[%d]: %w", i, err) }
	}
	return nil
}
func EqMatrixList(a, b MatrixList) bool { return slices.EqualFunc(a, b, EqMatrix) }
