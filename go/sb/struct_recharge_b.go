package sb

import (
	"bytes"
	"fmt"
	"math"
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

func GetRechargeB(buf *bytes.Buffer, s *RechargeB) error {
	if s == nil { return nil }
	bitSize := int(math.Ceil(float64(5) / 8.0))
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
	bits := make([]byte, uint8(math.Ceil(float64(5)/8.0)))
	body := bytes.NewBuffer(nil)
	if s.Id != 0 {
		if err := SetU32(body, s.Id); err != nil { return fmt.Errorf("SetRechargeB Id: %w", err) }
		SetBit(bits, uint8(0), true)
	}
	if len(s.Type) > 0 {
		if err := SetOrderStatusList(body, (OrderStatusList)(s.Type)); err != nil { return fmt.Errorf("SetRechargeB Type: %w", err) }
		SetBit(bits, uint8(1), true)
	}
	if len(s.Phone) > 0 {
		if err := SetTextList(body, s.Phone); err != nil { return fmt.Errorf("SetRechargeB Phone: %w", err) }
		SetBit(bits, uint8(2), true)
	}
	if s.Si != nil {
		if err := SetSimInfo(body, s.Si); err != nil { return fmt.Errorf("SetRechargeB Si: %w", err) }
		SetBit(bits, uint8(3), true)
	}
	if s.Bid != 0 {
		if err := SetU32(body, s.Bid); err != nil { return fmt.Errorf("SetRechargeB Bid: %w", err) }
		SetBit(bits, uint8(4), true)
	}

	if _, err := buf.Write(bits); err != nil { return fmt.Errorf("SetRechargeB write bitmask: %w", err) }
	_, err := body.WriteTo(buf); return err
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
func GetRechargeBList(buf *bytes.Buffer) (RechargeBList, error) { return getList[*RechargeB, RechargeBList](buf, ReadRechargeB) }
func SetRechargeBList(buf *bytes.Buffer, v RechargeBList) error { return setList(buf, v, SetRechargeB) }
func ValidateRechargeBList(v RechargeBList) error {
	for i, item := range v {
		if item == nil { return fmt.Errorf("RechargeBList[%d]: nil item", i) }
		if err := ValidateRechargeB(item); err != nil { return fmt.Errorf("RechargeBList[%d]: %w", i, err) }
	}
	return nil
}
func EqRechargeBList(a, b RechargeBList) bool { return slices.EqualFunc(a, b, EqRechargeB) }
