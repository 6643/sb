package sb

import (
	"bytes"
	"fmt"
	"math"
	"slices"
)

type BlobWrap struct {
	TextValue string  
	BinValue []byte  
	Nums []uint8  
	Level Level  
}

func GetBlobWrap(buf *bytes.Buffer, s *BlobWrap) error {
	if s == nil { return nil }
	bitSize := int(math.Ceil(float64(4) / 8.0))
	if buf.Len() < bitSize { return fmt.Errorf("GetBlobWrap bitmask: %d - %d", buf.Len(), bitSize) }
	bits := buf.Next(bitSize)
	if GetBit(bits, uint8(0)) {
		val, err := GetText(buf)
		if err != nil { return fmt.Errorf("GetBlobWrap TextValue: %w", err) }
		s.TextValue = val
	}
	if GetBit(bits, uint8(1)) {
		val, err := GetBin(buf)
		if err != nil { return fmt.Errorf("GetBlobWrap BinValue: %w", err) }
		s.BinValue = val
	}
	if GetBit(bits, uint8(2)) {
		val, err := GetU8List(buf)
		if err != nil { return fmt.Errorf("GetBlobWrap Nums: %w", err) }
		s.Nums = val
	}
	if GetBit(bits, uint8(3)) {
		val, err := GetLevel(buf)
		if err != nil { return fmt.Errorf("GetBlobWrap Level: %w", err) }
		s.Level = val
		if !IsLevel(s.Level) { return fmt.Errorf("GetBlobWrap Level: 非法枚举值: %d", s.Level) }
	}
	if err := ValidateBlobWrap(s); err != nil { return fmt.Errorf("ValidateBlobWrap: %w", err) }
	return nil
}

func SetBlobWrap(buf *bytes.Buffer, s *BlobWrap) error {
	if s == nil { return fmt.Errorf("SetBlobWrap: nil value") }
	if err := ValidateBlobWrap(s); err != nil { return fmt.Errorf("ValidateBlobWrap: %w", err) }
	bits := make([]byte, uint8(math.Ceil(float64(4)/8.0)))
	body := bytes.NewBuffer(nil)
	if s.TextValue != "" {
		if err := SetText(body, s.TextValue); err != nil { return fmt.Errorf("SetBlobWrap TextValue: %w", err) }
		SetBit(bits, uint8(0), true)
	}
	if s.BinValue != nil {
		if err := SetBin(body, s.BinValue); err != nil { return fmt.Errorf("SetBlobWrap BinValue: %w", err) }
		SetBit(bits, uint8(1), true)
	}
	if len(s.Nums) > 0 {
		if err := SetU8List(body, s.Nums); err != nil { return fmt.Errorf("SetBlobWrap Nums: %w", err) }
		SetBit(bits, uint8(2), true)
	}
	if s.Level != 0 {
		if err := SetLevel(body, s.Level); err != nil { return fmt.Errorf("SetBlobWrap Level: %w", err) }
		SetBit(bits, uint8(3), true)
	}

	if _, err := buf.Write(bits); err != nil { return fmt.Errorf("SetBlobWrap write bitmask: %w", err) }
	_, err := body.WriteTo(buf); return err
}

func ValidateBlobWrap(s *BlobWrap) error {
	if s == nil { return nil }
	if s.Level != 0 && !IsLevel(s.Level) { return fmt.Errorf("Level 非法枚举值: %d", s.Level) }
	return nil
}

func EqBlobWrap(a, b *BlobWrap) bool {
	if a == b { return true }
	if a == nil || b == nil { return false }
	if !EqText(a.TextValue, b.TextValue) { return false }
	if !EqBin(a.BinValue, b.BinValue) { return false }
	if !EqU8List(a.Nums, b.Nums) { return false }
	if a.Level != b.Level { return false }
	return true
}

// Standalone functions
func ReadBlobWrap(buf *bytes.Buffer) (*BlobWrap, error) {
	s := new(BlobWrap)
	return s, GetBlobWrap(buf, s)
}

type BlobWrapList []*BlobWrap
func GetBlobWrapList(buf *bytes.Buffer) (BlobWrapList, error) { return getList[*BlobWrap, BlobWrapList](buf, ReadBlobWrap) }
func SetBlobWrapList(buf *bytes.Buffer, v BlobWrapList) error { return setList(buf, v, SetBlobWrap) }
func ValidateBlobWrapList(v BlobWrapList) error {
	for i, item := range v {
		if item == nil { return fmt.Errorf("BlobWrapList[%d]: nil item", i) }
		if err := ValidateBlobWrap(item); err != nil { return fmt.Errorf("BlobWrapList[%d]: %w", i, err) }
	}
	return nil
}
func EqBlobWrapList(a, b BlobWrapList) bool { return slices.EqualFunc(a, b, EqBlobWrap) }
