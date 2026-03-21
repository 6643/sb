import * as rt from "./type"
import * as _ from "./_"

export interface Recharge extends rt.Serializable, rt.Deserializable {
    // abcd
    id: number;
    type: _.OrderStatus[];
    phone: string[];
    si: _.SimInfo | null;
}

const rechargeHeaderWidths = [1, 2, 2, 1] as const;

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
        if (err === undefined) Object.assign(s, next);
        return err;
    };
    return s;
};

export const isZeroRecharge = (s: Recharge | null | undefined): boolean => {
    if (s === null || s === undefined) return true;
    if (s.id !== 0) return false;
    if (!rt.isArrayValue(s.type) || s.type.length !== 0) return false;
    if (!rt.isArrayValue(s.phone) || s.phone.length !== 0) return false;
    if (!_.isZeroSimInfo(s.si as any)) return false;
    return true;
};

export const validateRecharge = (s: Recharge | null | undefined): rt.Err => {
    if (s === null || s === undefined) return undefined;
    if (!rt.isArrayValue(s.type)) return new Error(`Type: value is not array`);
    for (let i = 0; i < s.type.length; i++) {
        if (!_.IsAssignableOrderStatus(s.type[i] as any)) return new Error(`Type[${i}] 非法枚举值: ${s.type[i] as any}`);
    }
    { const [, err] = rt.listCountState(s.type.length); if (err !== null) return new Error(`Type: ${err.message}`); }
    if (!rt.isArrayValue(s.phone)) return new Error(`Phone: value is not array`);
    { const [, err] = rt.listCountState(s.phone.length); if (err !== null) return new Error(`Phone: ${err.message}`); }
    for (let i = 0; i < s.phone.length; i++) {
        const [, err] = rt.textState(s.phone[i]);
        if (err !== null) return new Error(`Phone[${i}]: ${err.message}`);
    }
    { const err = _.validateSimInfo(s.si); if (err !== undefined) return new Error(`Si: ${err.message}`); }
    return undefined;
};

export const getRecharge = (buf: rt.Buffer): [Recharge, rt.Err] => {
    const s = newRecharge();
    const headerBits = 6;
    const [header, errHeader] = buf.read(rt.headerSize(headerBits));
    if (errHeader !== null) return [s, new Error(`not enough data`)];
    const [headerStates, errHeaderState] = rt.readHeader(header, rechargeHeaderWidths, "Recharge header");
    if (errHeaderState !== undefined) return [s, errHeaderState];
    const idPresent = headerStates[0] === 1;
    const typeState = headerStates[1];
    const phoneState = headerStates[2];
    const siPresent = headerStates[3] === 1;
    if (idPresent) {
        const [value, err] = rt.getU32(buf);
        if (err !== null) return [s, new Error(`get Recharge Id: ${err.message}`)];
        s.id = value as any;
    }
    { const [value, err] = _.getOrderStatusListBody(buf, typeState); if (err !== undefined) return [s, new Error(`get Recharge Type: ${err.message}`)]; s.type = value; }
    { const [value, err] = rt.getTextList(buf, phoneState); if (err !== null) return [s, new Error(`get Recharge Phone: ${err.message}`)]; s.phone = value; }
    if (siPresent) {
        const [value, err] = _.readSimInfo(buf);
        if (err !== undefined) return [s, new Error(`get Recharge Si: ${err.message}`)];
        s.si = value;
    }
    const errValidate = validateRecharge(s);
    if (errValidate !== undefined) return [s, new Error(`validate failed: ${errValidate.message}`)];
    return [s, undefined];
};

export const setRecharge = (buf: rt.Buffer, s: Recharge): rt.Err => {
    if (s === null || s === undefined) return new Error(`set Recharge: value is null or undefined`);
    const errValidate = validateRecharge(s);
    if (errValidate !== undefined) return new Error(`validate Recharge: ${errValidate.message}`);
    const startOffset = buf.writeOffset;
    const [typeState, errTypeState] = rt.listCountState(s.type.length);
    if (errTypeState !== null) return errTypeState;
    const [phoneState, errPhoneState] = rt.listCountState(s.phone.length);
    if (errPhoneState !== null) return errPhoneState;
    const headerStates = [
        s.id !== 0 ? 1 : 0,
        typeState,
        phoneState,
        !_.isZeroSimInfo(s.si as any) ? 1 : 0,
    ];
    const [header, errHeader] = rt.writeHeader(rechargeHeaderWidths, headerStates);
    if (errHeader !== undefined) { buf.rewindWrite(startOffset); return new Error(`set header: ${errHeader.message}`); }
    const errHeaderWrite = buf.write(header);
    if (errHeaderWrite !== null) { buf.rewindWrite(startOffset); return errHeaderWrite; }
    if (s.id !== 0) {
        const err = rt.setU32(buf, s.id as any);
        if (err !== null) { buf.rewindWrite(startOffset); return new Error(`set Recharge Id: ${err.message}`); }
    }
    { const err = _.setOrderStatusListBody(buf, typeState, s.type); if (err !== undefined) { buf.rewindWrite(startOffset); return new Error(`set Recharge Type: ${err.message}`); } }
    { const err = rt.setTextList(buf, phoneState, s.phone); if (err !== null) { buf.rewindWrite(startOffset); return new Error(`set Recharge Phone: ${err.message}`); } }
    if (!_.isZeroSimInfo(s.si)) {
        const err = _.setSimInfo(buf, s.si!);
        if (err !== undefined) { buf.rewindWrite(startOffset); return new Error(`set Recharge Si: ${err.message}`); }
    }
    return undefined;
};

export const readRecharge = (buf: rt.Buffer): [Recharge, rt.Err] => getRecharge(buf);

export const eqRecharge = (a: Recharge | null | undefined, b: Recharge | null | undefined): boolean => {
    if (isZeroRecharge(a as any) && isZeroRecharge(b as any)) return true;
    if (a === null || a === undefined || b === null || b === undefined) return false;
    if (!rt.eqU32(a.id, b.id)) return false;
    if (!_.eqOrderStatusList(a.type as any, b.type as any)) return false;
    if (!rt.eqList(a.phone, b.phone, rt.eqText)) return false;
    if (!_.eqSimInfo(a.si, b.si)) return false;
    return true;
};

export const getRechargeListBody = (buf: rt.Buffer, state: number): [Recharge[], rt.Err] => {
    const [list, err] = rt.getDefaultList<Recharge>(
        buf,
        state,
        () => newRecharge(),
        (buf) => readRecharge(buf),
    );
    return [list, err];
};

export const setRechargeListBody = (buf: rt.Buffer, state: number, v: Recharge[]): rt.Err => {
    return rt.setDefaultList<Recharge>(
        buf,
        state,
        v,
        (item) => isZeroRecharge(item),
        (buf, item) => setRecharge(buf, item),
    );
}
