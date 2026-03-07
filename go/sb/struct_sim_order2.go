package sb

import (
	"bytes"
	"fmt"
	"slices"
)

type SimOrder2 struct {
	// SIM卡ID
	Id uint32 `bson:"id" json:"id"`
	// 办理人姓名
	Name string `bson:"name" json:"name"`
	// 联系电话
	Phone string `bson:"phone" json:"phone"`
	// 身份证号
	IdNo string `bson:"id_no" json:"id_no"`
	// 所在城市
	CityCode uint32 `bson:"city_code" json:"city_code"`
	// 详细地址
	Address string `bson:"address" json:"address"`
	// 新手机号码
	NewPhone string `bson:"new_phone" json:"new_phone"`
}

func sizeSimOrder2Body(s *SimOrder2, bits []byte) (int, error) {
	if s == nil { return 0, fmt.Errorf("sizeSimOrder2: nil value") }
	bodySize := 0
	if s.Id != 0 {
		bodySize += 4
		if bits != nil { SetBit(bits, uint8(0), true) }
	}
	if s.Name != "" {
		fieldSize, err := sizeText(s.Name)
		if err != nil { return 0, fmt.Errorf("sizeSimOrder2 Name: %w", err) }
		bodySize += fieldSize
		if bits != nil { SetBit(bits, uint8(1), true) }
	}
	if s.Phone != "" {
		fieldSize, err := sizeText(s.Phone)
		if err != nil { return 0, fmt.Errorf("sizeSimOrder2 Phone: %w", err) }
		bodySize += fieldSize
		if bits != nil { SetBit(bits, uint8(2), true) }
	}
	if s.IdNo != "" {
		fieldSize, err := sizeText(s.IdNo)
		if err != nil { return 0, fmt.Errorf("sizeSimOrder2 IdNo: %w", err) }
		bodySize += fieldSize
		if bits != nil { SetBit(bits, uint8(3), true) }
	}
	if s.CityCode != 0 {
		bodySize += 4
		if bits != nil { SetBit(bits, uint8(4), true) }
	}
	if s.Address != "" {
		fieldSize, err := sizeText(s.Address)
		if err != nil { return 0, fmt.Errorf("sizeSimOrder2 Address: %w", err) }
		bodySize += fieldSize
		if bits != nil { SetBit(bits, uint8(5), true) }
	}
	if s.NewPhone != "" {
		fieldSize, err := sizeText(s.NewPhone)
		if err != nil { return 0, fmt.Errorf("sizeSimOrder2 NewPhone: %w", err) }
		bodySize += fieldSize
		if bits != nil { SetBit(bits, uint8(6), true) }
	}
	return bodySize, nil
}

func sizeSimOrder2(s *SimOrder2) (int, error) {
	bodySize, err := sizeSimOrder2Body(s, nil)
	if err != nil { return 0, err }
	return 1 + bodySize, nil
}

func sizeSimOrder2List(v SimOrder2List) (int, error) {
	if len(v) > 65535 { return 0, fmt.Errorf("list length exceeds uint16 max") }
	totalSize := 2
	for i, item := range v {
		itemSize, err := sizeSimOrder2(item)
		if err != nil { return 0, fmt.Errorf("sizeSimOrder2List[%d]: %w", i, err) }
		totalSize += itemSize
	}
	return totalSize, nil
}

func GetSimOrder2(buf *bytes.Buffer, s *SimOrder2) error {
	if s == nil { return nil }
	const bitSize = 1
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
	const bitSize = 1
	var bits [1]byte
	bodySize, err := sizeSimOrder2Body(s, bits[:])
	if err != nil { return err }
	buf.Grow(bitSize + bodySize)
	if _, err := buf.Write(bits[:]); err != nil { return fmt.Errorf("SetSimOrder2 write bitmask: %w", err) }
	if s.Id != 0 {
		if err := SetU32(buf, s.Id); err != nil { return fmt.Errorf("SetSimOrder2 Id: %w", err) }
	}
	if s.Name != "" {
		if err := SetText(buf, s.Name); err != nil { return fmt.Errorf("SetSimOrder2 Name: %w", err) }
	}
	if s.Phone != "" {
		if err := SetText(buf, s.Phone); err != nil { return fmt.Errorf("SetSimOrder2 Phone: %w", err) }
	}
	if s.IdNo != "" {
		if err := SetText(buf, s.IdNo); err != nil { return fmt.Errorf("SetSimOrder2 IdNo: %w", err) }
	}
	if s.CityCode != 0 {
		if err := SetU32(buf, s.CityCode); err != nil { return fmt.Errorf("SetSimOrder2 CityCode: %w", err) }
	}
	if s.Address != "" {
		if err := SetText(buf, s.Address); err != nil { return fmt.Errorf("SetSimOrder2 Address: %w", err) }
	}
	if s.NewPhone != "" {
		if err := SetText(buf, s.NewPhone); err != nil { return fmt.Errorf("SetSimOrder2 NewPhone: %w", err) }
	}
	return nil
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
func GetSimOrder2List(buf *bytes.Buffer) (SimOrder2List, error) {
	count, err := GetU16(buf)
	if err != nil { return nil, err }
	list := make([]*SimOrder2, count)
	for i := range list {
		item, err := ReadSimOrder2(buf)
		if err != nil { return nil, err }
		list[i] = item
	}
	return SimOrder2List(list), nil
}
func SetSimOrder2List(buf *bytes.Buffer, v SimOrder2List) error {
	totalSize, err := sizeSimOrder2List(v)
	if err != nil { return err }
	buf.Grow(totalSize)
	if err := SetU16(buf, uint16(len(v))); err != nil { return err }
	for _, item := range v {
		if err := SetSimOrder2(buf, item); err != nil { return err }
	}
	return nil
}
func ValidateSimOrder2List(v SimOrder2List) error {
	for i, item := range v {
		if item == nil { return fmt.Errorf("SimOrder2List[%d]: nil item", i) }
		if err := ValidateSimOrder2(item); err != nil { return fmt.Errorf("SimOrder2List[%d]: %w", i, err) }
	}
	return nil
}
func EqSimOrder2List(a, b SimOrder2List) bool { return slices.EqualFunc(a, b, EqSimOrder2) }
