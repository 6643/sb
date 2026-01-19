import * as _ from "./_.ts"
import * as Enum from "./enum"

export interface SimOrder extends _.Serializable, _.Deserializable {
    id: number;
    accountId: number;
    itemId: number;
    name: string;
    phone: string;
    idNo: string;
    cityCode: number;
    address: string;
    newPhone: string;
    commission: number;
    status: Enum.OrderStatus;
}

export const newSimOrder = (): SimOrder => {
    const s = {
        id: 0,
        accountId: 0,
        itemId: 0,
        name: "",
        phone: "",
        idNo: "",
        cityCode: 0,
        address: "",
        newPhone: "",
        commission: 0,
        status: 0,
    } as any as SimOrder;
    s.set = (buf: _.Buffer) => setSimOrder(buf, s);
    s.get = (buf: _.Buffer) => {
        const [res, err] = getSimOrder(buf);
        if (err === null) Object.assign(s, res);
        return err;
    };
    return s;
}

export const eqSimOrder = (a: SimOrder, b: SimOrder): boolean => {
    if (a === b) return true;
    if (a === null || b === null) return false;
    if (!_.eqU32(a.id, b.id)) return false;
    if (!_.eqU32(a.accountId, b.accountId)) return false;
    if (!_.eqU32(a.itemId, b.itemId)) return false;
    if (!_.eqText(a.name, b.name)) return false;
    if (!_.eqText(a.phone, b.phone)) return false;
    if (!_.eqText(a.idNo, b.idNo)) return false;
    if (!_.eqU32(a.cityCode, b.cityCode)) return false;
    if (!_.eqText(a.address, b.address)) return false;
    if (!_.eqText(a.newPhone, b.newPhone)) return false;
    if (!_.eqU16(a.commission, b.commission)) return false;
    if (a.status !== b.status) return false;
    return true;
}

export const getSimOrder = (buf: _.Buffer): [SimOrder, Error | null] => {
    const s = newSimOrder();
    const bitmaskSize = Math.ceil(11 / 8);
    const [bits, err] = buf.read(bitmaskSize);
    if (err !== null) return [s, err];
    if (_.GetBit(bits, 0)) {
        const [v, err] = _.getU32(buf);
        if (err !== null) return [s, err];
        s.id = v;
    }
    if (_.GetBit(bits, 1)) {
        const [v, err] = _.getU32(buf);
        if (err !== null) return [s, err];
        s.accountId = v;
    }
    if (_.GetBit(bits, 2)) {
        const [v, err] = _.getU32(buf);
        if (err !== null) return [s, err];
        s.itemId = v;
    }
    if (_.GetBit(bits, 3)) {
        const [v, err] = _.getText(buf);
        if (err !== null) return [s, err];
        s.name = v;
    }
    if (_.GetBit(bits, 4)) {
        const [v, err] = _.getText(buf);
        if (err !== null) return [s, err];
        s.phone = v;
    }
    if (_.GetBit(bits, 5)) {
        const [v, err] = _.getText(buf);
        if (err !== null) return [s, err];
        s.idNo = v;
    }
    if (_.GetBit(bits, 6)) {
        const [v, err] = _.getU32(buf);
        if (err !== null) return [s, err];
        s.cityCode = v;
    }
    if (_.GetBit(bits, 7)) {
        const [v, err] = _.getText(buf);
        if (err !== null) return [s, err];
        s.address = v;
    }
    if (_.GetBit(bits, 8)) {
        const [v, err] = _.getText(buf);
        if (err !== null) return [s, err];
        s.newPhone = v;
    }
    if (_.GetBit(bits, 9)) {
        const [v, err] = _.getU16(buf);
        if (err !== null) return [s, err];
        s.commission = v;
    }
    if (_.GetBit(bits, 10)) {
        const [v, err] = _.getU8(buf);
        if (err !== null) return [s, err];
        s.status = v as any;
    }
    return [s, null];
}

export const setSimOrder = (buf: _.Buffer, s: SimOrder): Error | null => {
    if (s === null || s === undefined) return new Error(`set SimOrder: value is null or undefined`);
    const bits = new Uint8Array(Math.ceil(11 / 8));
    const body = new _.Buffer();
    if (!_.eqU32(s.id, 0)) {
        const err = _.setU32(body, s.id);
        if (err !== null) return err;
        _.SetBit(bits, 0, true);
    }
    if (!_.eqU32(s.accountId, 0)) {
        const err = _.setU32(body, s.accountId);
        if (err !== null) return err;
        _.SetBit(bits, 1, true);
    }
    if (!_.eqU32(s.itemId, 0)) {
        const err = _.setU32(body, s.itemId);
        if (err !== null) return err;
        _.SetBit(bits, 2, true);
    }
    if (!_.eqText(s.name, "")) {
        const err = _.setText(body, s.name);
        if (err !== null) return err;
        _.SetBit(bits, 3, true);
    }
    if (!_.eqText(s.phone, "")) {
        const err = _.setText(body, s.phone);
        if (err !== null) return err;
        _.SetBit(bits, 4, true);
    }
    if (!_.eqText(s.idNo, "")) {
        const err = _.setText(body, s.idNo);
        if (err !== null) return err;
        _.SetBit(bits, 5, true);
    }
    if (!_.eqU32(s.cityCode, 0)) {
        const err = _.setU32(body, s.cityCode);
        if (err !== null) return err;
        _.SetBit(bits, 6, true);
    }
    if (!_.eqText(s.address, "")) {
        const err = _.setText(body, s.address);
        if (err !== null) return err;
        _.SetBit(bits, 7, true);
    }
    if (!_.eqText(s.newPhone, "")) {
        const err = _.setText(body, s.newPhone);
        if (err !== null) return err;
        _.SetBit(bits, 8, true);
    }
    if (!_.eqU16(s.commission, 0)) {
        const err = _.setU16(body, s.commission);
        if (err !== null) return err;
        _.SetBit(bits, 9, true);
    }
    if ((s.status as any) !== 0) {
        const err = _.setU8(body, s.status as any);
        if (err !== null) return err;
        _.SetBit(bits, 10, true);
    }

    const errBits = buf.write(bits);
    if (errBits !== null) return errBits;
    return buf.write(body.bytes);
}

export const getSimOrderList = (buf: _.Buffer): [SimOrder[], Error | null] => _.getList(buf, getSimOrder);
export const setSimOrderList = (buf: _.Buffer, v: SimOrder[]): Error | null => _.setList(buf, v, setSimOrder);
export const eqSimOrderList = (a: SimOrder[], b: SimOrder[]): boolean => _.eqList(a, b, eqSimOrder);