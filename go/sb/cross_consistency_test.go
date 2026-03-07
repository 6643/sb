package sb

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	rt "sb/go/sb/rt"
)

type crossWireCase struct {
	name   string
	kind   string
	encode func() ([]byte, error)
	verify func([]byte) error
}

type crossWireInput struct {
	Name string `json:"name"`
	Kind string `json:"kind"`
	Hex  string `json:"hex"`
}

type crossWireOutput struct {
	Name string `json:"name"`
	Hex  string `json:"hex"`
}

func TestCrossLanguageWireConsistency(t *testing.T) {
	cases := buildCrossWireCases()
	inputs := make([]crossWireInput, 0, len(cases))
	expected := make(map[string]string, len(cases))

	for _, tc := range cases {
		wire, err := tc.encode()
		if err != nil {
			t.Fatalf("%s encode failed: %v", tc.name, err)
		}
		if err := tc.verify(bytes.Clone(wire)); err != nil {
			t.Fatalf("%s local roundtrip failed: %v", tc.name, err)
		}
		hexWire := hex.EncodeToString(wire)
		inputs = append(inputs, crossWireInput{Name: tc.name, Kind: tc.kind, Hex: hexWire})
		expected[tc.name] = hexWire
	}

	inputPath := filepath.Join(t.TempDir(), "cross_wire_cases.json")
	data, err := json.Marshal(inputs)
	if err != nil {
		t.Fatalf("marshal cross cases failed: %v", err)
	}
	if err := os.WriteFile(inputPath, data, 0o644); err != nil {
		t.Fatalf("write cross cases failed: %v", err)
	}

	repoRoot, err := filepath.Abs(filepath.Join("..", ".."))
	if err != nil {
		t.Fatalf("resolve repo root failed: %v", err)
	}

	cmd := exec.Command("bun", "ts/sb/cross_consistency.ts", inputPath)
	cmd.Dir = repoRoot
	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("bun cross consistency failed: %v\n%s", err, string(output))
	}

	var results []crossWireOutput
	if err := json.Unmarshal(output, &results); err != nil {
		t.Fatalf("parse bun result failed: %v\n%s", err, string(output))
	}
	if len(results) != len(expected) {
		t.Fatalf("bun result count mismatch: got %d want %d", len(results), len(expected))
	}

	seen := make(map[string]bool, len(results))
	for _, result := range results {
		want, ok := expected[result.Name]
		if !ok {
			t.Fatalf("unexpected bun case: %s", result.Name)
		}
		seen[result.Name] = true
		if result.Hex != want {
			t.Fatalf("%s wire mismatch:\n  go=%s\n  ts=%s", result.Name, want, result.Hex)
		}
	}
	for name := range expected {
		if !seen[name] {
			t.Fatalf("bun result missing case: %s", name)
		}
	}
}

