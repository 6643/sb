import { describe, expect, test } from "bun:test";

import { Buffer } from "./type";
import { eqSimInfo, getSimInfo, newSimInfo, setSimInfo, validateSimInfo } from "./struct_sim_info";

describe("struct smoke", () => {
    test("SimInfo public api remains usable", () => {
        const value = newSimInfo();
        value.id = 7;
        value.title = "starter";
        value.content = "runtime";
        value.a = true;
        value.c = true;
        value.zip = new Uint8Array([0x01, 0x02, 0x03]);

        expect(validateSimInfo(value)).toBeUndefined();

        const buf = new Buffer();
        expect(setSimInfo(buf, value)).toBeUndefined();

        const [decoded, err] = getSimInfo(new Buffer(buf.bytes));
        expect(err).toBeUndefined();
        expect(eqSimInfo(value, decoded)).toBe(true);
    });
});
