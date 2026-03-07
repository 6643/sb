package sb

import (
	"bytes"
	"fmt"
	"slices"
)

type RechargeA struct {
	// abcd
	Id uint32 `bson:"_id" json:"_id"`
	Type []OrderStatus `bson:"type" json:"type"`
	Phone []string `bson:"phone" json:"phone"`
	Si *SimInfo `bson:"si" json:"si"`
	Aid uint32 `bson:"aid" json:"aid"`
}

func sizeRechargeABody(s *RechargeA, bits []byte) (int, error) {
	if s == nil { return 0, fmt.Errorf("sizeRechargeA: nil value") }
	bodySize := 0
	if s.Id != 0 {
		bodySize += 4
		if bits != nil { SetBit(bits, uint8(0), true) }
	}
	if len(s.Type) > 0 {
		fieldSize, err := sizeFixedList(len(s.Type), 1)
		if err != nil { return 0, fmt.Errorf("sizeRechargeA Type: %w", err) }
		bodySize += fieldSize
		if bits != nil { SetBit(bits, uint8(1), true) }
	}
	if len(s.Phone) > 0 {
		fieldSize, err := sizeTextList(s.Phone)
		if err != nil { return 0, fmt.Errorf("sizeRechargeA Phone: %w", err) }
		bodySize += fieldSize
		if bits != nil { SetBit(bits, uint8(2), true) }
	}
	if s.Si != nil {
		fieldSize, err := sizeSimInfo(s.Si)
		if err != nil { return 0, fmt.Errorf("sizeRechargeA Si: %w", err) }
		bodySize += fieldSize
		if bits != nil { SetBit(bits, uint8(3), true) }
	}
	if s.Aid != 0 {
		bodySize += 4
		if bits != nil { SetBit(bits, uint8(4), true) }
	}
	return bodySize, nil
}

func sizeRechargeA(s *RechargeA) (int, error) {
	bodySize, err := sizeRechargeABody(s, nil)
	if err != nil { return 0, err }
	return 1 + bodySize, nil
}

func sizeRechargeAList(v RechargeAList) (int, error) {
	if len(v) > 65535 { return 0, fmt.Errorf("list length exceeds uint16 max") }
	totalSize := 2
	for i, item := range v {
		itemSize, err := sizeRechargeA(item)
		if err != nil { return 0, fmt.Errorf("sizeRechargeAList[%d]: %w", i, err) }
		totalSize += itemSize
	}
	return totalSize, nil
}

func GetRechargeA(buf *bytes.Buffer, s *RechargeA) error {
	if s == nil { return nil }
	const bitSize = 1
	if buf.Len() < bitSize { return fmt.Errorf("GetRechargeA bitmask: %d - %d", buf.Len(), bitSize) }
	bits := buf.Next(bitSize)
	if GetBit(bits, uint8(0)) {
		val, err := GetU32(buf)
		if err != nil { return fmt.Errorf("GetRechargeA Id: %w", err) }
		s.Id = val
	}
	if GetBit(bits, uint8(1)) {
		val, err := GetOrderStatusList(buf)
		if err != nil { return fmt.Errorf("GetRechargeA Type: %w", err) }
		s.Type = val
	}
	if GetBit(bits, uint8(2)) {
		val, err := GetTextList(buf)
		if err != nil { return fmt.Errorf("GetRechargeA Phone: %w", err) }
		s.Phone = val
	}
	if GetBit(bits, uint8(3)) {
		val, err := ReadSimInfo(buf)
		if err != nil { return fmt.Errorf("GetRechargeA Si: %w", err) }
		s.Si = val
	}
	if GetBit(bits, uint8(4)) {
		val, err := GetU32(buf)
		if err != nil { return fmt.Errorf("GetRechargeA Aid: %w", err) }
		s.Aid = val
	}
	if err := ValidateRechargeA(s); err != nil { return fmt.Errorf("ValidateRechargeA: %w", err) }
	return nil
}

