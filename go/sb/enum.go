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

func IsAccountStatusValid(v AccountStatus) bool {
	switch v {
	case AccountStatusOffline:
		return true
	case AccountStatusOnline:
		return true
	case AccountStatusDeleted:
		return true
	default:
		return false
	}
}

type AccountStatusList []AccountStatus
func (v AccountStatusList) Set(buf *bytes.Buffer) error { return SetU8List(buf, *(*[]uint8)(unsafe.Pointer(&v))) }
func (v *AccountStatusList) Get(buf *bytes.Buffer) error {
	val, err := GetU8List(buf)
	if err == nil { *v = *(*AccountStatusList)(unsafe.Pointer(&val)) }
	return err
}
func IsAccountStatusListValid(v AccountStatusList) bool {
	for _, item := range v {
		if !IsAccountStatusValid(item) { return false }
	}
	return true
}
func (v AccountStatusList) Eq(other AccountStatusList) bool { return slices.Equal(v, other) }

// Type 类型
type Type uint8

const (
	TypeSim Type = 0 
	TypeRecharge Type = 1 
)

func IsTypeValid(v Type) bool {
	switch v {
	case TypeSim:
		return true
	case TypeRecharge:
		return true
	default:
		return false
	}
}

type TypeList []Type
func (v TypeList) Set(buf *bytes.Buffer) error { return SetU8List(buf, *(*[]uint8)(unsafe.Pointer(&v))) }
func (v *TypeList) Get(buf *bytes.Buffer) error {
	val, err := GetU8List(buf)
	if err == nil { *v = *(*TypeList)(unsafe.Pointer(&val)) }
	return err
}
func IsTypeListValid(v TypeList) bool {
	for _, item := range v {
		if !IsTypeValid(item) { return false }
	}
	return true
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

func IsStatusValid(v Status) bool {
	switch v {
	case StatusOk:
		return true
	case StatusErr:
		return true
	case StatusTwo:
		return true
	case StatusThree:
		return true
	case StatusFour:
		return true
	case StatusFive:
		return true
	case StatusSix:
		return true
	case StatusSeven:
		return true
	case StatusOne:
		return true
	default:
		return false
	}
}

type StatusList []Status
func (v StatusList) Set(buf *bytes.Buffer) error { return SetU8List(buf, *(*[]uint8)(unsafe.Pointer(&v))) }
func (v *StatusList) Get(buf *bytes.Buffer) error {
	val, err := GetU8List(buf)
	if err == nil { *v = *(*StatusList)(unsafe.Pointer(&val)) }
	return err
}
func IsStatusListValid(v StatusList) bool {
	for _, item := range v {
		if !IsStatusValid(item) { return false }
	}
	return true
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

func IsStatusAValid(v StatusA) bool {
	switch v {
	case StatusAOk:
		return true
	case StatusAOne:
		return true
	case StatusATwo:
		return true
	case StatusAThree:
		return true
	case StatusAFour:
		return true
	case StatusAFive:
		return true
	case StatusASix:
		return true
	case StatusASeven:
		return true
	default:
		return false
	}
}

type StatusAList []StatusA
func (v StatusAList) Set(buf *bytes.Buffer) error { return SetU8List(buf, *(*[]uint8)(unsafe.Pointer(&v))) }
func (v *StatusAList) Get(buf *bytes.Buffer) error {
	val, err := GetU8List(buf)
	if err == nil { *v = *(*StatusAList)(unsafe.Pointer(&val)) }
	return err
}
func IsStatusAListValid(v StatusAList) bool {
	for _, item := range v {
		if !IsStatusAValid(item) { return false }
	}
	return true
}
func (v StatusAList) Eq(other StatusAList) bool { return slices.Equal(v, other) }

// ItemStatus 订单状态
type ItemStatus uint8

const (
	ItemStatusOffline ItemStatus = 0 
	ItemStatusOnline ItemStatus = 1 
)

func IsItemStatusValid(v ItemStatus) bool {
	switch v {
	case ItemStatusOffline:
		return true
	case ItemStatusOnline:
		return true
	default:
		return false
	}
}

type ItemStatusList []ItemStatus
func (v ItemStatusList) Set(buf *bytes.Buffer) error { return SetU8List(buf, *(*[]uint8)(unsafe.Pointer(&v))) }
func (v *ItemStatusList) Get(buf *bytes.Buffer) error {
	val, err := GetU8List(buf)
	if err == nil { *v = *(*ItemStatusList)(unsafe.Pointer(&val)) }
	return err
}
func IsItemStatusListValid(v ItemStatusList) bool {
	for _, item := range v {
		if !IsItemStatusValid(item) { return false }
	}
	return true
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

func IsSimPickPhoneValid(v SimPickPhone) bool {
	switch v {
	case SimPickPhoneNo:
		return true
	case SimPickPhoneYes:
		return true
	case SimPickPhoneActive:
		return true
	case SimPickPhoneAbcc:
		return true
	default:
		return false
	}
}

type SimPickPhoneList []SimPickPhone
func (v SimPickPhoneList) Set(buf *bytes.Buffer) error { return SetU8List(buf, *(*[]uint8)(unsafe.Pointer(&v))) }
func (v *SimPickPhoneList) Get(buf *bytes.Buffer) error {
	val, err := GetU8List(buf)
	if err == nil { *v = *(*SimPickPhoneList)(unsafe.Pointer(&val)) }
	return err
}
func IsSimPickPhoneListValid(v SimPickPhoneList) bool {
	for _, item := range v {
		if !IsSimPickPhoneValid(item) { return false }
	}
	return true
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

func IsSimOperatorValid(v SimOperator) bool {
	switch v {
	case SimOperatorZz:
		return true
	case SimOperatorLt:
		return true
	case SimOperatorYd:
		return true
	case SimOperatorDx:
		return true
	case SimOperatorGd:
		return true
	case SimOperatorXx:
		return true
	case SimOperatorA:
		return true
	case SimOperatorB:
		return true
	default:
		return false
	}
}

type SimOperatorList []SimOperator
func (v SimOperatorList) Set(buf *bytes.Buffer) error { return SetU8List(buf, *(*[]uint8)(unsafe.Pointer(&v))) }
func (v *SimOperatorList) Get(buf *bytes.Buffer) error {
	val, err := GetU8List(buf)
	if err == nil { *v = *(*SimOperatorList)(unsafe.Pointer(&val)) }
	return err
}
func IsSimOperatorListValid(v SimOperatorList) bool {
	for _, item := range v {
		if !IsSimOperatorValid(item) { return false }
	}
	return true
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

func IsOrderStatusValid(v OrderStatus) bool {
	switch v {
	case OrderStatusPending:
		return true
	case OrderStatusClosed:
		return true
	case OrderStatusCanceled:
		return true
	case OrderStatusShipped:
		return true
	case OrderStatusDelivered:
		return true
	case OrderStatusActived:
		return true
	case OrderStatusSettled:
		return true
	default:
		return false
	}
}

type OrderStatusList []OrderStatus
func (v OrderStatusList) Set(buf *bytes.Buffer) error { return SetU8List(buf, *(*[]uint8)(unsafe.Pointer(&v))) }
func (v *OrderStatusList) Get(buf *bytes.Buffer) error {
	val, err := GetU8List(buf)
	if err == nil { *v = *(*OrderStatusList)(unsafe.Pointer(&val)) }
	return err
}
func IsOrderStatusListValid(v OrderStatusList) bool {
	for _, item := range v {
		if !IsOrderStatusValid(item) { return false }
	}
	return true
}
func (v OrderStatusList) Eq(other OrderStatusList) bool { return slices.Equal(v, other) }

