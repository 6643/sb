import { afterEach, describe, expect, test } from "bun:test";

import { RpcClient, RpcErrCode } from "./rpc";

const originalFetch = globalThis.fetch;

afterEach(() => {
    globalThis.fetch = originalFetch;
});

describe("rpc fetch", () => {
    test("rejects oversized chunked response without arrayBuffer fallback", async () => {
        let arrayBufferCalled = false;
        globalThis.fetch = (async () => ({
            ok: true,
            status: 200,
            headers: new Headers(),
            body: new ReadableStream<Uint8Array>({
                start(controller) {
                    controller.enqueue(new Uint8Array([1, 2, 3]));
                    controller.enqueue(new Uint8Array([4, 5, 6]));
                    controller.close();
                },
            }),
            arrayBuffer: async () => {
                arrayBufferCalled = true;
                return new Uint8Array([1, 2, 3, 4, 5, 6]).buffer;
            },
        }) as any);

        const client = new RpcClient({ host: "http://example.com", maxRespBytes: 4 }) as any;
        const [bytes, status] = await client._fetch("get_count", new Uint8Array(0));

        expect(status).toBe(RpcErrCode.RespErr);
        expect(bytes).toBeNull();
        expect(arrayBufferCalled).toBeFalse();
    });

    test("preserves unknown http status", async () => {
        globalThis.fetch = (async () => new Response(null, { status: 429 })) as any;

        const client = new RpcClient({ host: "http://example.com" }) as any;
        const [bytes, status] = await client._fetch("get_count", new Uint8Array(0));

        expect(bytes).toBeNull();
        expect(status).toBe(429);
    });
});
