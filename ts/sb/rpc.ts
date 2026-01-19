import * as _ from "./_.ts"

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
}

export class RpcClient {
    private headers: Record<string, string> = {};
    private timeout: number;
    private retries: number;

    constructor(private config: RpcConfig) {
        if (config.headers) this.headers = { ...config.headers };
        this.timeout = config.timeout || 5000;
        this.retries = config.retries !== undefined ? config.retries : 3;
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
            const timeoutId = setTimeout(() => controller.abort(), this.timeout);
            try {
                const res = await fetch(`${this.config.host}/${path}`, {
                    method: "POST",
                    headers: { "Content-Type": "application/octet-stream", ...this.headers },
                    body: body as any,
                    signal: controller.signal
                });
                if (res.ok) return [new Uint8Array(await res.arrayBuffer()), RpcErrCode.Ok];
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
                clearTimeout(timeoutId);
            }
        }
        return [null, lastStatus];
    }

    /** 获取用户的id */
    public userGetAbc = async (): Promise<[_.OrderStatus, RpcErrCode]> => {
        const buf = new _.Buffer();

        const [bytes, status] = await this._fetch("user.get_abc", buf.bytes);
        if (status !== RpcErrCode.Ok || bytes === null) return [0 as _.OrderStatus, status];

        const [result, err] = _.getU8(new _.Buffer(bytes));
        if (err !== null) return [0 as _.OrderStatus, RpcErrCode.RespErr];
        return [result as any, RpcErrCode.Ok];
    };
    /** 获取abcd */
    public userGetAbcd = async (page: number, size: number): Promise<[_.OrderStatus, RpcErrCode]> => {
        const buf = new _.Buffer();
        if (_.setAll(buf, _.u8(page), _.u8(size)) !== null) return [0 as _.OrderStatus, RpcErrCode.ReqErr];

        const [bytes, status] = await this._fetch("user.get_abcd", buf.bytes);
        if (status !== RpcErrCode.Ok || bytes === null) return [0 as _.OrderStatus, status];

        const [result, err] = _.getU8(new _.Buffer(bytes));
        if (err !== null) return [0 as _.OrderStatus, RpcErrCode.RespErr];
        return [result as any, RpcErrCode.Ok];
    };
    /** 设置sim信息 */
    public userSetSimInfo = async (info: _.SimInfo): Promise<RpcErrCode> => {
        const buf = new _.Buffer();
        if (_.setAll(buf, info) !== null) return RpcErrCode.ReqErr;

        const [bytes, status] = await this._fetch("user.set_sim_info", buf.bytes);
        if (status !== RpcErrCode.Ok || bytes === null) return status;

        return RpcErrCode.Ok;
    };
    /** 获取数量 */
    public getCount = async (page: number): Promise<[number, RpcErrCode]> => {
        const buf = new _.Buffer();
        if (_.setAll(buf, _.u8(page)) !== null) return [0, RpcErrCode.ReqErr];

        const [bytes, status] = await this._fetch("get_count", buf.bytes);
        if (status !== RpcErrCode.Ok || bytes === null) return [0, status];

        const [result, err] = _.getU8(new _.Buffer(bytes));
        if (err !== null) return [0, RpcErrCode.RespErr];
        return [result as any, RpcErrCode.Ok];
    };
    /** 获取bin */
    public getBin = async (page: number): Promise<[Uint8Array, RpcErrCode]> => {
        const buf = new _.Buffer();
        if (_.setAll(buf, _.u8(page)) !== null) return [new Uint8Array(0), RpcErrCode.ReqErr];

        const [bytes, status] = await this._fetch("get_bin", buf.bytes);
        if (status !== RpcErrCode.Ok || bytes === null) return [new Uint8Array(0), status];

        const [result, err] = _.getBin(new _.Buffer(bytes));
        if (err !== null) return [new Uint8Array(0), RpcErrCode.RespErr];
        return [result as any, RpcErrCode.Ok];
    };
    
}
