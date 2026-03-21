import * as rt from "./type"

export interface SimOrder2 extends rt.Serializable, rt.Deserializable {
    // SIM卡ID
    id: number;
    // 办理人姓名
    name: string;
    // 联系电话
    phone: string;
    // 身份证号
    idNo: string;
    // 所在城市
    cityCode: number;
    // 详细地址
    address: string;
    // 新手机号码
    newPhone: string;
}

const simOrder2HeaderWidths = [1, 2, 2, 2, 1, 2, 2] as const;

const simOrder2Meta = rt.defineStruct<SimOrder2>({
    name: "SimOrder2",
    headerWidths: simOrder2HeaderWidths,
    create: () => ({
        id: 0,
        name: "",
        phone: "",
        idNo: "",
        cityCode: 0,
        address: "",
        newPhone: "",
    }) as any as SimOrder2,
    fields: [

        rt.scalarField<SimOrder2, number>("id", "Id", 0, rt.getU32, rt.setU32, rt.eqU32),
        rt.textField<SimOrder2>("name", "Name"),
        rt.textField<SimOrder2>("phone", "Phone"),
        rt.textField<SimOrder2>("idNo", "IdNo"),
        rt.scalarField<SimOrder2, number>("cityCode", "CityCode", 0, rt.getU32, rt.setU32, rt.eqU32),
        rt.textField<SimOrder2>("address", "Address"),
        rt.textField<SimOrder2>("newPhone", "NewPhone"),
    ],
});

export const newSimOrder2 = (): SimOrder2 => rt.newStruct(simOrder2Meta, getSimOrder2, setSimOrder2);
export const isZeroSimOrder2 = (s: SimOrder2 | null | undefined): boolean => rt.isZeroStruct(simOrder2Meta, s);
export const validateSimOrder2 = (s: SimOrder2 | null | undefined): rt.Err => rt.validateStruct(simOrder2Meta, s);
export const getSimOrder2 = (buf: rt.Buffer): [SimOrder2, rt.Err] => rt.getStruct(simOrder2Meta, buf);
export const setSimOrder2 = (buf: rt.Buffer, s: SimOrder2): rt.Err => rt.setStruct(simOrder2Meta, buf, s);
export const readSimOrder2 = (buf: rt.Buffer): [SimOrder2, rt.Err] => getSimOrder2(buf);
export const eqSimOrder2 = (a: SimOrder2 | null | undefined, b: SimOrder2 | null | undefined): boolean => rt.eqStruct(simOrder2Meta, a, b);

export const getSimOrder2ListBody = (buf: rt.Buffer, state: number): [SimOrder2[], rt.Err] => {
    const [list, err] = rt.getDefaultList<SimOrder2>(
        buf,
        state,
        () => newSimOrder2(),
        (buf) => readSimOrder2(buf),
    );
    return [list, err];
};

export const setSimOrder2ListBody = (buf: rt.Buffer, state: number, v: SimOrder2[]): rt.Err => {
    return rt.setDefaultList<SimOrder2>(
        buf,
        state,
        v,
        (item) => isZeroSimOrder2(item),
        (buf, item) => setSimOrder2(buf, item),
    );
}
