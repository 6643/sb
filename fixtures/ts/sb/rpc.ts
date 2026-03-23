import * as rt from "./runtime_core"
import { DefaultOrderStatus, IsAssignableOrderStatus, IsOrderStatus, NormalizeOrderStatus, OrderStatus, getOrderStatusListBody, setOrderStatusListBody } from "./enum"
import { SimInfo, getSimInfoListBody, newSimInfo, readSimInfo, setSimInfo, setSimInfoListBody } from "./struct_sim_info"

const defaultMaxRespBytes = 4 * 1024 * 1024;
const maxSafeRespBytes = Number.MAX_SAFE_INTEGER;
const maxTimeoutMs = 2147483647;

export enum RpcErrCode {
    Ok = 200,
    NoConn = 0,
    Timeout = 408,
    ReqErr = 400,
    RespErr = 500,
    NotAuth = 401,
    NotExist = 404,
}

export type RpcStatus = RpcErrCode | number;

export interface RpcConfig {
    host: string;
    headers?: Record<string, string>;
    timeout?: number;
    retries?: number;
    enableRetries?: boolean;
    maxRespBytes?: number;
}

export interface RpcClient {
    setHeader: (key: string, value: string) => void;
    getHeader: (key: string) => string | undefined;
    removeHeader: (key: string) => void;
    setAuthorization: (token: string) => void;
    getAuthorization: () => string | undefined;
    removeAuthorization: () => void;
    isAuthorized: () => boolean;
    userGetAbc: () => Promise<[OrderStatus, RpcStatus]>;
    userGetAbcd: (page: number, size: number) => Promise<[OrderStatus, RpcStatus]>;
    userSetSimInfo: (info: SimInfo) => Promise<RpcStatus>;
    getCount: (page: number) => Promise<[number, RpcStatus]>;
    getBin: (page: number) => Promise<[Uint8Array, RpcStatus]>;
}

type _RpcClientState = RpcClient & {
    config: RpcConfig;
    headers: Record<string, string>;
    timeout: number;
    retries: number;
    enableRetries: boolean;
    maxRespBytes: number;
    _fetch: (path: string, body: Uint8Array) => Promise<[Uint8Array | null, RpcStatus]>;
};
type RpcClientCtor = new (config: RpcConfig) => RpcClient;

function _RpcClientCtor(this: _RpcClientState, config: RpcConfig): void {
    this.config = { ...config, host: config.host.replace(/\/+$/, "") };
    this.headers = config.headers ? { ...config.headers } : {};
    const cfgTimeout = config.timeout;
    this.timeout = cfgTimeout !== undefined && Number.isFinite(cfgTimeout) && cfgTimeout >= 0 ? Math.min(Math.floor(cfgTimeout), maxTimeoutMs) : 5000;
    const cfgRetries = config.retries;
    this.retries = cfgRetries !== undefined && Number.isFinite(cfgRetries) && cfgRetries >= 0 ? Math.floor(cfgRetries) : 3;
    this.enableRetries = config.enableRetries === true;
    const cfgMaxRespBytes = config.maxRespBytes;
    this.maxRespBytes = cfgMaxRespBytes !== undefined && Number.isFinite(cfgMaxRespBytes) && cfgMaxRespBytes > 0 ? Math.min(Math.floor(cfgMaxRespBytes), maxSafeRespBytes) : defaultMaxRespBytes;
}
const _rpcClientProto = _RpcClientCtor.prototype as _RpcClientState;

_rpcClientProto.setHeader = function(this: _RpcClientState, key: string, value: string): void { this.headers[key] = value; };
_rpcClientProto.getHeader = function(this: _RpcClientState, key: string): string | undefined { return this.headers[key]; };
_rpcClientProto.removeHeader = function(this: _RpcClientState, key: string): void { delete this.headers[key]; };

