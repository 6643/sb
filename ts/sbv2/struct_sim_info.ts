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
        if (err === null) Object.assign(s, next);
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

export const validateSimInfo = (s: SimInfo | null | undefined): Error | null => {
    if (s === null || s === undefined) return null;
    { const [, err] = rt.textState(s.title); if (err !== null) return new Error(`Title: ${err.message}`); }
    { const [, err] = rt.textState(s.content); if (err !== null) return new Error(`Content: ${err.message}`); }
    { const [, err] = rt.binState(s.zip.byteLength); if (err !== null) return new Error(`Zip: ${err.message}`); }
    return null;
};

export const getSimInfo = (buf: rt.Buffer): [SimInfo, Error | null] => {
    const s = newSimInfo();
    const headerBits = 11;
    const [header, errHeader] = buf.read(rt.headerSize(headerBits));
    if (errHeader !== null) return [s, new Error(`not enough data`)];
    const errPadding = rt.validatePaddingZero(header, headerBits, "SimInfo header");
    if (errPadding !== null) return [s, errPadding];
    const reader = new rt.BitReader(header, headerBits);
    const [idPresent, errIdPresent] = reader.readBit();
    if (errIdPresent !== null) return [s, new Error(`get SimInfo Id header: ${errIdPresent.message}`)];
    const [titleState, errTitleState] = reader.readBits(2);
    if (errTitleState !== null) return [s, new Error(`get SimInfo Title header: ${errTitleState.message}`)];
    const [contentState, errContentState] = reader.readBits(2);
    if (errContentState !== null) return [s, new Error(`get SimInfo Content header: ${errContentState.message}`)];
    const [aState, errAState] = reader.readBit();
    if (errAState !== null) return [s, new Error(`get SimInfo A header: ${errAState.message}`)];
    const [bState, errBState] = reader.readBit();
    if (errBState !== null) return [s, new Error(`get SimInfo B header: ${errBState.message}`)];
    const [cState, errCState] = reader.readBit();
    if (errCState !== null) return [s, new Error(`get SimInfo C header: ${errCState.message}`)];
    const [dState, errDState] = reader.readBit();
    if (errDState !== null) return [s, new Error(`get SimInfo D header: ${errDState.message}`)];
    const [zipState, errZipState] = reader.readBits(2);
    if (errZipState !== null) return [s, new Error(`get SimInfo Zip header: ${errZipState.message}`)];
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
    if (errValidate !== null) return [s, new Error(`validate failed: ${errValidate.message}`)];
    return [s, null];
};

export const setSimInfo = (buf: rt.Buffer, s: SimInfo): Error | null => {
    if (s === null || s === undefined) return new Error(`set SimInfo: value is null or undefined`);
    const errValidate = validateSimInfo(s);
    if (errValidate !== null) return new Error(`validate SimInfo: ${errValidate.message}`);
    const startOffset = buf.writeOffset;
    const [titleState, errTitleState] = rt.textState(s.title);
    if (errTitleState !== null) return errTitleState;
    const [contentState, errContentState] = rt.textState(s.content);
    if (errContentState !== null) return errContentState;
    const [zipState, errZipState] = rt.binState(s.zip.byteLength);
    if (errZipState !== null) return errZipState;
    const header = new rt.BitWriter(11);
    header.writeBit(s.id !== 0);
    { const err = header.writeBits(titleState, 2); if (err !== null) { buf.rewindWrite(startOffset); return new Error(`set SimInfo Title header: ${err.message}`); } }
    { const err = header.writeBits(contentState, 2); if (err !== null) { buf.rewindWrite(startOffset); return new Error(`set SimInfo Content header: ${err.message}`); } }
    header.writeBit(s.a);
    header.writeBit(s.b);
    header.writeBit(s.c);
    header.writeBit(s.d);
    { const err = header.writeBits(zipState, 2); if (err !== null) { buf.rewindWrite(startOffset); return new Error(`set SimInfo Zip header: ${err.message}`); } }
    const errHeaderWrite = buf.write(header.bytes);
    if (errHeaderWrite !== null) { buf.rewindWrite(startOffset); return errHeaderWrite; }
    if (s.id !== 0) {
        const err = rt.setU32(buf, s.id as any);
        if (err !== null) { buf.rewindWrite(startOffset); return new Error(`set SimInfo Id: ${err.message}`); }
    }
    { const err = rt.setTextCompact(buf, titleState, s.title); if (err !== null) { buf.rewindWrite(startOffset); return new Error(`set SimInfo Title: ${err.message}`); } }
    { const err = rt.setTextCompact(buf, contentState, s.content); if (err !== null) { buf.rewindWrite(startOffset); return new Error(`set SimInfo Content: ${err.message}`); } }
    { const err = rt.setBinCompact(buf, zipState, s.zip); if (err !== null) { buf.rewindWrite(startOffset); return new Error(`set SimInfo Zip: ${err.message}`); } }
    return null;
};

export const readSimInfo = (buf: rt.Buffer): [SimInfo, Error | null] => getSimInfo(buf);

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

export const getSimInfoListBody = (buf: rt.Buffer, state: number): [SimInfo[], Error | null] =>
    rt.getBitmapListCompact<SimInfo>(
        buf,
        state,
        () => newSimInfo(),
        (buf) => readSimInfo(buf),
    );

export const setSimInfoListBody = (buf: rt.Buffer, state: number, v: SimInfo[]): Error | null =>
    rt.setBitmapListCompact<SimInfo>(
        buf,
        state,
        v,
        (item) => isZeroSimInfo(item),
        (buf, item) => setSimInfo(buf, item),
    );
