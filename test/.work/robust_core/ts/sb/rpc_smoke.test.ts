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
    test("method echo exists", () => {
        const client = new RpcClient({ host: baseUrl });
        expect(typeof (client as any).echo).toBe("function");
    });
    test("method echoTail exists", () => {
        const client = new RpcClient({ host: baseUrl });
        expect(typeof (client as any).echoTail).toBe("function");
    });
    test("method getColor exists", () => {
        const client = new RpcClient({ host: baseUrl });
        expect(typeof (client as any).getColor).toBe("function");
    });
    test("method getBadColor exists", () => {
        const client = new RpcClient({ host: baseUrl });
        expect(typeof (client as any).getBadColor).toBe("function");
    });
    test("method ping exists", () => {
        const client = new RpcClient({ host: baseUrl });
        expect(typeof (client as any).ping).toBe("function");
    });
    test("method pingJunk exists", () => {
        const client = new RpcClient({ host: baseUrl });
        expect(typeof (client as any).pingJunk).toBe("function");
    });
    test("method pick exists", () => {
        const client = new RpcClient({ host: baseUrl });
        expect(typeof (client as any).pick).toBe("function");
    });
    test("method pickBad exists", () => {
        const client = new RpcClient({ host: baseUrl });
        expect(typeof (client as any).pickBad).toBe("function");
    });
});