func buildCrossWireCases() []crossWireCase {
	simInfoZero := &SimInfo{}
	simInfoRich := &SimInfo{
		Id:      7,
		Title:   "标题",
		Content: "正文",
		A:       true,
		C:       true,
		Zip:     []byte{1, 2, 3, 4},
	}
	simInfoBoundary := &SimInfo{
		Id:      9,
		Title:   strings.Repeat("a", 255),
		Content: strings.Repeat("b", 256),
		B:       true,
		D:       true,
		Zip:     seqBytes(256),
	}

	rechargeZero := &Recharge{}
	rechargeRich := &Recharge{
		Id:    17,
		Type:  []OrderStatus{OrderStatusPending, OrderStatusClosed, OrderStatusSettled},
		Phone: []string{"13800138000", "", "客服"},
		Si:    simInfoRich,
	}
	rechargeBoundary := &Recharge{
		Id:    18,
		Type:  buildOrderStatusList255(),
		Phone: buildTextList255(),
		Si:    nil,
	}
	rechargeAZero := &RechargeA{}
	rechargeARich := &RechargeA{
		Id:    21,
		Type:  []OrderStatus{OrderStatusPending, OrderStatusDelivered},
		Phone: []string{"10086", "热线"},
		Si:    simInfoRich,
		Aid:   99,
	}
	rechargeBZero := &RechargeB{}
	rechargeBRich := &RechargeB{
		Id:    31,
		Type:  []OrderStatus{OrderStatusCanceled, OrderStatusActived},
		Phone: []string{"18800001111", ""},
		Si:    simInfoBoundary,
		Bid:   199,
	}
	simZero := &Sim{}
	simRich := &Sim{
		Id:                99,
		Type:              TypeRecharge,
		Status:            ItemStatusOnline,
		Commission:        55,
		Supplier:          77,
		Aff:               88,
		ContractDuration:  12,
		Name:              "测试套餐",
		Operator:          SimOperatorDx,
		Monthly:           29,
		FlowUniversal:     100,
		FlowDirectional:   20,
		CanMoveFlow:       true,
		CallMonth:         300,
		CallPrice:         1,
		SmsMonth:          50,
		SmsPrice:          2,
		MinAge:            18,
		MaxAge:            60,
		Attribution:       1,
		PickPhone:         []SimPickPhone{SimPickPhoneNo, SimPickPhoneYes, SimPickPhoneActive, SimPickPhoneAbcc},
		FirstChargeLink:   "https://pay.example/sim",
		FirstChargeMoney:  "100",
		FirstChargeReturn: "20",
		BanCity:           []uint32{1, 2, 0, 0, 3},
		Info:              []*SimInfo{nil, simInfoRich, {}},
		Snapshot:          []string{"", "snap-1", "截图"},
	}
	simOrder2Zero := &SimOrder2{}
	simOrder2Rich := &SimOrder2{
		Id:       101,
		Name:     "张三",
		Phone:    "13800138000",
		IdNo:     "320311770706001",
		CityCode: 310100,
		Address:  "上海市浦东新区测试路 1 号",
		NewPhone: "13900139000",
	}
	simOrderZero := &SimOrder{}
	simOrderRich := &SimOrder{
		Id:         202,
		AccountId:  301,
		ItemId:     401,
		Name:       "李四",
		Phone:      "13700137000",
		IdNo:       "110101199003076611",
		CityCode:   110100,
		Address:    "北京市朝阳区测试街道 2 号",
		NewPhone:   "13600136000",
		Commission: 88,
		Status:     OrderStatusDelivered,
	}

	getBinShort := []byte{9, 8, 7}
	getBinBoundary := seqBytes(256)

	return []crossWireCase{
		structWireCase("struct.recharge.zero", "struct.recharge", rechargeZero, SetRecharge, ReadRecharge, EqRecharge),
		structWireCase("struct.recharge.rich", "struct.recharge", rechargeRich, SetRecharge, ReadRecharge, EqRecharge),
		structWireCase("struct.recharge.list255", "struct.recharge", rechargeBoundary, SetRecharge, ReadRecharge, EqRecharge),
		structWireCase("struct.rechargeA.zero", "struct.rechargeA", rechargeAZero, SetRechargeA, ReadRechargeA, EqRechargeA),
		structWireCase("struct.rechargeA.rich", "struct.rechargeA", rechargeARich, SetRechargeA, ReadRechargeA, EqRechargeA),
		structWireCase("struct.rechargeB.zero", "struct.rechargeB", rechargeBZero, SetRechargeB, ReadRechargeB, EqRechargeB),
		structWireCase("struct.rechargeB.rich", "struct.rechargeB", rechargeBRich, SetRechargeB, ReadRechargeB, EqRechargeB),
		structWireCase("struct.simInfo.zero", "struct.simInfo", simInfoZero, SetSimInfo, ReadSimInfo, EqSimInfo),
		structWireCase("struct.simInfo.rich", "struct.simInfo", simInfoRich, SetSimInfo, ReadSimInfo, EqSimInfo),
		structWireCase("struct.simInfo.boundary", "struct.simInfo", simInfoBoundary, SetSimInfo, ReadSimInfo, EqSimInfo),
		structWireCase("struct.sim.zero", "struct.sim", simZero, SetSim, ReadSim, EqSim),
		structWireCase("struct.sim.rich", "struct.sim", simRich, SetSim, ReadSim, EqSim),
		structWireCase("struct.simOrder2.zero", "struct.simOrder2", simOrder2Zero, SetSimOrder2, ReadSimOrder2, EqSimOrder2),
		structWireCase("struct.simOrder2.rich", "struct.simOrder2", simOrder2Rich, SetSimOrder2, ReadSimOrder2, EqSimOrder2),
		structWireCase("struct.simOrder.zero", "struct.simOrder", simOrderZero, SetSimOrder, ReadSimOrder, EqSimOrder),
		structWireCase("struct.simOrder.rich", "struct.simOrder", simOrderRich, SetSimOrder, ReadSimOrder, EqSimOrder),
		emptyWireCase("api.user.get_abc.req.empty", "api.user.get_abc.req"),
		orderStatusWireCase("api.user.get_abc.resp.default", "api.user.get_abc.resp", OrderStatusPending),
		orderStatusWireCase("api.user.get_abc.resp.rich", "api.user.get_abc.resp", OrderStatusSettled),
		u8PairWireCase("api.user.get_abcd.req", "api.user.get_abcd.req", 7, 9),
		orderStatusWireCase("api.user.get_abcd.resp", "api.user.get_abcd.resp", OrderStatusActived),
		structWireCase("api.user.set_sim_info.req", "api.user.set_sim_info.req", simInfoBoundary, SetSimInfo, ReadSimInfo, EqSimInfo),
		emptyWireCase("api.user.set_sim_info.resp.empty", "api.user.set_sim_info.resp"),
		u8WireCase("api.get_count.req", "api.get_count.req", 77),
		u8WireCase("api.get_count.resp", "api.get_count.resp", 123),
		u8WireCase("api.get_bin.req", "api.get_bin.req", 12),
		binRespWireCase("api.get_bin.resp.empty", "api.get_bin.resp", nil),
		binRespWireCase("api.get_bin.resp.u8", "api.get_bin.resp", getBinShort),
		binRespWireCase("api.get_bin.resp.u16", "api.get_bin.resp", getBinBoundary),
	}
}

