package sb

import (
	"bytes"
	"fmt"
	"math"
	"slices"
)

type Sim struct {
	// SIM卡ID
	Id uint32 `bson:"_id" json:"_id"`
	Type Type `bson:"type" json:"type"`
	Status ItemStatus `bson:"status" json:"status"`
	// 佣金
	Commission uint16 `bson:"commission" json:"commission"`
	// 供应商ID
	Supplier uint32 `bson:"supplier" json:"supplier"`
	// 推广员ID
	Aff uint32 `bson:"aff" json:"aff"`
	// 合约期(月), 0:长期
	ContractDuration uint8 `bson:"contract_duration" json:"contract_duration"`
	Name string `bson:"name" json:"name"`
	// 运营商
	Operator SimOperator `bson:"operator" json:"operator"`
	// 月租
	Monthly uint16 `bson:"monthly" json:"monthly"`
	// 通用流量
	FlowUniversal uint16 `bson:"flow_universal" json:"flow_universal"`
	// 定向流量
	FlowDirectional uint16 `bson:"flow_directional" json:"flow_directional"`
	// 流量是否结转
	CanMoveFlow bool `bson:"can_move_flow" json:"can_move_flow"`
	// 每月通话(分钟)
	CallMonth uint16 `bson:"call_month" json:"call_month"`
	CallPrice uint16 `bson:"call_price" json:"call_price"`
	// 每月短信(条)
	SmsMonth uint16 `bson:"sms_month" json:"sms_month"`
	SmsPrice uint16 `bson:"sms_price" json:"sms_price"`
	MinAge uint8 `bson:"min_age" json:"min_age"`
	MaxAge uint8 `bson:"max_age" json:"max_age"`
	// 归属地, 0:随机, 1:收货地
	Attribution uint32 `bson:"attribution" json:"attribution"`
	// 选号
	PickPhone []SimPickPhone `bson:"pick_phone" json:"pick_phone"`
	// 首充渠道
	FirstChargeLink string `bson:"first_charge_link" json:"first_charge_link"`
	// 首充金额
	FirstChargeMoney string `bson:"first_charge_money" json:"first_charge_money"`
	// 首充返额
	FirstChargeReturn string `bson:"first_charge_return" json:"first_charge_return"`
	// 禁发区域
	BanCity []uint32 `bson:"ban_city" json:"ban_city"`
	Info []*SimInfo `bson:"info" json:"info"`
	// 套餐截图
	Snapshot []string `bson:"snapshot" json:"snapshot"`
}

