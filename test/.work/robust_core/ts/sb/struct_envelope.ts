import * as _ from "./_"
import * as Enum from "./enum"

export interface Envelope extends _.Serializable, _.Deserializable {
    item: _.Item | null;
    items: _.Item[];
    note: string;
}

export const newEnvelope = (): Envelope => {
    const s = {
        item: null,
        items: [],
        note: "",
    } as any as Envelope;
    s.set = (buf: _.Buffer) => setEnvelope(buf, s);
    s.get = (buf: _.Buffer) => {
        const [res, err] = getEnvelope(buf);
        if (err === null) Object.assign(s, res);
        return err;
    };
    return s;
}

export const eqEnvelope = (a: Envelope, b: Envelope): boolean => {
    if (a === b) return true;
    if (a === null || b === null) return false;
    if (!_.eqItem(a.item, b.item)) return false;
    if (!_.eqItemList(a.items, b.items)) return false;
    if (!_.eqText(a.note, b.note)) return false;
    return true;
}

export const getEnvelope = (buf: _.Buffer): [Envelope, Error | null] => {
    const s = newEnvelope();
    const bitmaskSize = Math.ceil(3 / 8);
    const [bits, err] = buf.read(bitmaskSize);
    if (err !== null) return [s, err];
    if (_.GetBit(bits, 0)) {
        const [v, err] = _.getItem(buf);
        if (err !== null) return [s, err];
        s.item = v;
    }
    if (_.GetBit(bits, 1)) {
        const [v, err] = _.getItemList(buf);
        if (err !== null) return [s, err];
        s.items = v;
    }
    if (_.GetBit(bits, 2)) {
        const [v, err] = _.getText(buf);
        if (err !== null) return [s, err];
        s.note = v;
    }
    return [s, null];
}

export const setEnvelope = (buf: _.Buffer, s: Envelope): Error | null => {
    if (s === null || s === undefined) return new Error(`set Envelope: value is null or undefined`);
    const bits = new Uint8Array(Math.ceil(3 / 8));
    const body = new _.Buffer();
    if (s.item !== null && s.item !== undefined) {
        const err = _.setItem(body, s.item);
        if (err !== null) return err;
        _.SetBit(bits, 0, true);
    }
    if (s.items && s.items.length > 0) {
        const err = _.setItemList(body, s.items);
        if (err !== null) return err;
        _.SetBit(bits, 1, true);
    }
    if (!_.eqText(s.note, "")) {
        const err = _.setText(body, s.note);
        if (err !== null) return err;
        _.SetBit(bits, 2, true);
    }

    const errBits = buf.write(bits);
    if (errBits !== null) return errBits;
    return buf.write(body.bytes);
}

export const getEnvelopeList = (buf: _.Buffer): [Envelope[], Error | null] => _.getList(buf, getEnvelope);
export const setEnvelopeList = (buf: _.Buffer, v: Envelope[]): Error | null => _.setList(buf, v, setEnvelope);
export const eqEnvelopeList = (a: Envelope[], b: Envelope[]): boolean => _.eqList(a, b, eqEnvelope);
