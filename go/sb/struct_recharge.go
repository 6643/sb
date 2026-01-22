package sb

import (
	"bytes"
	"fmt"
	"math"
	"slices"
	"unsafe"
)

type Recharge struct {
	Id uint32 `bson:"_id" json:"_id"` // abcd
	Type []OrderStatus `bson:"type" json:"type"` 
	Phone []string `bson:"phone" json:"phone"` 
	Si *SimInfo `bson:"si" json:"si"` 
}

func (s *Recharge) Get(buf *bytes.Buffer) error {
	if buf.Len() == 0 { return nil }
	bitSize := int(math.Ceil(float64(4) / 8.0))
	if buf.Len() < bitSize { return fmt.Errorf("GetRecharge bitmask: %d - %d", buf.Len(), bitSize) }
	bits := buf.Next(bitSize)
	if GetBit(bits, uint8(0)) {
		val, err := GetU32(buf)
		if err != nil { return fmt.Errorf("GetRecharge Id: %w", err) }
		s.Id = val
	}
	if GetBit(bits, uint8(1)) {
		val, err := GetU8List(buf)
		if err != nil { return fmt.Errorf("GetRecharge Type: %w", err) }
		s.Type = *(*[]OrderStatus)(unsafe.Pointer(&val))
	}
	if GetBit(bits, uint8(2)) {
		val, err := GetTextList(buf)
		if err != nil { return fmt.Errorf("GetRecharge Phone: %w", err) }
		s.Phone = val
	}
	if GetBit(bits, uint8(3)) {
		if s.Si == nil { s.Si = new(SimInfo) }
		if err := s.Si.Get(buf); err != nil { return fmt.Errorf("GetRecharge Si: %w", err) }
	}
	return nil
}

func (s *Recharge) Set(buf *bytes.Buffer) error {
	if s == nil { return nil }
	bits := make([]byte, uint8(math.Ceil(float64(4)/8.0)))
	body := bytes.NewBuffer(nil)
	if s.Id != 0 {
		if err := SetU32(body, s.Id); err != nil { return fmt.Errorf("SetRecharge Id: %w", err) }
		SetBit(bits, uint8(0), true)
	}
	if len(s.Type) > 0 {
		if err := SetU8List(body, *(*[]uint8)(unsafe.Pointer(&s.Type))); err != nil { return fmt.Errorf("SetRecharge Type: %w", err) }
		SetBit(bits, uint8(1), true)
	}
	if len(s.Phone) > 0 {
		if err := SetTextList(body, s.Phone); err != nil { return fmt.Errorf("SetRecharge Phone: %w", err) }
		SetBit(bits, uint8(2), true)
	}
	if s.Si != nil {
		if err := s.Si.Set(body); err != nil { return fmt.Errorf("SetRecharge Si: %w", err) }
		SetBit(bits, uint8(3), true)
	}

	if _, err := buf.Write(bits); err != nil { return fmt.Errorf("SetRecharge write bitmask: %w", err) }
	_, err := body.WriteTo(buf); return err
}

func (s *Recharge) Eq(other *Recharge) bool {
	if s == other { return true }
	if s == nil || other == nil { return false }
	if !EqU32(s.Id, other.Id) { return false }
	if !slices.Equal(s.Type, other.Type) { return false }
	if !EqTextList(s.Phone, other.Phone) { return false }
	if !s.Si.Eq(other.Si) { return false }
	return true
}

// Standalone functions for compatibility
func GetRecharge(buf *bytes.Buffer) (*Recharge, error) {
	s := new(Recharge); return s, s.Get(buf)
}
func SetRecharge(buf *bytes.Buffer, s *Recharge) error { return s.Set(buf) }
func EqRecharge(a, b *Recharge) bool { return a.Eq(b) }

type RechargeList []*Recharge
func (v RechargeList) Set(buf *bytes.Buffer) error { return setList(buf, v, SetRecharge) }
func (v *RechargeList) Get(buf *bytes.Buffer) error {
	val, err := getList[*Recharge, RechargeList](buf, GetRecharge)
	if err == nil { *v = val }; return err
}
func (v RechargeList) Eq(other RechargeList) bool { return slices.EqualFunc(v, other, EqRecharge) }
