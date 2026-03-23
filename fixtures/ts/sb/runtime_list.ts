import { Buffer, StateU16, StateU24, StateU8, StateZero, bitmapSize, errU, itemHeaderSize, resultU, setU8, getU8 } from "./runtime_base"
import type { Err } from "./runtime_base"
import { decodeBitmap, decodeStateBlock, encodeBitmap, encodeStateBlock } from "./runtime_header"
import { getBin, setBin, binState } from "./runtime_bin"
import { getText, setText, textState } from "./runtime_text"

export const listCountState = (count: number): [number, Error | null] => {
    if (count < 0) return [0, new Error(`negative list count: ${count}`)];
    if (count === 0) return [StateZero, null];
    if (count <= 0xFF) return [StateU8, null];
    return [0, new Error(`list count exceeds u8 max: ${count}`)];
};

export const getListCount = (buf: Buffer, state: number): [number, Error | null] => {
    switch (state) {
    case StateZero:
        return [0, null];
    case StateU8: {
        const [v, err] = getU8(buf);
        if (err !== null) return [0, err];
        if (v === 0) return [0, new Error(`list count state ${state} encoded zero count`)];
        return [v, null];
    }
    case StateU16:
    case StateU24:
        return [0, new Error(`list count state ${state} is illegal`)];
    default:
        return [0, new Error(`invalid list count state: ${state}`)];
    }
};
export const setListCount = (buf: Buffer, state: number, count: number): Error | null => {
    const [canonical, err] = listCountState(count);
    if (err !== null) return err;
    if (state !== canonical) return new Error(`list count state ${state} is not canonical for count ${count}`);
    switch (state) {
    case StateZero:
        return null;
    case StateU8:
        return setU8(buf, count);
    default:
        return new Error(`invalid list count state: ${state}`);
    }
};

export const getBoolList = (buf: Buffer, state: number): [boolean[], Error | null] => {
    const [count, err] = getListCount(buf, state);
    if (err !== null) return [[], err];
    const size = bitmapSize(count);
    const [data, err2] = buf.read(size);
    if (err2 !== null) return [[], new Error(`bool list bitmap not enough data: ${buf.len} - ${size}`)];
    return decodeBitmap(data, count);
};
export const setBoolList = (buf: Buffer, state: number, v: boolean[]): Error | null => {
    const [canonical, err] = listCountState(v.length);
    if (err !== null) return err;
    if (state !== canonical) return new Error(`bool list state ${state} is not canonical for count ${v.length}`);
    const err2 = setListCount(buf, state, v.length);
    if (err2 !== null || v.length === 0) return err2;
    return buf.write(encodeBitmap(v));
};

export const getDefaultList = <T>(
    buf: Buffer,
    state: number,
    defaultItem: () => T,
    getItem: (buf: Buffer) => [T, Err | null],
): [T[], Err] => {
    const [count, err] = getListCount(buf, state);
    if (err !== null) return [[], err];
    const size = bitmapSize(count);
    const [bitmap, err2] = buf.read(size);
    if (err2 !== null) return [[], new Error(`bitmap list not enough data: ${buf.len} - ${size}`)];
    const [bits, err3] = decodeBitmap(bitmap, count);
    if (err3 !== null) return [[], err3];
    const list = new Array<T>(count);
    for (let i = 0; i < count; i++) {
        if (!bits[i]) {
            list[i] = defaultItem();
            continue;
        }
        const [item, err4] = getItem(buf);
        if (err4 !== undefined && err4 !== null) return [[], new Error(`bitmap list[${i}]: ${err4.message}`)];
        list[i] = item;
    }
    return [list, undefined];
};
export const getZeroList = <T>(
    buf: Buffer,
    state: number,
    zeroValue: T,
    getItem: (buf: Buffer) => [T, Error | null],
): [T[], Err] => getDefaultList(buf, state, () => zeroValue, (next) => resultU(...getItem(next)));
export const setDefaultList = <T>(
    buf: Buffer,
    state: number,
    v: T[],
    isDefault: (item: T) => boolean,
    setItem: (buf: Buffer, item: T) => Err | null,
): Err => {
    const [canonical, err] = listCountState(v.length);
    if (err !== null) return err;
    if (state !== canonical) return new Error(`bitmap list state ${state} is not canonical for count ${v.length}`);
    const err2 = setListCount(buf, state, v.length);
    if (err2 !== null) return err2;
    if (v.length === 0) return undefined;
    const bits = v.map(item => !isDefault(item));
    const err3 = buf.write(encodeBitmap(bits));
    if (err3 !== null) return err3;
    for (let i = 0; i < v.length; i++) {
        if (!bits[i]) continue;
        const err4 = setItem(buf, v[i]);
        if (err4 !== undefined && err4 !== null) return new Error(`bitmap list[${i}]: ${err4.message}`);
    }
    return undefined;
};
export const setZeroList = <T>(
    buf: Buffer,
    state: number,
    v: T[],
    zeroValue: T,
    setItem: (buf: Buffer, item: T) => Error | null,
): Err => setDefaultList(buf, state, v, (item) => item === zeroValue, (next, item) => errU(setItem(next, item)));

