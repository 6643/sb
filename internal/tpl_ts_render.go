package internal

import (
	"fmt"
	"strings"
)

func (g *TsGenerator) renderTsEnumFile(enums []TplEnum) string {
	var w sourceWriter
	for _, enum := range enums {
		enumName := PascalCase(enum.Name)
		w.Blank()
		w.WriteLineComment("// ", enum.Note)
		w.Linef("export enum %s {", enumName)
		for _, child := range enum.Children {
			w.WriteLineComment("    // ", child.Note)
			w.Linef("    %s = %d,", PascalCase(child.Name), child.ID)
		}
		w.Line("}")
		w.Blank()
		w.Linef("export const Is%s = (v: %s): boolean => {", enumName, enumName)
		w.Line("    switch (v) {")
		for _, child := range enum.Children {
			w.Linef("    case %s.%s:", enumName, PascalCase(child.Name))
		}
		w.Line("        return true;")
		w.Line("    default:")
		w.Line("        return false;")
		w.Line("    }")
		w.Line("}")
		w.Blank()
		w.Linef("export const Is%sList = (v: %s[]): boolean => {", enumName, enumName)
		w.Line("    for (const item of v) {")
		w.Linef("        if (!Is%s(item)) return false;", enumName)
		w.Line("    }")
		w.Line("    return true;")
		w.Line("}")
	}
	return w.String()
}

func (g *TsGenerator) renderTsStructFile(st TplStruct) string {
	name := PascalCase(st.Name)
	bitmaskSize := goBitmaskSize(len(st.Fields))
	var w sourceWriter
	w.Line("import * as _ from \"./_\"")
	w.Line("import * as TplEnum from \"./enum\"")
	w.Blank()
	w.Linef("export interface %s extends _.Serializable, _.Deserializable {", name)
	for _, field := range st.Fields {
		w.Linef("    %s: %s;", CamelCase(field.Name), g.tsStructFieldType(field.Type))
	}
	w.Line("}")
	w.Blank()
	w.Linef("export const new%s = (): %s => {", name, name)
	w.Line("    const s = {")
	for _, field := range st.Fields {
		w.Linef("        %s: %s,", CamelCase(field.Name), g.tsStructDefaultValue(field.Type))
	}
	w.Linef("    } as any as %s;", name)
	w.Linef("    s.set = (buf: _.Buffer) => set%s(buf, s);", name)
	w.Line("    s.get = (buf: _.Buffer) => {")
	w.Linef("        const [res, err] = get%s(buf);", name)
	w.Line("        if (err === null) Object.assign(s, res);")
	w.Line("        return err;")
	w.Line("    };")
	w.Line("    return s;")
	w.Line("}")
	w.Blank()
	w.Linef("export const eq%s = (a: %s, b: %s): boolean => {", name, name, name)
	w.Line("    if (a === b) return true;")
	w.Line("    if (a === null || b === null) return false;")
	for _, field := range st.Fields {
		left := "a." + CamelCase(field.Name)
		right := "b." + CamelCase(field.Name)
		for _, line := range g.tsEqLines(field.Type, left, right) {
			w.Line(line)
		}
	}
	w.Line("    return true;")
	w.Line("}")
	w.Blank()
	w.Linef("export const get%s = (buf: _.Buffer): [%s, Error | null] => {", name, name)
	w.Linef("    const s = new%s();", name)
	w.Linef("    const bitmaskSize = %d;", bitmaskSize)
	w.Line("    const [bits, err] = buf.read(bitmaskSize);")
	w.Line("    if (err !== null) return [s, err];")
	for i, field := range st.Fields {
		fieldRef := "s." + CamelCase(field.Name)
		fieldName := PascalCase(field.Name)
		if field.Type.Name == "bool" {
			w.Linef("    %s = _.GetBit(bits, %d);", fieldRef, i)
			continue
		}
		w.Linef("    if (_.GetBit(bits, %d)) {", i)
		for _, line := range g.tsGetLines(name, fieldName, field.Type, fieldRef) {
			w.Line(line)
		}
		w.Line("    }")
	}
	w.Line("    return [s, null];")
	w.Line("}")
	w.Blank()
	w.Linef("export const set%s = (buf: _.Buffer, s: %s): Error | null => {", name, name)
	w.Linef("    if (s === null || s === undefined) return new Error(`set %s: value is null or undefined`);", name)
	w.Line("    const startOffset = buf.write_offset;")
	w.Linef("    const bits = new Uint8Array(%d);", bitmaskSize)
	for i, field := range st.Fields {
		fieldRef := "s." + CamelCase(field.Name)
		fieldName := PascalCase(field.Name)
		for _, line := range g.tsSetBitLines(name, fieldName, field.Type, fieldRef, i) {
			w.Line(line)
		}
	}
	w.Blank()
	w.Line("    const errBits = buf.write(bits);")
	w.Line("    if (errBits !== null) return errBits;")
	for _, field := range st.Fields {
		fieldRef := "s." + CamelCase(field.Name)
		for _, line := range g.tsWriteLines(field.Type, fieldRef) {
			w.Line(line)
		}
	}
	w.Line("    return null;")
	w.Line("}")
	w.Blank()
	w.Linef("export const get%sList = (buf: _.Buffer): [%s[], Error | null] => {", name, name)
	w.Line("    const [count, err] = _.getU16(buf);")
	w.Line("    if (err !== null) return [[], err];")
	w.Linef("    const list: %s[] = new Array(count);", name)
	w.Line("    for (let i = 0; i < count; i++) {")
	w.Linef("        const [item, err2] = get%s(buf);", name)
	w.Line("        if (err2 !== null) return [[], err2];")
	w.Line("        list[i] = item;")
	w.Line("    }")
	w.Line("    return [list, null];")
	w.Line("}")
	w.Linef("export const set%sList = (buf: _.Buffer, v: %s[]): Error | null => {", name, name)
	w.Line("    if (v.length > 65535) return new Error(`list length ${v.length} exceeds u16 max`);")
	w.Line("    const err = _.setU16(buf, v.length);")
	w.Line("    if (err !== null) return err;")
	w.Line("    for (const item of v) {")
	w.Linef("        const err2 = set%s(buf, item);", name)
	w.Line("        if (err2 !== null) return err2;")
	w.Line("    }")
	w.Line("    return null;")
	w.Line("}")
	w.Linef("export const eq%sList = (a: %s[], b: %s[]): boolean => _.eqList(a, b, eq%s);", name, name, name, name)
	return w.String()
}

