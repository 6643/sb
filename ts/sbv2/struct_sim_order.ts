import * as rt from "./type"
import * as _ from "./_"

export interface SimOrder extends rt.Serializable, rt.Deserializable {
    id: number;
    accountId: number;
    itemId: number;
    // 办理人姓名
    name: string;
    // 联系电话
    phone: string;
    // 身份证号
    idNo: string;
    // 所在城市
    cityCode: number;
    // 详细地址
    address: string;
    // 新手机号码
    newPhone: string;
    // 佣金
    commission: number;
    status: _.OrderStatus;
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
        status: _.DefaultOrderStatus(),
    } as any as SimOrder;
    s.set = (buf: rt.Buffer) => setSimOrder(buf, s);
    s.get = (buf: rt.Buffer) => {
        const [next, err] = getSimOrder(buf);
        if (err === null) Object.assign(s, next);
        return err;
    };
    return s;
};

export const isZeroSimOrder = (s: SimOrder | null | undefined): boolean => {
    if (s === null || s === undefined) return true;
    if (s.id !== 0) return false;
    if (s.accountId !== 0) return false;
    if (s.itemId !== 0) return false;
    if (s.name !== "") return false;
    if (s.phone !== "") return false;
    if (s.idNo !== "") return false;
    if (s.cityCode !== 0) return false;
    if (s.address !== "") return false;
    if (s.newPhone !== "") return false;
    if (s.commission !== 0) return false;
    if (!_.IsDefaultOrderStatus(s.status as any)) return false;
    return true;
};

export const validateSimOrder = (s: SimOrder | null | undefined): Error | null => {
    if (s === null || s === undefined) return null;
    { const [, err] = rt.textState(s.name); if (err !== null) return new Error(`Name: ${err.message}`); }
    { const [, err] = rt.textState(s.phone); if (err !== null) return new Error(`Phone: ${err.message}`); }
    { const [, err] = rt.textState(s.idNo); if (err !== null) return new Error(`IdNo: ${err.message}`); }
    { const [, err] = rt.textState(s.address); if (err !== null) return new Error(`Address: ${err.message}`); }
    { const [, err] = rt.textState(s.newPhone); if (err !== null) return new Error(`NewPhone: ${err.message}`); }
    if (!_.IsAssignableOrderStatus(s.status as any)) return new Error(`Status 非法枚举值: ${s.status as any}`);
    return null;
};

