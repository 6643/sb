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

func IsAccountStatus(v AccountStatus) bool {
	switch v {
	case AccountStatusOffline, AccountStatusOnline, AccountStatusDeleted:
		return true
	default:
		return false
	}
}

func GetAccountStatus(buf *bytes.Buffer) (AccountStatus, error) {
	val, err := GetU8(buf)
	return AccountStatus(val), err
}

func SetAccountStatus(buf *bytes.Buffer, v AccountStatus) error { return SetU8(buf, uint8(v)) }

type AccountStatusList []AccountStatus
func GetAccountStatusList(buf *bytes.Buffer) (AccountStatusList, error) {
	val, err := GetU8List(buf)
	if err != nil { return nil, err }
	return *(*AccountStatusList)(unsafe.Pointer(&val)), nil
}
func SetAccountStatusList(buf *bytes.Buffer, v AccountStatusList) error { return SetU8List(buf, *(*[]uint8)(unsafe.Pointer(&v))) }
func IsAccountStatusList(v AccountStatusList) bool {
	for _, item := range v {
		if !IsAccountStatus(item) { return false }
	}
	return true
}
func EqAccountStatusList(a, b AccountStatusList) bool { return slices.Equal(a, b) }

// Type 类型
type Type uint8

const (
	TypeSim Type = 0 
	TypeRecharge Type = 1 
)

func IsType(v Type) bool {
	switch v {
	case TypeSim, TypeRecharge:
		return true
	default:
		return false
	}
}

func GetType(buf *bytes.Buffer) (Type, error) {
	val, err := GetU8(buf)
	return Type(val), err
}

func SetType(buf *bytes.Buffer, v Type) error { return SetU8(buf, uint8(v)) }

type TypeList []Type
func GetTypeList(buf *bytes.Buffer) (TypeList, error) {
	val, err := GetU8List(buf)
	if err != nil { return nil, err }
	return *(*TypeList)(unsafe.Pointer(&val)), nil
}
func SetTypeList(buf *bytes.Buffer, v TypeList) error { return SetU8List(buf, *(*[]uint8)(unsafe.Pointer(&v))) }
func IsTypeList(v TypeList) bool {
	for _, item := range v {
		if !IsType(item) { return false }
	}
	return true
}
func EqTypeList(a, b TypeList) bool { return slices.Equal(a, b) }

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

func IsStatus(v Status) bool {
	switch v {
	case StatusOk, StatusErr, StatusTwo, StatusThree, StatusFour, StatusFive, StatusSix, StatusSeven, StatusOne:
		return true
	default:
		return false
	}
}

func GetStatus(buf *bytes.Buffer) (Status, error) {
	val, err := GetU8(buf)
	return Status(val), err
}

func SetStatus(buf *bytes.Buffer, v Status) error { return SetU8(buf, uint8(v)) }

type StatusList []Status
func GetStatusList(buf *bytes.Buffer) (StatusList, error) {
	val, err := GetU8List(buf)
	if err != nil { return nil, err }
	return *(*StatusList)(unsafe.Pointer(&val)), nil
}
func SetStatusList(buf *bytes.Buffer, v StatusList) error { return SetU8List(buf, *(*[]uint8)(unsafe.Pointer(&v))) }
func IsStatusList(v StatusList) bool {
	for _, item := range v {
		if !IsStatus(item) { return false }
	}
	return true
}
func EqStatusList(a, b StatusList) bool { return slices.Equal(a, b) }

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

func IsStatusA(v StatusA) bool {
	switch v {
	case StatusAOk, StatusAOne, StatusATwo, StatusAThree, StatusAFour, StatusAFive, StatusASix, StatusASeven:
		return true
	default:
		return false
	}
}

func GetStatusA(buf *bytes.Buffer) (StatusA, error) {
	val, err := GetU8(buf)
	return StatusA(val), err
}

func SetStatusA(buf *bytes.Buffer, v StatusA) error { return SetU8(buf, uint8(v)) }

type StatusAList []StatusA
func GetStatusAList(buf *bytes.Buffer) (StatusAList, error) {
	val, err := GetU8List(buf)
	if err != nil { return nil, err }
	return *(*StatusAList)(unsafe.Pointer(&val)), nil
}
func SetStatusAList(buf *bytes.Buffer, v StatusAList) error { return SetU8List(buf, *(*[]uint8)(unsafe.Pointer(&v))) }
func IsStatusAList(v StatusAList) bool {
	for _, item := range v {
		if !IsStatusA(item) { return false }
	}
	return true
}
func EqStatusAList(a, b StatusAList) bool { return slices.Equal(a, b) }

// ItemStatus 订单状态
type ItemStatus uint8

const (
	ItemStatusOffline ItemStatus = 0 
	ItemStatusOnline ItemStatus = 1 
)

func IsItemStatus(v ItemStatus) bool {
	switch v {
	case ItemStatusOffline, ItemStatusOnline:
		return true
	default:
		return false
	}
}

func GetItemStatus(buf *bytes.Buffer) (ItemStatus, error) {
	val, err := GetU8(buf)
	return ItemStatus(val), err
}

