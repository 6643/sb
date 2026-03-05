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

    {{range .Apis}}
    {{- $resData := .Result -}}
    {{- $hasRet := ne $resData.Name "nil" -}}
    {{- $retType := TsLogicType $resData -}}
    {{- $defaultVal := "null" -}}
    {{- if $hasRet -}}
        {{- if $resData.IsList -}}{{$defaultVal = "[]"}}
        {{- else if IsBaseType $resData -}}{{$defaultVal = TsValue $resData.Name}}
        {{- else if IsEnum $resData -}}{{$defaultVal = printf "0 as _.%s" (PascalCase $resData.Name)}}
        {{- else -}}{{$defaultVal = printf "_.new%s()" (PascalCase $resData.Name)}}
        {{- end -}}
    {{- end -}}
    /** {{.Note}} */
    public {{.Name | CamelCase}} = async ({{range $i, $arg := .Args}}{{if $i}}, {{end}}{{$arg.Name}}: {{if not (IsBaseType .Type)}}_.{{end}}{{TsLogicType $arg.Type}}{{end}}): Promise<{{if $hasRet}}[{{if not (IsBaseType $resData)}}_.{{end}}{{$retType}}, RpcErrCode]{{else}}RpcErrCode{{end}}> => {
        const buf = new _.Buffer();
        {{- range .Args}}
        {{- if IsEnum .Type}}
        {{- if .Type.IsList}}
        if (!_.Is{{.Type.Name | PascalCase}}List({{.Name}} as any)) return {{if $hasRet}}[{{$defaultVal}}, RpcErrCode.ReqErr]{{else}}RpcErrCode.ReqErr{{end}};
        {{- else}}
        if (!_.Is{{.Type.Name | PascalCase}}({{.Name}} as any)) return {{if $hasRet}}[{{$defaultVal}}, RpcErrCode.ReqErr]{{else}}RpcErrCode.ReqErr{{end}};
        {{- end}}
        {{- end}}
        {{- end}}
        {{- if .Args}}
        if (_.setAll(buf, {{range $i, $arg := .Args}}{{if $i}}, {{end}}{{if IsBaseType .Type}}(buf) => _.set{{.Type.Name | PascalCase}}{{if .Type.IsList}}List{{end}}(buf, {{$arg.Name}}){{else if IsEnum .Type}}(buf) => _.setU8{{if .Type.IsList}}List{{end}}(buf, {{$arg.Name}} as any){{else}}(buf) => _.set{{.Type.Name | PascalCase}}{{if .Type.IsList}}List{{end}}(buf, {{$arg.Name}} as any){{end}}{{end}}) !== null) return {{if $hasRet}}[{{$defaultVal}}, RpcErrCode.ReqErr]{{else}}RpcErrCode.ReqErr{{end}};
        {{- end}}

        const [bytes, status] = await this._fetch("{{.Name}}", buf.bytes);
        if (status !== RpcErrCode.Ok || bytes === null) return {{if $hasRet}}[{{$defaultVal}}, status]{{else}}status{{end}};

        {{if $hasRet -}}
        const respBuf = new _.Buffer(bytes);
        {{- if IsEnum $resData -}}
        {{- if $resData.IsList -}}
        const [result, err] = _.getU8List(respBuf);
        if (err === null && !_.Is{{$resData.Name | PascalCase}}List(result as any)) return [{{$defaultVal}}, RpcErrCode.RespErr];
        {{- else -}}
        const [result, err] = _.getU8(respBuf);
        if (err === null && !_.Is{{$resData.Name | PascalCase}}(result as any)) return [{{$defaultVal}}, RpcErrCode.RespErr];
        {{- end -}}
        {{- else -}}
        const [result, err] = _.get{{$resData.Name | PascalCase}}{{if $resData.IsList}}List{{end}}(respBuf);
        {{- end}}
        if (err !== null) return [{{$defaultVal}}, RpcErrCode.RespErr];
        if (respBuf.len !== 0) return [{{$defaultVal}}, RpcErrCode.RespErr];
        return [result as any, RpcErrCode.Ok];
        {{- else -}}
        if (bytes.byteLength !== 0) return RpcErrCode.RespErr;
        return RpcErrCode.Ok;
        {{- end}}
    };
    {{end}}
}
