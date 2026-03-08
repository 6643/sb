const _textEncoder = new TextEncoder();
const _textDecoder = new TextDecoder();

export const StateZero = 0;
export const StateU8 = 1;
export const StateU16 = 2;
export const StateU24 = 3;

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

_bufferProto.ensureCapacity = function(this: _BufferState, needed: number): void {
    const required = this._writeOffset + needed;
    if (required <= this._bytes.length) return;
    let newCapacity = this._bytes.length === 0 ? 128 : this._bytes.length;
    while (newCapacity < required) newCapacity *= 2;
    const next = new Uint8Array(newCapacity);
    next.set(this._bytes.subarray(0, this._writeOffset));
    this._bytes = next;
    this._view = new DataView(next.buffer);
};

_bufferProto.read = function(this: _BufferState, byteLength: number): [Uint8Array, Error | null] {
    if (byteLength < 0) return [new Uint8Array(0), new Error(`invalid read length ${byteLength}`)];
    if (this._readOffset + byteLength > this._writeOffset) return [new Uint8Array(0), new Error("not enough data")];
    const slice = this._bytes.subarray(this._readOffset, this._readOffset + byteLength);
    this._readOffset += byteLength;
    return [slice, null];
};

_bufferProto.write = function(this: _BufferState, data: Uint8Array): Error | null {
    this.ensureCapacity(data.byteLength);
    this.writeUnsafe(data);
    return null;
};

_bufferProto.writeUnsafe = function(this: _BufferState, data: Uint8Array): void {
    this._bytes.set(data, this._writeOffset);
    this._writeOffset += data.byteLength;
};

_bufferProto.rewindWrite = function(this: _BufferState, offset: number): void {
    if (offset < 0) offset = 0;
    if (offset > this._writeOffset) offset = this._writeOffset;
    this._writeOffset = offset;
    if (this._readOffset > this._writeOffset) this._readOffset = this._writeOffset;
};

Object.defineProperty(_bufferProto, "bytes", { get(this: _BufferState): Uint8Array { return this._bytes.subarray(0, this._writeOffset); } });
Object.defineProperty(_bufferProto, "view", { get(this: _BufferState): DataView { return this._view; } });
Object.defineProperty(_bufferProto, "len", { get(this: _BufferState): number { return this._writeOffset - this._readOffset; } });
Object.defineProperty(_bufferProto, "readOffset", { get(this: _BufferState): number { return this._readOffset; } });
Object.defineProperty(_bufferProto, "writeOffset", { get(this: _BufferState): number { return this._writeOffset; } });

export const Buffer = _BufferCtor as unknown as BufferCtor;

export type Err = Error | undefined;
export const errU = (err: Error | null): Err => err === null ? undefined : err;
export const resultU = <T>(value: T, err: Error | null): [T, Err] => err === null ? [value, undefined] : [value, err];

export interface Serializable {
    set: (buf: Buffer) => Err;
}

export interface Deserializable {
    get: (buf: Buffer) => Err;
}

export const headerSize = (bitCount: number): number => bitCount <= 0 ? 0 : Math.floor((bitCount + 7) / 8);
export const bitmapSize = (count: number): number => count <= 0 ? 0 : Math.floor((count + 7) / 8);
export const itemHeaderSize = (count: number): number => count <= 0 ? 0 : Math.floor((count * 2 + 7) / 8);

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
export const eqF32 = (a: number, b: number): boolean => Math.abs(a - b) < 1e-6;
export const eqF64 = (a: number, b: number): boolean => Math.abs(a - b) < 1e-9;
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

