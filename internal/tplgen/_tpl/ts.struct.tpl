import * as _ from "./_.ts"
import * as Enum from "./enum"

export interface {{.Name | PascalCase}} extends _.Serializable, _.Deserializable {
    {{- range .Fields}}
    {{.Name | CamelCase}}: {{if .Type.IsList}}{{if IsEnum .Type}}Enum.{{end}}{{if not (IsBaseType .Type)}}{{if not (IsEnum .Type)}}_.{{end}}{{end}}{{TsType .Type}}[]{{else}}{{if IsEnum .Type}}Enum.{{end}}{{if not (IsBaseType .Type)}}{{if not (IsEnum .Type)}}_.{{end}}{{end}}{{TsType .Type}}{{end}};
    {{- end}}
}

export const new{{.Name | PascalCase}} = (): {{.Name | PascalCase}} => {
    const s = {
        {{- range .Fields}}
        {{.Name | CamelCase}}: {{if .Type.IsList}}[]{{else}}{{if IsBaseType .Type}}{{TsValue .Type.Name}}{{else if IsEnum .Type}}0{{else}}_.new{{TsType .Type}}(){{end}}{{end}},
        {{- end}}
    } as any as {{.Name | PascalCase}};
    s.set = (buf: _.Buffer) => set{{.Name | PascalCase}}(buf, s);
    s.get = (buf: _.Buffer) => {
        const [res, err] = get{{.Name | PascalCase}}(buf);
        if (err === null) Object.assign(s, res);
        return err;
    };
    return s;
}

export const eq{{.Name | PascalCase}} = (a: {{.Name | PascalCase}}, b: {{.Name | PascalCase}}): boolean => {
    if (a === b) return true;
    if (a === null || b === null) return false;
    {{- range .Fields}}
    {{- if IsBaseType .Type}}
    if (!_.eq{{.Type.Name | PascalCase}}{{if .Type.IsList}}List{{end}}(a.{{.Name | CamelCase}}, b.{{.Name | CamelCase}})) return false;
    {{- else if IsEnum .Type}}
    {{- if .Type.IsList}}
    if (!_.eqU8List(a.{{.Name | CamelCase}} as any, b.{{.Name | CamelCase}} as any)) return false;
    {{- else}}
    if (a.{{.Name | CamelCase}} !== b.{{.Name | CamelCase}}) return false;
    {{- end}}
    {{- else}}
    {{- if .Type.IsList}}
    if (!_.eq{{.Type.Name | PascalCase}}List(a.{{.Name | CamelCase}}, b.{{.Name | CamelCase}})) return false;
    {{- else}}
    if (!_.eq{{.Type.Name | PascalCase}}(a.{{.Name | CamelCase}}, b.{{.Name | CamelCase}})) return false;
    {{- end}}
    {{- end}}
    {{- end}}
    return true;
}

export const get{{.Name | PascalCase}} = (buf: _.Buffer): [{{.Name | PascalCase}}, Error | null] => {
    const s = new{{.Name | PascalCase}}();
    const bitmaskSize = Math.ceil({{len .Fields}} / 8);
    const [bits, err] = buf.read(bitmaskSize);
    if (err !== null) return [s, err];

    {{- range $i, $field := .Fields}}
    {{- if eq .Type.Name "bool"}}
    s.{{$field.Name | CamelCase}} = _.GetBit(bits, {{$i}});
    {{- else}}
    if (_.GetBit(bits, {{$i}})) {
        {{- if IsBaseType .Type}}
        const [v, err] = _.get{{.Type.Name | PascalCase}}{{if .Type.IsList}}List{{end}}(buf);
        if (err !== null) return [s, err];
        s.{{$field.Name | CamelCase}} = v;
        {{- else if IsEnum .Type}}
        {{- if .Type.IsList}}
        const [v, err] = _.getU8List(buf);
        if (err !== null) return [s, err];
        s.{{$field.Name | CamelCase}} = v as any;
        {{- else}}
        const [v, err] = _.getU8(buf);
        if (err !== null) return [s, err];
        s.{{$field.Name | CamelCase}} = v as any;
        {{- end}}
        {{- else}}
        {{- if .Type.IsList}}
        const [v, err] = _.get{{.Type.Name | PascalCase}}List(buf);
        if (err !== null) return [s, err];
        s.{{$field.Name | CamelCase}} = v;
        {{- else}}
        const [v, err] = _.get{{.Type.Name | PascalCase}}(buf);
        if (err !== null) return [s, err];
        s.{{$field.Name | CamelCase}} = v;
        {{- end}}
        {{- end}}
    }
    {{- end}}
    {{- end}}
    return [s, null];
}

