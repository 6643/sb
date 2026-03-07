import * as _ from "./_"

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

export interface RpcConfig {
    host: string;
    headers?: Record<string, string>;
    timeout?: number;
    retries?: number;
    maxRespBytes?: number;
}

export class RpcClient {
    private headers: Record<string, string> = {};
    private timeout: number;
    private retries: number;
    private maxRespBytes: number;

    constructor(private config: RpcConfig) {
        this.config = { ...config, host: config.host.replace(/\/+$/, "") };
        if (config.headers) this.headers = { ...config.headers };
        const cfgTimeout = config.timeout;
        this.timeout = cfgTimeout !== undefined && Number.isFinite(cfgTimeout) && cfgTimeout >= 0 ? Math.min(Math.floor(cfgTimeout), maxTimeoutMs) : 5000;
        const cfgRetries = config.retries;
        this.retries = cfgRetries !== undefined && Number.isFinite(cfgRetries) && cfgRetries >= 0 ? Math.floor(cfgRetries) : 3;
        const cfgMaxRespBytes = config.maxRespBytes;
        this.maxRespBytes = cfgMaxRespBytes !== undefined && Number.isFinite(cfgMaxRespBytes) && cfgMaxRespBytes > 0 ? Math.min(Math.floor(cfgMaxRespBytes), maxSafeRespBytes) : defaultMaxRespBytes;
    }

    public setHeader = (key: string, value: string): void => { this.headers[key] = value; };
    public getHeader = (key: string): string | undefined => this.headers[key];
    public removeHeader = (key: string): void => { delete this.headers[key]; };

    public setAuthorization = (token: string): void => { this.setHeader("Authorization", `Bearer ${token}`); };
    public getAuthorization = (): string | undefined => this.getHeader("Authorization");
    public removeAuthorization = (): void => { this.removeHeader("Authorization"); };
    public isAuthorized = (): boolean => !!this.getAuthorization();

