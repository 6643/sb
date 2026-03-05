import * as _ from "./_"
import * as Enum from "./enum"

export interface Child extends _.Serializable, _.Deserializable {
    id: number;
}

export const newChild = (): Child => {
    const s = {
        id: 0,
    } as any as Child;
    s.set = (buf: _.Buffer) => setChild(buf, s);
    s.get = (buf: _.Buffer) => {
        const [res, err] = getChild(buf);
        if (err === null) Object.assign(s, res);
        return err;
    };
    return s;
}

export const eqChild = (a: Child, b: Child): boolean => {
    if (a === b) return true;
    if (a === null || b === null) return false;
    if (!_.eqU32(a.id, b.id)) return false;
    return true;
}

export const getChild = (buf: _.Buffer): [Child, Error | null] => {
    const s = newChild();
    const bitmaskSize = Math.ceil(1 / 8);
    const [bits, err] = buf.read(bitmaskSize);
    if (err !== null) return [s, err];
    if (_.GetBit(bits, 0)) {
        const [v, err] = _.getU32(buf);
        if (err !== null) return [s, err];
        s.id = v;
    }
    return [s, null];
}

export const setChild = (buf: _.Buffer, s: Child): Error | null => {
    if (s === null || s === undefined) return new Error(`set Child: value is null or undefined`);
    const bits = new Uint8Array(Math.ceil(1 / 8));
    const body = new _.Buffer();
    if (!_.eqU32(s.id, 0)) {
        const err = _.setU32(body, s.id);
        if (err !== null) return err;
        _.SetBit(bits, 0, true);
    }

    const errBits = buf.write(bits);
    if (errBits !== null) return errBits;
    return buf.write(body.bytes);
}

export const getChildList = (buf: _.Buffer): [Child[], Error | null] => _.getList(buf, getChild);
export const setChildList = (buf: _.Buffer, v: Child[]): Error | null => _.setList(buf, v, setChild);
export const eqChildList = (a: Child[], b: Child[]): boolean => _.eqList(a, b, eqChild);
