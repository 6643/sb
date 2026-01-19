export class Buffer {
    private _bytes: Uint8Array;
    private _view: DataView;
    private _read_offset: number;
    private _write_offset: number;

    constructor(bytes?: Uint8Array) {
        if (bytes) {
            this._bytes = bytes;
            this._view = new DataView(bytes.buffer, bytes.byteOffset, bytes.byteLength);
            this._read_offset = 0;
            this._write_offset = bytes.length;
        } else {
            const initialCapacity = 128;
            this._bytes = new Uint8Array(initialCapacity);
            this._view = new DataView(this._bytes.buffer);
            this._read_offset = 0;
            this._write_offset = 0;
        }
    }

    public ensureCapacity = (needed: number): void => {
        const required = this._write_offset + needed;
        if (required > this._bytes.length) {
            let newCapacity = this._bytes.length;
            if (newCapacity === 0) newCapacity = 128;
            while (newCapacity < required) {
                newCapacity *= 2;
            }
            const newBytes = new Uint8Array(newCapacity);
            newBytes.set(this._bytes);
            this._bytes = newBytes;
            this._view = new DataView(newBytes.buffer);
        }
    };

    public read = (byteLength: number): [Uint8Array, Error | null] => {
        if (this._read_offset + byteLength > this._write_offset) {
            return [new Uint8Array(0), new Error("not enough data")];
        }
        const slice = this._bytes.subarray(this._read_offset, this._read_offset + byteLength);
        this._read_offset += byteLength;
        return [slice, null];
    }

    public write = (data: Uint8Array): Error | null => {
        this.ensureCapacity(data.length);
        this._bytes.set(data, this._write_offset);
        this._write_offset += data.length;
        return null;
    }

    get read_offset(): number { return this._read_offset; }
    get write_offset(): number { return this._write_offset; }
    get view(): DataView { return this._view; }
    get bytes(): Uint8Array { return this._bytes.subarray(0, this._write_offset); }
    get len(): number { return this._write_offset - this._read_offset; }
}

export interface Serializable {
    set: (buf: Buffer) => Error | null;
}

export interface Deserializable {
    get: (buf: Buffer) => Error | null;
}

// Bit Operations
export const GetBit = (bits: Uint8Array, i: number): boolean => {
    const byteIndex = Math.floor(i / 8);
    if (byteIndex >= bits.length) return false;
    const bitIndex = i % 8;
    return (bits[byteIndex] & (1 << bitIndex)) !== 0;
};

export const SetBit = (bits: Uint8Array, i: number, value: boolean): void => {
    const byteIndex = Math.floor(i / 8);
    if (byteIndex >= bits.length) return;
    const bitIndex = i % 8;
    if (value) {
        bits[byteIndex] |= (1 << bitIndex);
    } else {
        bits[byteIndex] &= ~(1 << bitIndex);
    }
};

const _setNum = (buf: Buffer, byteLength: number, value: number | bigint, setter: string): void => {
    buf.ensureCapacity(byteLength);
    (buf.view as any)[setter](buf.write_offset, value, true);
    (buf as any)._write_offset += byteLength;
};

const _checkRead = (buf: Buffer, len: number): Error | null => {
    if (buf.len < len) return new Error("not enough data");
    return null;
};

// List Helpers
export const getList = <T>(buf: Buffer, getter: (buf: Buffer) => [T, Error | null]): [T[], Error | null] => {
    const [count, err] = getU8(buf);
    if (err !== null) return [[], err];
    const list: T[] = new Array(count);
    for (let i = 0; i < count; i++) {
        const [item, err2] = getter(buf);
        if (err2 !== null) return [[], err2];
        list[i] = item;
    }
    return [list, null];
};

export const setList = <T>(buf: Buffer, list: T[], setter: (buf: Buffer, val: T) => Error | null): Error | null => {
    if (list.length > 255) return new Error(`list length ${list.length} exceeds u8 max`);
    const err = setU8(buf, list.length);
    if (err !== null) return err;
    for (const item of list) {
        const err2 = setter(buf, item);
        if (err2 !== null) return err2;
    }
    return null;
};

export const eqList = <T>(a: T[], b: T[], eq: (a: T, b: T) => boolean): boolean => {
    if (a === b) return true;
    if (a === null || b === null) return false;
    if (a.length !== b.length) return false;
    for (let i = 0; i < a.length; i++) {
        if (!eq(a[i], b[i])) return false;
    }
    return true;
};

