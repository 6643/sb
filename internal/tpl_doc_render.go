package internal

func generateDoc(schema *TplSchema, cfg Config) error {
	content := renderDoc(schema)
	if err := writeDocFile(cfg.GoDir, []byte(content)); err != nil {
		return err
	}
	if err := writeDocFile(cfg.TsDir, []byte(content)); err != nil {
		return err
	}
	return nil
}

func renderDoc(schema *TplSchema) string {
	var w sourceWriter
	w.Line("# API Documentation")
	w.Blank()
	if schema.Note != "" {
		w.Line(schema.Note)
		w.Blank()
	}
	renderDocAPIList(&w, schema.Apis)
	renderDocRPCStatusCodes(&w)
	renderDocUsageDemos(&w, schema)
	renderDocTypes(&w, schema)
	return w.String()
}