func structWireCase[T any](name string, kind string, value *T, set func(*bytes.Buffer, *T) error, read func(*bytes.Buffer) (*T, error), eq func(*T, *T) bool) crossWireCase {
	return crossWireCase{
		name: name,
		kind: kind,
		encode: func() ([]byte, error) {
			var buf bytes.Buffer
			if err := set(&buf, value); err != nil {
				return nil, err
			}
			return bytes.Clone(buf.Bytes()), nil
		},
		verify: func(data []byte) error {
			buf := bytes.NewBuffer(data)
			got, err := read(buf)
			if err != nil {
				return err
			}
			if buf.Len() != 0 {
				return fmt.Errorf("leftover bytes: %d", buf.Len())
			}
			if !eq(value, got) {
				return fmt.Errorf("decoded value mismatch")
			}
			var re bytes.Buffer
			if err := set(&re, got); err != nil {
				return err
			}
			if !bytes.Equal(data, re.Bytes()) {
				return fmt.Errorf("re-encoded bytes mismatch")
			}
			return nil
		},
	}
}

func emptyWireCase(name string, kind string) crossWireCase {
	return crossWireCase{
		name: name,
		kind: kind,
		encode: func() ([]byte, error) {
			return nil, nil
		},
		verify: func(data []byte) error {
			if len(data) != 0 {
				return fmt.Errorf("expected empty payload, got %d bytes", len(data))
			}
			return nil
		},
	}
}

func u8WireCase(name string, kind string, value uint8) crossWireCase {
	return crossWireCase{
		name: name,
		kind: kind,
		encode: func() ([]byte, error) {
			var buf bytes.Buffer
			if err := rt.SetU8(&buf, value); err != nil {
				return nil, err
			}
			return bytes.Clone(buf.Bytes()), nil
		},
		verify: func(data []byte) error {
			buf := bytes.NewBuffer(data)
			got, err := rt.GetU8(buf)
			if err != nil {
				return err
			}
			if buf.Len() != 0 {
				return fmt.Errorf("leftover bytes: %d", buf.Len())
			}
			if got != value {
				return fmt.Errorf("decoded value mismatch: got %d want %d", got, value)
			}
			var re bytes.Buffer
			if err := rt.SetU8(&re, got); err != nil {
				return err
			}
			if !bytes.Equal(data, re.Bytes()) {
				return fmt.Errorf("re-encoded bytes mismatch")
			}
			return nil
		},
	}
}

func u8PairWireCase(name string, kind string, first uint8, second uint8) crossWireCase {
	return crossWireCase{
		name: name,
		kind: kind,
		encode: func() ([]byte, error) {
			var buf bytes.Buffer
			if err := rt.SetU8(&buf, first); err != nil {
				return nil, err
			}
			if err := rt.SetU8(&buf, second); err != nil {
				return nil, err
			}
			return bytes.Clone(buf.Bytes()), nil
		},
		verify: func(data []byte) error {
			buf := bytes.NewBuffer(data)
			gotFirst, err := rt.GetU8(buf)
			if err != nil {
				return err
			}
			gotSecond, err := rt.GetU8(buf)
			if err != nil {
				return err
			}
			if buf.Len() != 0 {
				return fmt.Errorf("leftover bytes: %d", buf.Len())
			}
			if gotFirst != first || gotSecond != second {
				return fmt.Errorf("decoded pair mismatch: got (%d,%d) want (%d,%d)", gotFirst, gotSecond, first, second)
			}
			var re bytes.Buffer
			if err := rt.SetU8(&re, gotFirst); err != nil {
				return err
			}
			if err := rt.SetU8(&re, gotSecond); err != nil {
				return err
			}
			if !bytes.Equal(data, re.Bytes()) {
				return fmt.Errorf("re-encoded bytes mismatch")
			}
			return nil
		},
	}
}

