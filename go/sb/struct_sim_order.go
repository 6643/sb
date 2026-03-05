package sb

import (
	"bytes"
	"fmt"
	"math"
	"slices"
)

type SimOrder struct {
	Id uint32 `bson:"id" json:"id"` 
	AccountId uint32 `bson:"account_id" json:"account_id"` 
	ItemId uint32 `bson:"item_id" json:"item_id"` 
	Name string `bson:"name" json:"name"` // 办理人姓名
	Phone string `bson:"phone" json:"phone"` // 联系电话
	IdNo string `bson:"id_no" json:"id_no"` // 身份证号
	CityCode uint32 `bson:"city_code" json:"city_code"` // 所在城市
	Address string `bson:"address" json:"address"` // 详细地址
	NewPhone string `bson:"new_phone" json:"new_phone"` // 新手机号码
	Commission uint16 `bson:"commission" json:"commission"` // 佣金
	Status OrderStatus `bson:"status" json:"status"` 
}

func GetSimOrder(buf *bytes.Buffer, s *SimOrder) error {
	if s == nil { return nil }
	bitSize := int(math.Ceil(float64(11) / 8.0))
	if buf.Len() < bitSize { return fmt.Errorf("GetSimOrder bitmask: %d - %d", buf.Len(), bitSize) }
	bits := buf.Next(bitSize)
	if GetBit(bits, uint8(0)) {
		val, err := GetU32(buf)
		if err != nil { return fmt.Errorf("GetSimOrder Id: %w", err) }
		s.Id = val
	}
	if GetBit(bits, uint8(1)) {
		val, err := GetU32(buf)
		if err != nil { return fmt.Errorf("GetSimOrder AccountId: %w", err) }
		s.AccountId = val
	}
	if GetBit(bits, uint8(2)) {
		val, err := GetU32(buf)
		if err != nil { return fmt.Errorf("GetSimOrder ItemId: %w", err) }
		s.ItemId = val
	}
	if GetBit(bits, uint8(3)) {
		val, err := GetText(buf)
		if err != nil { return fmt.Errorf("GetSimOrder Name: %w", err) }
		s.Name = val
	}
	if GetBit(bits, uint8(4)) {
		val, err := GetText(buf)
		if err != nil { return fmt.Errorf("GetSimOrder Phone: %w", err) }
		s.Phone = val
	}
	if GetBit(bits, uint8(5)) {
		val, err := GetText(buf)
		if err != nil { return fmt.Errorf("GetSimOrder IdNo: %w", err) }
		s.IdNo = val
	}
	if GetBit(bits, uint8(6)) {
		val, err := GetU32(buf)
		if err != nil { return fmt.Errorf("GetSimOrder CityCode: %w", err) }
		s.CityCode = val
	}
	if GetBit(bits, uint8(7)) {
		val, err := GetText(buf)
		if err != nil { return fmt.Errorf("GetSimOrder Address: %w", err) }
		s.Address = val
	}
	if GetBit(bits, uint8(8)) {
		val, err := GetText(buf)
		if err != nil { return fmt.Errorf("GetSimOrder NewPhone: %w", err) }
		s.NewPhone = val
	}
	if GetBit(bits, uint8(9)) {
		val, err := GetU16(buf)
		if err != nil { return fmt.Errorf("GetSimOrder Commission: %w", err) }
		s.Commission = val
	}
	if GetBit(bits, uint8(10)) {
		val, err := GetOrderStatus(buf)
		if err != nil { return fmt.Errorf("GetSimOrder Status: %w", err) }
		s.Status = val
		if !IsOrderStatus(s.Status) { return fmt.Errorf("GetSimOrder Status: 非法枚举值: %d", s.Status) }
	}
	if err := ValidateSimOrder(s); err != nil { return fmt.Errorf("ValidateSimOrder: %w", err) }
	return nil
}