func (g *TsGenerator) renderTsIndexFile(files []string) string {
	var w sourceWriter
	w.Line("export * from \"./type\"")
	for _, file := range files {
		w.Linef("export * from \"./%s\"", file)
	}
	return w.String()
}

func (g *TsGenerator) renderTsRPCFile(apis []TplApi) string {
	var w sourceWriter
	w.Line("import * as _ from \"./_\"")
	w.Blank()
	w.Line("const defaultMaxRespBytes = 4 * 1024 * 1024;")
	w.Line("const maxSafeRespBytes = Number.MAX_SAFE_INTEGER;")
	w.Line("const maxTimeoutMs = 2147483647;")
	w.Blank()
	w.Line("export enum RpcErrCode {")
	w.Line("    Ok = 200,")
	w.Line("    NoConn = 0,")
	w.Line("    Timeout = 408,")
	w.Line("    ReqErr = 400,")
	w.Line("    RespErr = 500,")
	w.Line("    NotAuth = 401,")
	w.Line("    NotExist = 404,")
	w.Line("}")
	w.Blank()
	w.Line("export interface RpcConfig {")
	w.Line("    host: string;")
	w.Line("    headers?: Record<string, string>;")
	w.Line("    timeout?: number;")
	w.Line("    retries?: number;")
	w.Line("    maxRespBytes?: number;")
	w.Line("}")
	w.Blank()
	w.Line("export class RpcClient {")
	w.Line("    private headers: Record<string, string> = {};")
	w.Line("    private timeout: number;")
	w.Line("    private retries: number;")
	w.Line("    private maxRespBytes: number;")
	w.Blank()
	w.Line("    constructor(private config: RpcConfig) {")
	w.Line("        this.config = { ...config, host: config.host.replace(/\\/+$/, \"\") };")
	w.Line("        if (config.headers) this.headers = { ...config.headers };")
	w.Line("        const cfgTimeout = config.timeout;")
	w.Line("        if (cfgTimeout !== undefined && Number.isFinite(cfgTimeout) && cfgTimeout >= 0) {")
	w.Line("            this.timeout = Math.min(Math.floor(cfgTimeout), maxTimeoutMs);")
	w.Line("        } else {")
	w.Line("            this.timeout = 5000;")
	w.Line("        }")
	w.Line("        const cfgRetries = config.retries;")
	w.Line("        if (cfgRetries !== undefined && Number.isFinite(cfgRetries) && cfgRetries >= 0) {")
	w.Line("            this.retries = Math.floor(cfgRetries);")
	w.Line("        } else {")
	w.Line("            this.retries = 3;")
	w.Line("        }")
	w.Line("        const cfgMaxRespBytes = config.maxRespBytes;")
	w.Line("        if (cfgMaxRespBytes !== undefined && Number.isFinite(cfgMaxRespBytes) && cfgMaxRespBytes > 0) {")
	w.Line("            this.maxRespBytes = Math.min(Math.floor(cfgMaxRespBytes), maxSafeRespBytes);")
	w.Line("        } else {")
	w.Line("            this.maxRespBytes = defaultMaxRespBytes;")
	w.Line("        }")
	w.Line("    }")
	w.Blank()
	w.Line("    public setHeader = (key: string, value: string): void => { this.headers[key] = value; };")
	w.Line("    public getHeader = (key: string): string | undefined => this.headers[key];")
	w.Line("    public removeHeader = (key: string): void => { delete this.headers[key]; };")
	w.Blank()
	w.Line("    public setAuthorization = (token: string): void => { this.setHeader(\"Authorization\", `Bearer ${token}`); };")
	w.Line("    public getAuthorization = (): string | undefined => this.getHeader(\"Authorization\");")
	w.Line("    public removeAuthorization = (): void => { this.removeHeader(\"Authorization\"); };")
	w.Line("    public isAuthorized = (): boolean => !!this.getAuthorization();")
	w.Blank()
	w.Line("    private async _fetch(path: string, body: Uint8Array): Promise<[Uint8Array | null, RpcErrCode]> {")
	w.Line("        let lastStatus = RpcErrCode.NoConn;")
	w.Line("        for (let i = 0; i <= this.retries; i++) {")
	w.Line("            if (i > 0) await new Promise(res => setTimeout(res, i * 1000));")
	w.Line("            const controller = new AbortController();")
	w.Line("            let timeoutId: ReturnType<typeof setTimeout> | null = null;")
	w.Line("            if (this.timeout > 0) timeoutId = setTimeout(() => controller.abort(), this.timeout);")
	w.Line("            try {")
	w.Line("                const res = await fetch(`${this.config.host}/${path}`, {")
	w.Line("                    method: \"POST\",")
	w.Line("                    headers: { \"Content-TplType\": \"application/octet-stream\", ...this.headers },")
	w.Line("                    body: body as any,")
	w.Line("                    signal: controller.signal")
	w.Line("                });")
	w.Line("                if (res.ok) {")
	w.Line("                    const contentLength = res.headers.get(\"content-length\");")
	w.Line("                    if (contentLength !== null) {")
	w.Line("                        const size = Number(contentLength);")
	w.Line("                        if (Number.isFinite(size) && size > this.maxRespBytes) return [null, RpcErrCode.RespErr];")
	w.Line("                    }")
	w.Line("                    const bytes = new Uint8Array(await res.arrayBuffer());")
	w.Line("                    if (bytes.byteLength > this.maxRespBytes) return [null, RpcErrCode.RespErr];")
	w.Line("                    return [bytes, RpcErrCode.Ok];")
	w.Line("                }")
	w.Line("                lastStatus = res.status as RpcErrCode;")
	w.Line("                if (res.status === 408 && i < this.retries) continue;")
	w.Line("                return [null, res.status as RpcErrCode];")
	w.Line("            } catch (e: any) {")
	w.Line("                if (e.name === \"AbortError\") {")
	w.Line("                    lastStatus = RpcErrCode.Timeout;")
	w.Line("                } else {")
	w.Line("                    lastStatus = RpcErrCode.NoConn;")
	w.Line("                }")
	w.Line("                if (i < this.retries) continue;")
	w.Line("            } finally {")
	w.Line("                if (timeoutId !== null) clearTimeout(timeoutId);")
	w.Line("            }")
	w.Line("        }")
	w.Line("        return [null, lastStatus];")
	w.Line("    }")
	for _, api := range apis {
		methodName := CamelCase(api.Name)
		callName := PascalCase(api.Result.Name)
		defaultVal := g.tsRPCDefaultValue(api.Result)
		w.Blank()
		w.WriteLineComment("    // ", api.Note)
		w.Linef("    public %s = async (%s): Promise<%s> => {", methodName, g.tsRPCArgList(api.Args), g.tsRPCPromiseType(api.Result))
		w.Line("        const buf = new _.Buffer();")
		for _, arg := range api.Args {
			if arg.Type.Kind != TplKindEnum {
				continue
			}
			name := arg.Name
			typeName := PascalCase(arg.Type.Name)
			if arg.Type.IsList {
				w.Linef("        if (!_.Is%sList(%s as any)) return %s;", typeName, name, g.tsRPCReqErrReturn(api.Result, defaultVal))
			} else {
				w.Linef("        if (!_.Is%s(%s as any)) return %s;", typeName, name, g.tsRPCReqErrReturn(api.Result, defaultVal))
			}
			}
			if len(api.Args) > 0 {
				w.Line("        let encodeErr: Error | null = null;")
				for _, arg := range api.Args {
					w.Linef("        if ((encodeErr = %s) !== null) return %s;", g.tsRPCSetterCall(arg, "buf"), g.tsRPCReqErrReturn(api.Result, defaultVal))
				}
			}
		w.Blank()
		w.Linef("        const [bytes, status] = await this._fetch(\"%s\", buf.bytes);", api.Name)
		w.Linef("        if (status !== RpcErrCode.Ok || bytes === null) return %s;", g.tsRPCStatusReturn(api.Result, defaultVal))
		w.Blank()
		if api.Result.Name == "nil" {
			w.Line("        if (bytes.byteLength !== 0) return RpcErrCode.RespErr;")
			w.Line("        return RpcErrCode.Ok;")
			w.Line("    };")
			continue
		}
		w.Line("        const respBuf = new _.Buffer(bytes);")
		w.Blank()
		switch {
		case api.Result.Kind == TplKindEnum && api.Result.IsList:
			w.Line("        const [result, err] = _.getU8List(respBuf);")
			w.Linef("        if (err === null && !_.Is%sList(result as any)) return [%s, RpcErrCode.RespErr];", callName, defaultVal)
		case api.Result.Kind == TplKindEnum:
			w.Line("        const [result, err] = _.getU8(respBuf);")
			w.Linef("        if (err === null && !_.Is%s(result as any)) return [%s, RpcErrCode.RespErr];", callName, defaultVal)
		default:
			w.Linef("        const [result, err] = _.get%s%s(respBuf);", callName, tsListSuffix(api.Result))
		}
		w.Blank()
		w.Linef("        if (err !== null) return [%s, RpcErrCode.RespErr];", defaultVal)
		w.Linef("        if (respBuf.len !== 0) return [%s, RpcErrCode.RespErr];", defaultVal)
		w.Line("        return [result as any, RpcErrCode.Ok];")
		w.Line("    };")
	}
	w.Blank()
	w.Line("}")
	return w.String()
}

