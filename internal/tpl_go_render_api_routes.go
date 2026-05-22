package internal

func (g *GoGenerator) renderAPIRouteRegistrations(w *sourceWriter, groups map[string][]TplApi) {
	modules := orderedGroupKeys(groups)
	for _, module := range modules {
		funcName := "Register" + PascalCase(module) + "Api"
		if PascalCase(module) == "Api" {
			funcName = "RegisterApi"
		}
		w.Blank()
		w.Linef("func %s(mux *http.ServeMux, mws ...Middleware) {", funcName)
		w.Line("\tmw := composeMiddleware(mws...)")
		for _, api := range groups[module] {
			w.Linef("\tmux.HandleFunc(\"POST /%s\", mw(%sHandler))", api.Name, PascalCase(api.Name))
		}
		w.Line("}")
	}
	w.Blank()
}
