import * as rt from "./type"
import * as _ from "./_"

export interface SimInfo extends rt.Serializable, rt.Deserializable {
    id: number;
    title: string;
    content: string;
    a: boolean;
    b: boolean;
    c: boolean;
    d: boolean;
    zip: Uint8Array;
}

const simInfoHeaderWidths = [1, 2, 2, 1, 1, 1, 1, 2] as const;

export const newSimInfo = (): SimInfo => {
    const s = {
        id: 0,
        title: "",
        content: "",
        a: false,
        b: false,
        c: false,
        d: false,
        zip: new Uint8Array(0),
    } as any as SimInfo;
    s.set = (buf: rt.Buffer) => setSimInfo(buf, s);
    s.get = (buf: rt.Buffer) => {
        const [next, err] = getSimInfo(buf);
        if (err === undefined) Object.assign(s, next);
        return err;
    };
    return s;
};

export const isZeroSimInfo = (s: SimInfo | null | undefined): boolean => {
    if (s === null || s === undefined) return true;
    if (s.id !== 0) return false;
    if (s.title !== "") return false;
    if (s.content !== "") return false;
    if (s.a) return false;
    if (s.b) return false;
    if (s.c) return false;
    if (s.d) return false;
    if (s.zip.byteLength !== 0) return false;
    return true;
};

export const validateSimInfo = (s: SimInfo | null | undefined): rt.Err => {
    if (s === null || s === undefined) return undefined;
    { const [, err] = rt.textState(s.title); if (err !== null) return new Error(`Title: ${err.message}`); }
    { const [, err] = rt.textState(s.content); if (err !== null) return new Error(`Content: ${err.message}`); }
    { const [, err] = rt.binState(s.zip.byteLength); if (err !== null) return new Error(`Zip: ${err.message}`); }
    return undefined;
};

export const getSimInfo = (buf: rt.Buffer): [SimInfo, rt.Err] => {
    const s = newSimInfo();
    const headerBits = 11;
    const [header, errHeader] = buf.read(rt.headerSize(headerBits));
    if (errHeader !== null) return [s, new Error(`not enough data`)];
    const [headerStates, errHeaderState] = rt.readHeader(header, simInfoHeaderWidths, "SimInfo header");
    if (errHeaderState !== undefined) return [s, errHeaderState];
    const idPresent = headerStates[0] === 1;
    const titleState = headerStates[1];
    const contentState = headerStates[2];
    const aState = headerStates[3] === 1;
    const bState = headerStates[4] === 1;
    const cState = headerStates[5] === 1;
    const dState = headerStates[6] === 1;
    const zipState = headerStates[7];
    if (idPresent) {
        const [value, err] = rt.getU32(buf);
        if (err !== null) return [s, new Error(`get SimInfo Id: ${err.message}`)];
        s.id = value as any;
    }
    { const [value, err] = rt.getTextCompact(buf, titleState); if (err !== null) return [s, new Error(`get SimInfo Title: ${err.message}`)]; s.title = value; }
    { const [value, err] = rt.getTextCompact(buf, contentState); if (err !== null) return [s, new Error(`get SimInfo Content: ${err.message}`)]; s.content = value; }
    s.a = aState;
    s.b = bState;
    s.c = cState;
    s.d = dState;
    { const [value, err] = rt.getBinCompact(buf, zipState); if (err !== null) return [s, new Error(`get SimInfo Zip: ${err.message}`)]; s.zip = value; }
    const errValidate = validateSimInfo(s);
    if (errValidate !== undefined) return [s, new Error(`validate failed: ${errValidate.message}`)];
    return [s, undefined];
};

export const setSimInfo = (buf: rt.Buffer, s: SimInfo): rt.Err => {
    if (s === null || s === undefined) return new Error(`set SimInfo: value is null or undefined`);
    const errValidate = validateSimInfo(s);
    if (errValidate !== undefined) return new Error(`validate SimInfo: ${errValidate.message}`);
    const startOffset = buf.writeOffset;
    const [titleState, errTitleState] = rt.textState(s.title);
    if (errTitleState !== null) return errTitleState;
    const [contentState, errContentState] = rt.textState(s.content);
    if (errContentState !== null) return errContentState;
    const [zipState, errZipState] = rt.binState(s.zip.byteLength);
    if (errZipState !== null) return errZipState;
    const headerStates = [
        s.id !== 0 ? 1 : 0,
        titleState,
        contentState,
        s.a ? 1 : 0,
        s.b ? 1 : 0,
        s.c ? 1 : 0,
        s.d ? 1 : 0,
        zipState,
    ];
    const [header, errHeader] = rt.writeHeader(simInfoHeaderWidths, headerStates);
    if (errHeader !== undefined) { buf.rewindWrite(startOffset); return new Error(`set header: ${errHeader.message}`); }
    const errHeaderWrite = buf.write(header);
    if (errHeaderWrite !== null) { buf.rewindWrite(startOffset); return errHeaderWrite; }
    if (s.id !== 0) {
        const err = rt.setU32(buf, s.id as any);
        if (err !== null) { buf.rewindWrite(startOffset); return new Error(`set SimInfo Id: ${err.message}`); }
    }
    { const err = rt.setTextCompact(buf, titleState, s.title); if (err !== null) { buf.rewindWrite(startOffset); return new Error(`set SimInfo Title: ${err.message}`); } }
    { const err = rt.setTextCompact(buf, contentState, s.content); if (err !== null) { buf.rewindWrite(startOffset); return new Error(`set SimInfo Content: ${err.message}`); } }
    { const err = rt.setBinCompact(buf, zipState, s.zip); if (err !== null) { buf.rewindWrite(startOffset); return new Error(`set SimInfo Zip: ${err.message}`); } }
    return undefined;
};

export const readSimInfo = (buf: rt.Buffer): [SimInfo, rt.Err] => getSimInfo(buf);

export const eqSimInfo = (a: SimInfo | null | undefined, b: SimInfo | null | undefined): boolean => {
    if (isZeroSimInfo(a as any) && isZeroSimInfo(b as any)) return true;
    if (a === null || a === undefined || b === null || b === undefined) return false;
    if (!rt.eqU32(a.id, b.id)) return false;
    if (!rt.eqText(a.title, b.title)) return false;
    if (!rt.eqText(a.content, b.content)) return false;
    if (!rt.eqBool(a.a, b.a)) return false;
    if (!rt.eqBool(a.b, b.b)) return false;
    if (!rt.eqBool(a.c, b.c)) return false;
    if (!rt.eqBool(a.d, b.d)) return false;
    if (!rt.eqBin(a.zip, b.zip)) return false;
    return true;
};

export const getSimInfoListBody = (buf: rt.Buffer, state: number): [SimInfo[], rt.Err] => {
    const [list, err] = rt.getBitmapListCompact<SimInfo>(
        buf,
        state,
        () => newSimInfo(),
        (buf) => readSimInfo(buf),
    );
    return [list, err];
};

export const setSimInfoListBody = (buf: rt.Buffer, state: number, v: SimInfo[]): rt.Err => {
    return rt.setBitmapListCompact<SimInfo>(
        buf,
        state,
        v,
        (item) => isZeroSimInfo(item),
        (buf, item) => setSimInfo(buf, item),
    );
}
