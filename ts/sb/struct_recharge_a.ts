import * as rt from "./type"
import * as _ from "./_"

export interface RechargeA extends rt.Serializable, rt.Deserializable {
    // abcd
    id: number;
    type: _.OrderStatus[];
    phone: string[];
    si: _.SimInfo | null;
    aid: number;
}

const rechargeAHeaderWidths = [1, 2, 2, 1, 1] as const;

export const newRechargeA = (): RechargeA => {
    const s = {
        id: 0,
        type: [],
        phone: [],
        si: null,
        aid: 0,
    } as any as RechargeA;
    s.set = (buf: rt.Buffer) => setRechargeA(buf, s);
    s.get = (buf: rt.Buffer) => {
        const [next, err] = getRechargeA(buf);
        if (err === undefined) Object.assign(s, next);
        return err;
    };
    return s;
};

export const isZeroRechargeA = (s: RechargeA | null | undefined): boolean => {
    if (s === null || s === undefined) return true;
    if (s.id !== 0) return false;
    if (s.type.length !== 0) return false;
    if (s.phone.length !== 0) return false;
    if (!_.isZeroSimInfo(s.si)) return false;
    if (s.aid !== 0) return false;
    return true;
};

export const validateRechargeA = (s: RechargeA | null | undefined): rt.Err => {
    if (s === null || s === undefined) return undefined;
    for (let i = 0; i < s.type.length; i++) {
        if (!_.IsAssignableOrderStatus(s.type[i] as any)) return new Error(`Type[${i}] 非法枚举值: ${s.type[i] as any}`);
    }
    { const [, err] = rt.listCountState(s.type.length); if (err !== null) return new Error(`Type: ${err.message}`); }
    { const [, err] = rt.listCountState(s.phone.length); if (err !== null) return new Error(`Phone: ${err.message}`); }
    for (let i = 0; i < s.phone.length; i++) {
        const [, err] = rt.textState(s.phone[i]);
        if (err !== null) return new Error(`Phone[${i}]: ${err.message}`);
    }
    { const err = _.validateSimInfo(s.si); if (err !== undefined) return new Error(`Si: ${err.message}`); }
    return undefined;
};

export const getRechargeA = (buf: rt.Buffer): [RechargeA, rt.Err] => {
    const s = newRechargeA();
    const headerBits = 7;
    const [header, errHeader] = buf.read(rt.headerSize(headerBits));
    if (errHeader !== null) return [s, new Error(`not enough data`)];
    const [headerStates, errHeaderState] = rt.readHeader(header, rechargeAHeaderWidths, "RechargeA header");
    if (errHeaderState !== undefined) return [s, errHeaderState];
    const idPresent = headerStates[0] === 1;
    const typeState = headerStates[1];
    const phoneState = headerStates[2];
    const siPresent = headerStates[3] === 1;
    const aidPresent = headerStates[4] === 1;
    if (idPresent) {
        const [value, err] = rt.getU32(buf);
        if (err !== null) return [s, new Error(`get RechargeA Id: ${err.message}`)];
        s.id = value as any;
    }
    { const [value, err] = _.getOrderStatusListBody(buf, typeState); if (err !== undefined) return [s, new Error(`get RechargeA Type: ${err.message}`)]; s.type = value; }
    { const [value, err] = rt.getTextListCompact(buf, phoneState); if (err !== null) return [s, new Error(`get RechargeA Phone: ${err.message}`)]; s.phone = value; }
    if (siPresent) {
        const [value, err] = _.readSimInfo(buf);
        if (err !== undefined) return [s, new Error(`get RechargeA Si: ${err.message}`)];
        s.si = value;
    }
    if (aidPresent) {
        const [value, err] = rt.getU32(buf);
        if (err !== null) return [s, new Error(`get RechargeA Aid: ${err.message}`)];
        s.aid = value as any;
    }
    const errValidate = validateRechargeA(s);
    if (errValidate !== undefined) return [s, new Error(`validate failed: ${errValidate.message}`)];
    return [s, undefined];
};