export const set{{.Name | PascalCase}} = (buf: _.Buffer, s: {{.Name | PascalCase}}): Error | null => {
    if (s === null || s === undefined) return new Error(`set {{.Name | PascalCase}}: value is null or undefined`);
    const bits = new Uint8Array(Math.ceil({{len .Fields}} / 8));
    const body = new _.Buffer();

    {{- range $i, $field := .Fields}}
    {{- if eq .Type.Name "bool"}}
    _.SetBit(bits, {{$i}}, s.{{$field.Name | CamelCase}} as boolean);
    {{- else if IsBaseType .Type}}
    {{- if .Type.IsList}}
    if (s.{{$field.Name | CamelCase}} && s.{{$field.Name | CamelCase}}.length > 0) {
        const err = _.set{{.Type.Name | PascalCase}}List(body, s.{{$field.Name | CamelCase}});
        if (err !== null) return err;
        _.SetBit(bits, {{$i}}, true);
    }
    {{- else}}
    if (!_.eq{{.Type.Name | PascalCase}}(s.{{$field.Name | CamelCase}}, {{TsValue .Type.Name}})) {
        const err = _.set{{.Type.Name | PascalCase}}(body, s.{{$field.Name | CamelCase}});
        if (err !== null) return err;
        _.SetBit(bits, {{$i}}, true);
    }
    {{- end}}
    {{- else if IsEnum .Type}}
    {{- if .Type.IsList}}
    if (s.{{$field.Name | CamelCase}} && s.{{$field.Name | CamelCase}}.length > 0) {
        const err = _.setU8List(body, s.{{$field.Name | CamelCase}} as any);
        if (err !== null) return err;
        _.SetBit(bits, {{$i}}, true);
    }
    {{- else}}
    if ((s.{{$field.Name | CamelCase}} as any) !== 0) {
        const err = _.setU8(body, s.{{$field.Name | CamelCase}} as any);
        if (err !== null) return err;
        _.SetBit(bits, {{$i}}, true);
    }
    {{- end}}
    {{- else}}
    {{- if .Type.IsList}}
    if (s.{{$field.Name | CamelCase}} && s.{{$field.Name | CamelCase}}.length > 0) {
        const err = _.set{{.Type.Name | PascalCase}}List(body, s.{{$field.Name | CamelCase}});
        if (err !== null) return err;
        _.SetBit(bits, {{$i}}, true);
    }
    {{- else}}
    if (s.{{$field.Name | CamelCase}} !== null) {
        const err = _.set{{.Type.Name | PascalCase}}(body, s.{{$field.Name | CamelCase}});
        if (err !== null) return err;
        _.SetBit(bits, {{$i}}, true);
    }
    {{- end}}
    {{- end}}
    {{- end}}

    const errBits = buf.write(bits);
    if (errBits !== null) return errBits;
    return buf.write(body.bytes);
}

export const get{{.Name | PascalCase}}List = (buf: _.Buffer): [{{.Name | PascalCase}}[], Error | null] => _.getList(buf, get{{.Name | PascalCase}});
export const set{{.Name | PascalCase}}List = (buf: _.Buffer, v: {{.Name | PascalCase}}[]): Error | null => _.setList(buf, v, set{{.Name | PascalCase}});
export const eq{{.Name | PascalCase}}List = (a: {{.Name | PascalCase}}[], b: {{.Name | PascalCase}}[]): boolean => _.eqList(a, b, eq{{.Name | PascalCase}});