import * as rtc from "./runtime_core"

type Buffer = rtc.Buffer
type Err = rtc.Err
type Serializable = rtc.Serializable
type Deserializable = rtc.Deserializable

export interface StructField<T> {
    key: keyof T & string;
    label: string;
    width: 1 | 2;
    defaultValue: () => any;
    headerState: (value: any) => [number, Err];
    read: (buf: Buffer, state: number) => [any, Err];
    write: (buf: Buffer, state: number, value: any) => Err;
    validate: (value: any) => Err;
    isZero: (value: any) => boolean;
    eq: (left: any, right: any) => boolean;
}

export interface StructMeta<T> {
    name: string;
    headerWidths: readonly number[];
    create: () => T;
    fields: readonly StructField<T>[];
}

export const defineStruct = <T>(meta: StructMeta<T>): StructMeta<T> => meta;

const _structHeaderBits = (meta: StructMeta<any>): number => {
    let total = 0;
    for (let i = 0; i < meta.headerWidths.length; i++) total += meta.headerWidths[i];
    return total;
};

const _structValue = <T>(value: T, key: keyof T & string): any => (value as any)[key];
const _structAssign = <T>(value: T, key: keyof T & string, next: any): void => {
    (value as any)[key] = next;
};

export const newStruct = <T extends Serializable & Deserializable>(
    meta: StructMeta<T>,
    getter: (buf: Buffer) => [T, Err],
    setter: (buf: Buffer, value: T) => Err,
): T => {
    const s = meta.create();
    s.set = (buf: Buffer) => setter(buf, s);
    s.get = (buf: Buffer) => {
        const [next, err] = getter(buf);
        if (err === undefined) Object.assign(s, next);
        return err;
    };
    return s;
};

export const isZeroStruct = <T>(meta: StructMeta<T>, value: T | null | undefined): boolean => {
    if (value === null || value === undefined) return true;
    for (let i = 0; i < meta.fields.length; i++) {
        const field = meta.fields[i];
        if (!field.isZero(_structValue(value, field.key))) return false;
    }
    return true;
};

export const validateStruct = <T>(meta: StructMeta<T>, value: T | null | undefined): Err => {
    if (value === null || value === undefined) return undefined;
    for (let i = 0; i < meta.fields.length; i++) {
        const field = meta.fields[i];
        const err = field.validate(_structValue(value, field.key));
        if (err !== undefined) return new Error(`${field.label}: ${err.message}`);
    }
    return undefined;
};

export const eqStruct = <T>(meta: StructMeta<T>, a: T | null | undefined, b: T | null | undefined): boolean => {
    if (isZeroStruct(meta, a) && isZeroStruct(meta, b)) return true;
    if (a === null || a === undefined || b === null || b === undefined) return false;
    for (let i = 0; i < meta.fields.length; i++) {
        const field = meta.fields[i];
        if (!field.eq(_structValue(a, field.key), _structValue(b, field.key))) return false;
    }
    return true;
};

export const getStruct = <T>(meta: StructMeta<T>, buf: Buffer): [T, Err] => {
    const value = meta.create();
    const [header, errHeader] = buf.read(rtc.headerSize(_structHeaderBits(meta)));
    if (errHeader !== null) return [value, new Error("not enough data")];
    const [headerStates, errHeaderState] = rtc.readHeader(header, meta.headerWidths, `${meta.name} header`);
    if (errHeaderState !== undefined) return [value, errHeaderState];
    for (let i = 0; i < meta.fields.length; i++) {
        const field = meta.fields[i];
        const [next, err] = field.read(buf, headerStates[i]);
        if (err !== undefined) return [value, new Error(`get ${meta.name} ${field.label}: ${err.message}`)];
        _structAssign(value, field.key, next);
    }
    const errValidate = validateStruct(meta, value);
    if (errValidate !== undefined) return [value, new Error(`validate failed: ${errValidate.message}`)];
    return [value, undefined];
};

