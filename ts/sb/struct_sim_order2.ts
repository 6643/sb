import * as _ from "./_.ts"
import * as Enum from "./enum"

export interface SimOrder2 extends _.Serializable, _.Deserializable {
    id: number;
    name: string;
    phone: string;
    idNo: string;
    cityCode: number;
    address: string;
    newPhone: string;
}

export const newSimOrder2 = (): SimOrder2 => {
    const s = {
        id: 0,
        name: "",
        phone: "",
        idNo: "",
        cityCode: 0,
        address: "",
        newPhone: "",
    } as any as SimOrder2;
    s.set = (buf: _.Buffer) => setSimOrder2(buf, s);
    s.get = (buf: _.Buffer) => {
        const [res, err] = getSimOrder2(buf);
        if (err === null) Object.assign(s, res);
        return err;
    };
    return s;
}

export const eqSimOrder2 = (a: SimOrder2, b: SimOrder2): boolean => {
    if (a === b) return true;
    if (a === null || b === null) return false;
    if (!_.eqU32(a.id, b.id)) return false;
    if (!_.eqText(a.name, b.name)) return false;
    if (!_.eqText(a.phone, b.phone)) return false;
    if (!_.eqText(a.idNo, b.idNo)) return false;
    if (!_.eqU32(a.cityCode, b.cityCode)) return false;
    if (!_.eqText(a.address, b.address)) return false;
    if (!_.eqText(a.newPhone, b.newPhone)) return false;
    return true;
}

export const getSimOrder2 = (buf: _.Buffer): [SimOrder2, Error | null] => {
    const s = newSimOrder2();
    const bitmaskSize = Math.ceil(7 / 8);
    const [bits, err] = buf.read(bitmaskSize);
    if (err !== null) return [s, err];
    if (_.GetBit(bits, 0)) {
        const [v, err] = _.getU32(buf);
        if (err !== null) return [s, err];
        s.id = v;
    }
    if (_.GetBit(bits, 1)) {
        const [v, err] = _.getText(buf);
        if (err !== null) return [s, err];
        s.name = v;
    }
    if (_.GetBit(bits, 2)) {
        const [v, err] = _.getText(buf);
        if (err !== null) return [s, err];
        s.phone = v;
    }
    if (_.GetBit(bits, 3)) {
        const [v, err] = _.getText(buf);
        if (err !== null) return [s, err];
        s.idNo = v;
    }
    if (_.GetBit(bits, 4)) {
        const [v, err] = _.getU32(buf);
        if (err !== null) return [s, err];
        s.cityCode = v;
    }
    if (_.GetBit(bits, 5)) {
        const [v, err] = _.getText(buf);
        if (err !== null) return [s, err];
        s.address = v;
    }
    if (_.GetBit(bits, 6)) {
        const [v, err] = _.getText(buf);
        if (err !== null) return [s, err];
        s.newPhone = v;
    }
    return [s, null];
}

export const setSimOrder2 = (buf: _.Buffer, s: SimOrder2): Error | null => {
    if (s === null || s === undefined) return new Error(`set SimOrder2: value is null or undefined`);
    const bits = new Uint8Array(Math.ceil(7 / 8));
    const body = new _.Buffer();
    if (!_.eqU32(s.id, 0)) {
        const err = _.setU32(body, s.id);
        if (err !== null) return err;
        _.SetBit(bits, 0, true);
    }
    if (!_.eqText(s.name, "")) {
        const err = _.setText(body, s.name);
        if (err !== null) return err;
        _.SetBit(bits, 1, true);
    }
    if (!_.eqText(s.phone, "")) {
        const err = _.setText(body, s.phone);
        if (err !== null) return err;
        _.SetBit(bits, 2, true);
    }
    if (!_.eqText(s.idNo, "")) {
        const err = _.setText(body, s.idNo);
        if (err !== null) return err;
        _.SetBit(bits, 3, true);
    }
    if (!_.eqU32(s.cityCode, 0)) {
        const err = _.setU32(body, s.cityCode);
        if (err !== null) return err;
        _.SetBit(bits, 4, true);
    }
    if (!_.eqText(s.address, "")) {
        const err = _.setText(body, s.address);
        if (err !== null) return err;
        _.SetBit(bits, 5, true);
    }
    if (!_.eqText(s.newPhone, "")) {
        const err = _.setText(body, s.newPhone);
        if (err !== null) return err;
        _.SetBit(bits, 6, true);
    }

    const errBits = buf.write(bits);
    if (errBits !== null) return errBits;
    return buf.write(body.bytes);
}

export const getSimOrder2List = (buf: _.Buffer): [SimOrder2[], Error | null] => _.getList(buf, getSimOrder2);
export const setSimOrder2List = (buf: _.Buffer, v: SimOrder2[]): Error | null => _.setList(buf, v, setSimOrder2);
export const eqSimOrder2List = (a: SimOrder2[], b: SimOrder2[]): boolean => _.eqList(a, b, eqSimOrder2);