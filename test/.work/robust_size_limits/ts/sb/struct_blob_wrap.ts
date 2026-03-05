import * as _ from "./_"
import * as Enum from "./enum"

export interface BlobWrap extends _.Serializable, _.Deserializable {
    textValue: string;
    binValue: Uint8Array;
    nums: number[];
    level: Enum.Level;
}

export const newBlobWrap = (): BlobWrap => {
    const s = {
        textValue: "",
        binValue: new Uint8Array(0),
        nums: [],
        level: 0,
    } as any as BlobWrap;
    s.set = (buf: _.Buffer) => setBlobWrap(buf, s);
    s.get = (buf: _.Buffer) => {
        const [res, err] = getBlobWrap(buf);
        if (err === null) Object.assign(s, res);
        return err;
    };
    return s;
}

export const eqBlobWrap = (a: BlobWrap, b: BlobWrap): boolean => {
    if (a === b) return true;
    if (a === null || b === null) return false;
    if (!_.eqText(a.textValue, b.textValue)) return false;
    if (!_.eqBin(a.binValue, b.binValue)) return false;
    if (!_.eqU8List(a.nums, b.nums)) return false;
    if (a.level !== b.level) return false;
    return true;
}

export const getBlobWrap = (buf: _.Buffer): [BlobWrap, Error | null] => {
    const s = newBlobWrap();
    const bitmaskSize = Math.ceil(4 / 8);
    const [bits, err] = buf.read(bitmaskSize);
    if (err !== null) return [s, err];
    if (_.GetBit(bits, 0)) {
        const [v, err] = _.getText(buf);
        if (err !== null) return [s, err];
        s.textValue = v;
    }
    if (_.GetBit(bits, 1)) {
        const [v, err] = _.getBin(buf);
        if (err !== null) return [s, err];
        s.binValue = v;
    }
    if (_.GetBit(bits, 2)) {
        const [v, err] = _.getU8List(buf);
        if (err !== null) return [s, err];
        s.nums = v;
    }
    if (_.GetBit(bits, 3)) {
        const [v, err] = _.getU8(buf);
        if (err !== null) return [s, err];
        s.level = v as any;
        if (!_.IsLevel(s.level as any)) return [s, new Error("get BlobWrap Level: invalid enum value")];
    }
    return [s, null];
}

export const setBlobWrap = (buf: _.Buffer, s: BlobWrap): Error | null => {
    if (s === null || s === undefined) return new Error(`set BlobWrap: value is null or undefined`);
    const bits = new Uint8Array(Math.ceil(4 / 8));
    const body = new _.Buffer();
    if (!_.eqText(s.textValue, "")) {
        const err = _.setText(body, s.textValue);
        if (err !== null) return err;
        _.SetBit(bits, 0, true);
    }
    if (!_.eqBin(s.binValue, new Uint8Array(0))) {
        const err = _.setBin(body, s.binValue);
        if (err !== null) return err;
        _.SetBit(bits, 1, true);
    }
    if (s.nums && s.nums.length > 0) {
        const err = _.setU8List(body, s.nums);
        if (err !== null) return err;
        _.SetBit(bits, 2, true);
    }
    if ((s.level as any) !== 0) {
        if (!_.IsLevel(s.level as any)) return new Error("set BlobWrap Level: invalid enum value");
        const err = _.setU8(body, s.level as any);
        if (err !== null) return err;
        _.SetBit(bits, 3, true);
    }

    const errBits = buf.write(bits);
    if (errBits !== null) return errBits;
    return buf.write(body.bytes);
}

export const getBlobWrapList = (buf: _.Buffer): [BlobWrap[], Error | null] => _.getList(buf, getBlobWrap);
export const setBlobWrapList = (buf: _.Buffer, v: BlobWrap[]): Error | null => _.setList(buf, v, setBlobWrap);
export const eqBlobWrapList = (a: BlobWrap[], b: BlobWrap[]): boolean => _.eqList(a, b, eqBlobWrap);