_rpcClientProto.setAuthorization = function(this: _RpcClientState, token: string): void { this.setHeader("Authorization", `Bearer ${token}`); };
_rpcClientProto.getAuthorization = function(this: _RpcClientState): string | undefined { return this.getHeader("Authorization"); };
_rpcClientProto.removeAuthorization = function(this: _RpcClientState): void { this.removeHeader("Authorization"); };
_rpcClientProto.isAuthorized = function(this: _RpcClientState): boolean { return !!this.getAuthorization(); };

async function readResponseBytes(res: Response, maxRespBytes: number): Promise<[Uint8Array | null, RpcErrCode | null]> {
    const contentLength = res.headers.get("content-length");
    if (contentLength !== null) {
        const size = Number(contentLength);
        if (Number.isFinite(size) && size > maxRespBytes) return [null, RpcErrCode.RespErr];
    }

    if (res.body !== null) {
        const reader = res.body.getReader();
        const chunks: Uint8Array[] = [];
        let total = 0;
        while (true) {
            const { done, value } = await reader.read();
            if (done) break;
            if (value === undefined) continue;
            total += value.byteLength;
            if (total > maxRespBytes) {
                try {
                    await reader.cancel();
                } catch {
                }
                return [null, RpcErrCode.RespErr];
            }
            chunks.push(value);
        }

        const bytes = new Uint8Array(total);
        let offset = 0;
        for (const chunk of chunks) {
            bytes.set(chunk, offset);
            offset += chunk.byteLength;
        }
        return [bytes, null];
    }

    const bytes = new Uint8Array(await res.arrayBuffer());
    if (bytes.byteLength > maxRespBytes) return [null, RpcErrCode.RespErr];
    return [bytes, null];
}

_rpcClientProto._fetch = async function(this: _RpcClientState, path: string, body: Uint8Array): Promise<[Uint8Array | null, RpcStatus]> {
    const maxRetries = this.enableRetries ? this.retries : 0;
    let lastStatus = RpcErrCode.NoConn;
    for (let i = 0; i <= maxRetries; i++) {
        if (i > 0) await new Promise((res) => setTimeout(res, i * 1000));
        const controller = new AbortController();
        let timeoutId: ReturnType<typeof setTimeout> | null = null;
        if (this.timeout > 0) timeoutId = setTimeout(() => controller.abort(), this.timeout);
        try {
            const res = await fetch(`${this.config.host}/${path}`, {
                method: "POST",
                headers: { "Content-Type": "application/octet-stream", ...this.headers },
                body: body as any,
                signal: controller.signal,
            });
            if (res.ok) {
                const [bytes, readErr] = await readResponseBytes(res, this.maxRespBytes);
                if (readErr !== null) return [null, readErr];
                return [bytes, RpcErrCode.Ok];
            }
            lastStatus = res.status;
            if (res.status === 408 && i < maxRetries) continue;
            return [null, res.status];
        } catch (e: any) {
            lastStatus = e && e.name === "AbortError" ? RpcErrCode.Timeout : RpcErrCode.NoConn;
            if (i < maxRetries) continue;
        } finally {
            if (timeoutId !== null) clearTimeout(timeoutId);
        }
    }
    return [null, lastStatus];
};

// 获取用户的id
_rpcClientProto.userGetAbc = async function(this: _RpcClientState, ): Promise<[OrderStatus, RpcStatus]> {
    const buf = new rt.Buffer();
    const [bytes, status] = await this._fetch("user.get_abc", buf.bytes);
    if (status !== RpcErrCode.Ok || bytes === null) return [DefaultOrderStatus(), status];
    const respBuf = new rt.Buffer(bytes);
    let result = DefaultOrderStatus() as any;
        {
            const [value, err] = rt.getU8(respBuf);
            if (err !== null) return [DefaultOrderStatus(), RpcErrCode.RespErr];
            const item = value as OrderStatus;
            if (!IsOrderStatus(item)) return [DefaultOrderStatus(), RpcErrCode.RespErr];
            result = item as any;
        }
    if (respBuf.len !== 0) return [DefaultOrderStatus(), RpcErrCode.RespErr];
    return [result as any, RpcErrCode.Ok];
};

