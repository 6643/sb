package sb

import (
	"bytes"
	"fmt"
	"slices"
)

type Recharge struct {
	// abcd
	Id uint32 `bson:"_id" json:"_id"`
	Type []OrderStatus `bson:"type" json:"type"`
	Phone []string `bson:"phone" json:"phone"`
	Si *SimInfo `bson:"si" json:"si"`
}

func sizeRechargeBody(s *Recharge, bits []byte) (int, error) {
	if s == nil { return 0, fmt.Errorf("sizeRecharge: nil value") }
	bodySize := 0
	if s.Id != 0 {
		bodySize += 4
		if bits != nil { SetBit(bits, uint8(0), true) }
	}
	if len(s.Type) > 0 {
		fieldSize, err := sizeFixedList(len(s.Type), 1)
		if err != nil { return 0, fmt.Errorf("sizeRecharge Type: %w", err) }
		bodySize += fieldSize
		if bits != nil { SetBit(bits, uint8(1), true) }
	}
	if len(s.Phone) > 0 {
		fieldSize, err := sizeTextList(s.Phone)
		if err != nil { return 0, fmt.Errorf("sizeRecharge Phone: %w", err) }
		bodySize += fieldSize
		if bits != nil { SetBit(bits, uint8(2), true) }
	}
	if s.Si != nil {
		fieldSize, err := sizeSimInfo(s.Si)
		if err != nil { return 0, fmt.Errorf("sizeRecharge Si: %w", err) }
		bodySize += fieldSize
		if bits != nil { SetBit(bits, uint8(3), true) }
	}
	return bodySize, nil
}

func sizeRecharge(s *Recharge) (int, error) {
	bodySize, err := sizeRechargeBody(s, nil)
	if err != nil { return 0, err }
	return 1 + bodySize, nil
}

func sizeRechargeList(v RechargeList) (int, error) {
	if len(v) > 65535 { return 0, fmt.Errorf("list length exceeds uint16 max") }
	totalSize := 2
	for i, item := range v {
		itemSize, err := sizeRecharge(item)
		if err != nil { return 0, fmt.Errorf("sizeRechargeList[%d]: %w", i, err) }
		totalSize += itemSize
	}
	return totalSize, nil
}

func GetRecharge(buf *bytes.Buffer, s *Recharge) error {
	if s == nil { return nil }
	const bitSize = 1
	if buf.Len() < bitSize { return fmt.Errorf("GetRecharge bitmask: %d - %d", buf.Len(), bitSize) }
	bits := buf.Next(bitSize)
	if GetBit(bits, uint8(0)) {
		val, err := GetU32(buf)
		if err != nil { return fmt.Errorf("GetRecharge Id: %w", err) }
		s.Id = val
	}
	if GetBit(bits, uint8(1)) {
		val, err := GetOrderStatusList(buf)
		if err != nil { return fmt.Errorf("GetRecharge Type: %w", err) }
		s.Type = val
	}
	if GetBit(bits, uint8(2)) {
		val, err := GetTextList(buf)
		if err != nil { return fmt.Errorf("GetRecharge Phone: %w", err) }
		s.Phone = val
	}
	if GetBit(bits, uint8(3)) {
		val, err := ReadSimInfo(buf)
		if err != nil { return fmt.Errorf("GetRecharge Si: %w", err) }
		s.Si = val
	}
	if err := ValidateRecharge(s); err != nil { return fmt.Errorf("ValidateRecharge: %w", err) }
	return nil
}

func SetRecharge(buf *bytes.Buffer, s *Recharge) error {
	if s == nil { return fmt.Errorf("SetRecharge: nil value") }
	if err := ValidateRecharge(s); err != nil { return fmt.Errorf("ValidateRecharge: %w", err) }
	const bitSize = 1
	var bits [1]byte
	bodySize, err := sizeRechargeBody(s, bits[:])
	if err != nil { return err }
	buf.Grow(bitSize + bodySize)
	if _, err := buf.Write(bits[:]); err != nil { return fmt.Errorf("SetRecharge write bitmask: %w", err) }
	if s.Id != 0 {
		if err := SetU32(buf, s.Id); err != nil { return fmt.Errorf("SetRecharge Id: %w", err) }
	}
	if len(s.Type) > 0 {
		if err := SetOrderStatusList(buf, (OrderStatusList)(s.Type)); err != nil { return fmt.Errorf("SetRecharge Type: %w", err) }
	}
	if len(s.Phone) > 0 {
		if err := SetTextList(buf, s.Phone); err != nil { return fmt.Errorf("SetRecharge Phone: %w", err) }
	}
	if s.Si != nil {
		if err := SetSimInfo(buf, s.Si); err != nil { return fmt.Errorf("SetRecharge Si: %w", err) }
	}
	return nil
}

func ValidateRecharge(s *Recharge) error {
	if s == nil { return nil }
	for i, item := range s.Type {
		if !IsOrderStatus(item) { return fmt.Errorf("Type[%d] 非法枚举值: %d", i, item) }
	}
	if s.Si != nil {
		if err := ValidateSimInfo(s.Si); err != nil { return fmt.Errorf("Si: %w", err) }
	}
	return nil
}

func EqRecharge(a, b *Recharge) bool {
	if a == b { return true }
	if a == nil || b == nil { return false }
	if !EqU32(a.Id, b.Id) { return false }
	if !EqOrderStatusList(a.Type, b.Type) { return false }
	if !EqTextList(a.Phone, b.Phone) { return false }
	if !EqSimInfo(a.Si, b.Si) { return false }
	return true
}

// Standalone functions
func ReadRecharge(buf *bytes.Buffer) (*Recharge, error) {
	s := new(Recharge)
	return s, GetRecharge(buf, s)
}

type RechargeList []*Recharge
func GetRechargeList(buf *bytes.Buffer) (RechargeList, error) {
	count, err := GetU16(buf)
	if err != nil { return nil, err }
	list := make([]*Recharge, count)
	for i := range list {
		item, err := ReadRecharge(buf)
		if err != nil { return nil, err }
		list[i] = item
	}
	return RechargeList(list), nil
}
func SetRechargeList(buf *bytes.Buffer, v RechargeList) error {
	totalSize, err := sizeRechargeList(v)
	if err != nil { return err }
	buf.Grow(totalSize)
	if err := SetU16(buf, uint16(len(v))); err != nil { return err }
	for _, item := range v {
		if err := SetRecharge(buf, item); err != nil { return err }
	}
	return nil
}
func ValidateRechargeList(v RechargeList) error {
	for i, item := range v {
		if item == nil { return fmt.Errorf("RechargeList[%d]: nil item", i) }
		if err := ValidateRecharge(item); err != nil { return fmt.Errorf("RechargeList[%d]: %w", i, err) }
	}
	return nil
}
func EqRechargeList(a, b RechargeList) bool { return slices.EqualFunc(a, b, EqRecharge) }
