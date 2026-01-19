package sb

import (
	"bytes"
	"slices"
	"unsafe"
)


// AccountStatus 账户状态
type AccountStatus uint8

const (
	AccountStatusOffline AccountStatus = 0 
	AccountStatusOnline AccountStatus = 1 
	AccountStatusDeleted AccountStatus = 2 
)

type AccountStatusList []AccountStatus
func (v AccountStatusList) Set(buf *bytes.Buffer) error { return SetU8List(buf, *(*[]uint8)(unsafe.Pointer(&v))) }
func (v *AccountStatusList) Get(buf *bytes.Buffer) error {
	val, err := GetU8List(buf)
	if err == nil { *v = *(*AccountStatusList)(unsafe.Pointer(&val)) }
	return err
}
func (v AccountStatusList) Eq(other AccountStatusList) bool { return slices.Equal(v, other) }

// Type 类型
type Type uint8

const (
	TypeSim Type = 0 
	TypeRecharge Type = 1 
)

type TypeList []Type
func (v TypeList) Set(buf *bytes.Buffer) error { return SetU8List(buf, *(*[]uint8)(unsafe.Pointer(&v))) }
func (v *TypeList) Get(buf *bytes.Buffer) error {
	val, err := GetU8List(buf)
	if err == nil { *v = *(*TypeList)(unsafe.Pointer(&val)) }
	return err
}
func (v TypeList) Eq(other TypeList) bool { return slices.Equal(v, other) }

// Status 错误码
type Status uint8

const (
	StatusOk Status = 0 
	StatusErr Status = 1 
	StatusTwo Status = 2 
	StatusThree Status = 3 
	StatusFour Status = 4 
	StatusFive Status = 5 
	StatusSix Status = 6 
	StatusSeven Status = 7 
	StatusOne Status = 11 
)

type StatusList []Status
func (v StatusList) Set(buf *bytes.Buffer) error { return SetU8List(buf, *(*[]uint8)(unsafe.Pointer(&v))) }
func (v *StatusList) Get(buf *bytes.Buffer) error {
	val, err := GetU8List(buf)
	if err == nil { *v = *(*StatusList)(unsafe.Pointer(&val)) }
	return err
}
func (v StatusList) Eq(other StatusList) bool { return slices.Equal(v, other) }

// StatusA 状态A
type StatusA uint8

const (
	StatusAOk StatusA = 0 
	StatusAOne StatusA = 1 
	StatusATwo StatusA = 2 
	StatusAThree StatusA = 3 
	StatusAFour StatusA = 4 
	StatusAFive StatusA = 5 
	StatusASix StatusA = 6 
	StatusASeven StatusA = 7 
)

type StatusAList []StatusA
func (v StatusAList) Set(buf *bytes.Buffer) error { return SetU8List(buf, *(*[]uint8)(unsafe.Pointer(&v))) }
func (v *StatusAList) Get(buf *bytes.Buffer) error {
	val, err := GetU8List(buf)
	if err == nil { *v = *(*StatusAList)(unsafe.Pointer(&val)) }
	return err
}
func (v StatusAList) Eq(other StatusAList) bool { return slices.Equal(v, other) }

// ItemStatus 订单状态
type ItemStatus uint8

const (
	ItemStatusOffline ItemStatus = 0 
	ItemStatusOnline ItemStatus = 1 
)

type ItemStatusList []ItemStatus
func (v ItemStatusList) Set(buf *bytes.Buffer) error { return SetU8List(buf, *(*[]uint8)(unsafe.Pointer(&v))) }
func (v *ItemStatusList) Get(buf *bytes.Buffer) error {
	val, err := GetU8List(buf)
	if err == nil { *v = *(*ItemStatusList)(unsafe.Pointer(&val)) }
	return err
}
func (v ItemStatusList) Eq(other ItemStatusList) bool { return slices.Equal(v, other) }

// SimPickPhone 可否选号
type SimPickPhone uint8

const (
	SimPickPhoneNo SimPickPhone = 0 
	SimPickPhoneYes SimPickPhone = 1 
	SimPickPhoneActive SimPickPhone = 3 
	SimPickPhoneAbcc SimPickPhone = 4 
)

type SimPickPhoneList []SimPickPhone
func (v SimPickPhoneList) Set(buf *bytes.Buffer) error { return SetU8List(buf, *(*[]uint8)(unsafe.Pointer(&v))) }
func (v *SimPickPhoneList) Get(buf *bytes.Buffer) error {
	val, err := GetU8List(buf)
	if err == nil { *v = *(*SimPickPhoneList)(unsafe.Pointer(&val)) }
	return err
}
func (v SimPickPhoneList) Eq(other SimPickPhoneList) bool { return slices.Equal(v, other) }

// SimOperator 运营商
type SimOperator uint8

const (
	SimOperatorZz SimOperator = 2 
	SimOperatorLt SimOperator = 3 
	SimOperatorYd SimOperator = 4 
	SimOperatorDx SimOperator = 5 
	SimOperatorGd SimOperator = 6 
	SimOperatorXx SimOperator = 7 
	SimOperatorA SimOperator = 11 
	SimOperatorB SimOperator = 12 
)

type SimOperatorList []SimOperator
func (v SimOperatorList) Set(buf *bytes.Buffer) error { return SetU8List(buf, *(*[]uint8)(unsafe.Pointer(&v))) }
func (v *SimOperatorList) Get(buf *bytes.Buffer) error {
	val, err := GetU8List(buf)
	if err == nil { *v = *(*SimOperatorList)(unsafe.Pointer(&val)) }
	return err
}
func (v SimOperatorList) Eq(other SimOperatorList) bool { return slices.Equal(v, other) }

// OrderStatus 订单状态
type OrderStatus uint8

const (
	OrderStatusPending OrderStatus = 0 // 待处理
	OrderStatusClosed OrderStatus = 1 // 已关闭
	OrderStatusCanceled OrderStatus = 2 // 已取消
	OrderStatusShipped OrderStatus = 3 // 已发货
	OrderStatusDelivered OrderStatus = 4 // 已送达
	OrderStatusActived OrderStatus = 5 // 已激活
	OrderStatusSettled OrderStatus = 6 // 已结算
)

type OrderStatusList []OrderStatus
func (v OrderStatusList) Set(buf *bytes.Buffer) error { return SetU8List(buf, *(*[]uint8)(unsafe.Pointer(&v))) }
func (v *OrderStatusList) Get(buf *bytes.Buffer) error {
	val, err := GetU8List(buf)
	if err == nil { *v = *(*OrderStatusList)(unsafe.Pointer(&val)) }
	return err
}
func (v OrderStatusList) Eq(other OrderStatusList) bool { return slices.Equal(v, other) }
