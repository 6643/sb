package sb

import (
	"bytes"
	"fmt"
	"slices"
)

type SimOrder struct {
	Id uint32 `bson:"id" json:"id"`
	AccountId uint32 `bson:"account_id" json:"account_id"`
	ItemId uint32 `bson:"item_id" json:"item_id"`
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
	// 佣金
	Commission uint16 `bson:"commission" json:"commission"`
	Status OrderStatus `bson:"status" json:"status"`
}

func sizeSimOrderBody(s *SimOrder, bits []byte) (int, error) {
	if s == nil { return 0, fmt.Errorf("sizeSimOrder: nil value") }
	bodySize := 0
	if s.Id != 0 {
		bodySize += 4
		if bits != nil { SetBit(bits, uint8(0), true) }
	}
	if s.AccountId != 0 {
		bodySize += 4
		if bits != nil { SetBit(bits, uint8(1), true) }
	}
	if s.ItemId != 0 {
		bodySize += 4
		if bits != nil { SetBit(bits, uint8(2), true) }
	}
	if s.Name != "" {
		fieldSize, err := sizeText(s.Name)
		if err != nil { return 0, fmt.Errorf("sizeSimOrder Name: %w", err) }
		bodySize += fieldSize
		if bits != nil { SetBit(bits, uint8(3), true) }
	}
	if s.Phone != "" {
		fieldSize, err := sizeText(s.Phone)
		if err != nil { return 0, fmt.Errorf("sizeSimOrder Phone: %w", err) }
		bodySize += fieldSize
		if bits != nil { SetBit(bits, uint8(4), true) }
	}
	if s.IdNo != "" {
		fieldSize, err := sizeText(s.IdNo)
		if err != nil { return 0, fmt.Errorf("sizeSimOrder IdNo: %w", err) }
		bodySize += fieldSize
		if bits != nil { SetBit(bits, uint8(5), true) }
	}
	if s.CityCode != 0 {
		bodySize += 4
		if bits != nil { SetBit(bits, uint8(6), true) }
	}
	if s.Address != "" {
		fieldSize, err := sizeText(s.Address)
		if err != nil { return 0, fmt.Errorf("sizeSimOrder Address: %w", err) }
		bodySize += fieldSize
		if bits != nil { SetBit(bits, uint8(7), true) }
	}
	if s.NewPhone != "" {
		fieldSize, err := sizeText(s.NewPhone)
		if err != nil { return 0, fmt.Errorf("sizeSimOrder NewPhone: %w", err) }
		bodySize += fieldSize
		if bits != nil { SetBit(bits, uint8(8), true) }
	}
	if s.Commission != 0 {
		bodySize += 2
		if bits != nil { SetBit(bits, uint8(9), true) }
	}
	if s.Status != 0 {
		bodySize += 1
		if bits != nil { SetBit(bits, uint8(10), true) }
	}
	return bodySize, nil
}

func sizeSimOrder(s *SimOrder) (int, error) {
	bodySize, err := sizeSimOrderBody(s, nil)
	if err != nil { return 0, err }
	return 2 + bodySize, nil
}

func sizeSimOrderList(v SimOrderList) (int, error) {
	if len(v) > 65535 { return 0, fmt.Errorf("list length exceeds uint16 max") }
	totalSize := 2
	for i, item := range v {
		itemSize, err := sizeSimOrder(item)
		if err != nil { return 0, fmt.Errorf("sizeSimOrderList[%d]: %w", i, err) }
		totalSize += itemSize
	}
	return totalSize, nil
}

func GetSimOrder(buf *bytes.Buffer, s *SimOrder) error {
	if s == nil { return nil }
	const bitSize = 2
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
	const bitSize = 2
	var bits [2]byte
	bodySize, err := sizeSimOrderBody(s, bits[:])
	if err != nil { return err }
	buf.Grow(bitSize + bodySize)
	if _, err := buf.Write(bits[:]); err != nil { return fmt.Errorf("SetSimOrder write bitmask: %w", err) }
	if s.Id != 0 {
		if err := SetU32(buf, s.Id); err != nil { return fmt.Errorf("SetSimOrder Id: %w", err) }
	}
	if s.AccountId != 0 {
		if err := SetU32(buf, s.AccountId); err != nil { return fmt.Errorf("SetSimOrder AccountId: %w", err) }
	}
	if s.ItemId != 0 {
		if err := SetU32(buf, s.ItemId); err != nil { return fmt.Errorf("SetSimOrder ItemId: %w", err) }
	}
	if s.Name != "" {
		if err := SetText(buf, s.Name); err != nil { return fmt.Errorf("SetSimOrder Name: %w", err) }
	}
	if s.Phone != "" {
		if err := SetText(buf, s.Phone); err != nil { return fmt.Errorf("SetSimOrder Phone: %w", err) }
	}
	if s.IdNo != "" {
		if err := SetText(buf, s.IdNo); err != nil { return fmt.Errorf("SetSimOrder IdNo: %w", err) }
	}
	if s.CityCode != 0 {
		if err := SetU32(buf, s.CityCode); err != nil { return fmt.Errorf("SetSimOrder CityCode: %w", err) }
	}
	if s.Address != "" {
		if err := SetText(buf, s.Address); err != nil { return fmt.Errorf("SetSimOrder Address: %w", err) }
	}
	if s.NewPhone != "" {
		if err := SetText(buf, s.NewPhone); err != nil { return fmt.Errorf("SetSimOrder NewPhone: %w", err) }
	}
	if s.Commission != 0 {
		if err := SetU16(buf, s.Commission); err != nil { return fmt.Errorf("SetSimOrder Commission: %w", err) }
	}
	if s.Status != 0 {
		if err := SetOrderStatus(buf, s.Status); err != nil { return fmt.Errorf("SetSimOrder Status: %w", err) }
	}
	return nil
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
	if !(a.Status == b.Status) { return false }
	return true
}

// Standalone functions
func ReadSimOrder(buf *bytes.Buffer) (*SimOrder, error) {
	s := new(SimOrder)
	return s, GetSimOrder(buf, s)
}

type SimOrderList []*SimOrder
func GetSimOrderList(buf *bytes.Buffer) (SimOrderList, error) {
	count, err := GetU16(buf)
	if err != nil { return nil, err }
	list := make([]*SimOrder, count)
	for i := range list {
		item, err := ReadSimOrder(buf)
		if err != nil { return nil, err }
		list[i] = item
	}
	return SimOrderList(list), nil
}
func SetSimOrderList(buf *bytes.Buffer, v SimOrderList) error {
	totalSize, err := sizeSimOrderList(v)
	if err != nil { return err }
	buf.Grow(totalSize)
	if err := SetU16(buf, uint16(len(v))); err != nil { return err }
	for _, item := range v {
		if err := SetSimOrder(buf, item); err != nil { return err }
	}
	return nil
}
func ValidateSimOrderList(v SimOrderList) error {
	for i, item := range v {
		if item == nil { return fmt.Errorf("SimOrderList[%d]: nil item", i) }
		if err := ValidateSimOrder(item); err != nil { return fmt.Errorf("SimOrderList[%d]: %w", i, err) }
	}
	return nil
}
func EqSimOrderList(a, b SimOrderList) bool { return slices.EqualFunc(a, b, EqSimOrder) }
