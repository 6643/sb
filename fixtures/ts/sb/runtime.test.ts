import { describe, expect, test } from "bun:test";

import {
    Buffer,
    decodeBitmap,
    encodeBitmap,
    encodeStateBlock,
    decodeStateBlock,
    getF32,
    getF64,
    getBin,
    getText,
    getListCount,
    headerSize,
    readHeader,
    StateU16,
    setF32,
    setF64,
    setI16,
    setI32,
    setI64,
    setI8,
    setBin,
    setText,
    setU16,
    setU24,
    setU32,
    setU64,
    setU8,
    StateU8,
    validatePaddingZero,
    writeHeader,
} from "./type";

const encodeGameWireForTest = (id: number, name: string): Uint8Array => {
    const [nameState, nameErr] = name.length === 0
        ? [0, null]
        : name.length <= 0xff
            ? [1, null]
            : name.length <= 0xffff
                ? [2, null]
                : [0, new Error(`text length exceeds u16 max: ${name.length}`)];
    if (nameErr !== null) throw nameErr;

    const [header, headerErr] = writeHeader([1, 2], [id !== 0 ? 1 : 0, nameState]);
    if (headerErr !== undefined) throw headerErr;

    const buf = new Buffer();
    expect(buf.write(header)).toBeNull();
    if (id !== 0) expect(setU32(buf, id)).toBeNull();
    expect(setText(buf, nameState, name)).toBeNull();
    return buf.bytes;
};

describe("runtime header helpers", () => {
    test("./type keeps representative runtime exports usable", () => {
        const [header, headerErr] = writeHeader([1, 1], [1, 0]);
        expect(headerErr).toBeUndefined();
        expect(Array.from(header)).toEqual([0x80]);

        const textBuf = new Buffer(new Uint8Array([0x02, 0x68, 0x69]));
        const [text, textErr] = getText(textBuf, StateU8);
        expect(textErr).toBeNull();
        expect(text).toBe("hi");

        const binBuf = new Buffer();
        expect(setBin(binBuf, StateU8, new Uint8Array([0x01, 0x02]))).toBeNull();
        expect(binBuf.bytes).toEqual(new Uint8Array([0x02, 0x01, 0x02]));
    });

    test("protocol Game vectors match spec bytes", () => {
        expect(encodeGameWireForTest(0, "lol")).toEqual(new Uint8Array([0x20, 0x03, 0x6c, 0x6f, 0x6c]));
        expect(encodeGameWireForTest(7, "")).toEqual(new Uint8Array([0x80, 0x07, 0x00, 0x00, 0x00]));
        expect(encodeGameWireForTest(7, "lol")).toEqual(new Uint8Array([0xa0, 0x07, 0x00, 0x00, 0x00, 0x03, 0x6c, 0x6f, 0x6c]));
    });

    test("readHeader rejects oversized header bytes with Go-compatible error", () => {
        const [, err] = readHeader(new Uint8Array([0x80, 0x00]), [1, 2, 1], "demo header");

        expect(err?.message).toBe("demo header invalid header size: 2 != 1");
    });

    test("readHeader rejects undersized header bytes with Go-compatible error", () => {
        const [, err] = readHeader(new Uint8Array(0), [1, 2, 1], "demo header");

        expect(err?.message).toBe("demo header invalid header size: 0 != 1");
    });

    test("readHeader rejects invalid field width with Go-compatible error", () => {
        const [, err] = readHeader(new Uint8Array([0x00]), [1, 3], "demo header");

        expect(err?.message).toBe("demo header: header field[1] invalid width: 3");
    });

    test("writeHeader rejects invalid field width with Go-compatible error", () => {
        const [, err] = writeHeader([1, 3], [1, 0]);

        expect(err?.message).toBe("header field[1] invalid width: 3");
    });

    test("writeHeader returns exact-sized bytes that readHeader roundtrips", () => {
        const widths = [2, 2, 2, 2, 1, 1] as const;
        const values = [1, 2, 3, 0, 1, 0] as const;
        const [header, writeErr] = writeHeader(widths, values);

        expect(writeErr).toBeUndefined();
        expect(header.byteLength).toBe(headerSize(widths.reduce((sum, width) => sum + width, 0)));

        const [decoded, readErr] = readHeader(header, widths, "demo header");
        expect(readErr).toBeUndefined();
        expect(decoded).toEqual([...values]);
    });

	test("protocol demo cross-byte header matches spec bytes", () => {
		const [header, writeErr] = writeHeader([1, 2, 2, 2, 2], [1, 1, 2, 0, 1]);

		expect(writeErr).toBeUndefined();
		expect(Array.from(header)).toEqual([0xB0, 0x80]);

		const [decoded, readErr] = readHeader(header, [1, 2, 2, 2, 2], "Demo header");
		expect(readErr).toBeUndefined();
		expect(decoded).toEqual([1, 1, 2, 0, 1]);
	});
});

