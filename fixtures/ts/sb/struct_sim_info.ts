import * as rt from "./runtime_core"
import * as rm from "./runtime_meta"

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

const simInfoMeta = rm.defineStruct<SimInfo>({
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

        rm.scalarField<SimInfo, number>("id", "Id", 0, rt.getU32, rt.setU32, rt.eqU32),
        rm.textField<SimInfo>("title", "Title"),
        rm.textField<SimInfo>("content", "Content"),
        rm.boolField<SimInfo>("a", "A"),
        rm.boolField<SimInfo>("b", "B"),
        rm.boolField<SimInfo>("c", "C"),
        rm.boolField<SimInfo>("d", "D"),
        rm.binField<SimInfo>("zip", "Zip"),
    ],
});

export const newSimInfo = (): SimInfo => rm.newStruct(simInfoMeta, getSimInfo, setSimInfo);
export const isZeroSimInfo = (s: SimInfo | null | undefined): boolean => rm.isZeroStruct(simInfoMeta, s);
export const validateSimInfo = (s: SimInfo | null | undefined): rt.Err => rm.validateStruct(simInfoMeta, s);
export const getSimInfo = (buf: rt.Buffer): [SimInfo, rt.Err] => rm.getStruct(simInfoMeta, buf);
export const setSimInfo = (buf: rt.Buffer, s: SimInfo): rt.Err => rm.setStruct(simInfoMeta, buf, s);
export const readSimInfo = (buf: rt.Buffer): [SimInfo, rt.Err] => getSimInfo(buf);
export const eqSimInfo = (a: SimInfo | null | undefined, b: SimInfo | null | undefined): boolean => rm.eqStruct(simInfoMeta, a, b);

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
