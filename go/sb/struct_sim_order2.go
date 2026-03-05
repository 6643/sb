package sb

import (
	"bytes"
	"fmt"
	"math"
	"slices"
)

type SimOrder2 struct {
	Id uint32 `bson:"id" json:"id"` // SIM卡ID
	Name string `bson:"name" json:"name"` // 办理人姓名
	Phone string `bson:"phone" json:"phone"` // 联系电话
	IdNo string `bson:"id_no" json:"id_no"` // 身份证号
	CityCode uint32 `bson:"city_code" json:"city_code"` // 所在城市
	Address string `bson:"address" json:"address"` // 详细地址
	NewPhone string `bson:"new_phone" json:"new_phone"` // 新手机号码
}

func GetSimOrder2(buf *bytes.Buffer, s *SimOrder2) error {
	if s == nil { return nil }
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
	if err := ValidateSimOrder2(s); err != nil { return fmt.Errorf("ValidateSimOrder2: %w", err) }
	return nil
}

func SetSimOrder2(buf *bytes.Buffer, s *SimOrder2) error {
	if s == nil { return fmt.Errorf("SetSimOrder2: nil value") }
	if err := ValidateSimOrder2(s); err != nil { return fmt.Errorf("ValidateSimOrder2: %w", err) }
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

func ValidateSimOrder2(s *SimOrder2) error {
	if s == nil { return nil }
	return nil
}

func EqSimOrder2(a, b *SimOrder2) bool {
	if a == b { return true }
	if a == nil || b == nil { return false }
	if !EqU32(a.Id, b.Id) { return false }
	if !EqText(a.Name, b.Name) { return false }
	if !EqText(a.Phone, b.Phone) { return false }
	if !EqText(a.IdNo, b.IdNo) { return false }
	if !EqU32(a.CityCode, b.CityCode) { return false }
	if !EqText(a.Address, b.Address) { return false }
	if !EqText(a.NewPhone, b.NewPhone) { return false }
	return true
}

// Standalone functions
func ReadSimOrder2(buf *bytes.Buffer) (*SimOrder2, error) {
	s := new(SimOrder2)
	return s, GetSimOrder2(buf, s)
}

type SimOrder2List []*SimOrder2
func GetSimOrder2List(buf *bytes.Buffer) (SimOrder2List, error) { return getList[*SimOrder2, SimOrder2List](buf, ReadSimOrder2) }
func SetSimOrder2List(buf *bytes.Buffer, v SimOrder2List) error { return setList(buf, v, SetSimOrder2) }
func ValidateSimOrder2List(v SimOrder2List) error {
	for i, item := range v {
		if item == nil { return fmt.Errorf("SimOrder2List[%d]: nil item", i) }
		if err := ValidateSimOrder2(item); err != nil { return fmt.Errorf("SimOrder2List[%d]: %w", i, err) }
	}
	return nil
}
func EqSimOrder2List(a, b SimOrder2List) bool { return slices.EqualFunc(a, b, EqSimOrder2) }