describe("runtime list and block guards", () => {
    test("decodeBitmap rejects negative count with Error instead of RangeError", () => {
        const [bits, err] = decodeBitmap(new Uint8Array(0), -1);

        expect(bits).toEqual([]);
        expect(err).toBeInstanceOf(Error);
        expect(err).not.toBeInstanceOf(RangeError);
        expect(err?.message).toBe("bitmap invalid count: -1");
    });

    test("decodeStateBlock rejects negative count with Error instead of RangeError", () => {
        const [states, err] = decodeStateBlock(new Uint8Array(0), -1);

        expect(states).toEqual([]);
        expect(err).toBeInstanceOf(Error);
        expect(err).not.toBeInstanceOf(RangeError);
        expect(err?.message).toBe("state block invalid count: -1");
    });

    test("getListCount reports encoded zero count", () => {
        const buf = new Buffer(new Uint8Array([0x00]));

        const [, err] = getListCount(buf, StateU8);
        expect(err?.message).toBe("list count state 1 encoded zero count");
    });

    test("getListCount rejects illegal list state", () => {
        const buf = new Buffer(new Uint8Array(0));

        const [, err] = getListCount(buf, 2);
        expect(err?.message).toBe("list count state 2 is illegal");
    });

    test("validatePaddingZero rejects negative usedBits", () => {
        const err = validatePaddingZero(new Uint8Array([0x00]), -1, "bitmap");

        expect(err?.message).toBe("bitmap invalid used bits: -1 not in [0,8]");
    });

    test("validatePaddingZero rejects oversized usedBits", () => {
        const err = validatePaddingZero(new Uint8Array([0x00]), 9, "bitmap");

        expect(err?.message).toBe("bitmap invalid used bits: 9 not in [0,8]");
    });

	test("protocol text list item header block matches spec bytes", () => {
		const [header, err] = encodeStateBlock([1, 0, 2, 1, 0]);

		expect(err).toBeNull();
		expect(Array.from(header)).toEqual([0x49, 0x00]);

		const [states, decodeErr] = decodeStateBlock(header, 5);
		expect(decodeErr).toBeNull();
		expect(states).toEqual([1, 0, 2, 1, 0]);
	});

	test("protocol bitmap examples match spec bytes", () => {
		expect(Array.from(encodeBitmap([true, true, false, false, true]))).toEqual([0xC8]);
		expect(Array.from(encodeBitmap([false, true, false, true]))).toEqual([0x50]);

		const [bits, err] = decodeBitmap(new Uint8Array([0xC8]), 5);
		expect(err).toBeNull();
		expect(bits).toEqual([true, true, false, false, true]);
	});
});

