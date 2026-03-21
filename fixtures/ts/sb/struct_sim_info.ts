import * as rt from "./type"

export interface SimInfo extends rt.Serializable, rt.Deserializable {
    id: number;
    title: string;
    content: string;
    a: boolean;
    b: boolean;
    c: boolean;
    d: boolean;
    zip: Uint8Array;
}

const simInfoHeaderWidths = [1, 2, 2, 1, 1, 1, 1, 2] as const;

const simInfoMeta = rt.defineStruct<SimInfo>({
    name: "SimInfo",
    headerWidths: simInfoHeaderWidths,
    create: () => ({
        id: 0,
        title: "",
        content: "",
        a: false,
        b: false,
        c: false,
        d: false,
        zip: new Uint8Array(0),
    }) as any as SimInfo,
    fields: [

        rt.scalarField<SimInfo, number>("id", "Id", 0, rt.getU32, rt.setU32, rt.eqU32),
        rt.textField<SimInfo>("title", "Title"),
        rt.textField<SimInfo>("content", "Content"),
        rt.boolField<SimInfo>("a", "A"),
        rt.boolField<SimInfo>("b", "B"),
        rt.boolField<SimInfo>("c", "C"),
        rt.boolField<SimInfo>("d", "D"),
        rt.binField<SimInfo>("zip", "Zip"),
    ],
});

export const newSimInfo = (): SimInfo => rt.newStruct(simInfoMeta, getSimInfo, setSimInfo);
export const isZeroSimInfo = (s: SimInfo | null | undefined): boolean => rt.isZeroStruct(simInfoMeta, s);
export const validateSimInfo = (s: SimInfo | null | undefined): rt.Err => rt.validateStruct(simInfoMeta, s);
export const getSimInfo = (buf: rt.Buffer): [SimInfo, rt.Err] => rt.getStruct(simInfoMeta, buf);
export const setSimInfo = (buf: rt.Buffer, s: SimInfo): rt.Err => rt.setStruct(simInfoMeta, buf, s);
export const readSimInfo = (buf: rt.Buffer): [SimInfo, rt.Err] => getSimInfo(buf);
export const eqSimInfo = (a: SimInfo | null | undefined, b: SimInfo | null | undefined): boolean => rt.eqStruct(simInfoMeta, a, b);

export const getSimInfoListBody = (buf: rt.Buffer, state: number): [SimInfo[], rt.Err] => {
    const [list, err] = rt.getDefaultList<SimInfo>(
        buf,
        state,
        () => newSimInfo(),
        (buf) => readSimInfo(buf),
    );
    return [list, err];
};

export const setSimInfoListBody = (buf: rt.Buffer, state: number, v: SimInfo[]): rt.Err => {
    return rt.setDefaultList<SimInfo>(
        buf,
        state,
        v,
        (item) => isZeroSimInfo(item),
        (buf, item) => setSimInfo(buf, item),
    );
}
