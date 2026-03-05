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
        if (cfgTimeout !== undefined && Number.isFinite(cfgTimeout) && cfgTimeout >= 0) {
            this.timeout = Math.min(Math.floor(cfgTimeout), maxTimeoutMs);
        } else {
            this.timeout = 5000;
        }
        const cfgRetries = config.retries;
        if (cfgRetries !== undefined && Number.isFinite(cfgRetries) && cfgRetries >= 0) {
            this.retries = Math.floor(cfgRetries);
        } else {
            this.retries = 3;
        }
        const cfgMaxRespBytes = config.maxRespBytes;
        if (cfgMaxRespBytes !== undefined && Number.isFinite(cfgMaxRespBytes) && cfgMaxRespBytes > 0) {
            this.maxRespBytes = Math.min(Math.floor(cfgMaxRespBytes), maxSafeRespBytes);
        } else {
            this.maxRespBytes = defaultMaxRespBytes;
        }
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
            if (i > 0) await new Promise(res => setTimeout(res, i * 1000));
            const controller = new AbortController();
            let timeoutId: ReturnType<typeof setTimeout> | null = null;
            if (this.timeout > 0) timeoutId = setTimeout(() => controller.abort(), this.timeout);
            try {
                const res = await fetch(`${this.config.host}/${path}`, {
                    method: "POST",
                    headers: { "Content-Type": "application/octet-stream", ...this.headers },
                    body: body as any,
                    signal: controller.signal
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
                if (e.name === "AbortError") {
                    lastStatus = RpcErrCode.Timeout;
                } else {
                    lastStatus = RpcErrCode.NoConn;
                }
                if (i < this.retries) continue;
            } finally {
                if (timeoutId !== null) clearTimeout(timeoutId);
            }
        }
        return [null, lastStatus];
    }

    /**  */
    public echo = async (env: _.Envelope): Promise<[_.Envelope, RpcErrCode]> => {
        const buf = new _.Buffer();
        if (_.setAll(buf, (buf) => _.setEnvelope(buf, env as any)) !== null) return [_.newEnvelope(), RpcErrCode.ReqErr];

        const [bytes, status] = await this._fetch("echo", buf.bytes);
        if (status !== RpcErrCode.Ok || bytes === null) return [_.newEnvelope(), status];

        const respBuf = new _.Buffer(bytes);const [result, err] = _.getEnvelope(respBuf);
        if (err !== null) return [_.newEnvelope(), RpcErrCode.RespErr];
        if (respBuf.len !== 0) return [_.newEnvelope(), RpcErrCode.RespErr];
        return [result as any, RpcErrCode.Ok];
    };
    /**  */
    public echoTail = async (env: _.Envelope): Promise<[_.Envelope, RpcErrCode]> => {
        const buf = new _.Buffer();
        if (_.setAll(buf, (buf) => _.setEnvelope(buf, env as any)) !== null) return [_.newEnvelope(), RpcErrCode.ReqErr];

        const [bytes, status] = await this._fetch("echo_tail", buf.bytes);
        if (status !== RpcErrCode.Ok || bytes === null) return [_.newEnvelope(), status];

        const respBuf = new _.Buffer(bytes);const [result, err] = _.getEnvelope(respBuf);
        if (err !== null) return [_.newEnvelope(), RpcErrCode.RespErr];
        if (respBuf.len !== 0) return [_.newEnvelope(), RpcErrCode.RespErr];
        return [result as any, RpcErrCode.Ok];
    };
    /**  */
    public getColor = async (): Promise<[_.Color, RpcErrCode]> => {
        const buf = new _.Buffer();

        const [bytes, status] = await this._fetch("get_color", buf.bytes);
        if (status !== RpcErrCode.Ok || bytes === null) return [0 as _.Color, status];

        const respBuf = new _.Buffer(bytes);const [result, err] = _.getU8(respBuf);
        if (err === null && !_.IsColor(result as any)) return [0 as _.Color, RpcErrCode.RespErr];
        if (err !== null) return [0 as _.Color, RpcErrCode.RespErr];
        if (respBuf.len !== 0) return [0 as _.Color, RpcErrCode.RespErr];
        return [result as any, RpcErrCode.Ok];
    };
    /**  */
    public getBadColor = async (): Promise<[_.Color, RpcErrCode]> => {
        const buf = new _.Buffer();

        const [bytes, status] = await this._fetch("get_bad_color", buf.bytes);
        if (status !== RpcErrCode.Ok || bytes === null) return [0 as _.Color, status];

        const respBuf = new _.Buffer(bytes);const [result, err] = _.getU8(respBuf);
        if (err === null && !_.IsColor(result as any)) return [0 as _.Color, RpcErrCode.RespErr];
        if (err !== null) return [0 as _.Color, RpcErrCode.RespErr];
        if (respBuf.len !== 0) return [0 as _.Color, RpcErrCode.RespErr];
        return [result as any, RpcErrCode.Ok];
    };
    /**  */
    public ping = async (): Promise<RpcErrCode> => {
        const buf = new _.Buffer();

        const [bytes, status] = await this._fetch("ping", buf.bytes);
        if (status !== RpcErrCode.Ok || bytes === null) return status;

        if (bytes.byteLength !== 0) return RpcErrCode.RespErr;
        return RpcErrCode.Ok;
    };
    /**  */
    public pingJunk = async (): Promise<RpcErrCode> => {
        const buf = new _.Buffer();

        const [bytes, status] = await this._fetch("ping_junk", buf.bytes);
        if (status !== RpcErrCode.Ok || bytes === null) return status;

        if (bytes.byteLength !== 0) return RpcErrCode.RespErr;
        return RpcErrCode.Ok;
    };
    /**  */
    public pick = async (color: _.Color): Promise<[_.Color, RpcErrCode]> => {
        const buf = new _.Buffer();
        if (!_.IsColor(color as any)) return [0 as _.Color, RpcErrCode.ReqErr];
        if (_.setAll(buf, (buf) => _.setU8(buf, color as any)) !== null) return [0 as _.Color, RpcErrCode.ReqErr];

        const [bytes, status] = await this._fetch("pick", buf.bytes);
        if (status !== RpcErrCode.Ok || bytes === null) return [0 as _.Color, status];

        const respBuf = new _.Buffer(bytes);const [result, err] = _.getU8(respBuf);
        if (err === null && !_.IsColor(result as any)) return [0 as _.Color, RpcErrCode.RespErr];
        if (err !== null) return [0 as _.Color, RpcErrCode.RespErr];
        if (respBuf.len !== 0) return [0 as _.Color, RpcErrCode.RespErr];
        return [result as any, RpcErrCode.Ok];
    };
    /**  */
    public pickBad = async (color: _.Color): Promise<[_.Color, RpcErrCode]> => {
        const buf = new _.Buffer();
        if (!_.IsColor(color as any)) return [0 as _.Color, RpcErrCode.ReqErr];
        if (_.setAll(buf, (buf) => _.setU8(buf, color as any)) !== null) return [0 as _.Color, RpcErrCode.ReqErr];

        const [bytes, status] = await this._fetch("pick_bad", buf.bytes);
        if (status !== RpcErrCode.Ok || bytes === null) return [0 as _.Color, status];

        const respBuf = new _.Buffer(bytes);const [result, err] = _.getU8(respBuf);
        if (err === null && !_.IsColor(result as any)) return [0 as _.Color, RpcErrCode.RespErr];
        if (err !== null) return [0 as _.Color, RpcErrCode.RespErr];
        if (respBuf.len !== 0) return [0 as _.Color, RpcErrCode.RespErr];
        return [result as any, RpcErrCode.Ok];
    };
    
}
