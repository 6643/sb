package sb

import (
	"bytes"
	"fmt"
	"math"
	"slices"
	"unsafe"
)

type RechargeB struct {
	Id uint32  // abcd
	Type []OrderStatus  
	Phone []string  
	Si *SimInfo  
	Bid uint32  
}

func (s *RechargeB) Get(buf *bytes.Buffer) error {
	if buf.Len() == 0 { return nil }
	bitSize := int(math.Ceil(float64(5) / 8.0))
	if buf.Len() < bitSize { return fmt.Errorf("GetRechargeB bitmask: %d - %d", buf.Len(), bitSize) }
	bits := buf.Next(bitSize)
	if GetBit(bits, uint8(0)) {
		val, err := GetU32(buf)
		if err != nil { return fmt.Errorf("GetRechargeB Id: %w", err) }
		s.Id = val
	}
	if GetBit(bits, uint8(1)) {
		val, err := GetU8List(buf)
		if err != nil { return fmt.Errorf("GetRechargeB Type: %w", err) }
		s.Type = *(*[]OrderStatus)(unsafe.Pointer(&val))
	}
	if GetBit(bits, uint8(2)) {
		val, err := GetTextList(buf)
		if err != nil { return fmt.Errorf("GetRechargeB Phone: %w", err) }
		s.Phone = val
	}
	if GetBit(bits, uint8(3)) {
		if s.Si == nil { s.Si = new(SimInfo) }
		if err := s.Si.Get(buf); err != nil { return fmt.Errorf("GetRechargeB Si: %w", err) }
	}
	if GetBit(bits, uint8(4)) {
		val, err := GetU32(buf)
		if err != nil { return fmt.Errorf("GetRechargeB Bid: %w", err) }
		s.Bid = val
	}
	return nil
}

func (s *RechargeB) Set(buf *bytes.Buffer) error {
	if s == nil { return nil }
	bits := make([]byte, uint8(math.Ceil(float64(5)/8.0)))
	body := bytes.NewBuffer(nil)
	if s.Id != 0 {
		if err := SetU32(body, s.Id); err != nil { return fmt.Errorf("SetRechargeB Id: %w", err) }
		SetBit(bits, uint8(0), true)
	}
	if len(s.Type) > 0 {
		if err := SetU8List(body, *(*[]uint8)(unsafe.Pointer(&s.Type))); err != nil { return fmt.Errorf("SetRechargeB Type: %w", err) }
		SetBit(bits, uint8(1), true)
	}
	if len(s.Phone) > 0 {
		if err := SetTextList(body, s.Phone); err != nil { return fmt.Errorf("SetRechargeB Phone: %w", err) }
		SetBit(bits, uint8(2), true)
	}
	if s.Si != nil {
		if err := s.Si.Set(body); err != nil { return fmt.Errorf("SetRechargeB Si: %w", err) }
		SetBit(bits, uint8(3), true)
	}
	if s.Bid != 0 {
		if err := SetU32(body, s.Bid); err != nil { return fmt.Errorf("SetRechargeB Bid: %w", err) }
		SetBit(bits, uint8(4), true)
	}

	if _, err := buf.Write(bits); err != nil { return fmt.Errorf("SetRechargeB write bitmask: %w", err) }
	_, err := body.WriteTo(buf); return err
}

func (s *RechargeB) Eq(other *RechargeB) bool {
	if s == other { return true }
	if s == nil || other == nil { return false }
	if !EqU32(s.Id, other.Id) { return false }
	if !slices.Equal(s.Type, other.Type) { return false }
	if !EqTextList(s.Phone, other.Phone) { return false }
	if !s.Si.Eq(other.Si) { return false }
	if !EqU32(s.Bid, other.Bid) { return false }
	return true
}

// Standalone functions for compatibility
func GetRechargeB(buf *bytes.Buffer) (*RechargeB, error) {
	s := new(RechargeB); return s, s.Get(buf)
}
func SetRechargeB(buf *bytes.Buffer, s *RechargeB) error { return s.Set(buf) }
func EqRechargeB(a, b *RechargeB) bool { return a.Eq(b) }

type RechargeBList []*RechargeB
func (v RechargeBList) Set(buf *bytes.Buffer) error { return setList(buf, v, SetRechargeB) }
func (v *RechargeBList) Get(buf *bytes.Buffer) error {
	val, err := getList[*RechargeB, RechargeBList](buf, GetRechargeB)
	if err == nil { *v = val }; return err
}
func (v RechargeBList) Eq(other RechargeBList) bool { return slices.EqualFunc(v, other, EqRechargeB) }