export const getU8 = (buf: Buffer): [number, Error | null] => {
    const err = _checkRead(buf, 1);
    if (err !== null) return [0, err];
    const v = buf.view.getUint8(buf.readOffset);
    (buf as any)._readOffset += 1;
    return [v, null];
};
export const setU8 = (buf: Buffer, v: number): Error | null => { _setNum(buf, 1, (view, offset) => view.setUint8(offset, v)); return null; };
export const getI8 = (buf: Buffer): [number, Error | null] => {
    const err = _checkRead(buf, 1);
    if (err !== null) return [0, err];
    const v = buf.view.getInt8(buf.readOffset);
    (buf as any)._readOffset += 1;
    return [v, null];
};
export const setI8 = (buf: Buffer, v: number): Error | null => { _setNum(buf, 1, (view, offset) => view.setInt8(offset, v)); return null; };
export const getU16 = (buf: Buffer): [number, Error | null] => {
    const err = _checkRead(buf, 2);
    if (err !== null) return [0, err];
    const v = buf.view.getUint16(buf.readOffset, true);
    (buf as any)._readOffset += 2;
    return [v, null];
};
export const setU16 = (buf: Buffer, v: number): Error | null => { _setNum(buf, 2, (view, offset) => view.setUint16(offset, v, true)); return null; };
export const getI16 = (buf: Buffer): [number, Error | null] => {
    const err = _checkRead(buf, 2);
    if (err !== null) return [0, err];
    const v = buf.view.getInt16(buf.readOffset, true);
    (buf as any)._readOffset += 2;
    return [v, null];
};
export const setI16 = (buf: Buffer, v: number): Error | null => { _setNum(buf, 2, (view, offset) => view.setInt16(offset, v, true)); return null; };
export const getU24 = (buf: Buffer): [number, Error | null] => {
    const [bytes, err] = buf.read(3);
    if (err !== null) return [0, err];
    return [bytes[0] | (bytes[1] << 8) | (bytes[2] << 16), null];
};
export const setU24 = (buf: Buffer, v: number): Error | null => {
    if (v < 0 || v > 0xFFFFFF) return new Error(`u24 out of range: ${v}`);
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
export const setU32 = (buf: Buffer, v: number): Error | null => { _setNum(buf, 4, (view, offset) => view.setUint32(offset, v, true)); return null; };
export const getI32 = (buf: Buffer): [number, Error | null] => {
    const err = _checkRead(buf, 4);
    if (err !== null) return [0, err];
    const v = buf.view.getInt32(buf.readOffset, true);
    (buf as any)._readOffset += 4;
    return [v, null];
};
export const setI32 = (buf: Buffer, v: number): Error | null => { _setNum(buf, 4, (view, offset) => view.setInt32(offset, v, true)); return null; };
export const getU64 = (buf: Buffer): [bigint, Error | null] => {
    const err = _checkRead(buf, 8);
    if (err !== null) return [0n, err];
    const v = buf.view.getBigUint64(buf.readOffset, true);
    (buf as any)._readOffset += 8;
    return [v, null];
};
export const setU64 = (buf: Buffer, v: bigint): Error | null => { _setNum(buf, 8, (view, offset) => view.setBigUint64(offset, v, true)); return null; };
export const getI64 = (buf: Buffer): [bigint, Error | null] => {
    const err = _checkRead(buf, 8);
    if (err !== null) return [0n, err];
    const v = buf.view.getBigInt64(buf.readOffset, true);
    (buf as any)._readOffset += 8;
    return [v, null];
};
export const setI64 = (buf: Buffer, v: bigint): Error | null => { _setNum(buf, 8, (view, offset) => view.setBigInt64(offset, v, true)); return null; };
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

export const encodeBitmap = (bits: boolean[]): Uint8Array => {
    const writer = new BitWriter(bits.length);
    for (const bit of bits) writer.writeBit(bit);
    return writer.bytes;
};
export const decodeBitmap = (data: Uint8Array, count: number): [boolean[], Error | null] => {
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
export const encodeStateBlock = (states: number[]): [Uint8Array, Error | null] => {
    const writer = new BitWriter(states.length * 2);
    for (let i = 0; i < states.length; i++) {
        const err = writer.writeBits(states[i], 2);
        if (err !== null) return [new Uint8Array(0), new Error(`state[${i}]: ${err.message}`)];
    }
    return [writer.bytes, null];
};
export const decodeStateBlock = (data: Uint8Array, count: number): [number[], Error | null] => {
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

export const validatePaddingZero = (data: Uint8Array, usedBits: number, kind: string): Error | null => {
    for (let bitOffset = usedBits; bitOffset < data.byteLength * 8; bitOffset++) {
        const byteIndex = Math.floor(bitOffset / 8);
        const bitIndex = 7 - (bitOffset % 8);
        if ((data[byteIndex] & (1 << bitIndex)) !== 0) return new Error(`${kind} padding bit ${bitOffset - usedBits} is not zero`);
    }
    return null;
};

const _headerBitCount = (widths: readonly number[]): [number, Error | null] => {
    let total = 0;
    for (let i = 0; i < widths.length; i++) {
        const width = widths[i];
        if (width !== 1 && width !== 2) return [0, new Error(`header width[${i}] is invalid: ${width}`)];
        total += width;
    }
    return [total, null];
};
export const readHeader = (data: Uint8Array, widths: readonly number[], kind: string): [number[], Err] => {
    const [usedBits, err] = _headerBitCount(widths);
    if (err !== null) return [[], new Error(`${kind}: ${err.message}`)];
    const expectBytes = headerSize(usedBits);
    if (data.byteLength !== expectBytes) return [[], new Error(`${kind}: size mismatch ${data.byteLength} !== ${expectBytes}`)];
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

const _textBytes = (v: string): Uint8Array => _textEncoder.encode(v);
const _textLengthState = (length: number): [number, Error | null] => {
    if (length < 0) return [0, new Error(`negative text length: ${length}`)];
    if (length === 0) return [StateZero, null];
    if (length <= 0xFF) return [StateU8, null];
    if (length <= 0xFFFF) return [StateU16, null];
    return [0, new Error(`text length exceeds u16 max: ${length}`)];
};
export const textState = (v: string): [number, Error | null] => _textLengthState(_textBytes(v).byteLength);
export const binState = (length: number): [number, Error | null] => {
    if (length < 0) return [0, new Error(`negative bin length: ${length}`)];
    if (length === 0) return [StateZero, null];
    if (length <= 0xFF) return [StateU8, null];
    if (length <= 0xFFFF) return [StateU16, null];
    if (length <= 0xFFFFFF) return [StateU24, null];
    return [0, new Error(`bin length exceeds u24 max: ${length}`)];
};
export const listCountState = (count: number): [number, Error | null] => {
    if (count < 0) return [0, new Error(`negative list count: ${count}`)];
    if (count === 0) return [StateZero, null];
    if (count <= 0xFF) return [StateU8, null];
    return [0, new Error(`list count exceeds u8 max: ${count}`)];
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

export const getText = (buf: Buffer, state: number): [string, Error | null] => {
    const [length, err] = _getCompactLength(buf, state, 0xFFFF, "text");
    if (err !== null) return ["", err];
    if (length === 0) return ["", null];
    const [data, err2] = buf.read(length);
    if (err2 !== null) return ["", new Error(`text body not enough data: ${buf.len} - ${length}`)];
    return [_textDecoder.decode(data), null];
};
export const setText = (buf: Buffer, state: number, v: string): Error | null => {
    const bytes = _textBytes(v);
    const [canonical, err] = _textLengthState(bytes.byteLength);
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

export const getBin = (buf: Buffer, state: number): [Uint8Array, Error | null] => {
    const [length, err] = _getCompactLength(buf, state, 0xFFFFFF, "bin");
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

export const getListCount = (buf: Buffer, state: number): [number, Error | null] => {
    switch (state) {
    case StateZero:
        return [0, null];
    case StateU8: {
        const [v, err] = getU8(buf);
        if (err !== null) return [0, err];
        if (v === 0) return [0, new Error(`list count state ${state} encoded zero length`)];
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