func orderStatusWireCase(name string, kind string, value OrderStatus) crossWireCase {
	return crossWireCase{
		name: name,
		kind: kind,
		encode: func() ([]byte, error) {
			if !isAssignableOrderStatus(value) {
				return nil, fmt.Errorf("非法枚举值: %d", value)
			}
			var buf bytes.Buffer
			if err := rt.SetU8(&buf, uint8(normalizeOrderStatus(value))); err != nil {
				return nil, err
			}
			return bytes.Clone(buf.Bytes()), nil
		},
		verify: func(data []byte) error {
			buf := bytes.NewBuffer(data)
			raw, err := rt.GetU8(buf)
			if err != nil {
				return err
			}
			if buf.Len() != 0 {
				return fmt.Errorf("leftover bytes: %d", buf.Len())
			}
			got := OrderStatus(raw)
			if !isOrderStatus(got) {
				return fmt.Errorf("非法枚举值: %d", got)
			}
			if !eqOrderStatusValue(got, value) {
				return fmt.Errorf("decoded value mismatch: got %d want %d", got, value)
			}
			var re bytes.Buffer
			if err := rt.SetU8(&re, uint8(normalizeOrderStatus(got))); err != nil {
				return err
			}
			if !bytes.Equal(data, re.Bytes()) {
				return fmt.Errorf("re-encoded bytes mismatch")
			}
			return nil
		},
	}
}

func binRespWireCase(name string, kind string, value []byte) crossWireCase {
	return crossWireCase{
		name: name,
		kind: kind,
		encode: func() ([]byte, error) {
			state, err := rt.BinState(len(value))
			if err != nil {
				return nil, err
			}
			var buf bytes.Buffer
			if err := rt.SetU8(&buf, state); err != nil {
				return nil, err
			}
			if err := rt.SetBinCompact(&buf, state, value); err != nil {
				return nil, err
			}
			return bytes.Clone(buf.Bytes()), nil
		},
		verify: func(data []byte) error {
			buf := bytes.NewBuffer(data)
			state, err := rt.GetU8(buf)
			if err != nil {
				return err
			}
			got, err := rt.GetBinCompactInto(buf, state, nil)
			if err != nil {
				return err
			}
			if buf.Len() != 0 {
				return fmt.Errorf("leftover bytes: %d", buf.Len())
			}
			if !bytes.Equal(got, value) {
				return fmt.Errorf("decoded value mismatch")
			}
			canonical, err := rt.BinState(len(got))
			if err != nil {
				return err
			}
			if state != canonical {
				return fmt.Errorf("non-canonical bin state: got %d want %d", state, canonical)
			}
			var re bytes.Buffer
			if err := rt.SetU8(&re, canonical); err != nil {
				return err
			}
			if err := rt.SetBinCompact(&re, canonical, got); err != nil {
				return err
			}
			if !bytes.Equal(data, re.Bytes()) {
				return fmt.Errorf("re-encoded bytes mismatch")
			}
			return nil
		},
	}
}

func seqBytes(n int) []byte {
	out := make([]byte, n)
	for i := range out {
		out[i] = byte(i % 251)
	}
	return out
}

func buildOrderStatusList255() []OrderStatus {
	out := make([]OrderStatus, 255)
	for i := range out {
		switch i % 4 {
		case 0:
			out[i] = OrderStatusPending
		case 1:
			out[i] = OrderStatusClosed
		case 2:
			out[i] = OrderStatusActived
		default:
			out[i] = OrderStatusSettled
		}
	}
	return out
}

func buildTextList255() []string {
	out := make([]string, 255)
	for i := range out {
		switch i % 4 {
		case 0:
			out[i] = ""
		case 1:
			out[i] = fmt.Sprintf("v%03d", i)
		case 2:
			out[i] = "中文"
		default:
			out[i] = "x"
		}
	}
	return out
}