export const setStruct = <T>(meta: StructMeta<T>, buf: Buffer, value: T | null | undefined): Err => {
    if (value === null || value === undefined) return new Error(`set ${meta.name}: value is null or undefined`);
    const errValidate = validateStruct(meta, value);
    if (errValidate !== undefined) return new Error(`validate ${meta.name}: ${errValidate.message}`);
    const startOffset = buf.writeOffset;
    const headerStates = new Array<number>(meta.fields.length);
    for (let i = 0; i < meta.fields.length; i++) {
        const field = meta.fields[i];
        const [state, err] = field.headerState(_structValue(value, field.key));
        if (err !== undefined) return err;
        headerStates[i] = state;
    }
    const [header, errHeader] = rtc.writeHeader(meta.headerWidths, headerStates);
    if (errHeader !== undefined) { buf.rewindWrite(startOffset); return new Error(`set header: ${errHeader.message}`); }
    const errHeaderWrite = buf.write(header);
    if (errHeaderWrite !== null) { buf.rewindWrite(startOffset); return errHeaderWrite; }
    for (let i = 0; i < meta.fields.length; i++) {
        const field = meta.fields[i];
        const err = field.write(buf, headerStates[i], _structValue(value, field.key));
        if (err !== undefined) { buf.rewindWrite(startOffset); return new Error(`set ${meta.name} ${field.label}: ${err.message}`); }
    }
    return undefined;
};

export const boolField = <T>(key: keyof T & string, label: string): StructField<T> => ({
    key,
    label,
    width: 1,
    defaultValue: () => false,
    headerState: (value) => [value ? 1 : 0, undefined],
    read: (_buf, state) => [state === 1, undefined],
    write: () => undefined,
    validate: () => undefined,
    isZero: (value) => value !== true,
    eq: (left, right) => rtc.eqBool(left as any, right as any),
});

export const scalarField = <T, V>(
    key: keyof T & string,
    label: string,
    zeroValue: V,
    getItem: (buf: Buffer) => [V, Error | null],
    setItem: (buf: Buffer, value: V) => Error | null,
    eq: (left: V, right: V) => boolean,
): StructField<T> => ({
    key,
    label,
    width: 1,
    defaultValue: () => zeroValue,
    headerState: (value) => [eq(value as V, zeroValue) ? 0 : 1, undefined],
    read: (buf, state) => state === 0 ? [zeroValue, undefined] : rtc.resultU(...getItem(buf)),
    write: (buf, state, value) => state === 0 ? undefined : rtc.errU(setItem(buf, value as V)),
    validate: () => undefined,
    isZero: (value) => eq(value as V, zeroValue),
    eq: (left, right) => eq(left as V, right as V),
});

export const textField = <T>(key: keyof T & string, label: string): StructField<T> => ({
    key,
    label,
    width: 2,
    defaultValue: () => "",
    headerState: (value) => { const [state, err] = rtc.textState(value as any); return [state, rtc.errU(err)]; },
    read: (buf, state) => rtc.resultU(...rtc.getText(buf, state)),
    write: (buf, state, value) => rtc.errU(rtc.setText(buf, state, value as string)),
    validate: (value) => { const [, err] = rtc.textState(value as any); return rtc.errU(err); },
    isZero: (value) => rtc.isStringValue(value) && value === "",
    eq: (left, right) => rtc.eqText(left as any, right as any),
});

export const binField = <T>(key: keyof T & string, label: string): StructField<T> => ({
    key,
    label,
    width: 2,
    defaultValue: () => new Uint8Array(0),
    headerState: (value) => {
        if (!rtc.isBinValue(value)) return [0, new Error("value is not Uint8Array")];
        const [state, err] = rtc.binState(value.byteLength);
        return [state, rtc.errU(err)];
    },
    read: (buf, state) => rtc.resultU(...rtc.getBin(buf, state)),
    write: (buf, state, value) => rtc.errU(rtc.setBin(buf, state, value as Uint8Array)),
    validate: (value) => {
        if (!rtc.isBinValue(value)) return new Error("value is not Uint8Array");
        const [, err] = rtc.binState(value.byteLength);
        return rtc.errU(err);
    },
    isZero: (value) => rtc.isBinValue(value) && value.byteLength === 0,
    eq: (left, right) => rtc.eqBin(left as any, right as any),
});