func GetSim(buf *bytes.Buffer, s *Sim) error {
	if s == nil { return nil }
	bitSize := int(math.Ceil(float64(27) / 8.0))
	if buf.Len() < bitSize { return fmt.Errorf("GetSim bitmask: %d - %d", buf.Len(), bitSize) }
	bits := buf.Next(bitSize)
	if GetBit(bits, uint8(0)) {
		val, err := GetU32(buf)
		if err != nil { return fmt.Errorf("GetSim Id: %w", err) }
		s.Id = val
	}
	if GetBit(bits, uint8(1)) {
		val, err := GetType(buf)
		if err != nil { return fmt.Errorf("GetSim Type: %w", err) }
		s.Type = val
		if !IsType(s.Type) { return fmt.Errorf("GetSim Type: 非法枚举值: %d", s.Type) }
	}
	if GetBit(bits, uint8(2)) {
		val, err := GetItemStatus(buf)
		if err != nil { return fmt.Errorf("GetSim Status: %w", err) }
		s.Status = val
		if !IsItemStatus(s.Status) { return fmt.Errorf("GetSim Status: 非法枚举值: %d", s.Status) }
	}
	if GetBit(bits, uint8(3)) {
		val, err := GetU16(buf)
		if err != nil { return fmt.Errorf("GetSim Commission: %w", err) }
		s.Commission = val
	}
	if GetBit(bits, uint8(4)) {
		val, err := GetU32(buf)
		if err != nil { return fmt.Errorf("GetSim Supplier: %w", err) }
		s.Supplier = val
	}
	if GetBit(bits, uint8(5)) {
		val, err := GetU32(buf)
		if err != nil { return fmt.Errorf("GetSim Aff: %w", err) }
		s.Aff = val
	}
	if GetBit(bits, uint8(6)) {
		val, err := GetU8(buf)
		if err != nil { return fmt.Errorf("GetSim ContractDuration: %w", err) }
		s.ContractDuration = val
	}
	if GetBit(bits, uint8(7)) {
		val, err := GetText(buf)
		if err != nil { return fmt.Errorf("GetSim Name: %w", err) }
		s.Name = val
	}
	if GetBit(bits, uint8(8)) {
		val, err := GetSimOperator(buf)
		if err != nil { return fmt.Errorf("GetSim Operator: %w", err) }
		s.Operator = val
		if !IsSimOperator(s.Operator) { return fmt.Errorf("GetSim Operator: 非法枚举值: %d", s.Operator) }
	}
	if GetBit(bits, uint8(9)) {
		val, err := GetU16(buf)
		if err != nil { return fmt.Errorf("GetSim Monthly: %w", err) }
		s.Monthly = val
	}
	if GetBit(bits, uint8(10)) {
		val, err := GetU16(buf)
		if err != nil { return fmt.Errorf("GetSim FlowUniversal: %w", err) }
		s.FlowUniversal = val
	}
	if GetBit(bits, uint8(11)) {
		val, err := GetU16(buf)
		if err != nil { return fmt.Errorf("GetSim FlowDirectional: %w", err) }
		s.FlowDirectional = val
	}
	s.CanMoveFlow = GetBit(bits, uint8(12))
	if GetBit(bits, uint8(13)) {
		val, err := GetU16(buf)
		if err != nil { return fmt.Errorf("GetSim CallMonth: %w", err) }
		s.CallMonth = val
	}
	if GetBit(bits, uint8(14)) {
		val, err := GetU16(buf)
		if err != nil { return fmt.Errorf("GetSim CallPrice: %w", err) }
		s.CallPrice = val
	}
	if GetBit(bits, uint8(15)) {
		val, err := GetU16(buf)
		if err != nil { return fmt.Errorf("GetSim SmsMonth: %w", err) }
		s.SmsMonth = val
	}
	if GetBit(bits, uint8(16)) {
		val, err := GetU16(buf)
		if err != nil { return fmt.Errorf("GetSim SmsPrice: %w", err) }
		s.SmsPrice = val
	}
	if GetBit(bits, uint8(17)) {
		val, err := GetU8(buf)
		if err != nil { return fmt.Errorf("GetSim MinAge: %w", err) }
		s.MinAge = val
	}
	if GetBit(bits, uint8(18)) {
		val, err := GetU8(buf)
		if err != nil { return fmt.Errorf("GetSim MaxAge: %w", err) }
		s.MaxAge = val
	}
	if GetBit(bits, uint8(19)) {
		val, err := GetU32(buf)
		if err != nil { return fmt.Errorf("GetSim Attribution: %w", err) }
		s.Attribution = val
	}
	if GetBit(bits, uint8(20)) {
		val, err := GetSimPickPhoneList(buf)
		if err != nil { return fmt.Errorf("GetSim PickPhone: %w", err) }
		s.PickPhone = val
	}
	if GetBit(bits, uint8(21)) {
		val, err := GetText(buf)
		if err != nil { return fmt.Errorf("GetSim FirstChargeLink: %w", err) }
		s.FirstChargeLink = val
	}
	if GetBit(bits, uint8(22)) {
		val, err := GetText(buf)
		if err != nil { return fmt.Errorf("GetSim FirstChargeMoney: %w", err) }
		s.FirstChargeMoney = val
	}
	if GetBit(bits, uint8(23)) {
		val, err := GetText(buf)
		if err != nil { return fmt.Errorf("GetSim FirstChargeReturn: %w", err) }
		s.FirstChargeReturn = val
	}
	if GetBit(bits, uint8(24)) {
		val, err := GetU32List(buf)
		if err != nil { return fmt.Errorf("GetSim BanCity: %w", err) }
		s.BanCity = val
	}
	if GetBit(bits, uint8(25)) {
		val, err := GetSimInfoList(buf)
		if err != nil { return fmt.Errorf("GetSim Info: %w", err) }
		s.Info = val
	}
	if GetBit(bits, uint8(26)) {
		val, err := GetTextList(buf)
		if err != nil { return fmt.Errorf("GetSim Snapshot: %w", err) }
		s.Snapshot = val
	}
	if err := ValidateSim(s); err != nil { return fmt.Errorf("ValidateSim: %w", err) }
	return nil
}

