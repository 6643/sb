package sb

import (
	"bytes"
	"fmt"
	"math"
	"slices"
	"unsafe"
)

type RechargeA struct {
	Id uint32  // abcd
	Type []OrderStatus  
	Phone []string  
	Si *SimInfo  
	Aid uint32  
}

func (s *RechargeA) Get(buf *bytes.Buffer) error {
	if buf.Len() == 0 { return nil }
	bitSize := int(math.Ceil(float64(5) / 8.0))
	if buf.Len() < bitSize { return fmt.Errorf("GetRechargeA bitmask: %d - %d", buf.Len(), bitSize) }
	bits := buf.Next(bitSize)
	if GetBit(bits, uint8(0)) {
		val, err := GetU32(buf)
		if err != nil { return fmt.Errorf("GetRechargeA Id: %w", err) }
		s.Id = val
	}
	if GetBit(bits, uint8(1)) {
		val, err := GetU8List(buf)
		if err != nil { return fmt.Errorf("GetRechargeA Type: %w", err) }
		s.Type = *(*[]OrderStatus)(unsafe.Pointer(&val))
	}
	if GetBit(bits, uint8(2)) {
		val, err := GetTextList(buf)
		if err != nil { return fmt.Errorf("GetRechargeA Phone: %w", err) }
		s.Phone = val
	}
	if GetBit(bits, uint8(3)) {
		if s.Si == nil { s.Si = new(SimInfo) }
		if err := s.Si.Get(buf); err != nil { return fmt.Errorf("GetRechargeA Si: %w", err) }
	}
	if GetBit(bits, uint8(4)) {
		val, err := GetU32(buf)
		if err != nil { return fmt.Errorf("GetRechargeA Aid: %w", err) }
		s.Aid = val
	}
	return nil
}

func (s *RechargeA) Set(buf *bytes.Buffer) error {
	if s == nil { return nil }
	bits := make([]byte, uint8(math.Ceil(float64(5)/8.0)))
	body := bytes.NewBuffer(nil)
	if s.Id != 0 {
		if err := SetU32(body, s.Id); err != nil { return fmt.Errorf("SetRechargeA Id: %w", err) }
		SetBit(bits, uint8(0), true)
	}
	if len(s.Type) > 0 {
		if err := SetU8List(body, *(*[]uint8)(unsafe.Pointer(&s.Type))); err != nil { return fmt.Errorf("SetRechargeA Type: %w", err) }
		SetBit(bits, uint8(1), true)
	}
	if len(s.Phone) > 0 {
		if err := SetTextList(body, s.Phone); err != nil { return fmt.Errorf("SetRechargeA Phone: %w", err) }
		SetBit(bits, uint8(2), true)
	}
	if s.Si != nil {
		if err := s.Si.Set(body); err != nil { return fmt.Errorf("SetRechargeA Si: %w", err) }
		SetBit(bits, uint8(3), true)
	}
	if s.Aid != 0 {
		if err := SetU32(body, s.Aid); err != nil { return fmt.Errorf("SetRechargeA Aid: %w", err) }
		SetBit(bits, uint8(4), true)
	}

	if _, err := buf.Write(bits); err != nil { return fmt.Errorf("SetRechargeA write bitmask: %w", err) }
	_, err := body.WriteTo(buf); return err
}

func (s *RechargeA) Eq(other *RechargeA) bool {
	if s == other { return true }
	if s == nil || other == nil { return false }
	if !EqU32(s.Id, other.Id) { return false }
	if !slices.Equal(s.Type, other.Type) { return false }
	if !EqTextList(s.Phone, other.Phone) { return false }
	if !s.Si.Eq(other.Si) { return false }
	if !EqU32(s.Aid, other.Aid) { return false }
	return true
}

// Standalone functions for compatibility
func GetRechargeA(buf *bytes.Buffer) (*RechargeA, error) {
	s := new(RechargeA); return s, s.Get(buf)
}
func SetRechargeA(buf *bytes.Buffer, s *RechargeA) error { return s.Set(buf) }
func EqRechargeA(a, b *RechargeA) bool { return a.Eq(b) }

type RechargeAList []*RechargeA
func (v RechargeAList) Set(buf *bytes.Buffer) error { return setList(buf, v, SetRechargeA) }
func (v *RechargeAList) Get(buf *bytes.Buffer) error {
	val, err := getList[*RechargeA, RechargeAList](buf, GetRechargeA)
	if err == nil { *v = val }; return err
}
func (v RechargeAList) Eq(other RechargeAList) bool { return slices.EqualFunc(v, other, EqRechargeA) }