describe("runtime numeric setters", () => {
    test.each([
        ["setU8", setU8, -1, "u8 out of range: -1"],
        ["setU8", setU8, 1.5, "u8 requires finite integer: 1.5"],
        ["setU8", setU8, Number.NaN, "u8 requires finite integer: NaN"],
        ["setU8", setU8, Number.POSITIVE_INFINITY, "u8 requires finite integer: Infinity"],
        ["setU8", setU8, 256, "u8 out of range: 256"],
        ["setI8", setI8, -129, "i8 out of range: -129"],
        ["setI8", setI8, -1.5, "i8 requires finite integer: -1.5"],
        ["setI8", setI8, Number.NaN, "i8 requires finite integer: NaN"],
        ["setI8", setI8, Number.NEGATIVE_INFINITY, "i8 requires finite integer: -Infinity"],
        ["setI8", setI8, 128, "i8 out of range: 128"],
        ["setU16", setU16, -1, "u16 out of range: -1"],
        ["setU16", setU16, 1.5, "u16 requires finite integer: 1.5"],
        ["setU16", setU16, Number.NaN, "u16 requires finite integer: NaN"],
        ["setU16", setU16, Number.POSITIVE_INFINITY, "u16 requires finite integer: Infinity"],
        ["setU16", setU16, 65536, "u16 out of range: 65536"],
        ["setI16", setI16, -32769, "i16 out of range: -32769"],
        ["setI16", setI16, -1.5, "i16 requires finite integer: -1.5"],
        ["setI16", setI16, Number.NaN, "i16 requires finite integer: NaN"],
        ["setI16", setI16, Number.NEGATIVE_INFINITY, "i16 requires finite integer: -Infinity"],
        ["setI16", setI16, 32768, "i16 out of range: 32768"],
        ["setU24", setU24, -1, "u24 out of range: -1"],
        ["setU24", setU24, 1.5, "u24 requires finite integer: 1.5"],
        ["setU24", setU24, Number.NaN, "u24 requires finite integer: NaN"],
        ["setU24", setU24, Number.POSITIVE_INFINITY, "u24 requires finite integer: Infinity"],
        ["setU24", setU24, 0x1000000, "u24 out of range: 16777216"],
        ["setU32", setU32, -1, "u32 out of range: -1"],
        ["setU32", setU32, 1.5, "u32 requires finite integer: 1.5"],
        ["setU32", setU32, Number.NaN, "u32 requires finite integer: NaN"],
        ["setU32", setU32, Number.POSITIVE_INFINITY, "u32 requires finite integer: Infinity"],
        ["setU32", setU32, 0x100000000, "u32 out of range: 4294967296"],
        ["setI32", setI32, -2147483649, "i32 out of range: -2147483649"],
        ["setI32", setI32, -1.5, "i32 requires finite integer: -1.5"],
        ["setI32", setI32, Number.NaN, "i32 requires finite integer: NaN"],
        ["setI32", setI32, Number.NEGATIVE_INFINITY, "i32 requires finite integer: -Infinity"],
        ["setI32", setI32, 2147483648, "i32 out of range: 2147483648"],
    ])("%s rejects invalid numeric input %p", (_name, setter, value, message) => {
        const buf = new Buffer();

        const err = setter(buf, value);

        expect(err).toBeInstanceOf(Error);
        expect(err?.message).toBe(message);
        expect(buf.bytes).toEqual(new Uint8Array(0));
    });

    test.each([
        ["setF32", setF32, getF32],
        ["setF64", setF64, getF64],
    ])("%s accepts finite numbers, NaN, and Infinity", (_name, setter, getter) => {
        const finiteBuf = new Buffer();
        expect(setter(finiteBuf, 1.25)).toBeNull();
        expect(finiteBuf.bytes.byteLength).toBeGreaterThan(0);
        const [finiteValue, finiteErr] = getter(new Buffer(finiteBuf.bytes));
        expect(finiteErr).toBeNull();
        expect(finiteValue).toBe(1.25);

        const nanBuf = new Buffer();
        expect(setter(nanBuf, Number.NaN)).toBeNull();
        const [nanValue, nanErr] = getter(new Buffer(nanBuf.bytes));
        expect(nanErr).toBeNull();
        expect(Number.isNaN(nanValue)).toBeTrue();

        const infBuf = new Buffer();
        expect(setter(infBuf, Number.POSITIVE_INFINITY)).toBeNull();
        const [infValue, infErr] = getter(new Buffer(infBuf.bytes));
        expect(infErr).toBeNull();
        expect(infValue).toBe(Number.POSITIVE_INFINITY);
    });

    test.each([
        ["setU64", setU64, -1n, "u64 out of range: -1"],
        ["setU64", setU64, 0x1_0000_0000_0000_0000n, "u64 out of range: 18446744073709551616"],
        ["setI64", setI64, -0x8000_0000_0000_0001n, "i64 out of range: -9223372036854775809"],
        ["setI64", setI64, 0x8000_0000_0000_0000n, "i64 out of range: 9223372036854775808"],
    ])("%s returns Error instead of throwing for out-of-range bigint %p", (_name, setter, value, message) => {
        const buf = new Buffer();
        let thrown: unknown;
        let err: Error | null = null;

        try {
            err = setter(buf, value);
        } catch (cause) {
            thrown = cause;
        }

        expect(thrown).toBeUndefined();
        expect(err).toBeInstanceOf(Error);
        expect(err?.message).toBe(message);
        expect(buf.bytes).toEqual(new Uint8Array(0));
    });

    test.each([
        ["setU64", setU64, 1 as any, "u64 requires bigint: 1"],
        ["setI64", setI64, "1" as any, "i64 requires bigint: 1"],
    ])("%s returns Error instead of throwing for non-bigint input %p", (_name, setter, value, message) => {
        const buf = new Buffer();
        let thrown: unknown;
        let err: Error | null = null;

        try {
            err = setter(buf, value);
        } catch (cause) {
            thrown = cause;
        }

        expect(thrown).toBeUndefined();
        expect(err).toBeInstanceOf(Error);
        expect(err?.message).toBe(message);
        expect(buf.bytes).toEqual(new Uint8Array(0));
    });
});

describe("runtime text UTF-8 semantics", () => {
    test("getText rejects invalid UTF-8 bytes", () => {
        const buf = new Buffer(new Uint8Array([0x01, 0xff]));

        const [, err] = getText(buf, StateU8);

        expect(err).toBeInstanceOf(Error);
    });

    test("setText rejects invalid UTF-8 input", () => {
        const buf = new Buffer();

        const err = setText(buf, StateU8, "\ud800");

        expect(err).toBeInstanceOf(Error);
    });

    test("valid UTF-8 text roundtrips", () => {
        const writeBuf = new Buffer();
        const value = "你好, runtime";

        expect(setText(writeBuf, StateU8, value)).toBeNull();

        const [got, err] = getText(new Buffer(writeBuf.bytes), StateU8);
        expect(err).toBeNull();
        expect(got).toBe(value);
    });

	test("getText rejects non-canonical u16 state for short payload", () => {
		const [, err] = getText(new Buffer(new Uint8Array([0x03, 0x00, 0x6c, 0x6f, 0x6c])), StateU16);

		expect(err?.message).toBe("text state 2 is not canonical for length 3");
	});

	test("getBin rejects non-canonical u16 state for short payload", () => {
		const [, err] = getBin(new Buffer(new Uint8Array([0x03, 0x00, 0x09, 0x08, 0x07])), StateU16);

		expect(err?.message).toBe("bin state 2 is not canonical for length 3");
	});
});
