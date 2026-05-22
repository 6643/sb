package internal

func (g *GoGenerator) renderAPIHandler(w *sourceWriter, api TplApi) {
	callName := PascalCase(api.Name)
	w.Blank()
	w.Linef("func %sHandler(w http.ResponseWriter, r *http.Request) {", callName)
	for _, arg := range api.Args {
		w.Linef("\tvar %s %s", CamelCase(arg.Name), g.logicType(arg.Type))
	}
	if len(api.Args) > 0 {
		w.Blank()
	}
	if len(api.Args) == 0 {
		w.Line("\tif !parseEmptyRequest(w, r) { return }")
	} else {
		w.Line("\tbuf, ok := parseRequest(w, r)")
		w.Line("\tif !ok { return }")
		for _, arg := range api.Args {
			g.directRead(w, arg.Type, CamelCase(arg.Name), "buf", "w.WriteHeader(http.StatusBadRequest); return")
		}
		w.Line("\tif buf.Len() != 0 { w.WriteHeader(http.StatusBadRequest); return }")
		w.Blank()
	}
	if api.Result.Name == "nil" {
		w.Linef("\tstatus := %s(r.Context()%s)", SnakeCase(api.Name), g.callArgNames(api.Args))
		w.Line("\tstatus = normalizeHandlerStatus(r.Context(), status)")
		w.Line("\tif !checkStatus(w, status) { return }")
		w.Line("\tw.WriteHeader(http.StatusOK)")
		w.Line("}")
		return
	}
	w.Linef("\tresult, status := %s(r.Context()%s)", SnakeCase(api.Name), g.callArgNames(api.Args))
	w.Line("\tstatus = normalizeHandlerStatus(r.Context(), status)")
	w.Line("\tif !checkStatus(w, status) { return }")
	w.Line("\tvar body bytes.Buffer")
	g.directWrite(w, api.Result, "result", "&body", "w.WriteHeader(http.StatusInternalServerError); return")
	w.Line("\tsendResponse(w, body.Bytes())")
	w.Line("}")
}
