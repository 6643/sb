package sb

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
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

type crossWireRejectCase struct {
	name         string
	kind         string
	hex          string
	verifyLocal  func([]byte) error
	wantContains string
}

func TestCrossLanguageWireConsistency(t *testing.T) {
	runCrossLanguageWireConsistency(t, buildCrossWireCases())
}

func TestCrossLanguageWireConsistencyRandom(t *testing.T) {
	const rounds = 12
	const seed int64 = 20260308
	runCrossLanguageWireConsistency(t, buildRandomCrossWireCases(seed, rounds))
}

func TestCrossLanguageWireRejectsMalformedInputs(t *testing.T) {
	runCrossLanguageWireRejects(t, buildCrossWireRejectCases())
}

func runCrossLanguageWireConsistency(t *testing.T, cases []crossWireCase) {
	t.Helper()
	if _, err := exec.LookPath("bun"); err != nil {
		t.Skip("bun not found in PATH")
	}
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

func runCrossLanguageWireRejects(t *testing.T, cases []crossWireRejectCase) {
	t.Helper()
	if _, err := exec.LookPath("bun"); err != nil {
		t.Skip("bun not found in PATH")
	}
	repoRoot, err := filepath.Abs(filepath.Join("..", ".."))
	if err != nil {
		t.Fatalf("resolve repo root failed: %v", err)
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			wire, err := hex.DecodeString(tc.hex)
			if err != nil {
				t.Fatalf("decode hex failed: %v", err)
			}
			if err := tc.verifyLocal(bytes.Clone(wire)); err == nil {
				t.Fatalf("local decoder should reject malformed wire")
			}

			inputPath := filepath.Join(t.TempDir(), "cross_wire_reject_case.json")
			data, err := json.Marshal([]crossWireInput{{Name: tc.name, Kind: tc.kind, Hex: tc.hex}})
			if err != nil {
				t.Fatalf("marshal reject case failed: %v", err)
			}
			if err := os.WriteFile(inputPath, data, 0o644); err != nil {
				t.Fatalf("write reject case failed: %v", err)
			}

			cmd := exec.Command("bun", "ts/sb/cross_consistency.ts", inputPath)
			cmd.Dir = repoRoot
			output, err := cmd.CombinedOutput()
			if err == nil {
				t.Fatalf("bun cross consistency should reject malformed wire, output=%s", string(output))
			}
			if !strings.Contains(string(output), tc.wantContains) {
				t.Fatalf("bun error mismatch:\nwant substring: %q\noutput:\n%s", tc.wantContains, string(output))
			}
		})
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

func buildCrossWireRejectCases() []crossWireRejectCase {
	return []crossWireRejectCase{
		{
			name: "struct.simInfo.non_zero_padding",
			kind: "struct.simInfo",
			hex:  "0010",
			verifyLocal: func(data []byte) error {
				_, err := ReadSimInfo(bytes.NewBuffer(data))
				return err
			},
			wantContains: "SimInfo header padding bit 0 is not zero",
		},
		{
			name: "struct.recharge.illegal_list_state",
			kind: "struct.recharge",
			hex:  "40",
			verifyLocal: func(data []byte) error {
				_, err := ReadRecharge(bytes.NewBuffer(data))
				return err
			},
			wantContains: "struct.recharge.illegal_list_state decode: get Recharge Type: list count state 2 is illegal",
		},
		{
			name: "api.get_bin.resp.non_canonical_u16",
			kind: "api.get_bin.resp",
			hex:  "020300090807",
			verifyLocal: func(data []byte) error {
				buf := bytes.NewBuffer(data)
				state, err := GetU8(buf)
				if err != nil {
					return err
				}
				_, err = GetBinInto(buf, state, nil)
				return err
			},
			wantContains: "api.get_bin.resp.non_canonical_u16 decode body: bin state 2 is not canonical for length 3",
		},
	}
}

func buildRandomCrossWireCases(seed int64, rounds int) []crossWireCase {
	r := rand.New(rand.NewSource(seed))
	cases := make([]crossWireCase, 0, rounds*17)
	for i := 0; i < rounds; i++ {
		prefix := fmt.Sprintf("random.%02d", i)
		cases = append(cases,
			structWireCase(prefix+".struct.recharge", "struct.recharge", randRecharge(r), SetRecharge, ReadRecharge, EqRecharge),
			structWireCase(prefix+".struct.rechargeA", "struct.rechargeA", randRechargeA(r), SetRechargeA, ReadRechargeA, EqRechargeA),
			structWireCase(prefix+".struct.rechargeB", "struct.rechargeB", randRechargeB(r), SetRechargeB, ReadRechargeB, EqRechargeB),
			structWireCase(prefix+".struct.simInfo", "struct.simInfo", randSimInfo(r), SetSimInfo, ReadSimInfo, EqSimInfo),
			structWireCase(prefix+".struct.sim", "struct.sim", randSim(r), SetSim, ReadSim, EqSim),
			structWireCase(prefix+".struct.simOrder2", "struct.simOrder2", randSimOrder2(r), SetSimOrder2, ReadSimOrder2, EqSimOrder2),
			structWireCase(prefix+".struct.simOrder", "struct.simOrder", randSimOrder(r), SetSimOrder, ReadSimOrder, EqSimOrder),
			emptyWireCase(prefix+".api.user.get_abc.req.empty", "api.user.get_abc.req"),
			orderStatusWireCase(prefix+".api.user.get_abc.resp", "api.user.get_abc.resp", randOrderStatus(r)),
			u8PairWireCase(prefix+".api.user.get_abcd.req", "api.user.get_abcd.req", randU8(r), randU8(r)),
			orderStatusWireCase(prefix+".api.user.get_abcd.resp", "api.user.get_abcd.resp", randOrderStatus(r)),
			structWireCase(prefix+".api.user.set_sim_info.req", "api.user.set_sim_info.req", randSimInfo(r), SetSimInfo, ReadSimInfo, EqSimInfo),
			emptyWireCase(prefix+".api.user.set_sim_info.resp.empty", "api.user.set_sim_info.resp"),
			u8WireCase(prefix+".api.get_count.req", "api.get_count.req", randU8(r)),
			u8WireCase(prefix+".api.get_count.resp", "api.get_count.resp", randU8(r)),
			u8WireCase(prefix+".api.get_bin.req", "api.get_bin.req", randU8(r)),
			binRespWireCase(prefix+".api.get_bin.resp", "api.get_bin.resp", randBin(r, 384)),
		)
	}
	return cases
}

var randomASCII = []byte("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789-_ /")

func randInt(r *rand.Rand, max int, edges ...int) int {
	if max <= 0 {
		return 0
	}
	filtered := make([]int, 0, len(edges))
	for _, edge := range edges {
		if edge >= 0 && edge <= max {
			filtered = append(filtered, edge)
		}
	}
	if len(filtered) > 0 && r.Intn(4) == 0 {
		return filtered[r.Intn(len(filtered))]
	}
	return r.Intn(max + 1)
}

func randU8(r *rand.Rand) uint8 {
	return uint8(randInt(r, 0xFF, 0, 1, 2, 3, 7, 15, 31, 63, 127, 254, 255))
}

func randU16(r *rand.Rand) uint16 {
	return uint16(randInt(r, 0xFFFF, 0, 1, 2, 3, 7, 15, 31, 63, 127, 255, 256, 1023, 4095, 65535))
}

func randU32(r *rand.Rand) uint32 {
	return uint32(randInt(r, 1<<20, 0, 1, 2, 3, 7, 15, 31, 63, 127, 255, 256, 1023, 4095, 65535, 1<<20))
}

func randText(r *rand.Rand, maxLen int) string {
	n := randInt(r, maxLen, 0, 1, 2, 3, 7, 15, 31, 63, 127, 255, 256)
	if n == 0 {
		return ""
	}
	buf := make([]byte, n)
	for i := range buf {
		buf[i] = randomASCII[r.Intn(len(randomASCII))]
	}
	return string(buf)
}

func randBin(r *rand.Rand, maxLen int) []byte {
	n := randInt(r, maxLen, 0, 1, 2, 3, 7, 15, 31, 63, 127, 255, 256)
	if n == 0 {
		return nil
	}
	buf := make([]byte, n)
	for i := range buf {
		buf[i] = byte(r.Intn(256))
	}
	return buf
}

func randListCount(r *rand.Rand, max int) int {
	return randInt(r, max, 0, 1, 2, 3, 7, 15, 31)
}

func randOrderStatus(r *rand.Rand) OrderStatus {
	values := [...]OrderStatus{
		OrderStatusPending,
		OrderStatusClosed,
		OrderStatusCanceled,
		OrderStatusShipped,
		OrderStatusDelivered,
		OrderStatusActived,
		OrderStatusSettled,
	}
	return values[r.Intn(len(values))]
}

func randType(r *rand.Rand) Type {
	values := [...]Type{TypeSim, TypeRecharge}
	return values[r.Intn(len(values))]
}

func randItemStatus(r *rand.Rand) ItemStatus {
	values := [...]ItemStatus{ItemStatusOffline, ItemStatusOnline}
	return values[r.Intn(len(values))]
}

func randSimPickPhone(r *rand.Rand) SimPickPhone {
	values := [...]SimPickPhone{SimPickPhoneNo, SimPickPhoneYes, SimPickPhoneActive, SimPickPhoneAbcc}
	return values[r.Intn(len(values))]
}

func randSimOperator(r *rand.Rand) SimOperator {
	values := [...]SimOperator{
		SimOperatorZz,
		SimOperatorLt,
		SimOperatorYd,
		SimOperatorDx,
		SimOperatorGd,
		SimOperatorXx,
		SimOperatorA,
		SimOperatorB,
	}
	return values[r.Intn(len(values))]
}

func randOrderStatusList(r *rand.Rand, max int) []OrderStatus {
	count := randListCount(r, max)
	out := make([]OrderStatus, count)
	for i := range out {
		out[i] = randOrderStatus(r)
	}
	return out
}

func randSimPickPhoneList(r *rand.Rand, max int) []SimPickPhone {
	count := randListCount(r, max)
	out := make([]SimPickPhone, count)
	for i := range out {
		out[i] = randSimPickPhone(r)
	}
	return out
}

func randTextList(r *rand.Rand, maxCount int, maxItemLen int) []string {
	count := randListCount(r, maxCount)
	out := make([]string, count)
	for i := range out {
		out[i] = randText(r, maxItemLen)
	}
	return out
}

func randU32List(r *rand.Rand, maxCount int) []uint32 {
	count := randListCount(r, maxCount)
	out := make([]uint32, count)
	for i := range out {
		out[i] = randU32(r)
	}
	return out
}

func randMaybeSimInfo(r *rand.Rand) *SimInfo {
	switch r.Intn(4) {
	case 0:
		return nil
	case 1:
		return &SimInfo{}
	default:
		return randSimInfo(r)
	}
}

func randSimInfoList(r *rand.Rand, maxCount int) []*SimInfo {
	count := randListCount(r, maxCount)
	out := make([]*SimInfo, count)
	for i := range out {
		out[i] = randMaybeSimInfo(r)
	}
	return out
}

func randSimInfo(r *rand.Rand) *SimInfo {
	return &SimInfo{
		Id:      randU32(r),
		Title:   randText(r, 320),
		Content: randText(r, 320),
		A:       r.Intn(2) == 1,
		B:       r.Intn(2) == 1,
		C:       r.Intn(2) == 1,
		D:       r.Intn(2) == 1,
		Zip:     randBin(r, 384),
	}
}

func randRecharge(r *rand.Rand) *Recharge {
	return &Recharge{
		Id:    randU32(r),
		Type:  randOrderStatusList(r, 24),
		Phone: randTextList(r, 18, 96),
		Si:    randMaybeSimInfo(r),
	}
}

func randRechargeA(r *rand.Rand) *RechargeA {
	return &RechargeA{
		Id:    randU32(r),
		Type:  randOrderStatusList(r, 24),
		Phone: randTextList(r, 18, 96),
		Si:    randMaybeSimInfo(r),
		Aid:   randU32(r),
	}
}

func randRechargeB(r *rand.Rand) *RechargeB {
	return &RechargeB{
		Id:    randU32(r),
		Type:  randOrderStatusList(r, 24),
		Phone: randTextList(r, 18, 96),
		Si:    randMaybeSimInfo(r),
		Bid:   randU32(r),
	}
}

func randSim(r *rand.Rand) *Sim {
	return &Sim{
		Id:                randU32(r),
		Type:              randType(r),
		Status:            randItemStatus(r),
		Commission:        randU16(r),
		Supplier:          randU32(r),
		Aff:               randU32(r),
		ContractDuration:  randU8(r),
		Name:              randText(r, 128),
		Operator:          randSimOperator(r),
		Monthly:           randU16(r),
		FlowUniversal:     randU16(r),
		FlowDirectional:   randU16(r),
		CanMoveFlow:       r.Intn(2) == 1,
		CallMonth:         randU16(r),
		CallPrice:         randU16(r),
		SmsMonth:          randU16(r),
		SmsPrice:          randU16(r),
		MinAge:            randU8(r),
		MaxAge:            randU8(r),
		Attribution:       randU32(r),
		PickPhone:         randSimPickPhoneList(r, 20),
		FirstChargeLink:   randText(r, 128),
		FirstChargeMoney:  randText(r, 64),
		FirstChargeReturn: randText(r, 64),
		BanCity:           randU32List(r, 20),
		Info:              randSimInfoList(r, 12),
		Snapshot:          randTextList(r, 16, 96),
	}
}

func randSimOrder2(r *rand.Rand) *SimOrder2 {
	return &SimOrder2{
		Id:       randU32(r),
		Name:     randText(r, 96),
		Phone:    randText(r, 32),
		IdNo:     randText(r, 32),
		CityCode: randU32(r),
		Address:  randText(r, 192),
		NewPhone: randText(r, 32),
	}
}

func randSimOrder(r *rand.Rand) *SimOrder {
	return &SimOrder{
		Id:         randU32(r),
		AccountId:  randU32(r),
		ItemId:     randU32(r),
		Name:       randText(r, 96),
		Phone:      randText(r, 32),
		IdNo:       randText(r, 32),
		CityCode:   randU32(r),
		Address:    randText(r, 192),
		NewPhone:   randText(r, 32),
		Commission: randU16(r),
		Status:     randOrderStatus(r),
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
			if err := SetU8(&buf, value); err != nil {
				return nil, err
			}
			return bytes.Clone(buf.Bytes()), nil
		},
		verify: func(data []byte) error {
			buf := bytes.NewBuffer(data)
			got, err := GetU8(buf)
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
			if err := SetU8(&re, got); err != nil {
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
			if err := SetU8(&buf, first); err != nil {
				return nil, err
			}
			if err := SetU8(&buf, second); err != nil {
				return nil, err
			}
			return bytes.Clone(buf.Bytes()), nil
		},
		verify: func(data []byte) error {
			buf := bytes.NewBuffer(data)
			gotFirst, err := GetU8(buf)
			if err != nil {
				return err
			}
			gotSecond, err := GetU8(buf)
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
			if err := SetU8(&re, gotFirst); err != nil {
				return err
			}
			if err := SetU8(&re, gotSecond); err != nil {
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
			if err := SetU8(&buf, uint8(normalizeOrderStatus(value))); err != nil {
				return nil, err
			}
			return bytes.Clone(buf.Bytes()), nil
		},
		verify: func(data []byte) error {
			buf := bytes.NewBuffer(data)
			raw, err := GetU8(buf)
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
			if err := SetU8(&re, uint8(normalizeOrderStatus(got))); err != nil {
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
			state, err := BinState(len(value))
			if err != nil {
				return nil, err
			}
			var buf bytes.Buffer
			if err := SetU8(&buf, state); err != nil {
				return nil, err
			}
			if err := SetBin(&buf, state, value); err != nil {
				return nil, err
			}
			return bytes.Clone(buf.Bytes()), nil
		},
		verify: func(data []byte) error {
			buf := bytes.NewBuffer(data)
			state, err := GetU8(buf)
			if err != nil {
				return err
			}
			got, err := GetBinInto(buf, state, nil)
			if err != nil {
				return err
			}
			if buf.Len() != 0 {
				return fmt.Errorf("leftover bytes: %d", buf.Len())
			}
			if !bytes.Equal(got, value) {
				return fmt.Errorf("decoded value mismatch")
			}
			canonical, err := BinState(len(got))
			if err != nil {
				return err
			}
			if state != canonical {
				return fmt.Errorf("non-canonical bin state: got %d want %d", state, canonical)
			}
			var re bytes.Buffer
			if err := SetU8(&re, canonical); err != nil {
				return err
			}
			if err := SetBin(&re, canonical, got); err != nil {
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
