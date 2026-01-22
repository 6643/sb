package sb

import (
	"bytes"
	"fmt"
	"math"
	"slices"
	
)

type SimOrder2 struct {
	Id uint32  // SIM卡ID
	Name string  // 办理人姓名
	Phone string  // 联系电话
	IdNo string  // 身份证号
	CityCode uint32  // 所在城市
	Address string  // 详细地址
	NewPhone string  // 新手机号码
}

func (s *SimOrder2) Get(buf *bytes.Buffer) error {
	if buf.Len() == 0 { return nil }
	bitSize := int(math.Ceil(float64(7) / 8.0))
	if buf.Len() < bitSize { return fmt.Errorf("GetSimOrder2 bitmask: %d - %d", buf.Len(), bitSize) }
	bits := buf.Next(bitSize)
	if GetBit(bits, uint8(0)) {
		val, err := GetU32(buf)
		if err != nil { return fmt.Errorf("GetSimOrder2 Id: %w", err) }
		s.Id = val
	}
	if GetBit(bits, uint8(1)) {
		val, err := GetText(buf)
		if err != nil { return fmt.Errorf("GetSimOrder2 Name: %w", err) }
		s.Name = val
	}
	if GetBit(bits, uint8(2)) {
		val, err := GetText(buf)
		if err != nil { return fmt.Errorf("GetSimOrder2 Phone: %w", err) }
		s.Phone = val
	}
	if GetBit(bits, uint8(3)) {
		val, err := GetText(buf)
		if err != nil { return fmt.Errorf("GetSimOrder2 IdNo: %w", err) }
		s.IdNo = val
	}
	if GetBit(bits, uint8(4)) {
		val, err := GetU32(buf)
		if err != nil { return fmt.Errorf("GetSimOrder2 CityCode: %w", err) }
		s.CityCode = val
	}
	if GetBit(bits, uint8(5)) {
		val, err := GetText(buf)
		if err != nil { return fmt.Errorf("GetSimOrder2 Address: %w", err) }
		s.Address = val
	}
	if GetBit(bits, uint8(6)) {
		val, err := GetText(buf)
		if err != nil { return fmt.Errorf("GetSimOrder2 NewPhone: %w", err) }
		s.NewPhone = val
	}
	return nil
}

func (s *SimOrder2) Set(buf *bytes.Buffer) error {
	if s == nil { return nil }
	bits := make([]byte, uint8(math.Ceil(float64(7)/8.0)))
	body := bytes.NewBuffer(nil)
	if s.Id != 0 {
		if err := SetU32(body, s.Id); err != nil { return fmt.Errorf("SetSimOrder2 Id: %w", err) }
		SetBit(bits, uint8(0), true)
	}
	if s.Name != "" {
		if err := SetText(body, s.Name); err != nil { return fmt.Errorf("SetSimOrder2 Name: %w", err) }
		SetBit(bits, uint8(1), true)
	}
	if s.Phone != "" {
		if err := SetText(body, s.Phone); err != nil { return fmt.Errorf("SetSimOrder2 Phone: %w", err) }
		SetBit(bits, uint8(2), true)
	}
	if s.IdNo != "" {
		if err := SetText(body, s.IdNo); err != nil { return fmt.Errorf("SetSimOrder2 IdNo: %w", err) }
		SetBit(bits, uint8(3), true)
	}
	if s.CityCode != 0 {
		if err := SetU32(body, s.CityCode); err != nil { return fmt.Errorf("SetSimOrder2 CityCode: %w", err) }
		SetBit(bits, uint8(4), true)
	}
	if s.Address != "" {
		if err := SetText(body, s.Address); err != nil { return fmt.Errorf("SetSimOrder2 Address: %w", err) }
		SetBit(bits, uint8(5), true)
	}
	if s.NewPhone != "" {
		if err := SetText(body, s.NewPhone); err != nil { return fmt.Errorf("SetSimOrder2 NewPhone: %w", err) }
		SetBit(bits, uint8(6), true)
	}

	if _, err := buf.Write(bits); err != nil { return fmt.Errorf("SetSimOrder2 write bitmask: %w", err) }
	_, err := body.WriteTo(buf); return err
}

func (s *SimOrder2) Eq(other *SimOrder2) bool {
	if s == other { return true }
	if s == nil || other == nil { return false }
	if !EqU32(s.Id, other.Id) { return false }
	if !EqText(s.Name, other.Name) { return false }
	if !EqText(s.Phone, other.Phone) { return false }
	if !EqText(s.IdNo, other.IdNo) { return false }
	if !EqU32(s.CityCode, other.CityCode) { return false }
	if !EqText(s.Address, other.Address) { return false }
	if !EqText(s.NewPhone, other.NewPhone) { return false }
	return true
}

// Standalone functions for compatibility
func GetSimOrder2(buf *bytes.Buffer) (*SimOrder2, error) {
	s := new(SimOrder2); return s, s.Get(buf)
}
func SetSimOrder2(buf *bytes.Buffer, s *SimOrder2) error { return s.Set(buf) }
func EqSimOrder2(a, b *SimOrder2) bool { return a.Eq(b) }

type SimOrder2List []*SimOrder2
func (v SimOrder2List) Set(buf *bytes.Buffer) error { return setList(buf, v, SetSimOrder2) }
func (v *SimOrder2List) Get(buf *bytes.Buffer) error {
	val, err := getList[*SimOrder2, SimOrder2List](buf, GetSimOrder2)
	if err == nil { *v = val }; return err
}
func (v SimOrder2List) Eq(other SimOrder2List) bool { return slices.EqualFunc(v, other, EqSimOrder2) }
