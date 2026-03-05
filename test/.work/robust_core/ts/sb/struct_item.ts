import * as _ from "./_"
import * as Enum from "./enum"

export interface Item extends _.Serializable, _.Deserializable {
    id: number;
    color: Enum.Color;
    tags: string[];
    active: boolean;
}

export const newItem = (): Item => {
    const s = {
        id: 0,
        color: 0,
        tags: [],
        active: false,
    } as any as Item;
    s.set = (buf: _.Buffer) => setItem(buf, s);
    s.get = (buf: _.Buffer) => {
        const [res, err] = getItem(buf);
        if (err === null) Object.assign(s, res);
        return err;
    };
    return s;
}

export const eqItem = (a: Item, b: Item): boolean => {
    if (a === b) return true;
    if (a === null || b === null) return false;
    if (!_.eqU32(a.id, b.id)) return false;
    if (a.color !== b.color) return false;
    if (!_.eqTextList(a.tags, b.tags)) return false;
    if (!_.eqBool(a.active, b.active)) return false;
    return true;
}

export const getItem = (buf: _.Buffer): [Item, Error | null] => {
    const s = newItem();
    const bitmaskSize = Math.ceil(4 / 8);
    const [bits, err] = buf.read(bitmaskSize);
    if (err !== null) return [s, err];
    if (_.GetBit(bits, 0)) {
        const [v, err] = _.getU32(buf);
        if (err !== null) return [s, err];
        s.id = v;
    }
    if (_.GetBit(bits, 1)) {
        const [v, err] = _.getU8(buf);
        if (err !== null) return [s, err];
        s.color = v as any;
        if (!_.IsColor(s.color as any)) return [s, new Error("get Item Color: invalid enum value")];
    }
    if (_.GetBit(bits, 2)) {
        const [v, err] = _.getTextList(buf);
        if (err !== null) return [s, err];
        s.tags = v;
    }
    s.active = _.GetBit(bits, 3);
    return [s, null];
}

export const setItem = (buf: _.Buffer, s: Item): Error | null => {
    if (s === null || s === undefined) return new Error(`set Item: value is null or undefined`);
    const bits = new Uint8Array(Math.ceil(4 / 8));
    const body = new _.Buffer();
    if (!_.eqU32(s.id, 0)) {
        const err = _.setU32(body, s.id);
        if (err !== null) return err;
        _.SetBit(bits, 0, true);
    }
    if ((s.color as any) !== 0) {
        if (!_.IsColor(s.color as any)) return new Error("set Item Color: invalid enum value");
        const err = _.setU8(body, s.color as any);
        if (err !== null) return err;
        _.SetBit(bits, 1, true);
    }
    if (s.tags && s.tags.length > 0) {
        const err = _.setTextList(body, s.tags);
        if (err !== null) return err;
        _.SetBit(bits, 2, true);
    }
    _.SetBit(bits, 3, s.active as boolean);

    const errBits = buf.write(bits);
    if (errBits !== null) return errBits;
    return buf.write(body.bytes);
}

export const getItemList = (buf: _.Buffer): [Item[], Error | null] => _.getList(buf, getItem);
export const setItemList = (buf: _.Buffer, v: Item[]): Error | null => _.setList(buf, v, setItem);
export const eqItemList = (a: Item[], b: Item[]): boolean => _.eqList(a, b, eqItem);
