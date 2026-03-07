import * as rt from "./type"
import * as _ from "./_"

export interface SimOrder2 extends rt.Serializable, rt.Deserializable {
    // SIM卡ID
    id: number;
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
    s.set = (buf: rt.Buffer) => setSimOrder2(buf, s);
    s.get = (buf: rt.Buffer) => {
        const [next, err] = getSimOrder2(buf);
        if (err === null) Object.assign(s, next);
        return err;
    };
    return s;
};

export const isZeroSimOrder2 = (s: SimOrder2 | null | undefined): boolean => {
    if (s === null || s === undefined) return true;
    if (s.id !== 0) return false;
    if (s.name !== "") return false;
    if (s.phone !== "") return false;
    if (s.idNo !== "") return false;
    if (s.cityCode !== 0) return false;
    if (s.address !== "") return false;
    if (s.newPhone !== "") return false;
    return true;
};

export const validateSimOrder2 = (s: SimOrder2 | null | undefined): Error | null => {
    if (s === null || s === undefined) return null;
    { const [, err] = rt.textState(s.name); if (err !== null) return new Error(`Name: ${err.message}`); }
    { const [, err] = rt.textState(s.phone); if (err !== null) return new Error(`Phone: ${err.message}`); }
    { const [, err] = rt.textState(s.idNo); if (err !== null) return new Error(`IdNo: ${err.message}`); }
    { const [, err] = rt.textState(s.address); if (err !== null) return new Error(`Address: ${err.message}`); }
    { const [, err] = rt.textState(s.newPhone); if (err !== null) return new Error(`NewPhone: ${err.message}`); }
    return null;
};

export const getSimOrder2 = (buf: rt.Buffer): [SimOrder2, Error | null] => {
    const s = newSimOrder2();
    const headerBits = 12;
    const [header, errHeader] = buf.read(rt.headerSize(headerBits));
    if (errHeader !== null) return [s, new Error(`not enough data`)];
    const errPadding = rt.validatePaddingZero(header, headerBits, "SimOrder2 header");
    if (errPadding !== null) return [s, errPadding];
    const reader = new rt.BitReader(header, headerBits);
    const [idPresent, errIdPresent] = reader.readBit();
    if (errIdPresent !== null) return [s, new Error(`get SimOrder2 Id header: ${errIdPresent.message}`)];
    const [nameState, errNameState] = reader.readBits(2);
    if (errNameState !== null) return [s, new Error(`get SimOrder2 Name header: ${errNameState.message}`)];
    const [phoneState, errPhoneState] = reader.readBits(2);
    if (errPhoneState !== null) return [s, new Error(`get SimOrder2 Phone header: ${errPhoneState.message}`)];
    const [idNoState, errIdNoState] = reader.readBits(2);
    if (errIdNoState !== null) return [s, new Error(`get SimOrder2 IdNo header: ${errIdNoState.message}`)];
    const [cityCodePresent, errCityCodePresent] = reader.readBit();
    if (errCityCodePresent !== null) return [s, new Error(`get SimOrder2 CityCode header: ${errCityCodePresent.message}`)];
    const [addressState, errAddressState] = reader.readBits(2);
    if (errAddressState !== null) return [s, new Error(`get SimOrder2 Address header: ${errAddressState.message}`)];
    const [newPhoneState, errNewPhoneState] = reader.readBits(2);
    if (errNewPhoneState !== null) return [s, new Error(`get SimOrder2 NewPhone header: ${errNewPhoneState.message}`)];
    if (idPresent) {
        const [value, err] = rt.getU32(buf);
        if (err !== null) return [s, new Error(`get SimOrder2 Id: ${err.message}`)];
        s.id = value as any;
    }
    { const [value, err] = rt.getTextCompact(buf, nameState); if (err !== null) return [s, new Error(`get SimOrder2 Name: ${err.message}`)]; s.name = value; }
    { const [value, err] = rt.getTextCompact(buf, phoneState); if (err !== null) return [s, new Error(`get SimOrder2 Phone: ${err.message}`)]; s.phone = value; }
    { const [value, err] = rt.getTextCompact(buf, idNoState); if (err !== null) return [s, new Error(`get SimOrder2 IdNo: ${err.message}`)]; s.idNo = value; }
    if (cityCodePresent) {
        const [value, err] = rt.getU32(buf);
        if (err !== null) return [s, new Error(`get SimOrder2 CityCode: ${err.message}`)];
        s.cityCode = value as any;
    }
    { const [value, err] = rt.getTextCompact(buf, addressState); if (err !== null) return [s, new Error(`get SimOrder2 Address: ${err.message}`)]; s.address = value; }
    { const [value, err] = rt.getTextCompact(buf, newPhoneState); if (err !== null) return [s, new Error(`get SimOrder2 NewPhone: ${err.message}`)]; s.newPhone = value; }
    const errValidate = validateSimOrder2(s);
    if (errValidate !== null) return [s, new Error(`validate failed: ${errValidate.message}`)];
    return [s, null];
};