export const enumField = <T, V>(
    key: keyof T & string,
    label: string,
    defaultValue: () => V,
    isValue: (value: V) => boolean,
    isAssignable: (value: V) => boolean,
    normalize: (value: V) => V,
): StructField<T> => ({
    key,
    label,
    width: 1,
    defaultValue,
    headerState: (value) => [normalize(value as V) === defaultValue() ? 0 : 1, undefined],
    read: (buf, state) => {
        if (state === 0) return [defaultValue(), undefined];
        const [raw, err] = rtc.getU8(buf);
        if (err !== null) return [defaultValue(), err];
        const item = raw as any as V;
        if (!isValue(item)) return [defaultValue(), new Error(`非法枚举值: ${item}`)];
        return [item, undefined];
    },
    write: (buf, state, value) => state === 0 ? undefined : rtc.errU(rtc.setU8(buf, normalize(value as V) as any)),
    validate: (value) => isAssignable(value as V) ? undefined : new Error(`非法枚举值: ${value as any}`),
    isZero: (value) => normalize(value as V) === defaultValue(),
    eq: (left, right) => normalize(left as V) === normalize(right as V),
});

export const structField = <T, V>(
    key: keyof T & string,
    label: string,
    isZeroValue: (value: V | null | undefined) => boolean,
    readValue: (buf: Buffer) => [V, Err],
    setValue: (buf: Buffer, value: V) => Err,
    validateValue: (value: V | null | undefined) => Err,
    eqValue: (left: V | null | undefined, right: V | null | undefined) => boolean,
): StructField<T> => ({
    key,
    label,
    width: 1,
    defaultValue: () => null,
    headerState: (value) => [isZeroValue(value as any) ? 0 : 1, undefined],
    read: (buf, state) => state === 0 ? [null, undefined] : readValue(buf),
    write: (buf, state, value) => state === 0 ? undefined : setValue(buf, value as V),
    validate: (value) => validateValue(value as any),
    isZero: (value) => isZeroValue(value as any),
    eq: (left, right) => eqValue(left as any, right as any),
});

const _arrayState = (value: any): [number, Err] => {
    if (!rtc.isArrayValue(value)) return [0, new Error("value is not array")];
    const [state, err] = rtc.listCountState(value.length);
    return [state, rtc.errU(err)];
};

export const boolListField = <T>(key: keyof T & string, label: string): StructField<T> => ({
    key,
    label,
    width: 2,
    defaultValue: () => [],
    headerState: _arrayState,
    read: (buf, state) => rtc.resultU(...rtc.getBoolList(buf, state)),
    write: (buf, state, value) => rtc.errU(rtc.setBoolList(buf, state, value as boolean[])),
    validate: (value) => {
        if (!rtc.isArrayValue(value)) return new Error("value is not array");
        const [, err] = rtc.listCountState(value.length);
        return rtc.errU(err);
    },
    isZero: (value) => rtc.isArrayValue(value) && value.length === 0,
    eq: (left, right) => rtc.isArrayValue(left) && rtc.isArrayValue(right) && rtc.eqList(left as boolean[], right as boolean[], rtc.eqBool),
});

export const zeroListField = <T, V>(
    key: keyof T & string,
    label: string,
    zeroValue: V,
    getItem: (buf: Buffer) => [V, Error | null],
    setItem: (buf: Buffer, value: V) => Error | null,
    eq: (left: V, right: V) => boolean,
): StructField<T> => ({
    key,
    label,
    width: 2,
    defaultValue: () => [],
    headerState: _arrayState,
    read: (buf, state) => rtc.getZeroList(buf, state, zeroValue, getItem),
    write: (buf, state, value) => rtc.setZeroList(buf, state, value as V[], zeroValue, setItem),
    validate: (value) => {
        if (!rtc.isArrayValue(value)) return new Error("value is not array");
        const [, err] = rtc.listCountState(value.length);
        return rtc.errU(err);
    },
    isZero: (value) => rtc.isArrayValue(value) && value.length === 0,
    eq: (left, right) => rtc.isArrayValue(left) && rtc.isArrayValue(right) && rtc.eqList(left as V[], right as V[], eq),
});

