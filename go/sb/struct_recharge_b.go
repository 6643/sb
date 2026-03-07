package sb

import (
	"bytes"
	"fmt"
	"slices"
)

type RechargeB struct {
	// abcd
	Id uint32 `bson:"_id" json:"_id"`
	Type []OrderStatus `bson:"type" json:"type"`
	Phone []string `bson:"phone" json:"phone"`
	Si *SimInfo `bson:"si" json:"si"`
	Bid uint32 `bson:"bid" json:"bid"`
}

func sizeRechargeBBody(s *RechargeB, bits []byte) (int, error) {
	if s == nil { return 0, fmt.Errorf("sizeRechargeB: nil value") }
	bodySize := 0
	if s.Id != 0 {
		bodySize += 4
		if bits != nil { SetBit(bits, uint8(0), true) }
	}
	if len(s.Type) > 0 {
		fieldSize, err := sizeFixedList(len(s.Type), 1)
		if err != nil { return 0, fmt.Errorf("sizeRechargeB Type: %w", err) }
		bodySize += fieldSize
		if bits != nil { SetBit(bits, uint8(1), true) }
	}
	if len(s.Phone) > 0 {
		fieldSize, err := sizeTextList(s.Phone)
		if err != nil { return 0, fmt.Errorf("sizeRechargeB Phone: %w", err) }
		bodySize += fieldSize
		if bits != nil { SetBit(bits, uint8(2), true) }
	}
	if s.Si != nil {
		fieldSize, err := sizeSimInfo(s.Si)
		if err != nil { return 0, fmt.Errorf("sizeRechargeB Si: %w", err) }
		bodySize += fieldSize
		if bits != nil { SetBit(bits, uint8(3), true) }
	}
	if s.Bid != 0 {
		bodySize += 4
		if bits != nil { SetBit(bits, uint8(4), true) }
	}
	return bodySize, nil
}

func sizeRechargeB(s *RechargeB) (int, error) {
	bodySize, err := sizeRechargeBBody(s, nil)
	if err != nil { return 0, err }
	return 1 + bodySize, nil
}

func sizeRechargeBList(v RechargeBList) (int, error) {
	if len(v) > 65535 { return 0, fmt.Errorf("list length exceeds uint16 max") }
	totalSize := 2
	for i, item := range v {
		itemSize, err := sizeRechargeB(item)
		if err != nil { return 0, fmt.Errorf("sizeRechargeBList[%d]: %w", i, err) }
		totalSize += itemSize
	}
	return totalSize, nil
}

func GetRechargeB(buf *bytes.Buffer, s *RechargeB) error {
	if s == nil { return nil }
	const bitSize = 1
	if buf.Len() < bitSize { return fmt.Errorf("GetRechargeB bitmask: %d - %d", buf.Len(), bitSize) }
	bits := buf.Next(bitSize)
	if GetBit(bits, uint8(0)) {
		val, err := GetU32(buf)
		if err != nil { return fmt.Errorf("GetRechargeB Id: %w", err) }
		s.Id = val
	}
	if GetBit(bits, uint8(1)) {
		val, err := GetOrderStatusList(buf)
		if err != nil { return fmt.Errorf("GetRechargeB Type: %w", err) }
		s.Type = val
	}
	if GetBit(bits, uint8(2)) {
		val, err := GetTextList(buf)
		if err != nil { return fmt.Errorf("GetRechargeB Phone: %w", err) }
		s.Phone = val
	}
	if GetBit(bits, uint8(3)) {
		val, err := ReadSimInfo(buf)
		if err != nil { return fmt.Errorf("GetRechargeB Si: %w", err) }
		s.Si = val
	}
	if GetBit(bits, uint8(4)) {
		val, err := GetU32(buf)
		if err != nil { return fmt.Errorf("GetRechargeB Bid: %w", err) }
		s.Bid = val
	}
	if err := ValidateRechargeB(s); err != nil { return fmt.Errorf("ValidateRechargeB: %w", err) }
	return nil
}