export const setSimOrder2 = (buf: rt.Buffer, s: SimOrder2): Error | null => {
    if (s === null || s === undefined) return new Error(`set SimOrder2: value is null or undefined`);
    const errValidate = validateSimOrder2(s);
    if (errValidate !== null) return new Error(`validate SimOrder2: ${errValidate.message}`);
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
    const header = new rt.BitWriter(12);
    header.writeBit(s.id !== 0);
    { const err = header.writeBits(nameState, 2); if (err !== null) { buf.rewindWrite(startOffset); return new Error(`set SimOrder2 Name header: ${err.message}`); } }
    { const err = header.writeBits(phoneState, 2); if (err !== null) { buf.rewindWrite(startOffset); return new Error(`set SimOrder2 Phone header: ${err.message}`); } }
    { const err = header.writeBits(idNoState, 2); if (err !== null) { buf.rewindWrite(startOffset); return new Error(`set SimOrder2 IdNo header: ${err.message}`); } }
    header.writeBit(s.cityCode !== 0);
    { const err = header.writeBits(addressState, 2); if (err !== null) { buf.rewindWrite(startOffset); return new Error(`set SimOrder2 Address header: ${err.message}`); } }
    { const err = header.writeBits(newPhoneState, 2); if (err !== null) { buf.rewindWrite(startOffset); return new Error(`set SimOrder2 NewPhone header: ${err.message}`); } }
    const errHeaderWrite = buf.write(header.bytes);
    if (errHeaderWrite !== null) { buf.rewindWrite(startOffset); return errHeaderWrite; }
    if (s.id !== 0) {
        const err = rt.setU32(buf, s.id as any);
        if (err !== null) { buf.rewindWrite(startOffset); return new Error(`set SimOrder2 Id: ${err.message}`); }
    }
    { const err = rt.setTextCompact(buf, nameState, s.name); if (err !== null) { buf.rewindWrite(startOffset); return new Error(`set SimOrder2 Name: ${err.message}`); } }
    { const err = rt.setTextCompact(buf, phoneState, s.phone); if (err !== null) { buf.rewindWrite(startOffset); return new Error(`set SimOrder2 Phone: ${err.message}`); } }
    { const err = rt.setTextCompact(buf, idNoState, s.idNo); if (err !== null) { buf.rewindWrite(startOffset); return new Error(`set SimOrder2 IdNo: ${err.message}`); } }
    if (s.cityCode !== 0) {
        const err = rt.setU32(buf, s.cityCode as any);
        if (err !== null) { buf.rewindWrite(startOffset); return new Error(`set SimOrder2 CityCode: ${err.message}`); }
    }
    { const err = rt.setTextCompact(buf, addressState, s.address); if (err !== null) { buf.rewindWrite(startOffset); return new Error(`set SimOrder2 Address: ${err.message}`); } }
    { const err = rt.setTextCompact(buf, newPhoneState, s.newPhone); if (err !== null) { buf.rewindWrite(startOffset); return new Error(`set SimOrder2 NewPhone: ${err.message}`); } }
    return null;
};

export const readSimOrder2 = (buf: rt.Buffer): [SimOrder2, Error | null] => getSimOrder2(buf);

export const eqSimOrder2 = (a: SimOrder2 | null | undefined, b: SimOrder2 | null | undefined): boolean => {
    if (isZeroSimOrder2(a as any) && isZeroSimOrder2(b as any)) return true;
    if (a === null || a === undefined || b === null || b === undefined) return false;
    if (!rt.eqU32(a.id, b.id)) return false;
    if (!rt.eqText(a.name, b.name)) return false;
    if (!rt.eqText(a.phone, b.phone)) return false;
    if (!rt.eqText(a.idNo, b.idNo)) return false;
    if (!rt.eqU32(a.cityCode, b.cityCode)) return false;
    if (!rt.eqText(a.address, b.address)) return false;
    if (!rt.eqText(a.newPhone, b.newPhone)) return false;
    return true;
};

export const getSimOrder2ListBody = (buf: rt.Buffer, state: number): [SimOrder2[], Error | null] =>
    rt.getBitmapListCompact<SimOrder2>(
        buf,
        state,
        () => newSimOrder2(),
        (buf) => readSimOrder2(buf),
    );

export const setSimOrder2ListBody = (buf: rt.Buffer, state: number, v: SimOrder2[]): Error | null =>
    rt.setBitmapListCompact<SimOrder2>(
        buf,
        state,
        v,
        (item) => isZeroSimOrder2(item),
        (buf, item) => setSimOrder2(buf, item),
    );
