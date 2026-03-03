{{range .Enums}}
{{if .Note}}// {{.Note}}{{end}}
export enum {{.Name | PascalCase}} {
{{- range .Children}}
    {{.Name | PascalCase}} = {{.ID}}, {{if .Note}}// {{.Note}}{{end}}
{{- end}}
}
{{end}}