func SetSim(buf *bytes.Buffer, s *Sim) error {
	if s == nil { return fmt.Errorf("SetSim: nil value") }
	if err := ValidateSim(s); err != nil { return fmt.Errorf("ValidateSim: %w", err) }
	bits := make([]byte, uint8(math.Ceil(float64(27)/8.0)))
	body := bytes.NewBuffer(nil)
	if s.Id != 0 {
		if err := SetU32(body, s.Id); err != nil { return fmt.Errorf("SetSim Id: %w", err) }
		SetBit(bits, uint8(0), true)
	}
	if s.Type != 0 {
		if err := SetType(body, s.Type); err != nil { return fmt.Errorf("SetSim Type: %w", err) }
		SetBit(bits, uint8(1), true)
	}
	if s.Status != 0 {
		if err := SetItemStatus(body, s.Status); err != nil { return fmt.Errorf("SetSim Status: %w", err) }
		SetBit(bits, uint8(2), true)
	}
	if s.Commission != 0 {
		if err := SetU16(body, s.Commission); err != nil { return fmt.Errorf("SetSim Commission: %w", err) }
		SetBit(bits, uint8(3), true)
	}
	if s.Supplier != 0 {
		if err := SetU32(body, s.Supplier); err != nil { return fmt.Errorf("SetSim Supplier: %w", err) }
		SetBit(bits, uint8(4), true)
	}
	if s.Aff != 0 {
		if err := SetU32(body, s.Aff); err != nil { return fmt.Errorf("SetSim Aff: %w", err) }
		SetBit(bits, uint8(5), true)
	}
	if s.ContractDuration != 0 {
		if err := SetU8(body, s.ContractDuration); err != nil { return fmt.Errorf("SetSim ContractDuration: %w", err) }
		SetBit(bits, uint8(6), true)
	}
	if s.Name != "" {
		if err := SetText(body, s.Name); err != nil { return fmt.Errorf("SetSim Name: %w", err) }
		SetBit(bits, uint8(7), true)
	}
	if s.Operator != 0 {
		if err := SetSimOperator(body, s.Operator); err != nil { return fmt.Errorf("SetSim Operator: %w", err) }
		SetBit(bits, uint8(8), true)
	}
	if s.Monthly != 0 {
		if err := SetU16(body, s.Monthly); err != nil { return fmt.Errorf("SetSim Monthly: %w", err) }
		SetBit(bits, uint8(9), true)
	}
	if s.FlowUniversal != 0 {
		if err := SetU16(body, s.FlowUniversal); err != nil { return fmt.Errorf("SetSim FlowUniversal: %w", err) }
		SetBit(bits, uint8(10), true)
	}
	if s.FlowDirectional != 0 {
		if err := SetU16(body, s.FlowDirectional); err != nil { return fmt.Errorf("SetSim FlowDirectional: %w", err) }
		SetBit(bits, uint8(11), true)
	}
	SetBit(bits, uint8(12), s.CanMoveFlow)
	if s.CallMonth != 0 {
		if err := SetU16(body, s.CallMonth); err != nil { return fmt.Errorf("SetSim CallMonth: %w", err) }
		SetBit(bits, uint8(13), true)
	}
	if s.CallPrice != 0 {
		if err := SetU16(body, s.CallPrice); err != nil { return fmt.Errorf("SetSim CallPrice: %w", err) }
		SetBit(bits, uint8(14), true)
	}
	if s.SmsMonth != 0 {
		if err := SetU16(body, s.SmsMonth); err != nil { return fmt.Errorf("SetSim SmsMonth: %w", err) }
		SetBit(bits, uint8(15), true)
	}
	if s.SmsPrice != 0 {
		if err := SetU16(body, s.SmsPrice); err != nil { return fmt.Errorf("SetSim SmsPrice: %w", err) }
		SetBit(bits, uint8(16), true)
	}
	if s.MinAge != 0 {
		if err := SetU8(body, s.MinAge); err != nil { return fmt.Errorf("SetSim MinAge: %w", err) }
		SetBit(bits, uint8(17), true)
	}
	if s.MaxAge != 0 {
		if err := SetU8(body, s.MaxAge); err != nil { return fmt.Errorf("SetSim MaxAge: %w", err) }
		SetBit(bits, uint8(18), true)
	}
	if s.Attribution != 0 {
		if err := SetU32(body, s.Attribution); err != nil { return fmt.Errorf("SetSim Attribution: %w", err) }
		SetBit(bits, uint8(19), true)
	}
	if len(s.PickPhone) > 0 {
		if err := SetSimPickPhoneList(body, (SimPickPhoneList)(s.PickPhone)); err != nil { return fmt.Errorf("SetSim PickPhone: %w", err) }
		SetBit(bits, uint8(20), true)
	}
	if s.FirstChargeLink != "" {
		if err := SetText(body, s.FirstChargeLink); err != nil { return fmt.Errorf("SetSim FirstChargeLink: %w", err) }
		SetBit(bits, uint8(21), true)
	}
	if s.FirstChargeMoney != "" {
		if err := SetText(body, s.FirstChargeMoney); err != nil { return fmt.Errorf("SetSim FirstChargeMoney: %w", err) }
		SetBit(bits, uint8(22), true)
	}
	if s.FirstChargeReturn != "" {
		if err := SetText(body, s.FirstChargeReturn); err != nil { return fmt.Errorf("SetSim FirstChargeReturn: %w", err) }
		SetBit(bits, uint8(23), true)
	}
	if len(s.BanCity) > 0 {
		if err := SetU32List(body, s.BanCity); err != nil { return fmt.Errorf("SetSim BanCity: %w", err) }
		SetBit(bits, uint8(24), true)
	}
	if len(s.Info) > 0 {
		if err := SetSimInfoList(body, (SimInfoList)(s.Info)); err != nil { return fmt.Errorf("SetSim Info: %w", err) }
		SetBit(bits, uint8(25), true)
	}
	if len(s.Snapshot) > 0 {
		if err := SetTextList(body, s.Snapshot); err != nil { return fmt.Errorf("SetSim Snapshot: %w", err) }
		SetBit(bits, uint8(26), true)
	}

	if _, err := buf.Write(bits); err != nil { return fmt.Errorf("SetSim write bitmask: %w", err) }
	_, err := body.WriteTo(buf); return err
}