func SetSimOrder(buf *bytes.Buffer, s *SimOrder) error {
	if s == nil { return fmt.Errorf("SetSimOrder: nil value") }
	if err := ValidateSimOrder(s); err != nil { return fmt.Errorf("ValidateSimOrder: %w", err) }
	bits := make([]byte, uint8(math.Ceil(float64(11)/8.0)))
	body := bytes.NewBuffer(nil)
	if s.Id != 0 {
		if err := SetU32(body, s.Id); err != nil { return fmt.Errorf("SetSimOrder Id: %w", err) }
		SetBit(bits, uint8(0), true)
	}
	if s.AccountId != 0 {
		if err := SetU32(body, s.AccountId); err != nil { return fmt.Errorf("SetSimOrder AccountId: %w", err) }
		SetBit(bits, uint8(1), true)
	}
	if s.ItemId != 0 {
		if err := SetU32(body, s.ItemId); err != nil { return fmt.Errorf("SetSimOrder ItemId: %w", err) }
		SetBit(bits, uint8(2), true)
	}
	if s.Name != "" {
		if err := SetText(body, s.Name); err != nil { return fmt.Errorf("SetSimOrder Name: %w", err) }
		SetBit(bits, uint8(3), true)
	}
	if s.Phone != "" {
		if err := SetText(body, s.Phone); err != nil { return fmt.Errorf("SetSimOrder Phone: %w", err) }
		SetBit(bits, uint8(4), true)
	}
	if s.IdNo != "" {
		if err := SetText(body, s.IdNo); err != nil { return fmt.Errorf("SetSimOrder IdNo: %w", err) }
		SetBit(bits, uint8(5), true)
	}
	if s.CityCode != 0 {
		if err := SetU32(body, s.CityCode); err != nil { return fmt.Errorf("SetSimOrder CityCode: %w", err) }
		SetBit(bits, uint8(6), true)
	}
	if s.Address != "" {
		if err := SetText(body, s.Address); err != nil { return fmt.Errorf("SetSimOrder Address: %w", err) }
		SetBit(bits, uint8(7), true)
	}
	if s.NewPhone != "" {
		if err := SetText(body, s.NewPhone); err != nil { return fmt.Errorf("SetSimOrder NewPhone: %w", err) }
		SetBit(bits, uint8(8), true)
	}
	if s.Commission != 0 {
		if err := SetU16(body, s.Commission); err != nil { return fmt.Errorf("SetSimOrder Commission: %w", err) }
		SetBit(bits, uint8(9), true)
	}
	if s.Status != 0 {
		if err := SetOrderStatus(body, s.Status); err != nil { return fmt.Errorf("SetSimOrder Status: %w", err) }
		SetBit(bits, uint8(10), true)
	}

	if _, err := buf.Write(bits); err != nil { return fmt.Errorf("SetSimOrder write bitmask: %w", err) }
	_, err := body.WriteTo(buf); return err
}

func ValidateSimOrder(s *SimOrder) error {
	if s == nil { return nil }
	if s.Status != 0 && !IsOrderStatus(s.Status) { return fmt.Errorf("Status 非法枚举值: %d", s.Status) }
	return nil
}

func EqSimOrder(a, b *SimOrder) bool {
	if a == b { return true }
	if a == nil || b == nil { return false }
	if !EqU32(a.Id, b.Id) { return false }
	if !EqU32(a.AccountId, b.AccountId) { return false }
	if !EqU32(a.ItemId, b.ItemId) { return false }
	if !EqText(a.Name, b.Name) { return false }
	if !EqText(a.Phone, b.Phone) { return false }
	if !EqText(a.IdNo, b.IdNo) { return false }
	if !EqU32(a.CityCode, b.CityCode) { return false }
	if !EqText(a.Address, b.Address) { return false }
	if !EqText(a.NewPhone, b.NewPhone) { return false }
	if !EqU16(a.Commission, b.Commission) { return false }
	if a.Status != b.Status { return false }
	return true
}

// Standalone functions
func ReadSimOrder(buf *bytes.Buffer) (*SimOrder, error) {
	s := new(SimOrder)
	return s, GetSimOrder(buf, s)
}

type SimOrderList []*SimOrder
func GetSimOrderList(buf *bytes.Buffer) (SimOrderList, error) { return getList[*SimOrder, SimOrderList](buf, ReadSimOrder) }
func SetSimOrderList(buf *bytes.Buffer, v SimOrderList) error { return setList(buf, v, SetSimOrder) }
func ValidateSimOrderList(v SimOrderList) error {
	for i, item := range v {
		if item == nil { return fmt.Errorf("SimOrderList[%d]: nil item", i) }
		if err := ValidateSimOrder(item); err != nil { return fmt.Errorf("SimOrderList[%d]: %w", i, err) }
	}
	return nil
}
func EqSimOrderList(a, b SimOrderList) bool { return slices.EqualFunc(a, b, EqSimOrder) }