func SetItemStatus(buf *bytes.Buffer, v ItemStatus) error { return SetU8(buf, uint8(v)) }

type ItemStatusList []ItemStatus
func GetItemStatusList(buf *bytes.Buffer) (ItemStatusList, error) {
	val, err := GetU8List(buf)
	if err != nil { return nil, err }
	return *(*ItemStatusList)(unsafe.Pointer(&val)), nil
}
func SetItemStatusList(buf *bytes.Buffer, v ItemStatusList) error { return SetU8List(buf, *(*[]uint8)(unsafe.Pointer(&v))) }
func IsItemStatusList(v ItemStatusList) bool {
	for _, item := range v {
		if !IsItemStatus(item) { return false }
	}
	return true
}
func EqItemStatusList(a, b ItemStatusList) bool { return slices.Equal(a, b) }

// SimPickPhone 可否选号
type SimPickPhone uint8

const (
	SimPickPhoneNo SimPickPhone = 0 
	SimPickPhoneYes SimPickPhone = 1 
	SimPickPhoneActive SimPickPhone = 3 
	SimPickPhoneAbcc SimPickPhone = 4 
)

func IsSimPickPhone(v SimPickPhone) bool {
	switch v {
	case SimPickPhoneNo, SimPickPhoneYes, SimPickPhoneActive, SimPickPhoneAbcc:
		return true
	default:
		return false
	}
}

func GetSimPickPhone(buf *bytes.Buffer) (SimPickPhone, error) {
	val, err := GetU8(buf)
	return SimPickPhone(val), err
}

func SetSimPickPhone(buf *bytes.Buffer, v SimPickPhone) error { return SetU8(buf, uint8(v)) }

type SimPickPhoneList []SimPickPhone
func GetSimPickPhoneList(buf *bytes.Buffer) (SimPickPhoneList, error) {
	val, err := GetU8List(buf)
	if err != nil { return nil, err }
	return *(*SimPickPhoneList)(unsafe.Pointer(&val)), nil
}
func SetSimPickPhoneList(buf *bytes.Buffer, v SimPickPhoneList) error { return SetU8List(buf, *(*[]uint8)(unsafe.Pointer(&v))) }
func IsSimPickPhoneList(v SimPickPhoneList) bool {
	for _, item := range v {
		if !IsSimPickPhone(item) { return false }
	}
	return true
}
func EqSimPickPhoneList(a, b SimPickPhoneList) bool { return slices.Equal(a, b) }

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

func IsSimOperator(v SimOperator) bool {
	switch v {
	case SimOperatorZz, SimOperatorLt, SimOperatorYd, SimOperatorDx, SimOperatorGd, SimOperatorXx, SimOperatorA, SimOperatorB:
		return true
	default:
		return false
	}
}

func GetSimOperator(buf *bytes.Buffer) (SimOperator, error) {
	val, err := GetU8(buf)
	return SimOperator(val), err
}

func SetSimOperator(buf *bytes.Buffer, v SimOperator) error { return SetU8(buf, uint8(v)) }

type SimOperatorList []SimOperator
func GetSimOperatorList(buf *bytes.Buffer) (SimOperatorList, error) {
	val, err := GetU8List(buf)
	if err != nil { return nil, err }
	return *(*SimOperatorList)(unsafe.Pointer(&val)), nil
}
func SetSimOperatorList(buf *bytes.Buffer, v SimOperatorList) error { return SetU8List(buf, *(*[]uint8)(unsafe.Pointer(&v))) }
func IsSimOperatorList(v SimOperatorList) bool {
	for _, item := range v {
		if !IsSimOperator(item) { return false }
	}
	return true
}
func EqSimOperatorList(a, b SimOperatorList) bool { return slices.Equal(a, b) }

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

func IsOrderStatus(v OrderStatus) bool {
	switch v {
	case OrderStatusPending, OrderStatusClosed, OrderStatusCanceled, OrderStatusShipped, OrderStatusDelivered, OrderStatusActived, OrderStatusSettled:
		return true
	default:
		return false
	}
}

func GetOrderStatus(buf *bytes.Buffer) (OrderStatus, error) {
	val, err := GetU8(buf)
	return OrderStatus(val), err
}

func SetOrderStatus(buf *bytes.Buffer, v OrderStatus) error { return SetU8(buf, uint8(v)) }

type OrderStatusList []OrderStatus
func GetOrderStatusList(buf *bytes.Buffer) (OrderStatusList, error) {
	val, err := GetU8List(buf)
	if err != nil { return nil, err }
	return *(*OrderStatusList)(unsafe.Pointer(&val)), nil
}
func SetOrderStatusList(buf *bytes.Buffer, v OrderStatusList) error { return SetU8List(buf, *(*[]uint8)(unsafe.Pointer(&v))) }
func IsOrderStatusList(v OrderStatusList) bool {
	for _, item := range v {
		if !IsOrderStatus(item) { return false }
	}
	return true
}
func EqOrderStatusList(a, b OrderStatusList) bool { return slices.Equal(a, b) }