export const getSimOrder = (buf: rt.Buffer): [SimOrder, Error | null] => {
    const s = newSimOrder();
    const headerBits = 16;
    const [header, errHeader] = buf.read(rt.headerSize(headerBits));
    if (errHeader !== null) return [s, new Error(`not enough data`)];
    const errPadding = rt.validatePaddingZero(header, headerBits, "SimOrder header");
    if (errPadding !== null) return [s, errPadding];
    const reader = new rt.BitReader(header, headerBits);
    const [idPresent, errIdPresent] = reader.readBit();
    if (errIdPresent !== null) return [s, new Error(`get SimOrder Id header: ${errIdPresent.message}`)];
    const [accountIdPresent, errAccountIdPresent] = reader.readBit();
    if (errAccountIdPresent !== null) return [s, new Error(`get SimOrder AccountId header: ${errAccountIdPresent.message}`)];
    const [itemIdPresent, errItemIdPresent] = reader.readBit();
    if (errItemIdPresent !== null) return [s, new Error(`get SimOrder ItemId header: ${errItemIdPresent.message}`)];
    const [nameState, errNameState] = reader.readBits(2);
    if (errNameState !== null) return [s, new Error(`get SimOrder Name header: ${errNameState.message}`)];
    const [phoneState, errPhoneState] = reader.readBits(2);
    if (errPhoneState !== null) return [s, new Error(`get SimOrder Phone header: ${errPhoneState.message}`)];
    const [idNoState, errIdNoState] = reader.readBits(2);
    if (errIdNoState !== null) return [s, new Error(`get SimOrder IdNo header: ${errIdNoState.message}`)];
    const [cityCodePresent, errCityCodePresent] = reader.readBit();
    if (errCityCodePresent !== null) return [s, new Error(`get SimOrder CityCode header: ${errCityCodePresent.message}`)];
    const [addressState, errAddressState] = reader.readBits(2);
    if (errAddressState !== null) return [s, new Error(`get SimOrder Address header: ${errAddressState.message}`)];
    const [newPhoneState, errNewPhoneState] = reader.readBits(2);
    if (errNewPhoneState !== null) return [s, new Error(`get SimOrder NewPhone header: ${errNewPhoneState.message}`)];
    const [commissionPresent, errCommissionPresent] = reader.readBit();
    if (errCommissionPresent !== null) return [s, new Error(`get SimOrder Commission header: ${errCommissionPresent.message}`)];
    const [statusPresent, errStatusPresent] = reader.readBit();
    if (errStatusPresent !== null) return [s, new Error(`get SimOrder Status header: ${errStatusPresent.message}`)];
    if (idPresent) {
        const [value, err] = rt.getU32(buf);
        if (err !== null) return [s, new Error(`get SimOrder Id: ${err.message}`)];
        s.id = value as any;
    }
    if (accountIdPresent) {
        const [value, err] = rt.getU32(buf);
        if (err !== null) return [s, new Error(`get SimOrder AccountId: ${err.message}`)];
        s.accountId = value as any;
    }
    if (itemIdPresent) {
        const [value, err] = rt.getU32(buf);
        if (err !== null) return [s, new Error(`get SimOrder ItemId: ${err.message}`)];
        s.itemId = value as any;
    }
    { const [value, err] = rt.getTextCompact(buf, nameState); if (err !== null) return [s, new Error(`get SimOrder Name: ${err.message}`)]; s.name = value; }
    { const [value, err] = rt.getTextCompact(buf, phoneState); if (err !== null) return [s, new Error(`get SimOrder Phone: ${err.message}`)]; s.phone = value; }
    { const [value, err] = rt.getTextCompact(buf, idNoState); if (err !== null) return [s, new Error(`get SimOrder IdNo: ${err.message}`)]; s.idNo = value; }
    if (cityCodePresent) {
        const [value, err] = rt.getU32(buf);
        if (err !== null) return [s, new Error(`get SimOrder CityCode: ${err.message}`)];
        s.cityCode = value as any;
    }
    { const [value, err] = rt.getTextCompact(buf, addressState); if (err !== null) return [s, new Error(`get SimOrder Address: ${err.message}`)]; s.address = value; }
    { const [value, err] = rt.getTextCompact(buf, newPhoneState); if (err !== null) return [s, new Error(`get SimOrder NewPhone: ${err.message}`)]; s.newPhone = value; }
    if (commissionPresent) {
        const [value, err] = rt.getU16(buf);
        if (err !== null) return [s, new Error(`get SimOrder Commission: ${err.message}`)];
        s.commission = value as any;
    }
    s.status = _.DefaultOrderStatus();
    if (statusPresent) {
        const [value, err] = rt.getU8(buf);
        if (err !== null) return [s, new Error(`get SimOrder Status: ${err.message}`)];
        const item = value as _.OrderStatus;
        if (!_.IsOrderStatus(item)) return [s, new Error(`get SimOrder Status: 非法枚举值: ${item}`)];
        s.status = item;
    }
    const errValidate = validateSimOrder(s);
    if (errValidate !== null) return [s, new Error(`validate failed: ${errValidate.message}`)];
    return [s, null];
};

