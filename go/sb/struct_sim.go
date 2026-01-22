package sb

import (
	"bytes"
	"fmt"
	"math"
	"slices"
	"unsafe"
)

type Sim struct {
	Id uint32  // SIM卡ID
	Type Type  
	Status ItemStatus  
	Commission uint16  // 佣金
	Supplier uint32  // 供应商ID
	Aff uint32  // 推广员ID
	ContractDuration uint8  // 合约期(月), 0:长期
	Name string  
	Operator SimOperator  // 运营商
	Monthly uint16  // 月租
	FlowUniversal uint16  // 通用流量
	FlowDirectional uint16  // 定向流量
	CanMoveFlow bool  // 流量是否结转
	CallMonth uint16  // 每月通话(分钟)
	CallPrice uint16  
	SmsMonth uint16  // 每月短信(条)
	SmsPrice uint16  
	MinAge uint8  
	MaxAge uint8  
	Attribution uint32  // 归属地, 0:随机, 1:收货地
	PickPhone []SimPickPhone  // 选号
	FirstChargeLink string  // 首充渠道
	FirstChargeMoney string  // 首充金额
	FirstChargeReturn string  // 首充返额
	BanCity []uint32  // 禁发区域
	Info []*SimInfo  
	Snapshot []string  // 套餐截图
}

func (s *Sim) Get(buf *bytes.Buffer) error {
	if buf.Len() == 0 { return nil }
	bitSize := int(math.Ceil(float64(27) / 8.0))
	if buf.Len() < bitSize { return fmt.Errorf("GetSim bitmask: %d - %d", buf.Len(), bitSize) }
	bits := buf.Next(bitSize)
	if GetBit(bits, uint8(0)) {
		val, err := GetU32(buf)
		if err != nil { return fmt.Errorf("GetSim Id: %w", err) }
		s.Id = val
	}
	if GetBit(bits, uint8(1)) {
		val, err := GetU8(buf)
		if err != nil { return fmt.Errorf("GetSim Type: %w", err) }
		s.Type = Type(val)
	}
	if GetBit(bits, uint8(2)) {
		val, err := GetU8(buf)
		if err != nil { return fmt.Errorf("GetSim Status: %w", err) }
		s.Status = ItemStatus(val)
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
		val, err := GetU8(buf)
		if err != nil { return fmt.Errorf("GetSim Operator: %w", err) }
		s.Operator = SimOperator(val)
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
		val, err := GetU8List(buf)
		if err != nil { return fmt.Errorf("GetSim PickPhone: %w", err) }
		s.PickPhone = *(*[]SimPickPhone)(unsafe.Pointer(&val))
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
		var val SimInfoList
		if err := val.Get(buf); err != nil { return fmt.Errorf("GetSim Info: %w", err) }
		s.Info = val
	}
	if GetBit(bits, uint8(26)) {
		val, err := GetTextList(buf)
		if err != nil { return fmt.Errorf("GetSim Snapshot: %w", err) }
		s.Snapshot = val
	}
	return nil
}

func (s *Sim) Set(buf *bytes.Buffer) error {
	if s == nil { return nil }
	bits := make([]byte, uint8(math.Ceil(float64(27)/8.0)))
	body := bytes.NewBuffer(nil)
	if s.Id != 0 {
		if err := SetU32(body, s.Id); err != nil { return fmt.Errorf("SetSim Id: %w", err) }
		SetBit(bits, uint8(0), true)
	}
	if s.Type != 0 {
		if err := SetU8(body, uint8(s.Type)); err != nil { return fmt.Errorf("SetSim Type: %w", err) }
		SetBit(bits, uint8(1), true)
	}
	if s.Status != 0 {
		if err := SetU8(body, uint8(s.Status)); err != nil { return fmt.Errorf("SetSim Status: %w", err) }
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
		if err := SetU8(body, uint8(s.Operator)); err != nil { return fmt.Errorf("SetSim Operator: %w", err) }
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
		if err := SetU8List(body, *(*[]uint8)(unsafe.Pointer(&s.PickPhone))); err != nil { return fmt.Errorf("SetSim PickPhone: %w", err) }
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
		if err := (SimInfoList)(s.Info).Set(body); err != nil { return fmt.Errorf("SetSim Info: %w", err) }
		SetBit(bits, uint8(25), true)
	}
	if len(s.Snapshot) > 0 {
		if err := SetTextList(body, s.Snapshot); err != nil { return fmt.Errorf("SetSim Snapshot: %w", err) }
		SetBit(bits, uint8(26), true)
	}

	if _, err := buf.Write(bits); err != nil { return fmt.Errorf("SetSim write bitmask: %w", err) }
	_, err := body.WriteTo(buf); return err
}

func (s *Sim) Eq(other *Sim) bool {
	if s == other { return true }
	if s == nil || other == nil { return false }
	if !EqU32(s.Id, other.Id) { return false }
	if s.Type != other.Type { return false }
	if s.Status != other.Status { return false }
	if !EqU16(s.Commission, other.Commission) { return false }
	if !EqU32(s.Supplier, other.Supplier) { return false }
	if !EqU32(s.Aff, other.Aff) { return false }
	if !EqU8(s.ContractDuration, other.ContractDuration) { return false }
	if !EqText(s.Name, other.Name) { return false }
	if s.Operator != other.Operator { return false }
	if !EqU16(s.Monthly, other.Monthly) { return false }
	if !EqU16(s.FlowUniversal, other.FlowUniversal) { return false }
	if !EqU16(s.FlowDirectional, other.FlowDirectional) { return false }
	if !EqBool(s.CanMoveFlow, other.CanMoveFlow) { return false }
	if !EqU16(s.CallMonth, other.CallMonth) { return false }
	if !EqU16(s.CallPrice, other.CallPrice) { return false }
	if !EqU16(s.SmsMonth, other.SmsMonth) { return false }
	if !EqU16(s.SmsPrice, other.SmsPrice) { return false }
	if !EqU8(s.MinAge, other.MinAge) { return false }
	if !EqU8(s.MaxAge, other.MaxAge) { return false }
	if !EqU32(s.Attribution, other.Attribution) { return false }
	if !slices.Equal(s.PickPhone, other.PickPhone) { return false }
	if !EqText(s.FirstChargeLink, other.FirstChargeLink) { return false }
	if !EqText(s.FirstChargeMoney, other.FirstChargeMoney) { return false }
	if !EqText(s.FirstChargeReturn, other.FirstChargeReturn) { return false }
	if !EqU32List(s.BanCity, other.BanCity) { return false }
	if !(SimInfoList)(s.Info).Eq((SimInfoList)(other.Info)) { return false }
	if !EqTextList(s.Snapshot, other.Snapshot) { return false }
	return true
}

// Standalone functions for compatibility
func GetSim(buf *bytes.Buffer) (*Sim, error) {
	s := new(Sim); return s, s.Get(buf)
}
func SetSim(buf *bytes.Buffer, s *Sim) error { return s.Set(buf) }
func EqSim(a, b *Sim) bool { return a.Eq(b) }

type SimList []*Sim
func (v SimList) Set(buf *bytes.Buffer) error { return setList(buf, v, SetSim) }
func (v *SimList) Get(buf *bytes.Buffer) error {
	val, err := getList[*Sim, SimList](buf, GetSim)
	if err == nil { *v = val }; return err
}
func (v SimList) Eq(other SimList) bool { return slices.EqualFunc(v, other, EqSim) }