// Primitives
export const getU8 = (buf: Buffer): [number, Error | null] => {
    const err = _checkRead(buf, 1);
    if (err !== null) return [0, err];
    const v = buf.view.getUint8(buf.read_offset);
    (buf as any)._read_offset += 1;
    return [v, null];
};
export const setU8 = (buf: Buffer, v: number): Error | null => {
    _setNum(buf, 1, v, 'setUint8');
    return null;
};
export const eqU8 = (a: number, b: number): boolean => a === b;

export const getU16 = (buf: Buffer): [number, Error | null] => {
    const err = _checkRead(buf, 2);
    if (err !== null) return [0, err];
    const v = buf.view.getUint16(buf.read_offset, true);
    (buf as any)._read_offset += 2;
    return [v, null];
};
export const setU16 = (buf: Buffer, v: number): Error | null => {
    _setNum(buf, 2, v, 'setUint16');
    return null;
};
export const eqU16 = (a: number, b: number): boolean => a === b;

export const getU32 = (buf: Buffer): [number, Error | null] => {
    const err = _checkRead(buf, 4);
    if (err !== null) return [0, err];
    const v = buf.view.getUint32(buf.read_offset, true);
    (buf as any)._read_offset += 4;
    return [v, null];
};
export const setU32 = (buf: Buffer, v: number): Error | null => {
    _setNum(buf, 4, v, 'setUint32');
    return null;
};
export const eqU32 = (a: number, b: number): boolean => a === b;

export const getI8 = (buf: Buffer): [number, Error | null] => {
    const err = _checkRead(buf, 1);
    if (err !== null) return [0, err];
    const v = buf.view.getInt8(buf.read_offset);
    (buf as any)._read_offset += 1;
    return [v, null];
};
export const setI8 = (buf: Buffer, v: number): Error | null => {
    _setNum(buf, 1, v, 'setInt8');
    return null;
};
export const eqI8 = (a: number, b: number): boolean => a === b;

export const getI16 = (buf: Buffer): [number, Error | null] => {
    const err = _checkRead(buf, 2);
    if (err !== null) return [0, err];
    const v = buf.view.getInt16(buf.read_offset, true);
    (buf as any)._read_offset += 2;
    return [v, null];
};
export const setI16 = (buf: Buffer, v: number): Error | null => {
    _setNum(buf, 2, v, 'setInt16');
    return null;
};
export const eqI16 = (a: number, b: number): boolean => a === b;

export const getI32 = (buf: Buffer): [number, Error | null] => {
    const err = _checkRead(buf, 4);
    if (err !== null) return [0, err];
    const v = buf.view.getInt32(buf.read_offset, true);
    (buf as any)._read_offset += 4;
    return [v, null];
};
export const setI32 = (buf: Buffer, v: number): Error | null => {
    _setNum(buf, 4, v, 'setInt32');
    return null;
};
export const eqI32 = (a: number, b: number): boolean => a === b;

export const getI64 = (buf: Buffer): [bigint, Error | null] => {
    const err = _checkRead(buf, 8);
    if (err !== null) return [0n, err];
    const v = buf.view.getBigInt64(buf.read_offset, true);
    (buf as any)._read_offset += 8;
    return [v, null];
};
export const setI64 = (buf: Buffer, v: bigint): Error | null => {
    _setNum(buf, 8, v, 'setBigInt64');
    return null;
};
export const eqI64 = (a: bigint, b: bigint): boolean => a === b;

export const getU64 = (buf: Buffer): [bigint, Error | null] => {
    const err = _checkRead(buf, 8);
    if (err !== null) return [0n, err];
    const v = buf.view.getBigUint64(buf.read_offset, true);
    (buf as any)._read_offset += 8;
    return [v, null];
};
export const setU64 = (buf: Buffer, v: bigint): Error | null => {
    _setNum(buf, 8, v, 'setBigUint64');
    return null;
};
export const eqU64 = (a: bigint, b: bigint): boolean => a === b;

export const getF32 = (buf: Buffer): [number, Error | null] => {
    const err = _checkRead(buf, 4);
    if (err !== null) return [0, err];
    const v = buf.view.getFloat32(buf.read_offset, true);
    (buf as any)._read_offset += 4;
    return [v, null];
};
export const setF32 = (buf: Buffer, v: number): Error | null => {
    _setNum(buf, 4, v, 'setFloat32');
    return null;
};
export const eqF32 = (a: number, b: number): boolean => Math.abs(a - b) < 1e-6;