func (g *TsGenerator) renderTsSmokeTest(apis []TplApi) string {
	var w sourceWriter
	w.Line("import { describe, expect, test } from \"bun:test\";")
	w.Blank()
	w.Line("import * as _ from \"./_\";")
	w.Line("import { RpcClient, RpcErrCode } from \"./rpc\";")
	w.Blank()
	w.Line("const baseUrl = process.env.SB_BASE_URL || \"http://127.0.0.1:18080\";")
	w.Blank()
	w.Line("describe(\"rpc smoke\", () => {")
	w.Line("    test(\"client construction\", () => {")
	w.Line("        const client = new RpcClient({ host: baseUrl });")
	w.Line("        expect(client).toBeDefined();")
	w.Line("        expect(typeof RpcErrCode.Ok).toBe(\"number\");")
	w.Line("        expect(typeof _).toBe(\"object\");")
	w.Line("    });")
	for _, api := range apis {
		w.Linef("    test(\"method %s exists\", () => {", CamelCase(api.Name))
		w.Line("        const client = new RpcClient({ host: baseUrl });")
		w.Linef("        expect(typeof (client as any).%s).toBe(\"function\");", CamelCase(api.Name))
		w.Line("    });")
	}
	w.Line("});")
	return w.String()
}

func (g *TsGenerator) tsStructFieldType(t TplType) string {
	base := g.getTsType(TplType{Name: t.Name, Kind: t.Kind})
	switch t.Kind {
	case TplKindEnum:
		base = "TplEnum." + base
	case TplKindStruct:
		base = "_." + base
	}
	if t.IsList {
		return base + "[]"
	}
	if t.Kind == TplKindStruct {
		return base + " | null"
	}
	return base
}

