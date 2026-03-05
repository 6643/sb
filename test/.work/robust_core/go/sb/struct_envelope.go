package sb

import (
	"bytes"
	"fmt"
	"math"
	"slices"
)

type Envelope struct {
	Item *Item  
	Items []*Item  
	Note string  
}

func GetEnvelope(buf *bytes.Buffer, s *Envelope) error {
	if s == nil { return nil }
	bitSize := int(math.Ceil(float64(3) / 8.0))
	if buf.Len() < bitSize { return fmt.Errorf("GetEnvelope bitmask: %d - %d", buf.Len(), bitSize) }
	bits := buf.Next(bitSize)
	if GetBit(bits, uint8(0)) {
		val, err := ReadItem(buf)
		if err != nil { return fmt.Errorf("GetEnvelope Item: %w", err) }
		s.Item = val
	}
	if GetBit(bits, uint8(1)) {
		val, err := GetItemList(buf)
		if err != nil { return fmt.Errorf("GetEnvelope Items: %w", err) }
		s.Items = val
	}
	if GetBit(bits, uint8(2)) {
		val, err := GetText(buf)
		if err != nil { return fmt.Errorf("GetEnvelope Note: %w", err) }
		s.Note = val
	}
	if err := ValidateEnvelope(s); err != nil { return fmt.Errorf("ValidateEnvelope: %w", err) }
	return nil
}

func SetEnvelope(buf *bytes.Buffer, s *Envelope) error {
	if s == nil { return fmt.Errorf("SetEnvelope: nil value") }
	if err := ValidateEnvelope(s); err != nil { return fmt.Errorf("ValidateEnvelope: %w", err) }
	bits := make([]byte, uint8(math.Ceil(float64(3)/8.0)))
	body := bytes.NewBuffer(nil)
	if s.Item != nil {
		if err := SetItem(body, s.Item); err != nil { return fmt.Errorf("SetEnvelope Item: %w", err) }
		SetBit(bits, uint8(0), true)
	}
	if len(s.Items) > 0 {
		if err := SetItemList(body, (ItemList)(s.Items)); err != nil { return fmt.Errorf("SetEnvelope Items: %w", err) }
		SetBit(bits, uint8(1), true)
	}
	if s.Note != "" {
		if err := SetText(body, s.Note); err != nil { return fmt.Errorf("SetEnvelope Note: %w", err) }
		SetBit(bits, uint8(2), true)
	}

	if _, err := buf.Write(bits); err != nil { return fmt.Errorf("SetEnvelope write bitmask: %w", err) }
	_, err := body.WriteTo(buf); return err
}

func ValidateEnvelope(s *Envelope) error {
	if s == nil { return nil }
	if s.Item != nil {
		if err := ValidateItem(s.Item); err != nil { return fmt.Errorf("Item: %w", err) }
	}
	if err := ValidateItemList(s.Items); err != nil { return fmt.Errorf("Items: %w", err) }
	return nil
}

func EqEnvelope(a, b *Envelope) bool {
	if a == b { return true }
	if a == nil || b == nil { return false }
	if !EqItem(a.Item, b.Item) { return false }
	if !EqItemList((ItemList)(a.Items), (ItemList)(b.Items)) { return false }
	if !EqText(a.Note, b.Note) { return false }
	return true
}

// Standalone functions
func ReadEnvelope(buf *bytes.Buffer) (*Envelope, error) {
	s := new(Envelope)
	return s, GetEnvelope(buf, s)
}

type EnvelopeList []*Envelope
func GetEnvelopeList(buf *bytes.Buffer) (EnvelopeList, error) { return getList[*Envelope, EnvelopeList](buf, ReadEnvelope) }
func SetEnvelopeList(buf *bytes.Buffer, v EnvelopeList) error { return setList(buf, v, SetEnvelope) }
func ValidateEnvelopeList(v EnvelopeList) error {
	for i, item := range v {
		if item == nil { return fmt.Errorf("EnvelopeList[%d]: nil item", i) }
		if err := ValidateEnvelope(item); err != nil { return fmt.Errorf("EnvelopeList[%d]: %w", i, err) }
	}
	return nil
}
func EqEnvelopeList(a, b EnvelopeList) bool { return slices.EqualFunc(a, b, EqEnvelope) }
