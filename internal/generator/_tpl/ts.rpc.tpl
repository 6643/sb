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
        {{- if .Args}}
        if (_.setAll(buf, {{range $i, $arg := .Args}}{{if $i}}, {{end}}{{if IsBaseType .Type}}{{if .Type.IsList}}_.set{{.Type.Name | PascalCase}}List{{else}}_.{{.Type.Name | CamelCase}}{{end}}({{$arg.Name}}){{else if IsEnum .Type}}{{if .Type.IsList}}_.u8List({{$arg.Name}} as any){{else}}_.u8({{$arg.Name}} as any){{end}}{{else}}{{$arg.Name}}{{end}}{{end}}) !== null) return {{if $hasRet}}[{{$defaultVal}}, RpcErrCode.ReqErr]{{else}}RpcErrCode.ReqErr{{end}};
        {{- end}}

        const [bytes, status] = await this._fetch("{{.Name}}", buf.bytes);
        if (status !== RpcErrCode.Ok || bytes === null) return {{if $hasRet}}[{{$defaultVal}}, status]{{else}}status{{end}};

        {{if $hasRet -}}
        {{- if IsEnum $resData -}}
        {{- if $resData.IsList -}}
        const [result, err] = _.getU8List(new _.Buffer(bytes));
        {{- else -}}
        const [result, err] = _.getU8(new _.Buffer(bytes));
        {{- end -}}
        {{- else -}}
        const [result, err] = _.get{{$resData.Name | PascalCase}}{{if $resData.IsList}}List{{end}}(new _.Buffer(bytes));
        {{- end}}
        if (err !== null) return [{{$defaultVal}}, RpcErrCode.RespErr];
        return [result as any, RpcErrCode.Ok];
        {{- else -}}
        return RpcErrCode.Ok;
        {{- end}}
    };
    {{end}}
}