export const getF64 = (buf: Buffer): [number, Error | null] => {
    const err = _checkRead(buf, 8);
    if (err !== null) return [0, err];
    const v = buf.view.getFloat64(buf.read_offset, true);
    (buf as any)._read_offset += 8;
    return [v, null];
};
export const setF64 = (buf: Buffer, v: number): Error | null => {
    _setNum(buf, 8, v, 'setFloat64');
    return null;
};
export const eqF64 = (a: number, b: number): boolean => Math.abs(a - b) < 1e-9;

export const getBool = (buf: Buffer): [boolean, Error | null] => {
    const [v, err] = getU8(buf);
    if (err !== null) return [false, err];
    return [v === 1, null];
};
export const setBool = (buf: Buffer, value: boolean): Error | null => setU8(buf, value ? 1 : 0);
export const eqBool = (a: boolean, b: boolean): boolean => a === b;

export const getBin = (buf: Buffer): [Uint8Array, Error | null] => {
    const [len, err] = getU16(buf);
    if (err !== null) return [new Uint8Array(0), err];
    return buf.read(len);
};
export const setBin = (buf: Buffer, value: Uint8Array): Error | null => {
    const err = setU16(buf, value.length);
    if (err !== null) return err;
    return buf.write(value);
};
export const eqBin = (a: Uint8Array, b: Uint8Array): boolean => {
    if (a === b) return true;
    if (a === null || b === null) return false;
    if (a.length !== b.length) return false;
    for (let i = 0; i < a.length; i++) if (a[i] !== b[i]) return false;
    return true;
};

export const getText = (buf: Buffer): [string, Error | null] => {
    const [data, err] = getBin(buf);
    if (err !== null) return ["", err];
    return [new TextDecoder().decode(data), null];
};
export const setText = (buf: Buffer, value: string): Error | null => {
    const data = new TextEncoder().encode(value);
    return setBin(buf, data);
};
export const eqText = (a: string, b: string): boolean => a === b;

export const getBoolList = (buf: Buffer): [boolean[], Error | null] => getList(buf, getBool);
export const setBoolList = (buf: Buffer, v: boolean[]): Error | null => setList(buf, v, setBool);
export const eqBoolList = (a: boolean[], b: boolean[]): boolean => eqList(a, b, eqBool);

export const getI8List = (buf: Buffer): [number[], Error | null] => getList(buf, getI8);
export const setI8List = (buf: Buffer, v: number[]): Error | null => setList(buf, v, setI8);
export const eqI8List = (a: number[], b: number[]): boolean => eqList(a, b, eqI8);

export const getU8List = (buf: Buffer): [number[], Error | null] => getList(buf, getU8);
export const setU8List = (buf: Buffer, v: number[]): Error | null => setList(buf, v, setU8);
export const eqU8List = (a: number[], b: number[]): boolean => eqList(a, b, eqU8);

export const getI16List = (buf: Buffer): [number[], Error | null] => getList(buf, getI16);
export const setI16List = (buf: Buffer, v: number[]): Error | null => setList(buf, v, setI16);
export const eqI16List = (a: number[], b: number[]): boolean => eqList(a, b, eqI16);

export const getU16List = (buf: Buffer): [number[], Error | null] => getList(buf, getU16);
export const setU16List = (buf: Buffer, v: number[]): Error | null => setList(buf, v, setU16);
export const eqU16List = (a: number[], b: number[]): boolean => eqList(a, b, eqU16);

export const getI32List = (buf: Buffer): [number[], Error | null] => getList(buf, getI32);
export const setI32List = (buf: Buffer, v: number[]): Error | null => setList(buf, v, setI32);
export const eqI32List = (a: number[], b: number[]): boolean => eqList(a, b, eqI32);

export const getU32List = (buf: Buffer): [number[], Error | null] => getList(buf, getU32);
export const setU32List = (buf: Buffer, v: number[]): Error | null => setList(buf, v, setU32);
export const eqU32List = (a: number[], b: number[]): boolean => eqList(a, b, eqU32);

export const getI64List = (buf: Buffer): [bigint[], Error | null] => getList(buf, getI64);
export const setI64List = (buf: Buffer, v: bigint[]): Error | null => setList(buf, v, setI64);
export const eqI64List = (a: bigint[], b: bigint[]): boolean => eqList(a, b, eqI64);