func (g *TsGenerator) tsStructDefaultValue(t TplType) string {
	if t.IsList {
		return "[]"
	}
	if t.Kind == TplKindBase {
		return g.getTsValue(t.Name)
	}
	if t.Kind == TplKindEnum {
		return "0"
	}
	return "null"
}

func (g *TsGenerator) tsEqLines(t TplType, left, right string) []string {
	name := PascalCase(t.Name)
	switch t.Kind {
	case TplKindBase:
		if t.IsList {
			return []string{fmt.Sprintf("    if (!_.eq%sList(%s, %s)) return false;", name, left, right)}
		}
		return []string{fmt.Sprintf("    if (!_.eq%s(%s, %s)) return false;", name, left, right)}
	case TplKindEnum:
		if t.IsList {
			return []string{fmt.Sprintf("    if (!_.eqU8List(%s as any, %s as any)) return false;", left, right)}
		}
		return []string{fmt.Sprintf("    if (%s !== %s) return false;", left, right)}
	case TplKindStruct:
		if t.IsList {
			return []string{fmt.Sprintf("    if (!_.eq%sList(%s, %s)) return false;", name, left, right)}
		}
		return []string{
			fmt.Sprintf("    if ((%s === null) !== (%s === null)) return false;", left, right),
			fmt.Sprintf("    if (%s !== null && %s !== null && !_.eq%s(%s, %s)) return false;", left, right, name, left, right),
		}
	}
	return nil
}