export const setSimOrder = (buf: rt.Buffer, s: SimOrder): Error | null => {
    if (s === null || s === undefined) return new Error(`set SimOrder: value is null or undefined`);
    const errValidate = validateSimOrder(s);
    if (errValidate !== null) return new Error(`validate SimOrder: ${errValidate.message}`);
    const startOffset = buf.writeOffset;
    const [nameState, errNameState] = rt.textState(s.name);
    if (errNameState !== null) return errNameState;
    const [phoneState, errPhoneState] = rt.textState(s.phone);
    if (errPhoneState !== null) return errPhoneState;
    const [idNoState, errIdNoState] = rt.textState(s.idNo);
    if (errIdNoState !== null) return errIdNoState;
    const [addressState, errAddressState] = rt.textState(s.address);
    if (errAddressState !== null) return errAddressState;
    const [newPhoneState, errNewPhoneState] = rt.textState(s.newPhone);
    if (errNewPhoneState !== null) return errNewPhoneState;
    const header = new rt.BitWriter(16);
    header.writeBit(s.id !== 0);
    header.writeBit(s.accountId !== 0);
    header.writeBit(s.itemId !== 0);
    { const err = header.writeBits(nameState, 2); if (err !== null) { buf.rewindWrite(startOffset); return new Error(`set SimOrder Name header: ${err.message}`); } }
    { const err = header.writeBits(phoneState, 2); if (err !== null) { buf.rewindWrite(startOffset); return new Error(`set SimOrder Phone header: ${err.message}`); } }
    { const err = header.writeBits(idNoState, 2); if (err !== null) { buf.rewindWrite(startOffset); return new Error(`set SimOrder IdNo header: ${err.message}`); } }
    header.writeBit(s.cityCode !== 0);
    { const err = header.writeBits(addressState, 2); if (err !== null) { buf.rewindWrite(startOffset); return new Error(`set SimOrder Address header: ${err.message}`); } }
    { const err = header.writeBits(newPhoneState, 2); if (err !== null) { buf.rewindWrite(startOffset); return new Error(`set SimOrder NewPhone header: ${err.message}`); } }
    header.writeBit(s.commission !== 0);
    header.writeBit(!_.IsDefaultOrderStatus(s.status as any));
    const errHeaderWrite = buf.write(header.bytes);
    if (errHeaderWrite !== null) { buf.rewindWrite(startOffset); return errHeaderWrite; }
    if (s.id !== 0) {
        const err = rt.setU32(buf, s.id as any);
        if (err !== null) { buf.rewindWrite(startOffset); return new Error(`set SimOrder Id: ${err.message}`); }
    }
    if (s.accountId !== 0) {
        const err = rt.setU32(buf, s.accountId as any);
        if (err !== null) { buf.rewindWrite(startOffset); return new Error(`set SimOrder AccountId: ${err.message}`); }
    }
    if (s.itemId !== 0) {
        const err = rt.setU32(buf, s.itemId as any);
        if (err !== null) { buf.rewindWrite(startOffset); return new Error(`set SimOrder ItemId: ${err.message}`); }
    }
    { const err = rt.setTextCompact(buf, nameState, s.name); if (err !== null) { buf.rewindWrite(startOffset); return new Error(`set SimOrder Name: ${err.message}`); } }
    { const err = rt.setTextCompact(buf, phoneState, s.phone); if (err !== null) { buf.rewindWrite(startOffset); return new Error(`set SimOrder Phone: ${err.message}`); } }
    { const err = rt.setTextCompact(buf, idNoState, s.idNo); if (err !== null) { buf.rewindWrite(startOffset); return new Error(`set SimOrder IdNo: ${err.message}`); } }
    if (s.cityCode !== 0) {
        const err = rt.setU32(buf, s.cityCode as any);
        if (err !== null) { buf.rewindWrite(startOffset); return new Error(`set SimOrder CityCode: ${err.message}`); }
    }
    { const err = rt.setTextCompact(buf, addressState, s.address); if (err !== null) { buf.rewindWrite(startOffset); return new Error(`set SimOrder Address: ${err.message}`); } }
    { const err = rt.setTextCompact(buf, newPhoneState, s.newPhone); if (err !== null) { buf.rewindWrite(startOffset); return new Error(`set SimOrder NewPhone: ${err.message}`); } }
    if (s.commission !== 0) {
        const err = rt.setU16(buf, s.commission as any);
        if (err !== null) { buf.rewindWrite(startOffset); return new Error(`set SimOrder Commission: ${err.message}`); }
    }
    if (!_.IsDefaultOrderStatus(s.status as any)) {
        const err = rt.setU8(buf, _.NormalizeOrderStatus(s.status as any) as any);
        if (err !== null) { buf.rewindWrite(startOffset); return new Error(`set SimOrder Status: ${err.message}`); }
    }
    return null;
};

export const readSimOrder = (buf: rt.Buffer): [SimOrder, Error | null] => getSimOrder(buf);

export const eqSimOrder = (a: SimOrder | null | undefined, b: SimOrder | null | undefined): boolean => {
    if (isZeroSimOrder(a as any) && isZeroSimOrder(b as any)) return true;
    if (a === null || a === undefined || b === null || b === undefined) return false;
    if (!rt.eqU32(a.id, b.id)) return false;
    if (!rt.eqU32(a.accountId, b.accountId)) return false;
    if (!rt.eqU32(a.itemId, b.itemId)) return false;
    if (!rt.eqText(a.name, b.name)) return false;
    if (!rt.eqText(a.phone, b.phone)) return false;
    if (!rt.eqText(a.idNo, b.idNo)) return false;
    if (!rt.eqU32(a.cityCode, b.cityCode)) return false;
    if (!rt.eqText(a.address, b.address)) return false;
    if (!rt.eqText(a.newPhone, b.newPhone)) return false;
    if (!rt.eqU16(a.commission, b.commission)) return false;
    if (!_.eqOrderStatusValue(a.status as any, b.status as any)) return false;
    return true;
};

export const getSimOrderListBody = (buf: rt.Buffer, state: number): [SimOrder[], Error | null] =>
    rt.getBitmapListCompact<SimOrder>(
        buf,
        state,
        () => newSimOrder(),
        (buf) => readSimOrder(buf),
    );

export const setSimOrderListBody = (buf: rt.Buffer, state: number, v: SimOrder[]): Error | null =>
    rt.setBitmapListCompact<SimOrder>(
        buf,
        state,
        v,
        (item) => isZeroSimOrder(item),
        (buf, item) => setSimOrder(buf, item),
    );
