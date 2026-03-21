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

const simOrder2HeaderWidths = [1, 2, 2, 2, 1, 2, 2] as const;

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
        if (err === undefined) Object.assign(s, next);
        return err;
    };
    return s;
};

export const isZeroSimOrder2 = (s: SimOrder2 | null | undefined): boolean => {
    if (s === null || s === undefined) return true;
    if (s.id !== 0) return false;
    if (!rt.isStringValue(s.name) || s.name !== "") return false;
    if (!rt.isStringValue(s.phone) || s.phone !== "") return false;
    if (!rt.isStringValue(s.idNo) || s.idNo !== "") return false;
    if (s.cityCode !== 0) return false;
    if (!rt.isStringValue(s.address) || s.address !== "") return false;
    if (!rt.isStringValue(s.newPhone) || s.newPhone !== "") return false;
    return true;
};

export const validateSimOrder2 = (s: SimOrder2 | null | undefined): rt.Err => {
    if (s === null || s === undefined) return undefined;
    { const [, err] = rt.textState(s.name); if (err !== null) return new Error(`Name: ${err.message}`); }
    { const [, err] = rt.textState(s.phone); if (err !== null) return new Error(`Phone: ${err.message}`); }
    { const [, err] = rt.textState(s.idNo); if (err !== null) return new Error(`IdNo: ${err.message}`); }
    { const [, err] = rt.textState(s.address); if (err !== null) return new Error(`Address: ${err.message}`); }
    { const [, err] = rt.textState(s.newPhone); if (err !== null) return new Error(`NewPhone: ${err.message}`); }
    return undefined;
};

export const getSimOrder2 = (buf: rt.Buffer): [SimOrder2, rt.Err] => {
    const s = newSimOrder2();
    const headerBits = 12;
    const [header, errHeader] = buf.read(rt.headerSize(headerBits));
    if (errHeader !== null) return [s, new Error(`not enough data`)];
    const [headerStates, errHeaderState] = rt.readHeader(header, simOrder2HeaderWidths, "SimOrder2 header");
    if (errHeaderState !== undefined) return [s, errHeaderState];
    const idPresent = headerStates[0] === 1;
    const nameState = headerStates[1];
    const phoneState = headerStates[2];
    const idNoState = headerStates[3];
    const cityCodePresent = headerStates[4] === 1;
    const addressState = headerStates[5];
    const newPhoneState = headerStates[6];
    if (idPresent) {
        const [value, err] = rt.getU32(buf);
        if (err !== null) return [s, new Error(`get SimOrder2 Id: ${err.message}`)];
        s.id = value as any;
    }
    { const [value, err] = rt.getText(buf, nameState); if (err !== null) return [s, new Error(`get SimOrder2 Name: ${err.message}`)]; s.name = value; }
    { const [value, err] = rt.getText(buf, phoneState); if (err !== null) return [s, new Error(`get SimOrder2 Phone: ${err.message}`)]; s.phone = value; }
    { const [value, err] = rt.getText(buf, idNoState); if (err !== null) return [s, new Error(`get SimOrder2 IdNo: ${err.message}`)]; s.idNo = value; }
    if (cityCodePresent) {
        const [value, err] = rt.getU32(buf);
        if (err !== null) return [s, new Error(`get SimOrder2 CityCode: ${err.message}`)];
        s.cityCode = value as any;
    }
    { const [value, err] = rt.getText(buf, addressState); if (err !== null) return [s, new Error(`get SimOrder2 Address: ${err.message}`)]; s.address = value; }
    { const [value, err] = rt.getText(buf, newPhoneState); if (err !== null) return [s, new Error(`get SimOrder2 NewPhone: ${err.message}`)]; s.newPhone = value; }
    const errValidate = validateSimOrder2(s);
    if (errValidate !== undefined) return [s, new Error(`validate failed: ${errValidate.message}`)];
    return [s, undefined];
};

export const setSimOrder2 = (buf: rt.Buffer, s: SimOrder2): rt.Err => {
    if (s === null || s === undefined) return new Error(`set SimOrder2: value is null or undefined`);
    const errValidate = validateSimOrder2(s);
    if (errValidate !== undefined) return new Error(`validate SimOrder2: ${errValidate.message}`);
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
    const headerStates = [
        s.id !== 0 ? 1 : 0,
        nameState,
        phoneState,
        idNoState,
        s.cityCode !== 0 ? 1 : 0,
        addressState,
        newPhoneState,
    ];
    const [header, errHeader] = rt.writeHeader(simOrder2HeaderWidths, headerStates);
    if (errHeader !== undefined) { buf.rewindWrite(startOffset); return new Error(`set header: ${errHeader.message}`); }
    const errHeaderWrite = buf.write(header);
    if (errHeaderWrite !== null) { buf.rewindWrite(startOffset); return errHeaderWrite; }
    if (s.id !== 0) {
        const err = rt.setU32(buf, s.id as any);
        if (err !== null) { buf.rewindWrite(startOffset); return new Error(`set SimOrder2 Id: ${err.message}`); }
    }
    { const err = rt.setText(buf, nameState, s.name); if (err !== null) { buf.rewindWrite(startOffset); return new Error(`set SimOrder2 Name: ${err.message}`); } }
    { const err = rt.setText(buf, phoneState, s.phone); if (err !== null) { buf.rewindWrite(startOffset); return new Error(`set SimOrder2 Phone: ${err.message}`); } }
    { const err = rt.setText(buf, idNoState, s.idNo); if (err !== null) { buf.rewindWrite(startOffset); return new Error(`set SimOrder2 IdNo: ${err.message}`); } }
    if (s.cityCode !== 0) {
        const err = rt.setU32(buf, s.cityCode as any);
        if (err !== null) { buf.rewindWrite(startOffset); return new Error(`set SimOrder2 CityCode: ${err.message}`); }
    }
    { const err = rt.setText(buf, addressState, s.address); if (err !== null) { buf.rewindWrite(startOffset); return new Error(`set SimOrder2 Address: ${err.message}`); } }
    { const err = rt.setText(buf, newPhoneState, s.newPhone); if (err !== null) { buf.rewindWrite(startOffset); return new Error(`set SimOrder2 NewPhone: ${err.message}`); } }
    return undefined;
};

export const readSimOrder2 = (buf: rt.Buffer): [SimOrder2, rt.Err] => getSimOrder2(buf);

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

export const getSimOrder2ListBody = (buf: rt.Buffer, state: number): [SimOrder2[], rt.Err] => {
    const [list, err] = rt.getDefaultList<SimOrder2>(
        buf,
        state,
        () => newSimOrder2(),
        (buf) => readSimOrder2(buf),
    );
    return [list, err];
};

export const setSimOrder2ListBody = (buf: rt.Buffer, state: number, v: SimOrder2[]): rt.Err => {
    return rt.setDefaultList<SimOrder2>(
        buf,
        state,
        v,
        (item) => isZeroSimOrder2(item),
        (buf, item) => setSimOrder2(buf, item),
    );
}