func (g *TsGenerator) tsGetLines(structName, fieldName string, t TplType, target string) []string {
	name := PascalCase(t.Name)
	var lines []string
	switch {
	case t.Kind == TplKindBase:
		lines = append(lines, fmt.Sprintf("        const [v, err] = _.get%s%s(buf);", name, tsListSuffix(t)))
		lines = append(lines, "        if (err !== null) return [s, err];")
		lines = append(lines, fmt.Sprintf("        %s = v;", target))
	case t.Kind == TplKindEnum && t.IsList:
		lines = append(lines, "        const [v, err] = _.getU8List(buf);")
		lines = append(lines, "        if (err !== null) return [s, err];")
		lines = append(lines, fmt.Sprintf("        %s = v as any;", target))
		lines = append(lines, fmt.Sprintf("        if (!_.Is%sList(%s as any)) return [s, new Error(\"get %s %s: invalid enum value\")];", name, target, structName, fieldName))
	case t.Kind == TplKindEnum:
		lines = append(lines, "        const [v, err] = _.getU8(buf);")
		lines = append(lines, "        if (err !== null) return [s, err];")
		lines = append(lines, fmt.Sprintf("        %s = v as any;", target))
		lines = append(lines, fmt.Sprintf("        if (!_.Is%s(%s as any)) return [s, new Error(\"get %s %s: invalid enum value\")];", name, target, structName, fieldName))
	default:
		lines = append(lines, fmt.Sprintf("        const [v, err] = _.get%s%s(buf);", name, tsListSuffix(t)))
		lines = append(lines, "        if (err !== null) return [s, err];")
		lines = append(lines, fmt.Sprintf("        %s = v;", target))
	}
	return lines
}

