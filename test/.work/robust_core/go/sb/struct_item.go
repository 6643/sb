package sb

import (
	"bytes"
	"fmt"
	"math"
	"slices"
)

type Item struct {
	Id uint32  
	Color Color  
	Tags []string  
	Active bool  
}

func GetItem(buf *bytes.Buffer, s *Item) error {
	if s == nil { return nil }
	bitSize := int(math.Ceil(float64(4) / 8.0))
	if buf.Len() < bitSize { return fmt.Errorf("GetItem bitmask: %d - %d", buf.Len(), bitSize) }
	bits := buf.Next(bitSize)
	if GetBit(bits, uint8(0)) {
		val, err := GetU32(buf)
		if err != nil { return fmt.Errorf("GetItem Id: %w", err) }
		s.Id = val
	}
	if GetBit(bits, uint8(1)) {
		val, err := GetColor(buf)
		if err != nil { return fmt.Errorf("GetItem Color: %w", err) }
		s.Color = val
		if !IsColor(s.Color) { return fmt.Errorf("GetItem Color: 非法枚举值: %d", s.Color) }
	}
	if GetBit(bits, uint8(2)) {
		val, err := GetTextList(buf)
		if err != nil { return fmt.Errorf("GetItem Tags: %w", err) }
		s.Tags = val
	}
	s.Active = GetBit(bits, uint8(3))
	if err := ValidateItem(s); err != nil { return fmt.Errorf("ValidateItem: %w", err) }
	return nil
}

func SetItem(buf *bytes.Buffer, s *Item) error {
	if s == nil { return fmt.Errorf("SetItem: nil value") }
	if err := ValidateItem(s); err != nil { return fmt.Errorf("ValidateItem: %w", err) }
	bits := make([]byte, uint8(math.Ceil(float64(4)/8.0)))
	body := bytes.NewBuffer(nil)
	if s.Id != 0 {
		if err := SetU32(body, s.Id); err != nil { return fmt.Errorf("SetItem Id: %w", err) }
		SetBit(bits, uint8(0), true)
	}
	if s.Color != 0 {
		if err := SetColor(body, s.Color); err != nil { return fmt.Errorf("SetItem Color: %w", err) }
		SetBit(bits, uint8(1), true)
	}
	if len(s.Tags) > 0 {
		if err := SetTextList(body, s.Tags); err != nil { return fmt.Errorf("SetItem Tags: %w", err) }
		SetBit(bits, uint8(2), true)
	}
	SetBit(bits, uint8(3), s.Active)

	if _, err := buf.Write(bits); err != nil { return fmt.Errorf("SetItem write bitmask: %w", err) }
	_, err := body.WriteTo(buf); return err
}

func ValidateItem(s *Item) error {
	if s == nil { return nil }
	if s.Color != 0 && !IsColor(s.Color) { return fmt.Errorf("Color 非法枚举值: %d", s.Color) }
	return nil
}

func EqItem(a, b *Item) bool {
	if a == b { return true }
	if a == nil || b == nil { return false }
	if !EqU32(a.Id, b.Id) { return false }
	if a.Color != b.Color { return false }
	if !EqTextList(a.Tags, b.Tags) { return false }
	if !EqBool(a.Active, b.Active) { return false }
	return true
}

// Standalone functions
func ReadItem(buf *bytes.Buffer) (*Item, error) {
	s := new(Item)
	return s, GetItem(buf, s)
}

type ItemList []*Item
func GetItemList(buf *bytes.Buffer) (ItemList, error) { return getList[*Item, ItemList](buf, ReadItem) }
func SetItemList(buf *bytes.Buffer, v ItemList) error { return setList(buf, v, SetItem) }
func ValidateItemList(v ItemList) error {
	for i, item := range v {
		if item == nil { return fmt.Errorf("ItemList[%d]: nil item", i) }
		if err := ValidateItem(item); err != nil { return fmt.Errorf("ItemList[%d]: %w", i, err) }
	}
	return nil
}
func EqItemList(a, b ItemList) bool { return slices.EqualFunc(a, b, EqItem) }
