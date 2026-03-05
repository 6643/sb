import * as _ from "./_"
import * as Enum from "./enum"

export interface Matrix extends _.Serializable, _.Deserializable {
    single: _.Child | null;
    many: _.Child[];
    kind: Enum.Kind;
    kinds: Enum.Kind[];
}

export const newMatrix = (): Matrix => {
    const s = {
        single: null,
        many: [],
        kind: 0,
        kinds: [],
    } as any as Matrix;
    s.set = (buf: _.Buffer) => setMatrix(buf, s);
    s.get = (buf: _.Buffer) => {
        const [res, err] = getMatrix(buf);
        if (err === null) Object.assign(s, res);
        return err;
    };
    return s;
}

export const eqMatrix = (a: Matrix, b: Matrix): boolean => {
    if (a === b) return true;
    if (a === null || b === null) return false;
    if (!_.eqChild(a.single, b.single)) return false;
    if (!_.eqChildList(a.many, b.many)) return false;
    if (a.kind !== b.kind) return false;
    if (!_.eqU8List(a.kinds as any, b.kinds as any)) return false;
    return true;
}

export const getMatrix = (buf: _.Buffer): [Matrix, Error | null] => {
    const s = newMatrix();
    const bitmaskSize = Math.ceil(4 / 8);
    const [bits, err] = buf.read(bitmaskSize);
    if (err !== null) return [s, err];
    if (_.GetBit(bits, 0)) {
        const [v, err] = _.getChild(buf);
        if (err !== null) return [s, err];
        s.single = v;
    }
    if (_.GetBit(bits, 1)) {
        const [v, err] = _.getChildList(buf);
        if (err !== null) return [s, err];
        s.many = v;
    }
    if (_.GetBit(bits, 2)) {
        const [v, err] = _.getU8(buf);
        if (err !== null) return [s, err];
        s.kind = v as any;
        if (!_.IsKind(s.kind as any)) return [s, new Error("get Matrix Kind: invalid enum value")];
    }
    if (_.GetBit(bits, 3)) {
        const [v, err] = _.getU8List(buf);
        if (err !== null) return [s, err];
        s.kinds = v as any;
        if (!_.IsKindList(s.kinds as any)) return [s, new Error("get Matrix Kinds: invalid enum value")];
    }
    return [s, null];
}

export const setMatrix = (buf: _.Buffer, s: Matrix): Error | null => {
    if (s === null || s === undefined) return new Error(`set Matrix: value is null or undefined`);
    const bits = new Uint8Array(Math.ceil(4 / 8));
    const body = new _.Buffer();
    if (s.single !== null && s.single !== undefined) {
        const err = _.setChild(body, s.single);
        if (err !== null) return err;
        _.SetBit(bits, 0, true);
    }
    if (s.many && s.many.length > 0) {
        const err = _.setChildList(body, s.many);
        if (err !== null) return err;
        _.SetBit(bits, 1, true);
    }
    if ((s.kind as any) !== 0) {
        if (!_.IsKind(s.kind as any)) return new Error("set Matrix Kind: invalid enum value");
        const err = _.setU8(body, s.kind as any);
        if (err !== null) return err;
        _.SetBit(bits, 2, true);
    }
    if (s.kinds && s.kinds.length > 0) {
        if (!_.IsKindList(s.kinds as any)) return new Error("set Matrix Kinds: invalid enum value");
        const err = _.setU8List(body, s.kinds as any);
        if (err !== null) return err;
        _.SetBit(bits, 3, true);
    }

    const errBits = buf.write(bits);
    if (errBits !== null) return errBits;
    return buf.write(body.bytes);
}

export const getMatrixList = (buf: _.Buffer): [Matrix[], Error | null] => _.getList(buf, getMatrix);
export const setMatrixList = (buf: _.Buffer, v: Matrix[]): Error | null => _.setList(buf, v, setMatrix);
export const eqMatrixList = (a: Matrix[], b: Matrix[]): boolean => _.eqList(a, b, eqMatrix);