func SetRechargeA(buf *bytes.Buffer, s *RechargeA) error {
	if s == nil { return fmt.Errorf("SetRechargeA: nil value") }
	if err := ValidateRechargeA(s); err != nil { return fmt.Errorf("ValidateRechargeA: %w", err) }
	const bitSize = 1
	var bits [1]byte
	bodySize, err := sizeRechargeABody(s, bits[:])
	if err != nil { return err }
	buf.Grow(bitSize + bodySize)
	if _, err := buf.Write(bits[:]); err != nil { return fmt.Errorf("SetRechargeA write bitmask: %w", err) }
	if s.Id != 0 {
		if err := SetU32(buf, s.Id); err != nil { return fmt.Errorf("SetRechargeA Id: %w", err) }
	}
	if len(s.Type) > 0 {
		if err := SetOrderStatusList(buf, (OrderStatusList)(s.Type)); err != nil { return fmt.Errorf("SetRechargeA Type: %w", err) }
	}
	if len(s.Phone) > 0 {
		if err := SetTextList(buf, s.Phone); err != nil { return fmt.Errorf("SetRechargeA Phone: %w", err) }
	}
	if s.Si != nil {
		if err := SetSimInfo(buf, s.Si); err != nil { return fmt.Errorf("SetRechargeA Si: %w", err) }
	}
	if s.Aid != 0 {
		if err := SetU32(buf, s.Aid); err != nil { return fmt.Errorf("SetRechargeA Aid: %w", err) }
	}
	return nil
}

func ValidateRechargeA(s *RechargeA) error {
	if s == nil { return nil }
	for i, item := range s.Type {
		if !IsOrderStatus(item) { return fmt.Errorf("Type[%d] 非法枚举值: %d", i, item) }
	}
	if s.Si != nil {
		if err := ValidateSimInfo(s.Si); err != nil { return fmt.Errorf("Si: %w", err) }
	}
	return nil
}

func EqRechargeA(a, b *RechargeA) bool {
	if a == b { return true }
	if a == nil || b == nil { return false }
	if !EqU32(a.Id, b.Id) { return false }
	if !EqOrderStatusList(a.Type, b.Type) { return false }
	if !EqTextList(a.Phone, b.Phone) { return false }
	if !EqSimInfo(a.Si, b.Si) { return false }
	if !EqU32(a.Aid, b.Aid) { return false }
	return true
}

// Standalone functions
func ReadRechargeA(buf *bytes.Buffer) (*RechargeA, error) {
	s := new(RechargeA)
	return s, GetRechargeA(buf, s)
}

type RechargeAList []*RechargeA
func GetRechargeAList(buf *bytes.Buffer) (RechargeAList, error) {
	count, err := GetU16(buf)
	if err != nil { return nil, err }
	list := make([]*RechargeA, count)
	for i := range list {
		item, err := ReadRechargeA(buf)
		if err != nil { return nil, err }
		list[i] = item
	}
	return RechargeAList(list), nil
}
func SetRechargeAList(buf *bytes.Buffer, v RechargeAList) error {
	totalSize, err := sizeRechargeAList(v)
	if err != nil { return err }
	buf.Grow(totalSize)
	if err := SetU16(buf, uint16(len(v))); err != nil { return err }
	for _, item := range v {
		if err := SetRechargeA(buf, item); err != nil { return err }
	}
	return nil
}
func ValidateRechargeAList(v RechargeAList) error {
	for i, item := range v {
		if item == nil { return fmt.Errorf("RechargeAList[%d]: nil item", i) }
		if err := ValidateRechargeA(item); err != nil { return fmt.Errorf("RechargeAList[%d]: %w", i, err) }
	}
	return nil
}
func EqRechargeAList(a, b RechargeAList) bool { return slices.EqualFunc(a, b, EqRechargeA) }
