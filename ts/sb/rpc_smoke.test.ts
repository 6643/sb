import { describe, expect, test } from "bun:test";

import * as _ from "./_";
import { RpcClient, RpcErrCode } from "./rpc";

const baseUrl = process.env.SB_BASE_URL || "http://127.0.0.1:18080";

describe("rpc smoke", () => {
    test("client construction", () => {
        const client = new RpcClient({ host: baseUrl });
        expect(client).toBeDefined();
        expect(typeof RpcErrCode.Ok).toBe("number");
        expect(typeof _).toBe("object");
    });
    test("method userGetAbc exists", () => {
        const client = new RpcClient({ host: baseUrl });
        expect(typeof (client as any).userGetAbc).toBe("function");
    });
    test("method userGetAbcd exists", () => {
        const client = new RpcClient({ host: baseUrl });
        expect(typeof (client as any).userGetAbcd).toBe("function");
    });
    test("method userSetSimInfo exists", () => {
        const client = new RpcClient({ host: baseUrl });
        expect(typeof (client as any).userSetSimInfo).toBe("function");
    });
    test("method getCount exists", () => {
        const client = new RpcClient({ host: baseUrl });
        expect(typeof (client as any).getCount).toBe("function");
    });
    test("method getBin exists", () => {
        const client = new RpcClient({ host: baseUrl });
        expect(typeof (client as any).getBin).toBe("function");
    });
});
