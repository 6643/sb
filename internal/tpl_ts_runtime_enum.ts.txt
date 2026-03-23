import * as rtc from "./runtime_core"

type Buffer = rtc.Buffer
type Err = rtc.Err
type Serializable = rtc.Serializable
type Deserializable = rtc.Deserializable

export interface EnumMeta<T extends number> {
    defaultValue: T;
    values: readonly T[];
    zeroIsMember: boolean;
}

export const defineEnum = <T extends number>(defaultValue: T, values: readonly T[]): EnumMeta<T> => {
    let zeroIsMember = false;
    for (let i = 0; i < values.length; i++) {
        if ((values[i] as any) === 0) {
            zeroIsMember = true;
            break;
        }
    }
    return { defaultValue, values, zeroIsMember };
};

export const isEnum = <T extends number>(meta: EnumMeta<T>, value: T): boolean => {
    for (let i = 0; i < meta.values.length; i++) {
        if (meta.values[i] === value) return true;
    }
    return false;
};

export const normalizeEnum = <T extends number>(meta: EnumMeta<T>, value: T): T => {
    if (isEnum(meta, value)) return value;
    if ((value as any) === 0) return meta.defaultValue;
    return value;
};

export const isDefaultEnum = <T extends number>(meta: EnumMeta<T>, value: T): boolean => normalizeEnum(meta, value) === meta.defaultValue;
export const isAssignableEnum = <T extends number>(meta: EnumMeta<T>, value: T): boolean => isEnum(meta, value) || (((value as any) === 0) && !meta.zeroIsMember);
export const eqEnumValue = <T extends number>(meta: EnumMeta<T>, left: T, right: T): boolean => normalizeEnum(meta, left) === normalizeEnum(meta, right);
export const eqEnumList = <T extends number>(meta: EnumMeta<T>, left: T[], right: T[]): boolean => rtc.eqList(left, right, (a, b) => eqEnumValue(meta, a, b));

export const getEnumList = <T extends number>(meta: EnumMeta<T>, buf: Buffer, state: number): [T[], Err] => {
    return rtc.getDefaultList<T>(
        buf,
        state,
        () => meta.defaultValue,
        (buf) => {
            const [value, err] = rtc.getU8(buf);
            if (err !== null) return [meta.defaultValue, rtc.errU(err)];
            const item = value as T;
            if (!isEnum(meta, item)) return [meta.defaultValue, new Error(`非法枚举值: ${item}`)];
            return [item, undefined];
        },
    );
};

export const setEnumList = <T extends number>(meta: EnumMeta<T>, buf: Buffer, state: number, values: T[]): Err => {
    return rtc.setDefaultList<T>(
        buf,
        state,
        values,
        (item) => isDefaultEnum(meta, item),
        (buf, item) => {
            if (!isEnum(meta, item)) return new Error(`非法枚举值: ${item}`);
            return rtc.errU(rtc.setU8(buf, item as any));
        },
    );
};
