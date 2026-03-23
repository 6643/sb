export interface Buffer {
    ensureCapacity: (needed: number) => void;
    read: (byteLength: number) => [Uint8Array, Error | null];
    write: (data: Uint8Array) => Error | null;
    writeUnsafe: (data: Uint8Array) => void;
    rewindWrite: (offset: number) => void;
    readonly bytes: Uint8Array;
    readonly view: DataView;
    readonly len: number;
    readonly readOffset: number;
    readonly writeOffset: number;
}

export type Err = Error | undefined;

export interface Serializable {
    set: (buf: Buffer) => Err;
}

export interface Deserializable {
    get: (buf: Buffer) => Err;
}

export const attachBufferOps = (proto: any): void => {
    proto.ensureCapacity = function(this: any, needed: number): void {
        const required = this._writeOffset + needed;
        if (required <= this._bytes.length) return;
        let newCapacity = this._bytes.length === 0 ? 128 : this._bytes.length;
        while (newCapacity < required) newCapacity *= 2;
        const next = new Uint8Array(newCapacity);
        next.set(this._bytes.subarray(0, this._writeOffset));
        this._bytes = next;
        this._view = new DataView(next.buffer);
    };

    proto.read = function(this: any, byteLength: number): [Uint8Array, Error | null] {
        if (byteLength < 0) return [new Uint8Array(0), new Error(`invalid read length ${byteLength}`)];
        if (this._readOffset + byteLength > this._writeOffset) return [new Uint8Array(0), new Error("not enough data")];
        const slice = this._bytes.subarray(this._readOffset, this._readOffset + byteLength);
        this._readOffset += byteLength;
        return [slice, null];
    };

    proto.write = function(this: any, data: Uint8Array): Error | null {
        this.ensureCapacity(data.byteLength);
        this.writeUnsafe(data);
        return null;
    };

    proto.writeUnsafe = function(this: any, data: Uint8Array): void {
        this._bytes.set(data, this._writeOffset);
        this._writeOffset += data.byteLength;
    };

    proto.rewindWrite = function(this: any, offset: number): void {
        if (offset < 0) offset = 0;
        if (offset > this._writeOffset) offset = this._writeOffset;
        this._writeOffset = offset;
        if (this._readOffset > this._writeOffset) this._readOffset = this._writeOffset;
    };
};

export const errU = (err: Error | null): Err => err === null ? undefined : err;
export const resultU = <T>(value: T, err: Error | null): [T, Err] => err === null ? [value, undefined] : [value, err];

export const StateZero = 0;
export const StateU8 = 1;
export const StateU16 = 2;
export const StateU24 = 3;

export const headerSize = (bitCount: number): number => bitCount <= 0 ? 0 : Math.floor((bitCount + 7) / 8);
export const bitmapSize = (count: number): number => count <= 0 ? 0 : Math.floor((count + 7) / 8);
export const itemHeaderSize = (count: number): number => count <= 0 ? 0 : Math.floor((count * 2 + 7) / 8);

type _BufferState = Buffer & {
    _bytes: Uint8Array;
    _view: DataView;
    _readOffset: number;
    _writeOffset: number;
};
type BufferCtor = new (bytes?: Uint8Array) => Buffer;

function _BufferCtor(this: _BufferState, bytes?: Uint8Array): void {
    if (bytes) {
        this._bytes = bytes;
        this._view = new DataView(bytes.buffer, bytes.byteOffset, bytes.byteLength);
        this._readOffset = 0;
        this._writeOffset = bytes.byteLength;
        return;
    }
    this._bytes = new Uint8Array(128);
    this._view = new DataView(this._bytes.buffer);
    this._readOffset = 0;
    this._writeOffset = 0;
}
const _bufferProto = _BufferCtor.prototype as _BufferState;

attachBufferOps(_bufferProto);

Object.defineProperty(_bufferProto, "bytes", { get(this: _BufferState): Uint8Array { return this._bytes.subarray(0, this._writeOffset); } });
Object.defineProperty(_bufferProto, "view", { get(this: _BufferState): DataView { return this._view; } });
Object.defineProperty(_bufferProto, "len", { get(this: _BufferState): number { return this._writeOffset - this._readOffset; } });
Object.defineProperty(_bufferProto, "readOffset", { get(this: _BufferState): number { return this._readOffset; } });
Object.defineProperty(_bufferProto, "writeOffset", { get(this: _BufferState): number { return this._writeOffset; } });

export const Buffer = _BufferCtor as unknown as BufferCtor;