export const defaultListField = <T, V>(
    key: keyof T & string,
    label: string,
    defaultItem: () => V,
    getList: (buf: Buffer, state: number) => [V[], Err],
    setList: (buf: Buffer, state: number, value: V[]) => Err,
    validateItem: (value: V) => Err,
    isDefaultItem: (value: V) => boolean,
    eqItem: (left: V, right: V) => boolean,
): StructField<T> => ({
    key,
    label,
    width: 2,
    defaultValue: () => [],
    headerState: _arrayState,
    read: (buf, state) => getList(buf, state),
    write: (buf, state, value) => setList(buf, state, value as V[]),
    validate: (value) => {
        if (!rtc.isArrayValue(value)) return new Error("value is not array");
        const [, err] = rtc.listCountState(value.length);
        if (err !== null) return err;
        for (let i = 0; i < value.length; i++) {
            const itemErr = validateItem(value[i] as V);
            if (itemErr !== undefined) return new Error(`[${i}]: ${itemErr.message}`);
        }
        return undefined;
    },
    isZero: (value) => rtc.isArrayValue(value) && value.length === 0,
    eq: (left, right) => rtc.isArrayValue(left) && rtc.isArrayValue(right) && rtc.eqList(left as V[], right as V[], eqItem),
});

export const textListField = <T>(key: keyof T & string, label: string): StructField<T> => ({
    key,
    label,
    width: 2,
    defaultValue: () => [],
    headerState: _arrayState,
    read: (buf, state) => rtc.resultU(...rtc.getTextList(buf, state)),
    write: (buf, state, value) => rtc.errU(rtc.setTextList(buf, state, value as string[])),
    validate: (value) => {
        if (!rtc.isArrayValue(value)) return new Error("value is not array");
        const [, err] = rtc.listCountState(value.length);
        if (err !== null) return err;
        for (let i = 0; i < value.length; i++) {
            const [, itemErr] = rtc.textState(value[i] as any);
            if (itemErr !== null) return new Error(`[${i}]: ${itemErr.message}`);
        }
        return undefined;
    },
    isZero: (value) => rtc.isArrayValue(value) && value.length === 0,
    eq: (left, right) => rtc.isArrayValue(left) && rtc.isArrayValue(right) && rtc.eqList(left as string[], right as string[], rtc.eqText),
});

export const binListField = <T>(key: keyof T & string, label: string): StructField<T> => ({
    key,
    label,
    width: 2,
    defaultValue: () => [],
    headerState: _arrayState,
    read: (buf, state) => rtc.resultU(...rtc.getBinList(buf, state)),
    write: (buf, state, value) => rtc.errU(rtc.setBinList(buf, state, value as Uint8Array[])),
    validate: (value) => {
        if (!rtc.isArrayValue(value)) return new Error("value is not array");
        const [, err] = rtc.listCountState(value.length);
        if (err !== null) return err;
        for (let i = 0; i < value.length; i++) {
            if (!rtc.isBinValue(value[i])) return new Error(`[${i}]: value is not Uint8Array`);
            const [, itemErr] = rtc.binState((value[i] as Uint8Array).byteLength);
            if (itemErr !== null) return new Error(`[${i}]: ${itemErr.message}`);
        }
        return undefined;
    },
    isZero: (value) => rtc.isArrayValue(value) && value.length === 0,
    eq: (left, right) => rtc.isArrayValue(left) && rtc.isArrayValue(right) && rtc.eqList(left as Uint8Array[], right as Uint8Array[], rtc.eqBin),
});
