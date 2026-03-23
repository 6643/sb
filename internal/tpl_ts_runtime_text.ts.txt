import { Buffer, StateU16, StateU8, StateZero, getU16, getU8, isStringValue, setU16, setU8 } from "./runtime_base"

const _textEncoder = new TextEncoder();
const _textDecoderFatal = new TextDecoder("utf-8", { fatal: true });

export const textBytes = (v: string): Uint8Array => _textEncoder.encode(v);
export const textTypeError = (v: unknown): Error => new Error(`text value must be string: ${typeof v}`);

export const decodeTextStrict = (data: Uint8Array, kind: string): [string, Error | null] => {
    try {
        return [_textDecoderFatal.decode(data), null];
    } catch {
        return ["", new Error(`${kind} invalid utf-8`)];
    }
};

export const validateTextValue = (v: string): [Uint8Array, Error | null] => {
    if (!isStringValue(v)) return [new Uint8Array(0), textTypeError(v)];
    const bytes = textBytes(v);
    const [decoded, err] = decodeTextStrict(bytes, "text value");
    if (err !== null) return [new Uint8Array(0), err];
    if (decoded !== v) return [new Uint8Array(0), new Error("text value invalid utf-8")];
    return [bytes, null];
};

export const textLengthState = (length: number): [number, Error | null] => {
    if (length < 0) return [0, new Error(`negative text length: ${length}`)];
    if (length === 0) return [StateZero, null];
    if (length <= 0xFF) return [StateU8, null];
    if (length <= 0xFFFF) return [StateU16, null];
    return [0, new Error(`text length exceeds u16 max: ${length}`)];
};

export const textState = (v: string): [number, Error | null] => {
    if (!isStringValue(v)) return [0, textTypeError(v)];
    return textLengthState(textBytes(v).byteLength);
};

const _getCompactLength = (buf: Buffer, state: number, max: number, kind: string): [number, Error | null] => {
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
    default:
        void max;
        return [0, new Error(`invalid ${kind} state: ${state}`)];
    }
};

export const getText = (buf: Buffer, state: number): [string, Error | null] => {
    const [length, err] = _getCompactLength(buf, state, 0xFFFF, "text");
    if (err !== null) return ["", err];
    if (length === 0) return ["", null];
    const [data, err2] = buf.read(length);
    if (err2 !== null) return ["", new Error(`text body not enough data: ${buf.len} - ${length}`)];
    return decodeTextStrict(data, "text body");
};

export const setText = (buf: Buffer, state: number, v: string): Error | null => {
    const [bytes, err0] = validateTextValue(v);
    if (err0 !== null) return err0;
    const [canonical, err] = textLengthState(bytes.byteLength);
    if (err !== null) return err;
    if (state !== canonical) return new Error(`text state ${state} is not canonical for length ${bytes.byteLength}`);
    switch (state) {
    case StateZero:
        return null;
    case StateU8: {
        const err2 = setU8(buf, bytes.byteLength);
        if (err2 !== null) return err2;
        return buf.write(bytes);
    }
    case StateU16: {
        const err2 = setU16(buf, bytes.byteLength);
        if (err2 !== null) return err2;
        return buf.write(bytes);
    }
    default:
        return new Error(`invalid text state: ${state}`);
    }
};