func (g *TsGenerator) tsSetBitLines(structName, fieldName string, t TplType, ref string, bit int) []string {
	name := PascalCase(t.Name)
	var lines []string
	switch {
	case t.Name == "bool":
		return []string{fmt.Sprintf("    _.SetBit(bits, %d, %s as boolean);", bit, ref)}
	case t.Kind == TplKindBase && t.IsList:
		lines = append(lines, fmt.Sprintf("    if (%s && %s.length > 0) {", ref, ref))
	case t.Kind == TplKindBase:
		lines = append(lines, fmt.Sprintf("    if (!_.eq%s(%s, %s)) {", name, ref, g.getTsValue(t.Name)))
	case t.Kind == TplKindEnum && t.IsList:
		lines = append(lines, fmt.Sprintf("    if (%s && %s.length > 0) {", ref, ref))
		lines = append(lines, fmt.Sprintf("        if (!_.Is%sList(%s as any)) return new Error(\"set %s %s: invalid enum value\");", name, ref, structName, fieldName))
	case t.Kind == TplKindEnum:
		lines = append(lines, fmt.Sprintf("    if ((%s as any) !== 0) {", ref))
		lines = append(lines, fmt.Sprintf("        if (!_.Is%s(%s as any)) return new Error(\"set %s %s: invalid enum value\");", name, ref, structName, fieldName))
	case t.Kind == TplKindStruct && t.IsList:
		lines = append(lines, fmt.Sprintf("    if (%s && %s.length > 0) {", ref, ref))
	default:
		lines = append(lines, fmt.Sprintf("    if (%s !== null && %s !== undefined) {", ref, ref))
	}
	lines = append(lines, fmt.Sprintf("        _.SetBit(bits, %d, true);", bit))
	lines = append(lines, "    }")
	return lines
}

func (g *TsGenerator) tsWriteLines(t TplType, ref string) []string {
	name := PascalCase(t.Name)
	var lines []string
	switch {
	case t.Name == "bool":
		return nil
	case t.Kind == TplKindBase && t.IsList:
		lines = append(lines, fmt.Sprintf("    if (%s && %s.length > 0) {", ref, ref))
		lines = append(lines, fmt.Sprintf("        const err = _.set%sList(buf, %s);", name, ref))
	case t.Kind == TplKindBase:
		lines = append(lines, fmt.Sprintf("    if (!_.eq%s(%s, %s)) {", name, ref, g.getTsValue(t.Name)))
		lines = append(lines, fmt.Sprintf("        const err = _.set%s(buf, %s);", name, ref))
	case t.Kind == TplKindEnum && t.IsList:
		lines = append(lines, fmt.Sprintf("    if (%s && %s.length > 0) {", ref, ref))
		lines = append(lines, fmt.Sprintf("        const err = _.setU8List(buf, %s as any);", ref))
	case t.Kind == TplKindEnum:
		lines = append(lines, fmt.Sprintf("    if ((%s as any) !== 0) {", ref))
		lines = append(lines, fmt.Sprintf("        const err = _.setU8(buf, %s as any);", ref))
	case t.Kind == TplKindStruct && t.IsList:
		lines = append(lines, fmt.Sprintf("    if (%s && %s.length > 0) {", ref, ref))
		lines = append(lines, fmt.Sprintf("        const err = _.set%sList(buf, %s);", name, ref))
	default:
		lines = append(lines, fmt.Sprintf("    if (%s !== null && %s !== undefined) {", ref, ref))
		lines = append(lines, fmt.Sprintf("        const err = _.set%s(buf, %s);", name, ref))
	}
	lines = append(lines, "        if (err !== null) {")
	lines = append(lines, "            buf.rewindWrite(startOffset);")
	lines = append(lines, "            return err;")
	lines = append(lines, "        }")
	lines = append(lines, "    }")
	return lines
}

