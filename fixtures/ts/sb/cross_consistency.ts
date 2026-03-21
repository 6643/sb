import * as _ from "./_";

declare const Bun: {
    argv: string[];
    file: (path: string) => { text: () => Promise<string> };
};

declare const process: {
    stdout: { write: (text: string) => void };
};

type InputCase = {
    name: string;
    kind: string;
    hex: string;
};

type OutputCase = {
    name: string;
    hex: string;
};

const toHex = (bytes: Uint8Array): string => Array.from(bytes, b => b.toString(16).padStart(2, "0")).join("");

const fromHex = (hex: string): Uint8Array => {
    if (hex.length === 0) return new Uint8Array(0);
    if (hex.length % 2 !== 0) throw new Error(`invalid hex length: ${hex.length}`);
    const out = new Uint8Array(hex.length / 2);
    for (let i = 0; i < out.length; i++) out[i] = parseInt(hex.slice(i * 2, i * 2 + 2), 16);
    return out;
};

const ensureConsumed = (buf: _.Buffer, name: string): void => {
    if (buf.len !== 0) throw new Error(`${name}: leftover bytes ${buf.len}`);
};

const mustErr = (err: _.Err, name: string): void => {
    if (err !== undefined) throw new Error(`${name}: ${err.message}`);
};

const mustNullErr = (err: Error | null, name: string): void => {
    if (err !== null) throw new Error(`${name}: ${err.message}`);
};

const roundTripStruct = <T>(
    name: string,
    bytes: Uint8Array,
    read: (buf: _.Buffer) => [T, _.Err],
    write: (buf: _.Buffer, value: T) => _.Err,
): string => {
    const buf = new _.Buffer(bytes);
    const [value, err] = read(buf);
    mustErr(err, `${name} decode`);
    ensureConsumed(buf, `${name} decode`);
    const out = new _.Buffer();
    mustErr(write(out, value), `${name} encode`);
    return toHex(out.bytes);
};

const roundTripEmpty = (name: string, bytes: Uint8Array): string => {
    const buf = new _.Buffer(bytes);
    ensureConsumed(buf, `${name} decode`);
    return "";
};

const roundTripU8 = (name: string, bytes: Uint8Array): string => {
    const buf = new _.Buffer(bytes);
    const [value, err] = _.getU8(buf);
    mustNullErr(err, `${name} decode`);
    ensureConsumed(buf, `${name} decode`);
    const out = new _.Buffer();
    mustNullErr(_.setU8(out, value), `${name} encode`);
    return toHex(out.bytes);
};

const roundTripU8Pair = (name: string, bytes: Uint8Array): string => {
    const buf = new _.Buffer(bytes);
    const [first, errFirst] = _.getU8(buf);
    mustNullErr(errFirst, `${name} decode first`);
    const [second, errSecond] = _.getU8(buf);
    mustNullErr(errSecond, `${name} decode second`);
    ensureConsumed(buf, `${name} decode`);
    const out = new _.Buffer();
    mustNullErr(_.setU8(out, first), `${name} encode first`);
    mustNullErr(_.setU8(out, second), `${name} encode second`);
    return toHex(out.bytes);
};

const roundTripOrderStatus = (name: string, bytes: Uint8Array): string => {
    const buf = new _.Buffer(bytes);
    const [raw, err] = _.getU8(buf);
    mustNullErr(err, `${name} decode`);
    const value = raw as _.OrderStatus;
    if (!_.IsOrderStatus(value)) throw new Error(`${name}: invalid order status ${value}`);
    ensureConsumed(buf, `${name} decode`);
    const out = new _.Buffer();
    mustNullErr(_.setU8(out, value as any), `${name} encode`);
    return toHex(out.bytes);
};

const roundTripBinResp = (name: string, bytes: Uint8Array): string => {
    const buf = new _.Buffer(bytes);
    const [state, errState] = _.getU8(buf);
    mustNullErr(errState, `${name} decode state`);
    const [value, errValue] = _.getBin(buf, state);
    mustNullErr(errValue, `${name} decode body`);
    ensureConsumed(buf, `${name} decode`);
    const [canonical, errCanonical] = _.binState(value.byteLength);
    mustNullErr(errCanonical, `${name} canonical state`);
    if (canonical !== state) throw new Error(`${name}: non-canonical state ${state}, want ${canonical}`);
    const out = new _.Buffer();
    mustNullErr(_.setU8(out, canonical), `${name} encode state`);
    mustNullErr(_.setBin(out, canonical, value), `${name} encode body`);
    return toHex(out.bytes);
};

const roundTripCase = (entry: InputCase): string => {
    const bytes = fromHex(entry.hex);
    switch (entry.kind) {
    case "struct.recharge":
        return roundTripStruct(entry.name, bytes, _.readRecharge, _.setRecharge);
    case "struct.rechargeA":
        return roundTripStruct(entry.name, bytes, _.readRechargeA, _.setRechargeA);
    case "struct.rechargeB":
        return roundTripStruct(entry.name, bytes, _.readRechargeB, _.setRechargeB);
    case "struct.simInfo":
        return roundTripStruct(entry.name, bytes, _.readSimInfo, _.setSimInfo);
    case "struct.sim":
        return roundTripStruct(entry.name, bytes, _.readSim, _.setSim);
    case "struct.simOrder2":
        return roundTripStruct(entry.name, bytes, _.readSimOrder2, _.setSimOrder2);
    case "struct.simOrder":
        return roundTripStruct(entry.name, bytes, _.readSimOrder, _.setSimOrder);
    case "api.user.get_abc.req":
    case "api.user.set_sim_info.resp":
        return roundTripEmpty(entry.name, bytes);
    case "api.user.get_abc.resp":
    case "api.user.get_abcd.resp":
        return roundTripOrderStatus(entry.name, bytes);
    case "api.user.get_abcd.req":
        return roundTripU8Pair(entry.name, bytes);
    case "api.user.set_sim_info.req":
        return roundTripStruct(entry.name, bytes, _.readSimInfo, _.setSimInfo);
    case "api.get_count.req":
    case "api.get_count.resp":
    case "api.get_bin.req":
        return roundTripU8(entry.name, bytes);
    case "api.get_bin.resp":
        return roundTripBinResp(entry.name, bytes);
    default:
        throw new Error(`${entry.name}: unsupported kind ${entry.kind}`);
    }
};

const main = async (): Promise<void> => {
    const inputPath = Bun.argv[2];
    if (!inputPath) throw new Error("missing input path");
    const text = await Bun.file(inputPath).text();
    const cases = JSON.parse(text) as InputCase[];
    const output: OutputCase[] = cases.map((entry) => ({
        name: entry.name,
        hex: roundTripCase(entry),
    }));
    process.stdout.write(`${JSON.stringify(output)}\n`);
};

await main();
