import { BitReader, BitWriter, headerSize } from "./runtime_base"
import type { Err } from "./runtime_base"

export const encodeBitmap = (bits: boolean[]): Uint8Array => {
    const writer = new BitWriter(bits.length);
    for (const bit of bits) writer.writeBit(bit);
    return writer.bytes;
};
export const decodeBitmap = (data: Uint8Array, count: number): [boolean[], Error | null] => {
    if (count < 0) return [[], new Error(`bitmap invalid count: ${count}`)];
    const reader = new BitReader(data, count);
    const bits: boolean[] = new Array(count);
    for (let i = 0; i < count; i++) {
        const [bit, err] = reader.readBit();
        if (err !== null) return [[], err];
        bits[i] = bit;
    }
    const err = validatePaddingZero(data, count, "bitmap");
    if (err !== null) return [[], err];
    return [bits, null];
};

export const validatePaddingZero = (data: Uint8Array, usedBits: number, kind: string): Error | null => {
    const maxBits = data.byteLength * 8;
    if (usedBits < 0 || usedBits > maxBits) return new Error(`${kind} invalid used bits: ${usedBits} not in [0,${maxBits}]`);
    for (let bitOffset = usedBits; bitOffset < data.byteLength * 8; bitOffset++) {
        const byteIndex = Math.floor(bitOffset / 8);
        const bitIndex = 7 - (bitOffset % 8);
        if ((data[byteIndex] & (1 << bitIndex)) !== 0) return new Error(`${kind} padding bit ${bitOffset - usedBits} is not zero`);
    }
    return null;
};

export const encodeStateBlock = (states: number[]): [Uint8Array, Error | null] => {
    const writer = new BitWriter(states.length * 2);
    for (let i = 0; i < states.length; i++) {
        const err = writer.writeBits(states[i], 2);
        if (err !== null) return [new Uint8Array(0), new Error(`state[${i}]: ${err.message}`)];
    }
    return [writer.bytes, null];
};

export const decodeStateBlock = (data: Uint8Array, count: number): [number[], Error | null] => {
    if (count < 0) return [[], new Error(`state block invalid count: ${count}`)];
    const reader = new BitReader(data, count * 2);
    const states = new Array<number>(count);
    for (let i = 0; i < count; i++) {
        const [state, err] = reader.readBits(2);
        if (err !== null) return [[], err];
        states[i] = state;
    }
    const err = validatePaddingZero(data, count * 2, "state block");
    if (err !== null) return [[], err];
    return [states, null];
};

const _headerBitCount = (widths: readonly number[]): [number, Error | null] => {
    let total = 0;
    for (let i = 0; i < widths.length; i++) {
        const width = widths[i];
        if (width !== 1 && width !== 2) return [0, new Error(`header field[${i}] invalid width: ${width}`)];
        total += width;
    }
    return [total, null];
};
export const readHeader = (data: Uint8Array, widths: readonly number[], kind: string): [number[], Err] => {
    const [usedBits, err] = _headerBitCount(widths);
    if (err !== null) return [[], new Error(`${kind}: ${err.message}`)];
    const expectBytes = headerSize(usedBits);
    if (data.byteLength !== expectBytes) return [[], new Error(`${kind} invalid header size: ${data.byteLength} != ${expectBytes}`)];
    const errPadding = validatePaddingZero(data, usedBits, kind);
    if (errPadding !== null) return [[], errPadding];
    const reader = new BitReader(data, usedBits);
    const states = new Array<number>(widths.length);
    for (let i = 0; i < widths.length; i++) {
        const [state, errState] = reader.readBits(widths[i]);
        if (errState !== null) return [[], new Error(`${kind}[${i}]: ${errState.message}`)];
        states[i] = state;
    }
    return [states, undefined];
};
export const writeHeader = (widths: readonly number[], values: readonly number[]): [Uint8Array, Err] => {
    if (widths.length !== values.length) return [new Uint8Array(0), new Error(`header widths/values length mismatch: ${widths.length} !== ${values.length}`)];
    const [usedBits, err] = _headerBitCount(widths);
    if (err !== null) return [new Uint8Array(0), err];
    const writer = new BitWriter(usedBits);
    for (let i = 0; i < widths.length; i++) {
        const errState = writer.writeBits(values[i], widths[i]);
        if (errState !== null) return [new Uint8Array(0), new Error(`header[${i}]: ${errState.message}`)];
    }
    return [writer.bytes, undefined];
};