export const getU64List = (buf: Buffer): [bigint[], Error | null] => getList(buf, getU64);
export const setU64List = (buf: Buffer, v: bigint[]): Error | null => setList(buf, v, setU64);
export const eqU64List = (a: bigint[], b: bigint[]): boolean => eqList(a, b, eqU64);

export const getF32List = (buf: Buffer): [number[], Error | null] => getList(buf, getF32);
export const setF32List = (buf: Buffer, v: number[]): Error | null => setList(buf, v, setF32);
export const eqF32List = (a: number[], b: number[]): boolean => eqList(a, b, eqF32);

export const getF64List = (buf: Buffer): [number[], Error | null] => getList(buf, getF64);
export const setF64List = (buf: Buffer, v: number[]): Error | null => setList(buf, v, setF64);
export const eqF64List = (a: number[], b: number[]): boolean => eqList(a, b, eqF64);

export const getBinList = (buf: Buffer): [Uint8Array[], Error | null] => getList(buf, getBin);
export const setBinList = (buf: Buffer, v: Uint8Array[]): Error | null => setList(buf, v, setBin);
export const eqBinList = (a: Uint8Array[], b: Uint8Array[]): boolean => eqList(a, b, eqBin);

export const getTextList = (buf: Buffer): [string[], Error | null] => getList(buf, getText);
export const setTextList = (buf: Buffer, v: string[]): Error | null => setList(buf, v, setText);
export const eqTextList = (a: string[], b: string[]): boolean => eqList(a, b, eqText);

// SetAll performs batch set operations. Supports values or functions.
export const setAll = (buf: Buffer, ...args: any[]): Error | null => {
    for (const arg of args) {
        let err: Error | null = null;
        if (typeof arg === 'function') {
            err = arg(buf);
        } else if (arg && typeof arg === 'object' && 'set' in arg && typeof arg.set === 'function') {
            err = arg.set(buf);
        } else if (typeof arg === 'string') {
            err = setText(buf, arg);
        } else if (typeof arg === 'boolean') {
            err = setBool(buf, arg);
        } else if (typeof arg === 'bigint') {
            err = setI64(buf, arg);
        } else if (arg instanceof Uint8Array) {
            err = setBin(buf, arg);
        } else {
            return new Error(`setAll: unsupported type ${typeof arg}. Use wrappers for numbers.`);
        }
        if (err !== null) return err;
    }
    return null;
};

// GetAll performs batch get operations. Supports functions or objects with get().
export const getAll = (buf: Buffer, ...args: (any)[]): Error | null => {
    for (const arg of args) {
        let err: Error | null = null;
        if (typeof arg === 'function') {
            err = arg(buf);
        } else if (arg && typeof arg === 'object' && 'get' in arg && typeof arg.get === 'function') {
            err = arg.get(buf);
        } else {
            return new Error("getAll: argument must be a function or have a get() method");
        }
        if (err !== null) return err;
    }
    return null;
};

// Helpers to match Go's GetAny logic for primitives via closure update
export const into = <T>(target: { value: T }, getter: (buf: Buffer) => [T, Error | null]) => (buf: Buffer) => {
    const [v, e] = getter(buf);
    if (e === null) target.value = v;
    return e;
};

// Convenience wrappers for TS numbers
export const u8 = (v: number) => (buf: Buffer) => setU8(buf, v);
export const u8List = (v: number[]) => (buf: Buffer) => setU8List(buf, v);
export const u16 = (v: number) => (buf: Buffer) => setU16(buf, v);
export const u32 = (v: number) => (buf: Buffer) => setU32(buf, v);
export const i8 = (v: number) => (buf: Buffer) => setI8(buf, v);
export const i16 = (v: number) => (buf: Buffer) => setI16(buf, v);
export const i32 = (v: number) => (buf: Buffer) => setI32(buf, v);
export const f32 = (v: number) => (buf: Buffer) => setF32(buf, v);
export const f64 = (v: number) => (buf: Buffer) => setF64(buf, v);
export const bool = (v: boolean) => (buf: Buffer) => setBool(buf, v);
export const bin = (v: Uint8Array) => (buf: Buffer) => setBin(buf, v);
export const text = (v: string) => (buf: Buffer) => setText(buf, v);