export interface BitReader {
    readBit: () => [boolean, Error | null];
    readBits: (width: number) => [number, Error | null];
    readonly bitOffset: number;
}

type _BitReaderState = BitReader & {
    _data: Uint8Array;
    _bitLimit: number;
    _bitOffset: number;
};
type BitReaderCtor = new (data: Uint8Array, bitLimit: number) => BitReader;

function _BitReaderCtor(this: _BitReaderState, data: Uint8Array, bitLimit: number): void {
    const maxBits = data.byteLength * 8;
    this._data = data;
    this._bitLimit = bitLimit < 0 || bitLimit > maxBits ? maxBits : bitLimit;
    this._bitOffset = 0;
}
const _bitReaderProto = _BitReaderCtor.prototype as _BitReaderState;

_bitReaderProto.readBit = function(this: _BitReaderState): [boolean, Error | null] {
    if (this._bitOffset >= this._bitLimit) return [false, new Error("not enough bits")];
    const byteIndex = Math.floor(this._bitOffset / 8);
    const bitIndex = 7 - (this._bitOffset % 8);
    this._bitOffset++;
    return [((this._data[byteIndex] & (1 << bitIndex)) !== 0), null];
};

_bitReaderProto.readBits = function(this: _BitReaderState, width: number): [number, Error | null] {
    if (width <= 0 || width > 8) return [0, new Error(`invalid bit width: ${width}`)];
    let v = 0;
    for (let i = 0; i < width; i++) {
        const [bit, err] = this.readBit();
        if (err !== null) return [0, err];
        v <<= 1;
        if (bit) v |= 1;
    }
    return [v, null];
};

Object.defineProperty(_bitReaderProto, "bitOffset", { get(this: _BitReaderState): number { return this._bitOffset; } });

export const BitReader = _BitReaderCtor as unknown as BitReaderCtor;

export interface BitWriter {
    writeBit: (v: boolean) => void;
    writeBits: (v: number, width: number) => Error | null;
    readonly bitLen: number;
    readonly bytes: Uint8Array;
}

type _BitWriterState = BitWriter & {
    _data: Uint8Array;
    _bitOffset: number;
};
type BitWriterCtor = new (bitCapacity: number) => BitWriter;

function _BitWriterCtor(this: _BitWriterState, bitCapacity: number): void {
    this._data = new Uint8Array(headerSize(bitCapacity));
    this._bitOffset = 0;
}
const _bitWriterProto = _BitWriterCtor.prototype as _BitWriterState;

const _bitWriterEnsureBits = function(this: _BitWriterState, width: number): void {
    const requiredBits = this._bitOffset + width;
    const requiredBytes = headerSize(requiredBits);
    if (requiredBytes <= this._data.byteLength) return;
    const next = new Uint8Array(requiredBytes);
    next.set(this._data);
    this._data = next;
};

_bitWriterProto.writeBit = function(this: _BitWriterState, v: boolean): void {
    _bitWriterEnsureBits.call(this, 1);
    const byteIndex = Math.floor(this._bitOffset / 8);
    const bitIndex = 7 - (this._bitOffset % 8);
    if (v) this._data[byteIndex] |= 1 << bitIndex;
    this._bitOffset++;
};

_bitWriterProto.writeBits = function(this: _BitWriterState, v: number, width: number): Error | null {
    if (width <= 0 || width > 8) return new Error(`invalid bit width: ${width}`);
    const maxValue = width === 8 ? 0xff : (1 << width) - 1;
    if (v < 0 || v > maxValue) return new Error(`value ${v} exceeds ${width}-bit max ${maxValue}`);
    for (let shift = width - 1; shift >= 0; shift--) this.writeBit(((v >> shift) & 1) === 1);
    return null;
};

Object.defineProperty(_bitWriterProto, "bitLen", { get(this: _BitWriterState): number { return this._bitOffset; } });
Object.defineProperty(_bitWriterProto, "bytes", { get(this: _BitWriterState): Uint8Array { return this._data.slice(); } });

export const BitWriter = _BitWriterCtor as unknown as BitWriterCtor;

export const eqList = <T>(a: T[], b: T[], eq: (left: T, right: T) => boolean): boolean => {
    if (a === b) return true;
    if (a.length !== b.length) return false;
    for (let i = 0; i < a.length; i++) if (!eq(a[i], b[i])) return false;
    return true;
};