// 获取abcd
_rpcClientProto.userGetAbcd = async function(this: _RpcClientState, page: number, size: number): Promise<[OrderStatus, RpcStatus]> {
    const buf = new rt.Buffer();
        {
            const encodeErr = rt.setU8(buf, page as any);
            if (encodeErr !== null) return [DefaultOrderStatus(), RpcErrCode.ReqErr];
        }
        {
            const encodeErr = rt.setU8(buf, size as any);
            if (encodeErr !== null) return [DefaultOrderStatus(), RpcErrCode.ReqErr];
        }
    const [bytes, status] = await this._fetch("user.get_abcd", buf.bytes);
    if (status !== RpcErrCode.Ok || bytes === null) return [DefaultOrderStatus(), status];
    const respBuf = new rt.Buffer(bytes);
    let result = DefaultOrderStatus() as any;
        {
            const [value, err] = rt.getU8(respBuf);
            if (err !== null) return [DefaultOrderStatus(), RpcErrCode.RespErr];
            const item = value as OrderStatus;
            if (!IsOrderStatus(item)) return [DefaultOrderStatus(), RpcErrCode.RespErr];
            result = item as any;
        }
    if (respBuf.len !== 0) return [DefaultOrderStatus(), RpcErrCode.RespErr];
    return [result as any, RpcErrCode.Ok];
};

// 设置sim信息
// 无返回值
_rpcClientProto.userSetSimInfo = async function(this: _RpcClientState, info: SimInfo): Promise<RpcStatus> {
    const buf = new rt.Buffer();
        {
            const encodeErr = setSimInfo(buf, info);
            if (encodeErr !== undefined) return RpcErrCode.ReqErr;
        }
    const [bytes, status] = await this._fetch("user.set_sim_info", buf.bytes);
    if (status !== RpcErrCode.Ok || bytes === null) return status;
    if (bytes.byteLength !== 0) return RpcErrCode.RespErr;
    return RpcErrCode.Ok;
};

// 获取数量
_rpcClientProto.getCount = async function(this: _RpcClientState, page: number): Promise<[number, RpcStatus]> {
    const buf = new rt.Buffer();
        {
            const encodeErr = rt.setU8(buf, page as any);
            if (encodeErr !== null) return [0, RpcErrCode.ReqErr];
        }
    const [bytes, status] = await this._fetch("get_count", buf.bytes);
    if (status !== RpcErrCode.Ok || bytes === null) return [0, status];
    const respBuf = new rt.Buffer(bytes);
    let result = 0 as any;
        {
            const [value, err] = rt.getU8(respBuf);
            if (err !== null) return [0, RpcErrCode.RespErr];
            result = value as any;
        }
    if (respBuf.len !== 0) return [0, RpcErrCode.RespErr];
    return [result as any, RpcErrCode.Ok];
};

// 获取bin
_rpcClientProto.getBin = async function(this: _RpcClientState, page: number): Promise<[Uint8Array, RpcStatus]> {
    const buf = new rt.Buffer();
        {
            const encodeErr = rt.setU8(buf, page as any);
            if (encodeErr !== null) return [new Uint8Array(0), RpcErrCode.ReqErr];
        }
    const [bytes, status] = await this._fetch("get_bin", buf.bytes);
    if (status !== RpcErrCode.Ok || bytes === null) return [new Uint8Array(0), status];
    const respBuf = new rt.Buffer(bytes);
    let result = new Uint8Array(0) as any;
        {
            const [state, errState] = rt.getU8(respBuf);
            if (errState !== null) return [new Uint8Array(0), RpcErrCode.RespErr];
            const [value, err] = rt.getBin(respBuf, state);
            if (err !== null) return [new Uint8Array(0), RpcErrCode.RespErr];
            result = value as any;
        }
    if (respBuf.len !== 0) return [new Uint8Array(0), RpcErrCode.RespErr];
    return [result as any, RpcErrCode.Ok];
};

export const RpcClient = _RpcClientCtor as unknown as RpcClientCtor;