func (g *TsGenerator) tsRPCArgList(args []TplApiArg) string {
	parts := make([]string, 0, len(args))
	for _, arg := range args {
		prefix := ""
		if arg.Type.Kind != TplKindBase {
			prefix = "_."
		}
		parts = append(parts, fmt.Sprintf("%s: %s%s", arg.Name, prefix, g.getTsLogicType(arg.Type)))
	}
	return strings.Join(parts, ", ")
}

func (g *TsGenerator) tsRPCPromiseType(result TplType) string {
	if result.Name == "nil" {
		return "RpcErrCode"
	}
	prefix := ""
	if result.Kind != TplKindBase {
		prefix = "_."
	}
	return fmt.Sprintf("[%s%s, RpcErrCode]", prefix, g.getTsLogicType(result))
}

func (g *TsGenerator) tsRPCDefaultValue(result TplType) string {
	if result.IsList {
		return "[]"
	}
	if result.Kind == TplKindBase {
		return g.getTsValue(result.Name)
	}
	if result.Kind == TplKindEnum {
		return fmt.Sprintf("0 as _.%s", PascalCase(result.Name))
	}
	if result.Kind == TplKindStruct {
		return fmt.Sprintf("_.new%s()", PascalCase(result.Name))
	}
	return "null"
}

func (g *TsGenerator) tsRPCReqErrReturn(result TplType, defaultVal string) string {
	if result.Name == "nil" {
		return "RpcErrCode.ReqErr"
	}
	return fmt.Sprintf("[%s, RpcErrCode.ReqErr]", defaultVal)
}

func (g *TsGenerator) tsRPCStatusReturn(result TplType, defaultVal string) string {
	if result.Name == "nil" {
		return "status"
	}
	return fmt.Sprintf("[%s, status]", defaultVal)
}

func (g *TsGenerator) tsRPCSetters(args []TplApiArg) string {
	parts := make([]string, 0, len(args))
	for _, arg := range args {
		name := arg.Name
		switch arg.Type.Kind {
		case TplKindBase:
			parts = append(parts, fmt.Sprintf(", (buf: _.Buffer) => _.set%s%s(buf, %s)", PascalCase(arg.Type.Name), tsListSuffix(arg.Type), name))
		case TplKindEnum:
			parts = append(parts, fmt.Sprintf(", (buf: _.Buffer) => _.setU8%s(buf, %s as any)", tsListSuffix(arg.Type), name))
		case TplKindStruct:
			parts = append(parts, fmt.Sprintf(", (buf: _.Buffer) => _.set%s%s(buf, %s as any)", PascalCase(arg.Type.Name), tsListSuffix(arg.Type), name))
		}
	}
	return strings.Join(parts, "")
}

func (g *TsGenerator) tsRPCSetterCall(arg TplApiArg, bufVar string) string {
	name := arg.Name
	switch arg.Type.Kind {
	case TplKindBase:
		return fmt.Sprintf("_.set%s%s(%s, %s)", PascalCase(arg.Type.Name), tsListSuffix(arg.Type), bufVar, name)
	case TplKindEnum:
		return fmt.Sprintf("_.setU8%s(%s, %s as any)", tsListSuffix(arg.Type), bufVar, name)
	case TplKindStruct:
		return fmt.Sprintf("_.set%s%s(%s, %s as any)", PascalCase(arg.Type.Name), tsListSuffix(arg.Type), bufVar, name)
	default:
		return "null"
	}
}

func tsListSuffix(t TplType) string {
	if t.IsList {
		return "List"
	}
	return ""
}