export const eqBool = (a: boolean, b: boolean): boolean => a === b;
export const eqU8 = (a: number, b: number): boolean => a === b;
export const eqI8 = (a: number, b: number): boolean => a === b;
export const eqU16 = (a: number, b: number): boolean => a === b;
export const eqI16 = (a: number, b: number): boolean => a === b;
export const eqU32 = (a: number, b: number): boolean => a === b;
export const eqI32 = (a: number, b: number): boolean => a === b;
export const eqU64 = (a: bigint, b: bigint): boolean => a === b;
export const eqI64 = (a: bigint, b: bigint): boolean => a === b;
export const eqText = (a: string, b: string): boolean => a === b;
export const eqF32 = (a: number, b: number): boolean => a === b || (Number.isNaN(a) && Number.isNaN(b)) || Math.abs(a - b) < 1e-6;
export const eqF64 = (a: number, b: number): boolean => a === b || (Number.isNaN(a) && Number.isNaN(b)) || Math.abs(a - b) < 1e-9;
export const isStringValue = (v: unknown): v is string => typeof v === "string";
export const isBinValue = (v: unknown): v is Uint8Array => v instanceof Uint8Array;
export const isArrayValue = (v: unknown): v is unknown[] => Array.isArray(v);
export const eqBin = (a: Uint8Array, b: Uint8Array): boolean => {
    if (a === b) return true;
    if (a.byteLength !== b.byteLength) return false;
    for (let i = 0; i < a.byteLength; i++) if (a[i] !== b[i]) return false;
    return true;
};
export const eqBinList = (a: Uint8Array[], b: Uint8Array[]): boolean => eqList(a, b, eqBin);

const _checkRead = (buf: Buffer, len: number): Error | null => buf.len < len ? new Error("not enough data") : null;
const _setNum = (buf: Buffer, byteLength: number, setter: (view: DataView, offset: number) => void): void => {
    buf.ensureCapacity(byteLength);
    setter(buf.view, buf.writeOffset);
    (buf as any)._writeOffset += byteLength;
};
const _validateFiniteInteger = (kind: string, v: number, min: number, max: number): Error | null => {
    if (!Number.isFinite(v) || !Number.isInteger(v)) return new Error(`${kind} requires finite integer: ${v}`);
    if (v < min || v > max) return new Error(`${kind} out of range: ${v}`);
    return null;
};

export const getU8 = (buf: Buffer): [number, Error | null] => {
    const err = _checkRead(buf, 1);
    if (err !== null) return [0, err];
    const v = buf.view.getUint8(buf.readOffset);
    (buf as any)._readOffset += 1;
    return [v, null];
};
export const setU8 = (buf: Buffer, v: number): Error | null => {
    const err = _validateFiniteInteger("u8", v, 0, 0xFF);
    if (err !== null) return err;
    _setNum(buf, 1, (view, offset) => view.setUint8(offset, v));
    return null;
};
export const getI8 = (buf: Buffer): [number, Error | null] => {
    const err = _checkRead(buf, 1);
    if (err !== null) return [0, err];
    const v = buf.view.getInt8(buf.readOffset);
    (buf as any)._readOffset += 1;
    return [v, null];
};
export const setI8 = (buf: Buffer, v: number): Error | null => {
    const err = _validateFiniteInteger("i8", v, -0x80, 0x7F);
    if (err !== null) return err;
    _setNum(buf, 1, (view, offset) => view.setInt8(offset, v));
    return null;
};
export const getU16 = (buf: Buffer): [number, Error | null] => {
    const err = _checkRead(buf, 2);
    if (err !== null) return [0, err];
    const v = buf.view.getUint16(buf.readOffset, true);
    (buf as any)._readOffset += 2;
    return [v, null];
};
export const setU16 = (buf: Buffer, v: number): Error | null => {
    const err = _validateFiniteInteger("u16", v, 0, 0xFFFF);
    if (err !== null) return err;
    _setNum(buf, 2, (view, offset) => view.setUint16(offset, v, true));
    return null;
};
export const getI16 = (buf: Buffer): [number, Error | null] => {
    const err = _checkRead(buf, 2);
    if (err !== null) return [0, err];
    const v = buf.view.getInt16(buf.readOffset, true);
    (buf as any)._readOffset += 2;
    return [v, null];
};
export const setI16 = (buf: Buffer, v: number): Error | null => {
    const err = _validateFiniteInteger("i16", v, -0x8000, 0x7FFF);
    if (err !== null) return err;
    _setNum(buf, 2, (view, offset) => view.setInt16(offset, v, true));
    return null;
};
export const getU24 = (buf: Buffer): [number, Error | null] => {
    const [bytes, err] = buf.read(3);
    if (err !== null) return [0, err];
    return [bytes[0] | (bytes[1] << 8) | (bytes[2] << 16), null];
};
export const setU24 = (buf: Buffer, v: number): Error | null => {
    const err = _validateFiniteInteger("u24", v, 0, 0xFFFFFF);
    if (err !== null) return err;
    buf.ensureCapacity(3);
    const bytes = buf.bytes;
    (buf as any)._bytes[buf.writeOffset] = v & 0xff;
    (buf as any)._bytes[buf.writeOffset + 1] = (v >>> 8) & 0xff;
    (buf as any)._bytes[buf.writeOffset + 2] = (v >>> 16) & 0xff;
    (buf as any)._writeOffset += 3;
    void bytes;
    return null;
};
export const getU32 = (buf: Buffer): [number, Error | null] => {
    const err = _checkRead(buf, 4);
    if (err !== null) return [0, err];
    const v = buf.view.getUint32(buf.readOffset, true);
    (buf as any)._readOffset += 4;
    return [v, null];
};
export const setU32 = (buf: Buffer, v: number): Error | null => {
    const err = _validateFiniteInteger("u32", v, 0, 0xFFFFFFFF);
    if (err !== null) return err;
    _setNum(buf, 4, (view, offset) => view.setUint32(offset, v, true));
    return null;
};
export const getI32 = (buf: Buffer): [number, Error | null] => {
    const err = _checkRead(buf, 4);
    if (err !== null) return [0, err];
    const v = buf.view.getInt32(buf.readOffset, true);
    (buf as any)._readOffset += 4;
    return [v, null];
};
export const setI32 = (buf: Buffer, v: number): Error | null => {
    const err = _validateFiniteInteger("i32", v, -0x80000000, 0x7FFFFFFF);
    if (err !== null) return err;
    _setNum(buf, 4, (view, offset) => view.setInt32(offset, v, true));
    return null;
};