    private async _fetch(path: string, body: Uint8Array): Promise<[Uint8Array | null, RpcErrCode]> {
        let lastStatus = RpcErrCode.NoConn;
        for (let i = 0; i <= this.retries; i++) {
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
                    const contentLength = res.headers.get("content-length");
                    if (contentLength !== null) {
                        const size = Number(contentLength);
                        if (Number.isFinite(size) && size > this.maxRespBytes) return [null, RpcErrCode.RespErr];
                    }
                    const bytes = new Uint8Array(await res.arrayBuffer());
                    if (bytes.byteLength > this.maxRespBytes) return [null, RpcErrCode.RespErr];
                    return [bytes, RpcErrCode.Ok];
                }
                lastStatus = res.status as RpcErrCode;
                if (res.status === 408 && i < this.retries) continue;
                return [null, res.status as RpcErrCode];
            } catch (e: any) {
                lastStatus = e && e.name === "AbortError" ? RpcErrCode.Timeout : RpcErrCode.NoConn;
                if (i < this.retries) continue;
            } finally {
                if (timeoutId !== null) clearTimeout(timeoutId);
            }
        }
        return [null, lastStatus];
    }

    // 获取用户的id
    public userGetAbc = async (): Promise<[_.OrderStatus, RpcErrCode]> => {
        const buf = new _.Buffer();
        const [bytes, status] = await this._fetch("user.get_abc", buf.bytes);
        if (status !== RpcErrCode.Ok || bytes === null) return [_.DefaultOrderStatus(), status];
        const respBuf = new _.Buffer(bytes);
        let result = _.DefaultOrderStatus() as any;
        {
            const [value, err] = _.getU8(respBuf);
            if (err !== null) return [_.DefaultOrderStatus(), RpcErrCode.RespErr];
            const item = value as _.OrderStatus;
            if (!_.IsOrderStatus(item)) return [_.DefaultOrderStatus(), RpcErrCode.RespErr];
            result = item as any;
        }
        if (respBuf.len !== 0) return [_.DefaultOrderStatus(), RpcErrCode.RespErr];
        return [result as any, RpcErrCode.Ok];
    };

    // 获取abcd
    public userGetAbcd = async (page: number, size: number): Promise<[_.OrderStatus, RpcErrCode]> => {
        const buf = new _.Buffer();
        {
            const encodeErr = _.setU8(buf, page as any);
            if (encodeErr !== null) return [_.DefaultOrderStatus(), RpcErrCode.ReqErr];
        }
        {
            const encodeErr = _.setU8(buf, size as any);
            if (encodeErr !== null) return [_.DefaultOrderStatus(), RpcErrCode.ReqErr];
        }
        const [bytes, status] = await this._fetch("user.get_abcd", buf.bytes);
        if (status !== RpcErrCode.Ok || bytes === null) return [_.DefaultOrderStatus(), status];
        const respBuf = new _.Buffer(bytes);
        let result = _.DefaultOrderStatus() as any;
        {
            const [value, err] = _.getU8(respBuf);
            if (err !== null) return [_.DefaultOrderStatus(), RpcErrCode.RespErr];
            const item = value as _.OrderStatus;
            if (!_.IsOrderStatus(item)) return [_.DefaultOrderStatus(), RpcErrCode.RespErr];
            result = item as any;
        }
        if (respBuf.len !== 0) return [_.DefaultOrderStatus(), RpcErrCode.RespErr];
        return [result as any, RpcErrCode.Ok];
    };

    // 设置sim信息
    // 无返回值
    public userSetSimInfo = async (info: _.SimInfo): Promise<RpcErrCode> => {
        const buf = new _.Buffer();
        {
            const encodeErr = _.setSimInfo(buf, info);
            if (encodeErr !== undefined) return RpcErrCode.ReqErr;
        }
        const [bytes, status] = await this._fetch("user.set_sim_info", buf.bytes);
        if (status !== RpcErrCode.Ok || bytes === null) return status;
        if (bytes.byteLength !== 0) return RpcErrCode.RespErr;
        return RpcErrCode.Ok;
    };

    // 获取数量
    public getCount = async (page: number): Promise<[number, RpcErrCode]> => {
        const buf = new _.Buffer();
        {
            const encodeErr = _.setU8(buf, page as any);
            if (encodeErr !== null) return [0, RpcErrCode.ReqErr];
        }
        const [bytes, status] = await this._fetch("get_count", buf.bytes);
        if (status !== RpcErrCode.Ok || bytes === null) return [0, status];
        const respBuf = new _.Buffer(bytes);
        let result = 0 as any;
        {
            const [value, err] = _.getU8(respBuf);
            if (err !== null) return [0, RpcErrCode.RespErr];
            result = value as any;
        }
        if (respBuf.len !== 0) return [0, RpcErrCode.RespErr];
        return [result as any, RpcErrCode.Ok];
    };

    // 获取bin
    public getBin = async (page: number): Promise<[Uint8Array, RpcErrCode]> => {
        const buf = new _.Buffer();
        {
            const encodeErr = _.setU8(buf, page as any);
            if (encodeErr !== null) return [new Uint8Array(0), RpcErrCode.ReqErr];
        }
        const [bytes, status] = await this._fetch("get_bin", buf.bytes);
        if (status !== RpcErrCode.Ok || bytes === null) return [new Uint8Array(0), status];
        const respBuf = new _.Buffer(bytes);
        let result = new Uint8Array(0) as any;
        {
            const [state, errState] = _.getU8(respBuf);
            if (errState !== null) return [new Uint8Array(0), RpcErrCode.RespErr];
            const [value, err] = _.getBinCompact(respBuf, state);
            if (err !== null) return [new Uint8Array(0), RpcErrCode.RespErr];
            result = value as any;
        }
        if (respBuf.len !== 0) return [new Uint8Array(0), RpcErrCode.RespErr];
        return [result as any, RpcErrCode.Ok];
    };

}