export const setRechargeA = (buf: rt.Buffer, s: RechargeA): rt.Err => {
    if (s === null || s === undefined) return new Error(`set RechargeA: value is null or undefined`);
    const errValidate = validateRechargeA(s);
    if (errValidate !== undefined) return new Error(`validate RechargeA: ${errValidate.message}`);
    const startOffset = buf.writeOffset;
    const [typeState, errTypeState] = rt.listCountState(s.type.length);
    if (errTypeState !== null) return errTypeState;
    const [phoneState, errPhoneState] = rt.listCountState(s.phone.length);
    if (errPhoneState !== null) return errPhoneState;
    const headerStates = [
        s.id !== 0 ? 1 : 0,
        typeState,
        phoneState,
        !_.isZeroSimInfo(s.si) ? 1 : 0,
        s.aid !== 0 ? 1 : 0,
    ];
    const [header, errHeader] = rt.writeHeader(rechargeAHeaderWidths, headerStates);
    if (errHeader !== undefined) { buf.rewindWrite(startOffset); return new Error(`set header: ${errHeader.message}`); }
    const errHeaderWrite = buf.write(header);
    if (errHeaderWrite !== null) { buf.rewindWrite(startOffset); return errHeaderWrite; }
    if (s.id !== 0) {
        const err = rt.setU32(buf, s.id as any);
        if (err !== null) { buf.rewindWrite(startOffset); return new Error(`set RechargeA Id: ${err.message}`); }
    }
    { const err = _.setOrderStatusListBody(buf, typeState, s.type); if (err !== undefined) { buf.rewindWrite(startOffset); return new Error(`set RechargeA Type: ${err.message}`); } }
    { const err = rt.setTextListCompact(buf, phoneState, s.phone); if (err !== null) { buf.rewindWrite(startOffset); return new Error(`set RechargeA Phone: ${err.message}`); } }
    if (!_.isZeroSimInfo(s.si)) {
        const err = _.setSimInfo(buf, s.si!);
        if (err !== undefined) { buf.rewindWrite(startOffset); return new Error(`set RechargeA Si: ${err.message}`); }
    }
    if (s.aid !== 0) {
        const err = rt.setU32(buf, s.aid as any);
        if (err !== null) { buf.rewindWrite(startOffset); return new Error(`set RechargeA Aid: ${err.message}`); }
    }
    return undefined;
};

export const readRechargeA = (buf: rt.Buffer): [RechargeA, rt.Err] => getRechargeA(buf);

export const eqRechargeA = (a: RechargeA | null | undefined, b: RechargeA | null | undefined): boolean => {
    if (isZeroRechargeA(a as any) && isZeroRechargeA(b as any)) return true;
    if (a === null || a === undefined || b === null || b === undefined) return false;
    if (!rt.eqU32(a.id, b.id)) return false;
    if (!_.eqOrderStatusList(a.type as any, b.type as any)) return false;
    if (!rt.eqList(a.phone, b.phone, rt.eqText)) return false;
    if (!_.eqSimInfo(a.si, b.si)) return false;
    if (!rt.eqU32(a.aid, b.aid)) return false;
    return true;
};

export const getRechargeAListBody = (buf: rt.Buffer, state: number): [RechargeA[], rt.Err] => {
    const [list, err] = rt.getBitmapListCompact<RechargeA>(
        buf,
        state,
        () => newRechargeA(),
        (buf) => readRechargeA(buf),
    );
    return [list, err];
};

export const setRechargeAListBody = (buf: rt.Buffer, state: number, v: RechargeA[]): rt.Err => {
    return rt.setBitmapListCompact<RechargeA>(
        buf,
        state,
        v,
        (item) => isZeroRechargeA(item),
        (buf, item) => setRechargeA(buf, item),
    );
}