const _validateBigIntRange = (kind: string, v: unknown, min: bigint, max: bigint): Error | null => {
	if (typeof v !== "bigint") return new Error(`${kind} requires bigint: ${String(v)}`);
    if (v < min || v > max) return new Error(`${kind} out of range: ${v}`);
    return null;
};
export const getU64 = (buf: Buffer): [bigint, Error | null] => {
    const err = _checkRead(buf, 8);
    if (err !== null) return [0n, err];
    const v = buf.view.getBigUint64(buf.readOffset, true);
    (buf as any)._readOffset += 8;
    return [v, null];
};
export const setU64 = (buf: Buffer, v: bigint): Error | null => {
    const err = _validateBigIntRange("u64", v, 0n, 0xFFFF_FFFF_FFFF_FFFFn);
    if (err !== null) return err;
    _setNum(buf, 8, (view, offset) => view.setBigUint64(offset, v, true));
    return null;
};
export const getI64 = (buf: Buffer): [bigint, Error | null] => {
    const err = _checkRead(buf, 8);
    if (err !== null) return [0n, err];
    const v = buf.view.getBigInt64(buf.readOffset, true);
    (buf as any)._readOffset += 8;
    return [v, null];
};
export const setI64 = (buf: Buffer, v: bigint): Error | null => {
    const err = _validateBigIntRange("i64", v, -0x8000_0000_0000_0000n, 0x7FFF_FFFF_FFFF_FFFFn);
    if (err !== null) return err;
    _setNum(buf, 8, (view, offset) => view.setBigInt64(offset, v, true));
    return null;
};

export const getF32 = (buf: Buffer): [number, Error | null] => {
    const err = _checkRead(buf, 4);
    if (err !== null) return [0, err];
    const v = buf.view.getFloat32(buf.readOffset, true);
    (buf as any)._readOffset += 4;
    return [v, null];
};
export const setF32 = (buf: Buffer, v: number): Error | null => { _setNum(buf, 4, (view, offset) => view.setFloat32(offset, v, true)); return null; };
export const getF64 = (buf: Buffer): [number, Error | null] => {
    const err = _checkRead(buf, 8);
    if (err !== null) return [0, err];
    const v = buf.view.getFloat64(buf.readOffset, true);
    (buf as any)._readOffset += 8;
    return [v, null];
};
export const setF64 = (buf: Buffer, v: number): Error | null => { _setNum(buf, 8, (view, offset) => view.setFloat64(offset, v, true)); return null; };