export const getTextList = (buf: Buffer, state: number): [string[], Error | null] => {
    const [count, err] = getListCount(buf, state);
    if (err !== null) return [[], err];
    const size = itemHeaderSize(count);
    const [header, err2] = buf.read(size);
    if (err2 !== null) return [[], new Error(`text list state block not enough data: ${buf.len} - ${size}`)];
    const [states, err3] = decodeStateBlock(header, count);
    if (err3 !== null) return [[], err3];
    const list = new Array<string>(count);
    for (let i = 0; i < count; i++) {
        const [item, err4] = getText(buf, states[i]);
        if (err4 !== null) return [[], new Error(`text list[${i}]: ${err4.message}`)];
        list[i] = item;
    }
    return [list, null];
};
export const setTextList = (buf: Buffer, state: number, v: string[]): Error | null => {
    const [canonical, err] = listCountState(v.length);
    if (err !== null) return err;
    if (state !== canonical) return new Error(`text list state ${state} is not canonical for count ${v.length}`);
    const err2 = setListCount(buf, state, v.length);
    if (err2 !== null) return err2;
    if (v.length === 0) return null;
    const states = new Array<number>(v.length);
    for (let i = 0; i < v.length; i++) {
        const [itemState, err3] = textState(v[i]);
        if (err3 !== null) return new Error(`text list[${i}]: ${err3.message}`);
        states[i] = itemState;
    }
    const [header, err4] = encodeStateBlock(states);
    if (err4 !== null) return err4;
    const err5 = buf.write(header);
    if (err5 !== null) return err5;
    for (let i = 0; i < v.length; i++) {
        const err6 = setText(buf, states[i], v[i]);
        if (err6 !== null) return new Error(`text list[${i}]: ${err6.message}`);
    }
    return null;
};

export const getBinList = (buf: Buffer, state: number): [Uint8Array[], Error | null] => {
    const [count, err] = getListCount(buf, state);
    if (err !== null) return [[], err];
    const size = itemHeaderSize(count);
    const [header, err2] = buf.read(size);
    if (err2 !== null) return [[], new Error(`bin list state block not enough data: ${buf.len} - ${size}`)];
    const [states, err3] = decodeStateBlock(header, count);
    if (err3 !== null) return [[], err3];
    const list = new Array<Uint8Array>(count);
    for (let i = 0; i < count; i++) {
        const [item, err4] = getBin(buf, states[i]);
        if (err4 !== null) return [[], new Error(`bin list[${i}]: ${err4.message}`)];
        list[i] = item;
    }
    return [list, null];
};
export const setBinList = (buf: Buffer, state: number, v: Uint8Array[]): Error | null => {
    const [canonical, err] = listCountState(v.length);
    if (err !== null) return err;
    if (state !== canonical) return new Error(`bin list state ${state} is not canonical for count ${v.length}`);
    const err2 = setListCount(buf, state, v.length);
    if (err2 !== null) return err2;
    if (v.length === 0) return null;
    const states = new Array<number>(v.length);
    for (let i = 0; i < v.length; i++) {
        const [itemState, err3] = binState(v[i].byteLength);
        if (err3 !== null) return new Error(`bin list[${i}]: ${err3.message}`);
        states[i] = itemState;
    }
    const [header, err4] = encodeStateBlock(states);
    if (err4 !== null) return err4;
    const err5 = buf.write(header);
    if (err5 !== null) return err5;
    for (let i = 0; i < v.length; i++) {
        const err6 = setBin(buf, states[i], v[i]);
        if (err6 !== null) return new Error(`bin list[${i}]: ${err6.message}`);
    }
    return null;
};
