import * as _ from "./_"
import * as TplEnum from "./enum"

export interface SimInfo extends _.Serializable, _.Deserializable {
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
    s.set = (buf: _.Buffer) => setSimInfo(buf, s);
    s.get = (buf: _.Buffer) => {
        const [res, err] = getSimInfo(buf);
        if (err === null) Object.assign(s, res);
        return err;
    };
    return s;
}

export const eqSimInfo = (a: SimInfo, b: SimInfo): boolean => {
    if (a === b) return true;
    if (a === null || b === null) return false;
    if (!_.eqU32(a.id, b.id)) return false;
    if (!_.eqText(a.title, b.title)) return false;
    if (!_.eqText(a.content, b.content)) return false;
    if (!_.eqBool(a.a, b.a)) return false;
    if (!_.eqBool(a.b, b.b)) return false;
    if (!_.eqBool(a.c, b.c)) return false;
    if (!_.eqBool(a.d, b.d)) return false;
    if (!_.eqBin(a.zip, b.zip)) return false;
    return true;
}

export const getSimInfo = (buf: _.Buffer): [SimInfo, Error | null] => {
    const s = newSimInfo();
    const bitmaskSize = 1;
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
        s.title = v;
    }
    if (_.GetBit(bits, 2)) {
        const [v, err] = _.getText(buf);
        if (err !== null) return [s, err];
        s.content = v;
    }
    s.a = _.GetBit(bits, 3);
    s.b = _.GetBit(bits, 4);
    s.c = _.GetBit(bits, 5);
    s.d = _.GetBit(bits, 6);
    if (_.GetBit(bits, 7)) {
        const [v, err] = _.getBin(buf);
        if (err !== null) return [s, err];
        s.zip = v;
    }
    return [s, null];
}

export const setSimInfo = (buf: _.Buffer, s: SimInfo): Error | null => {
    if (s === null || s === undefined) return new Error(`set SimInfo: value is null or undefined`);
    const startOffset = buf.write_offset;
    const bits = new Uint8Array(1);
    if (!_.eqU32(s.id, 0)) {
        _.SetBit(bits, 0, true);
    }
    if (!_.eqText(s.title, "")) {
        _.SetBit(bits, 1, true);
    }
    if (!_.eqText(s.content, "")) {
        _.SetBit(bits, 2, true);
    }
    _.SetBit(bits, 3, s.a as boolean);
    _.SetBit(bits, 4, s.b as boolean);
    _.SetBit(bits, 5, s.c as boolean);
    _.SetBit(bits, 6, s.d as boolean);
    if (!_.eqBin(s.zip, new Uint8Array(0))) {
        _.SetBit(bits, 7, true);
    }

    const errBits = buf.write(bits);
    if (errBits !== null) return errBits;
    if (!_.eqU32(s.id, 0)) {
        const err = _.setU32(buf, s.id);
        if (err !== null) {
            buf.rewindWrite(startOffset);
            return err;
        }
    }
    if (!_.eqText(s.title, "")) {
        const err = _.setText(buf, s.title);
        if (err !== null) {
            buf.rewindWrite(startOffset);
            return err;
        }
    }
    if (!_.eqText(s.content, "")) {
        const err = _.setText(buf, s.content);
        if (err !== null) {
            buf.rewindWrite(startOffset);
            return err;
        }
    }
    if (!_.eqBin(s.zip, new Uint8Array(0))) {
        const err = _.setBin(buf, s.zip);
        if (err !== null) {
            buf.rewindWrite(startOffset);
            return err;
        }
    }
    return null;
}

export const getSimInfoList = (buf: _.Buffer): [SimInfo[], Error | null] => {
    const [count, err] = _.getU16(buf);
    if (err !== null) return [[], err];
    const list: SimInfo[] = new Array(count);
    for (let i = 0; i < count; i++) {
        const [item, err2] = getSimInfo(buf);
        if (err2 !== null) return [[], err2];
        list[i] = item;
    }
    return [list, null];
}
export const setSimInfoList = (buf: _.Buffer, v: SimInfo[]): Error | null => {
    if (v.length > 65535) return new Error(`list length ${v.length} exceeds u16 max`);
    const err = _.setU16(buf, v.length);
    if (err !== null) return err;
    for (const item of v) {
        const err2 = setSimInfo(buf, item);
        if (err2 !== null) return err2;
    }
    return null;
}
export const eqSimInfoList = (a: SimInfo[], b: SimInfo[]): boolean => _.eqList(a, b, eqSimInfo);