func SetRechargeB(buf *bytes.Buffer, s *RechargeB) error {
	if s == nil { return fmt.Errorf("SetRechargeB: nil value") }
	if err := ValidateRechargeB(s); err != nil { return fmt.Errorf("ValidateRechargeB: %w", err) }
	const bitSize = 1
	var bits [1]byte
	bodySize, err := sizeRechargeBBody(s, bits[:])
	if err != nil { return err }
	buf.Grow(bitSize + bodySize)
	if _, err := buf.Write(bits[:]); err != nil { return fmt.Errorf("SetRechargeB write bitmask: %w", err) }
	if s.Id != 0 {
		if err := SetU32(buf, s.Id); err != nil { return fmt.Errorf("SetRechargeB Id: %w", err) }
	}
	if len(s.Type) > 0 {
		if err := SetOrderStatusList(buf, (OrderStatusList)(s.Type)); err != nil { return fmt.Errorf("SetRechargeB Type: %w", err) }
	}
	if len(s.Phone) > 0 {
		if err := SetTextList(buf, s.Phone); err != nil { return fmt.Errorf("SetRechargeB Phone: %w", err) }
	}
	if s.Si != nil {
		if err := SetSimInfo(buf, s.Si); err != nil { return fmt.Errorf("SetRechargeB Si: %w", err) }
	}
	if s.Bid != 0 {
		if err := SetU32(buf, s.Bid); err != nil { return fmt.Errorf("SetRechargeB Bid: %w", err) }
	}
	return nil
}

func ValidateRechargeB(s *RechargeB) error {
	if s == nil { return nil }
	for i, item := range s.Type {
		if !IsOrderStatus(item) { return fmt.Errorf("Type[%d] 非法枚举值: %d", i, item) }
	}
	if s.Si != nil {
		if err := ValidateSimInfo(s.Si); err != nil { return fmt.Errorf("Si: %w", err) }
	}
	return nil
}

func EqRechargeB(a, b *RechargeB) bool {
	if a == b { return true }
	if a == nil || b == nil { return false }
	if !EqU32(a.Id, b.Id) { return false }
	if !EqOrderStatusList(a.Type, b.Type) { return false }
	if !EqTextList(a.Phone, b.Phone) { return false }
	if !EqSimInfo(a.Si, b.Si) { return false }
	if !EqU32(a.Bid, b.Bid) { return false }
	return true
}

// Standalone functions
func ReadRechargeB(buf *bytes.Buffer) (*RechargeB, error) {
	s := new(RechargeB)
	return s, GetRechargeB(buf, s)
}

type RechargeBList []*RechargeB
func GetRechargeBList(buf *bytes.Buffer) (RechargeBList, error) {
	count, err := GetU16(buf)
	if err != nil { return nil, err }
	list := make([]*RechargeB, count)
	for i := range list {
		item, err := ReadRechargeB(buf)
		if err != nil { return nil, err }
		list[i] = item
	}
	return RechargeBList(list), nil
}
func SetRechargeBList(buf *bytes.Buffer, v RechargeBList) error {
	totalSize, err := sizeRechargeBList(v)
	if err != nil { return err }
	buf.Grow(totalSize)
	if err := SetU16(buf, uint16(len(v))); err != nil { return err }
	for _, item := range v {
		if err := SetRechargeB(buf, item); err != nil { return err }
	}
	return nil
}
func ValidateRechargeBList(v RechargeBList) error {
	for i, item := range v {
		if item == nil { return fmt.Errorf("RechargeBList[%d]: nil item", i) }
		if err := ValidateRechargeB(item); err != nil { return fmt.Errorf("RechargeBList[%d]: %w", i, err) }
	}
	return nil
}
func EqRechargeBList(a, b RechargeBList) bool { return slices.EqualFunc(a, b, EqRechargeB) }
