import * as rt from "./type"
import * as _ from "./_"

export interface Recharge extends rt.Serializable, rt.Deserializable {
    // abcd
    id: number;
    type: _.OrderStatus[];
    phone: string[];
    si: _.SimInfo | null;
}

export const newRecharge = (): Recharge => {
    const s = {
        id: 0,
        type: [],
        phone: [],
        si: null,
    } as any as Recharge;
    s.set = (buf: rt.Buffer) => setRecharge(buf, s);
    s.get = (buf: rt.Buffer) => {
        const [next, err] = getRecharge(buf);
        if (err === null) Object.assign(s, next);
        return err;
    };
    return s;
};

export const isZeroRecharge = (s: Recharge | null | undefined): boolean => {
    if (s === null || s === undefined) return true;
    if (s.id !== 0) return false;
    if (s.type.length !== 0) return false;
    if (s.phone.length !== 0) return false;
    if (!_.isZeroSimInfo(s.si)) return false;
    return true;
};

export const validateRecharge = (s: Recharge | null | undefined): Error | null => {
    if (s === null || s === undefined) return null;
    for (let i = 0; i < s.type.length; i++) {
        if (!_.IsAssignableOrderStatus(s.type[i] as any)) return new Error(`Type[${i}] 非法枚举值: ${s.type[i] as any}`);
    }
    { const [, err] = rt.listCountState(s.type.length); if (err !== null) return new Error(`Type: ${err.message}`); }
    { const [, err] = rt.listCountState(s.phone.length); if (err !== null) return new Error(`Phone: ${err.message}`); }
    for (let i = 0; i < s.phone.length; i++) {
        const [, err] = rt.textState(s.phone[i]);
        if (err !== null) return new Error(`Phone[${i}]: ${err.message}`);
    }
    { const err = _.validateSimInfo(s.si); if (err !== null) return new Error(`Si: ${err.message}`); }
    return null;
};

export const getRecharge = (buf: rt.Buffer): [Recharge, Error | null] => {
    const s = newRecharge();
    const headerBits = 6;
    const [header, errHeader] = buf.read(rt.headerSize(headerBits));
    if (errHeader !== null) return [s, new Error(`not enough data`)];
    const errPadding = rt.validatePaddingZero(header, headerBits, "Recharge header");
    if (errPadding !== null) return [s, errPadding];
    const reader = new rt.BitReader(header, headerBits);
    const [idPresent, errIdPresent] = reader.readBit();
    if (errIdPresent !== null) return [s, new Error(`get Recharge Id header: ${errIdPresent.message}`)];
    const [typeState, errTypeState] = reader.readBits(2);
    if (errTypeState !== null) return [s, new Error(`get Recharge Type header: ${errTypeState.message}`)];
    const [phoneState, errPhoneState] = reader.readBits(2);
    if (errPhoneState !== null) return [s, new Error(`get Recharge Phone header: ${errPhoneState.message}`)];
    const [siPresent, errSiPresent] = reader.readBit();
    if (errSiPresent !== null) return [s, new Error(`get Recharge Si header: ${errSiPresent.message}`)];
    if (idPresent) {
        const [value, err] = rt.getU32(buf);
        if (err !== null) return [s, new Error(`get Recharge Id: ${err.message}`)];
        s.id = value as any;
    }
    { const [value, err] = _.getOrderStatusListBody(buf, typeState); if (err !== null) return [s, new Error(`get Recharge Type: ${err.message}`)]; s.type = value; }
    { const [value, err] = rt.getTextListCompact(buf, phoneState); if (err !== null) return [s, new Error(`get Recharge Phone: ${err.message}`)]; s.phone = value; }
    if (siPresent) {
        const [value, err] = _.readSimInfo(buf);
        if (err !== null) return [s, new Error(`get Recharge Si: ${err.message}`)];
        s.si = value;
    }
    const errValidate = validateRecharge(s);
    if (errValidate !== null) return [s, new Error(`validate failed: ${errValidate.message}`)];
    return [s, null];
};

export const setRecharge = (buf: rt.Buffer, s: Recharge): Error | null => {
    if (s === null || s === undefined) return new Error(`set Recharge: value is null or undefined`);
    const errValidate = validateRecharge(s);
    if (errValidate !== null) return new Error(`validate Recharge: ${errValidate.message}`);
    const startOffset = buf.writeOffset;
    const [typeState, errTypeState] = rt.listCountState(s.type.length);
    if (errTypeState !== null) return errTypeState;
    const [phoneState, errPhoneState] = rt.listCountState(s.phone.length);
    if (errPhoneState !== null) return errPhoneState;
    const header = new rt.BitWriter(6);
    header.writeBit(s.id !== 0);
    { const err = header.writeBits(typeState, 2); if (err !== null) { buf.rewindWrite(startOffset); return new Error(`set Recharge Type header: ${err.message}`); } }
    { const err = header.writeBits(phoneState, 2); if (err !== null) { buf.rewindWrite(startOffset); return new Error(`set Recharge Phone header: ${err.message}`); } }
    header.writeBit(!_.isZeroSimInfo(s.si));
    const errHeaderWrite = buf.write(header.bytes);
    if (errHeaderWrite !== null) { buf.rewindWrite(startOffset); return errHeaderWrite; }
    if (s.id !== 0) {
        const err = rt.setU32(buf, s.id as any);
        if (err !== null) { buf.rewindWrite(startOffset); return new Error(`set Recharge Id: ${err.message}`); }
    }
    { const err = _.setOrderStatusListBody(buf, typeState, s.type); if (err !== null) { buf.rewindWrite(startOffset); return new Error(`set Recharge Type: ${err.message}`); } }
    { const err = rt.setTextListCompact(buf, phoneState, s.phone); if (err !== null) { buf.rewindWrite(startOffset); return new Error(`set Recharge Phone: ${err.message}`); } }
    if (!_.isZeroSimInfo(s.si)) {
        const err = _.setSimInfo(buf, s.si!);
        if (err !== null) { buf.rewindWrite(startOffset); return new Error(`set Recharge Si: ${err.message}`); }
    }
    return null;
};

export const readRecharge = (buf: rt.Buffer): [Recharge, Error | null] => getRecharge(buf);

export const eqRecharge = (a: Recharge | null | undefined, b: Recharge | null | undefined): boolean => {
    if (isZeroRecharge(a as any) && isZeroRecharge(b as any)) return true;
    if (a === null || a === undefined || b === null || b === undefined) return false;
    if (!rt.eqU32(a.id, b.id)) return false;
    if (!_.eqOrderStatusList(a.type as any, b.type as any)) return false;
    if (!rt.eqList(a.phone, b.phone, rt.eqText)) return false;
    if (!_.eqSimInfo(a.si, b.si)) return false;
    return true;
};

export const getRechargeListBody = (buf: rt.Buffer, state: number): [Recharge[], Error | null] =>
    rt.getBitmapListCompact<Recharge>(
        buf,
        state,
        () => newRecharge(),
        (buf) => readRecharge(buf),
    );

export const setRechargeListBody = (buf: rt.Buffer, state: number, v: Recharge[]): Error | null =>
    rt.setBitmapListCompact<Recharge>(
        buf,
        state,
        v,
        (item) => isZeroRecharge(item),
        (buf, item) => setRecharge(buf, item),
    );
