import { Buffer, StateU16, StateU24, StateU8, StateZero, getU16, getU24, getU8, setU16, setU24, setU8 } from "./runtime_base"

export const binState = (length: number): [number, Error | null] => {
    if (length < 0) return [0, new Error(`negative bin length: ${length}`)];
    if (length === 0) return [StateZero, null];
    if (length <= 0xFF) return [StateU8, null];
    if (length <= 0xFFFF) return [StateU16, null];
    if (length <= 0xFFFFFF) return [StateU24, null];
    return [0, new Error(`bin length exceeds u24 max: ${length}`)];
};

export const getBinLength = (buf: Buffer, state: number, max: number, kind: string): [number, Error | null] => {
    switch (state) {
    case StateZero:
        return [0, null];
    case StateU8: {
        const [v, err] = getU8(buf);
        if (err !== null) return [0, err];
        if (v === 0) return [0, new Error(`${kind} state ${state} encoded zero length`)];
        return [v, null];
    }
    case StateU16: {
        const [v, err] = getU16(buf);
        if (err !== null) return [0, err];
        if (v === 0) return [0, new Error(`${kind} state ${state} encoded zero length`)];
        if (v <= 0xFF) return [0, new Error(`${kind} state ${state} is not canonical for length ${v}`)];
        return [v, null];
    }
    case StateU24: {
        if (max < 0x10000) return [0, new Error(`${kind} state ${state} is illegal`)];
        const [v, err] = getU24(buf);
        if (err !== null) return [0, err];
        if (v === 0) return [0, new Error(`${kind} state ${state} encoded zero length`)];
        if (v <= 0xFFFF) return [0, new Error(`${kind} state ${state} is not canonical for length ${v}`)];
        return [v, null];
    }
    default:
        return [0, new Error(`invalid ${kind} state: ${state}`)];
    }
};

export const getBin = (buf: Buffer, state: number): [Uint8Array, Error | null] => {
    const [length, err] = getBinLength(buf, state, 0xFFFFFF, "bin");
    if (err !== null) return [new Uint8Array(0), err];
    if (length === 0) return [new Uint8Array(0), null];
    const [data, err2] = buf.read(length);
    if (err2 !== null) return [new Uint8Array(0), new Error(`bin body not enough data: ${buf.len} - ${length}`)];
    return [new Uint8Array(data), null];
};

export const setBin = (buf: Buffer, state: number, v: Uint8Array): Error | null => {
    const [canonical, err] = binState(v.byteLength);
    if (err !== null) return err;
    if (state !== canonical) return new Error(`bin state ${state} is not canonical for length ${v.byteLength}`);
    switch (state) {
    case StateZero:
        return null;
    case StateU8: {
        const err2 = setU8(buf, v.byteLength);
        if (err2 !== null) return err2;
        return buf.write(v);
    }
    case StateU16: {
        const err2 = setU16(buf, v.byteLength);
        if (err2 !== null) return err2;
        return buf.write(v);
    }
    case StateU24: {
        const err2 = setU24(buf, v.byteLength);
        if (err2 !== null) return err2;
        return buf.write(v);
    }
    default:
        return new Error(`invalid bin state: ${state}`);
    }
};
