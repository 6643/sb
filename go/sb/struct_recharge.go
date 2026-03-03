package sb

import (
	"bytes"
	"fmt"
	"math"
	"slices"
)

type Recharge struct {
	Id uint32 `bson:"_id" json:"_id"` // abcd
	Type []OrderStatus `bson:"type" json:"type"` 
	Phone []string `bson:"phone" json:"phone"` 
	Si *SimInfo `bson:"si" json:"si"` 
}

func GetRecharge(buf *bytes.Buffer, s *Recharge) error {
	if s == nil { return nil }
	if buf.Len() == 0 { return nil }
	bitSize := int(math.Ceil(float64(4) / 8.0))
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
	if s == nil { return nil }
	if err := ValidateRecharge(s); err != nil { return fmt.Errorf("ValidateRecharge: %w", err) }
	bits := make([]byte, uint8(math.Ceil(float64(4)/8.0)))
	body := bytes.NewBuffer(nil)
	if s.Id != 0 {
		if err := SetU32(body, s.Id); err != nil { return fmt.Errorf("SetRecharge Id: %w", err) }
		SetBit(bits, uint8(0), true)
	}
	if len(s.Type) > 0 {
		if err := SetOrderStatusList(body, s.Type); err != nil { return fmt.Errorf("SetRecharge Type: %w", err) }
		SetBit(bits, uint8(1), true)
	}
	if len(s.Phone) > 0 {
		if err := SetTextList(body, s.Phone); err != nil { return fmt.Errorf("SetRecharge Phone: %w", err) }
		SetBit(bits, uint8(2), true)
	}
	if s.Si != nil {
		if err := SetSimInfo(body, s.Si); err != nil { return fmt.Errorf("SetRecharge Si: %w", err) }
		SetBit(bits, uint8(3), true)
	}

	if _, err := buf.Write(bits); err != nil { return fmt.Errorf("SetRecharge write bitmask: %w", err) }
	_, err := body.WriteTo(buf); return err
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
func GetRechargeList(buf *bytes.Buffer) (RechargeList, error) { return getList[*Recharge, RechargeList](buf, ReadRecharge) }
func SetRechargeList(buf *bytes.Buffer, v RechargeList) error { return setList(buf, v, SetRecharge) }
func ValidateRechargeList(v RechargeList) error {
	for i, item := range v {
		if item == nil { continue }
		if err := ValidateRecharge(item); err != nil { return fmt.Errorf("RechargeList[%d]: %w", i, err) }
	}
	return nil
}
func EqRechargeList(a, b RechargeList) bool { return slices.EqualFunc(a, b, EqRecharge) }
