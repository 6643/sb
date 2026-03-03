{{range .Enums}}
{{$enumName := .Name | PascalCase}}
{{if .Note}}// {{.Note}}{{end}}
export enum {{$enumName}} {
{{- range .Children}}
    {{.Name | PascalCase}} = {{.ID}}, {{if .Note}}// {{.Note}}{{end}}
{{- end}}
}

export const Is{{$enumName}} = (v: {{$enumName}}): boolean => {
    switch (v) {
    case {{range $i, $child := .Children}}{{if gt $i 0}}, {{end}}{{$enumName}}.{{$child.Name | PascalCase}}{{end}}:
        return true;
    default:
        return false;
    }
}

export const Is{{$enumName}}List = (v: {{$enumName}}[]): boolean => {
    for (const item of v) {
        if (!Is{{$enumName}}(item)) return false;
    }
    return true;
}
{{end}}
