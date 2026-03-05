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
{{- range .Apis}}
    test("method {{.Name | CamelCase}} exists", () => {
        const client = new RpcClient({ host: baseUrl });
        expect(typeof (client as any).{{.Name | CamelCase}}).toBe("function");
    });
{{- end}}
});
