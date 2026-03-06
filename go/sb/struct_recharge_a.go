package sb

import (
	"bytes"
	"fmt"
	"math"
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

func GetRechargeA(buf *bytes.Buffer, s *RechargeA) error {
	if s == nil { return nil }
	bitSize := int(math.Ceil(float64(5) / 8.0))
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
	bits := make([]byte, uint8(math.Ceil(float64(5)/8.0)))
	body := bytes.NewBuffer(nil)
	if s.Id != 0 {
		if err := SetU32(body, s.Id); err != nil { return fmt.Errorf("SetRechargeA Id: %w", err) }
		SetBit(bits, uint8(0), true)
	}
	if len(s.Type) > 0 {
		if err := SetOrderStatusList(body, (OrderStatusList)(s.Type)); err != nil { return fmt.Errorf("SetRechargeA Type: %w", err) }
		SetBit(bits, uint8(1), true)
	}
	if len(s.Phone) > 0 {
		if err := SetTextList(body, s.Phone); err != nil { return fmt.Errorf("SetRechargeA Phone: %w", err) }
		SetBit(bits, uint8(2), true)
	}
	if s.Si != nil {
		if err := SetSimInfo(body, s.Si); err != nil { return fmt.Errorf("SetRechargeA Si: %w", err) }
		SetBit(bits, uint8(3), true)
	}
	if s.Aid != 0 {
		if err := SetU32(body, s.Aid); err != nil { return fmt.Errorf("SetRechargeA Aid: %w", err) }
		SetBit(bits, uint8(4), true)
	}

	if _, err := buf.Write(bits); err != nil { return fmt.Errorf("SetRechargeA write bitmask: %w", err) }
	_, err := body.WriteTo(buf); return err
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
func GetRechargeAList(buf *bytes.Buffer) (RechargeAList, error) { return getList[*RechargeA, RechargeAList](buf, ReadRechargeA) }
func SetRechargeAList(buf *bytes.Buffer, v RechargeAList) error { return setList(buf, v, SetRechargeA) }
func ValidateRechargeAList(v RechargeAList) error {
	for i, item := range v {
		if item == nil { return fmt.Errorf("RechargeAList[%d]: nil item", i) }
		if err := ValidateRechargeA(item); err != nil { return fmt.Errorf("RechargeAList[%d]: %w", i, err) }
	}
	return nil
}
func EqRechargeAList(a, b RechargeAList) bool { return slices.EqualFunc(a, b, EqRechargeA) }