func ValidateSim(s *Sim) error {
	if s == nil { return nil }
	if s.Type != 0 && !IsType(s.Type) { return fmt.Errorf("Type 非法枚举值: %d", s.Type) }
	if s.Status != 0 && !IsItemStatus(s.Status) { return fmt.Errorf("Status 非法枚举值: %d", s.Status) }
	if s.Operator != 0 && !IsSimOperator(s.Operator) { return fmt.Errorf("Operator 非法枚举值: %d", s.Operator) }
	for i, item := range s.PickPhone {
		if !IsSimPickPhone(item) { return fmt.Errorf("PickPhone[%d] 非法枚举值: %d", i, item) }
	}
	if err := ValidateSimInfoList(s.Info); err != nil { return fmt.Errorf("Info: %w", err) }
	return nil
}

func EqSim(a, b *Sim) bool {
	if a == b { return true }
	if a == nil || b == nil { return false }
	if !EqU32(a.Id, b.Id) { return false }
	if !(a.Type == b.Type) { return false }
	if !(a.Status == b.Status) { return false }
	if !EqU16(a.Commission, b.Commission) { return false }
	if !EqU32(a.Supplier, b.Supplier) { return false }
	if !EqU32(a.Aff, b.Aff) { return false }
	if !EqU8(a.ContractDuration, b.ContractDuration) { return false }
	if !EqText(a.Name, b.Name) { return false }
	if !(a.Operator == b.Operator) { return false }
	if !EqU16(a.Monthly, b.Monthly) { return false }
	if !EqU16(a.FlowUniversal, b.FlowUniversal) { return false }
	if !EqU16(a.FlowDirectional, b.FlowDirectional) { return false }
	if !EqBool(a.CanMoveFlow, b.CanMoveFlow) { return false }
	if !EqU16(a.CallMonth, b.CallMonth) { return false }
	if !EqU16(a.CallPrice, b.CallPrice) { return false }
	if !EqU16(a.SmsMonth, b.SmsMonth) { return false }
	if !EqU16(a.SmsPrice, b.SmsPrice) { return false }
	if !EqU8(a.MinAge, b.MinAge) { return false }
	if !EqU8(a.MaxAge, b.MaxAge) { return false }
	if !EqU32(a.Attribution, b.Attribution) { return false }
	if !EqSimPickPhoneList(a.PickPhone, b.PickPhone) { return false }
	if !EqText(a.FirstChargeLink, b.FirstChargeLink) { return false }
	if !EqText(a.FirstChargeMoney, b.FirstChargeMoney) { return false }
	if !EqText(a.FirstChargeReturn, b.FirstChargeReturn) { return false }
	if !EqU32List(a.BanCity, b.BanCity) { return false }
	if !EqSimInfoList((SimInfoList)(a.Info), (SimInfoList)(b.Info)) { return false }
	if !EqTextList(a.Snapshot, b.Snapshot) { return false }
	return true
}

// Standalone functions
func ReadSim(buf *bytes.Buffer) (*Sim, error) {
	s := new(Sim)
	return s, GetSim(buf, s)
}

type SimList []*Sim
func GetSimList(buf *bytes.Buffer) (SimList, error) { return getList[*Sim, SimList](buf, ReadSim) }
func SetSimList(buf *bytes.Buffer, v SimList) error { return setList(buf, v, SetSim) }
func ValidateSimList(v SimList) error {
	for i, item := range v {
		if item == nil { return fmt.Errorf("SimList[%d]: nil item", i) }
		if err := ValidateSim(item); err != nil { return fmt.Errorf("SimList[%d]: %w", i, err) }
	}
	return nil
}
func EqSimList(a, b SimList) bool { return slices.EqualFunc(a, b, EqSim) }
