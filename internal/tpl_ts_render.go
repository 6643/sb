package internal

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func (g *TsGenerator) Generate(schema *TplSchema) error {
	targetDir := filepath.Join(g.Config.TsDir, "sb")
	dir, err := newGeneratedDir(targetDir)
	if err != nil {
		return err
	}
	if err := g.removeLegacyAPIWrapperStructFiles(targetDir); err != nil {
		return err
	}

	files := []generatedFile{
		{RelativePath: "type.ts", Data: []byte(renderTsRuntimeSource()), Perm: 0644},
		{RelativePath: "enum.ts", Data: []byte(g.renderEnumFile(schema.Enums)), Perm: 0644},
	}
	if len(schema.Enums) > 0 {
		files = append(files, generatedFile{RelativePath: "enum_smoke.test.ts", Data: []byte(g.renderEnumSmokeTest(schema.Enums)), Perm: 0644})
	}
	if err := dir.WriteAll(files...); err != nil {
		return err
	}

	structFiles := make([]string, 0, len(schema.Structs))
	for _, st := range schema.Structs {
		filename := "struct_" + SnakeCase(st.Name) + ".ts"
		structFiles = append(structFiles, strings.TrimSuffix(filename, ".ts"))
		if err := dir.Write(filename, []byte(g.renderStructFile(st)), 0644); err != nil {
			return err
		}
	}

	allFiles := append([]string{"enum"}, structFiles...)
	if err := dir.Write("_.ts", []byte(g.renderIndexFile(allFiles)), 0644); err != nil {
		return err
	}

	if len(schema.Apis) == 0 {
		return nil
	}
	return dir.WriteAll(
		generatedFile{RelativePath: "rpc.ts", Data: []byte(g.renderRPCFile(schema.Apis)), Perm: 0644},
		generatedFile{RelativePath: "rpc_smoke.test.ts", Data: []byte(g.renderSmokeTest(schema.Apis)), Perm: 0644},
	)
}

func (g *TsGenerator) removeLegacyAPIWrapperStructFiles(targetDir string) error {
	patterns := []string{
		filepath.Join(targetDir, "struct_api_*_req.ts"),
		filepath.Join(targetDir, "struct_api_*_resp.ts"),
	}
	for _, pattern := range patterns {
		matches, err := filepath.Glob(pattern)
		if err != nil {
			return err
		}
		for _, match := range matches {
			if err := os.Remove(match); err != nil && !os.IsNotExist(err) {
				return err
			}
		}
	}
	return nil
}

func (g *TsGenerator) renderEnumFile(enums []TplEnum) string {
	var w sourceWriter
	w.Line("import * as rt from \"./type\"")
	w.Blank()
	for _, enum := range enums {
		enumName := PascalCase(enum.Name)
		defaultName := enumName + PascalCase(enum.Children[0].Name)
		w.WriteLineComment("// ", enum.Note)
		w.Linef("export enum %s {", enumName)
		for _, child := range enum.Children {
			w.WriteLineComment("    // ", child.Note)
			w.Linef("    %s = %d,", PascalCase(child.Name), child.ID)
		}
		w.Line("}")
		w.Blank()
		w.Linef("export const Default%s = (): %s => %s.%s;", enumName, enumName, enumName, PascalCase(enum.Children[0].Name))
		w.Linef("export const Is%s = (v: %s): boolean => {", enumName, enumName)
		w.Line("    switch (v) {")
		for _, child := range enum.Children {
			w.Linef("    case %s.%s:", enumName, PascalCase(child.Name))
		}
		w.Line("        return true;")
		w.Line("    default:")
		w.Line("        return false;")
		w.Line("    }")
		w.Line("};")
		w.Blank()
		w.Linef("export const Normalize%s = (v: %s): %s => {", enumName, enumName, enumName)
		w.Linef("    if (Is%s(v)) return v;", enumName)
		w.Linef("    if ((v as any) === 0) return Default%s();", enumName)
		w.Line("    return v;")
		w.Line("};")
		w.Linef("export const IsDefault%s = (v: %s): boolean => Normalize%s(v) === Default%s();", enumName, enumName, enumName, enumName)
		w.Linef("export const IsAssignable%s = (v: %s): boolean => Is%s(v) || (((v as any) === 0) && !Is%s(0 as any));", enumName, enumName, enumName, enumName)
		w.Linef("export const eq%sValue = (a: %s, b: %s): boolean => Normalize%s(a) === Normalize%s(b);", enumName, enumName, enumName, enumName, enumName)
		w.Linef("export const eq%sList = (a: %s[], b: %s[]): boolean => rt.eqList(a, b, eq%sValue);", enumName, enumName, enumName, enumName)
		w.Blank()
		w.Linef("export const get%sListBody = (buf: rt.Buffer, state: number): [%s[], rt.Err] => {", enumName, enumName)
		w.Linef("    const [list, err] = rt.getDefaultList<%s>(", enumName)
		w.Line("        buf,")
		w.Line("        state,")
		w.Linef("        () => Default%s(),", enumName)
		w.Line("        (buf) => {")
		w.Line("            const [value, err] = rt.getU8(buf);")
		w.Linef("            if (err !== null) return [Default%s(), rt.errU(err)];", enumName)
		w.Linef("            const item = value as %s;", enumName)
		w.Linef("            if (!Is%s(item)) return [Default%s(), new Error(`非法枚举值: ${item}`)];", enumName, enumName)
		w.Line("            return [item, undefined];")
		w.Line("        },")
		w.Line("    );")
		w.Line("    return [list, err];")
		w.Line("};")
		w.Blank()
		w.Linef("export const set%sListBody = (buf: rt.Buffer, state: number, v: %s[]): rt.Err => {", enumName, enumName)
		w.Linef("    return rt.setDefaultList<%s>(", enumName)
		w.Line("        buf,")
		w.Line("        state,")
		w.Line("        v,")
		w.Linef("        (item) => IsDefault%s(item),", enumName)
		w.Line("        (buf, item) => {")
		w.Linef("            if (!Is%s(item)) return new Error(`非法枚举值: ${item}`);", enumName)
		w.Line("            return rt.setU8(buf, item as any);")
		w.Line("        },")
		w.Line("    );")
		w.Line("};")
		w.Blank()
		_ = defaultName
	}
	return w.String()
}

func (g *TsGenerator) renderStructFile(st TplStruct) string {
	name := PascalCase(st.Name)
	headerBits := g.headerBits(st)
	headerWidthsName := g.structHeaderWidthsName(st)

	var w sourceWriter
	w.Line("import * as rt from \"./type\"")
	w.Line("import * as _ from \"./_\"")
	w.Blank()
	w.WriteLineComment("// ", st.Note)
	w.Linef("export interface %s extends rt.Serializable, rt.Deserializable {", name)
	for _, field := range st.Fields {
		w.WriteLineComment("    // ", field.Note)
		w.Linef("    %s: %s;", CamelCase(field.Name), g.fieldType(field.Type))
	}
	w.Line("}")
	w.Blank()
	w.Linef("const %s = %s;", headerWidthsName, g.structHeaderWidthsLiteral(st))
	w.Blank()
	w.Linef("export const new%s = (): %s => {", name, name)
	w.Line("    const s = {")
	for _, field := range st.Fields {
		w.Linef("        %s: %s,", CamelCase(field.Name), g.defaultValue(field.Type))
	}
	w.Linef("    } as any as %s;", name)
	w.Linef("    s.set = (buf: rt.Buffer) => set%s(buf, s);", name)
	w.Line("    s.get = (buf: rt.Buffer) => {")
	w.Linef("        const [next, err] = get%s(buf);", name)
	w.Line("        if (err === undefined) Object.assign(s, next);")
	w.Line("        return err;")
	w.Line("    };")
	w.Line("    return s;")
	w.Line("};")
	w.Blank()
	w.Linef("export const isZero%s = (s: %s | null | undefined): boolean => {", name, name)
	w.Line("    if (s === null || s === undefined) return true;")
	for _, field := range st.Fields {
		ref := "s." + CamelCase(field.Name)
		w.Linef("    if (%s) return false;", g.nonDefaultExpr(field.Type, ref))
	}
	w.Line("    return true;")
	w.Line("};")
	w.Blank()
	w.Linef("export const validate%s = (s: %s | null | undefined): rt.Err => {", name, name)
	w.Line("    if (s === null || s === undefined) return undefined;")
	for _, field := range st.Fields {
		fieldName := PascalCase(field.Name)
		ref := "s." + CamelCase(field.Name)
		switch {
		case field.Type.Kind == TplKindEnum && !field.Type.IsList:
			w.Linef("    if (!_.IsAssignable%s(%s as any)) return new Error(`%s 非法枚举值: ${%s as any}`);", PascalCase(field.Type.Name), ref, fieldName, ref)
		case field.Type.Kind == TplKindEnum && field.Type.IsList:
			w.Linef("    for (let i = 0; i < %s.length; i++) {", ref)
			w.Linef("        if (!_.IsAssignable%s(%s[i] as any)) return new Error(`%s[${i}] 非法枚举值: ${%s[i] as any}`);", PascalCase(field.Type.Name), ref, fieldName, ref)
			w.Line("    }")
			w.Linef("    { const [, err] = rt.listCountState(%s.length); if (err !== null) return new Error(`%s: ${err.message}`); }", ref, fieldName)
		case field.Type.Kind == TplKindStruct && !field.Type.IsList:
			w.Linef("    { const err = _.validate%s(%s); if (err !== undefined) return new Error(`%s: ${err.message}`); }", PascalCase(field.Type.Name), ref, fieldName)
		case field.Type.Kind == TplKindStruct && field.Type.IsList:
			w.Linef("    { const [, err] = rt.listCountState(%s.length); if (err !== null) return new Error(`%s: ${err.message}`); }", ref, fieldName)
			w.Linef("    for (let i = 0; i < %s.length; i++) {", ref)
			w.Linef("        const err = _.validate%s(%s[i]);", PascalCase(field.Type.Name), ref)
			w.Linef("        if (err !== undefined) return new Error(`%s[${i}]: ${err.message}`);", fieldName)
			w.Line("    }")
		case field.Type.Kind == TplKindBase && !field.Type.IsList && field.Type.Name == "text":
			w.Linef("    { const [, err] = rt.textState(%s); if (err !== null) return new Error(`%s: ${err.message}`); }", ref, fieldName)
		case field.Type.Kind == TplKindBase && !field.Type.IsList && field.Type.Name == "bin":
			w.Linef("    { const [, err] = rt.binState(%s.byteLength); if (err !== null) return new Error(`%s: ${err.message}`); }", ref, fieldName)
		case field.Type.Kind == TplKindBase && field.Type.IsList && field.Type.Name == "text":
			w.Linef("    { const [, err] = rt.listCountState(%s.length); if (err !== null) return new Error(`%s: ${err.message}`); }", ref, fieldName)
			w.Linef("    for (let i = 0; i < %s.length; i++) {", ref)
			w.Linef("        const [, err] = rt.textState(%s[i]);", ref)
			w.Linef("        if (err !== null) return new Error(`%s[${i}]: ${err.message}`);", fieldName)
			w.Line("    }")
		case field.Type.Kind == TplKindBase && field.Type.IsList && field.Type.Name == "bin":
			w.Linef("    { const [, err] = rt.listCountState(%s.length); if (err !== null) return new Error(`%s: ${err.message}`); }", ref, fieldName)
			w.Linef("    for (let i = 0; i < %s.length; i++) {", ref)
			w.Linef("        const [, err] = rt.binState(%s[i].byteLength);", ref)
			w.Linef("        if (err !== null) return new Error(`%s[${i}]: ${err.message}`);", fieldName)
			w.Line("    }")
		case field.Type.Kind == TplKindBase && field.Type.IsList:
			w.Linef("    { const [, err] = rt.listCountState(%s.length); if (err !== null) return new Error(`%s: ${err.message}`); }", ref, fieldName)
		}
	}
	w.Line("    return undefined;")
	w.Line("};")
	w.Blank()
	w.Linef("export const get%s = (buf: rt.Buffer): [%s, rt.Err] => {", name, name)
	w.Linef("    const s = new%s();", name)
	w.Linef("    const headerBits = %d;", headerBits)
	w.Line("    const [header, errHeader] = buf.read(rt.headerSize(headerBits));")
	w.Line("    if (errHeader !== null) return [s, new Error(`not enough data`)];")
	w.Linef("    const [headerStates, errHeaderState] = rt.readHeader(header, %s, \"%s header\");", headerWidthsName, name)
	w.Line("    if (errHeaderState !== undefined) return [s, errHeaderState];")
	for i, field := range st.Fields {
		fieldName := CamelCase(PascalCase(field.Name))
		switch g.tagWidth(field.Type) {
		case 1:
			if field.Type.Name == "bool" {
				w.Linef("    const %sState = headerStates[%d] === 1;", fieldName, i)
			} else {
				w.Linef("    const %sPresent = headerStates[%d] === 1;", fieldName, i)
			}
		case 2:
			w.Linef("    const %sState = headerStates[%d];", fieldName, i)
		}
	}
	for _, field := range st.Fields {
		fieldName := PascalCase(field.Name)
		fieldVar := "s." + CamelCase(field.Name)
		stateVar := CamelCase(fieldName)
		switch {
		case field.Type.Name == "bool":
			w.Linef("    %s = %sState;", fieldVar, stateVar)
		case field.Type.Kind == TplKindBase && !field.Type.IsList && field.Type.Name == "text":
			w.Linef("    { const [value, err] = rt.getText(buf, %sState); if (err !== null) return [s, new Error(`get %s %s: ${err.message}`)]; %s = value; }", stateVar, name, fieldName, fieldVar)
		case field.Type.Kind == TplKindBase && !field.Type.IsList && field.Type.Name == "bin":
			w.Linef("    { const [value, err] = rt.getBin(buf, %sState); if (err !== null) return [s, new Error(`get %s %s: ${err.message}`)]; %s = value; }", stateVar, name, fieldName, fieldVar)
		case field.Type.Kind == TplKindBase && field.Type.IsList:
			if field.Type.Name == "bool" || field.Type.Name == "text" || field.Type.Name == "bin" {
				w.Linef("    { const [value, err] = %s; if (err !== null) return [s, new Error(`get %s %s: ${err.message}`)]; %s = value; }", g.listGetExpr(field.Type, stateVar+"State"), name, fieldName, fieldVar)
			} else {
				w.Linef("    { const [value, err] = %s; if (err !== undefined) return [s, new Error(`get %s %s: ${err.message}`)]; %s = value; }", g.listGetExpr(field.Type, stateVar+"State"), name, fieldName, fieldVar)
			}
		case field.Type.Kind == TplKindStruct && field.Type.IsList:
			w.Linef("    { const [value, err] = _.get%sListBody(buf, %sState); if (err !== undefined) return [s, new Error(`get %s %s: ${err.message}`)]; %s = value; }", PascalCase(field.Type.Name), stateVar, name, fieldName, fieldVar)
		case field.Type.Kind == TplKindEnum && field.Type.IsList:
			w.Linef("    { const [value, err] = _.get%sListBody(buf, %sState); if (err !== undefined) return [s, new Error(`get %s %s: ${err.message}`)]; %s = value; }", PascalCase(field.Type.Name), stateVar, name, fieldName, fieldVar)
		case field.Type.Kind == TplKindStruct:
			w.Linef("    if (%sPresent) {", stateVar)
			w.Linef("        const [value, err] = _.read%s(buf);", PascalCase(field.Type.Name))
			w.Linef("        if (err !== undefined) return [s, new Error(`get %s %s: ${err.message}`)];", name, fieldName)
			w.Linef("        %s = value;", fieldVar)
			w.Line("    }")
		case field.Type.Kind == TplKindEnum:
			w.Linef("    %s = _.Default%s();", fieldVar, PascalCase(field.Type.Name))
			w.Linef("    if (%sPresent) {", stateVar)
			w.Line("        const [value, err] = rt.getU8(buf);")
			w.Linef("        if (err !== null) return [s, new Error(`get %s %s: ${err.message}`)];", name, fieldName)
			w.Linef("        const item = value as _.%s;", PascalCase(field.Type.Name))
			w.Linef("        if (!_.Is%s(item)) return [s, new Error(`get %s %s: 非法枚举值: ${item}`)];", PascalCase(field.Type.Name), name, fieldName)
			w.Linef("        %s = item;", fieldVar)
			w.Line("    }")
		default:
			w.Linef("    if (%sPresent) {", stateVar)
			w.Linef("        const [value, err] = rt.%s(buf);", g.primitiveGetter(field.Type.Name))
			w.Linef("        if (err !== null) return [s, new Error(`get %s %s: ${err.message}`)];", name, fieldName)
			w.Linef("        %s = value as any;", fieldVar)
			w.Line("    }")
		}
	}
	w.Linef("    const errValidate = validate%s(s);", name)
	w.Line("    if (errValidate !== undefined) return [s, new Error(`validate failed: ${errValidate.message}`)];")
	w.Line("    return [s, undefined];")
	w.Line("};")
	w.Blank()
	w.Linef("export const set%s = (buf: rt.Buffer, s: %s): rt.Err => {", name, name)
	w.Linef("    if (s === null || s === undefined) return new Error(`set %s: value is null or undefined`);", name)
	w.Linef("    const errValidate = validate%s(s);", name)
	w.Linef("    if (errValidate !== undefined) return new Error(`validate %s: ${errValidate.message}`);", name)
	w.Line("    const startOffset = buf.writeOffset;")
	for _, field := range st.Fields {
		if g.tagWidth(field.Type) != 2 {
			continue
		}
		fieldName := PascalCase(field.Name)
		ref := "s." + CamelCase(field.Name)
		switch {
		case field.Type.Kind == TplKindBase && !field.Type.IsList && field.Type.Name == "text":
			w.Linef("    const [%sState, err%sState] = rt.textState(%s);", CamelCase(fieldName), fieldName, ref)
		case field.Type.Kind == TplKindBase && !field.Type.IsList && field.Type.Name == "bin":
			w.Linef("    const [%sState, err%sState] = rt.binState(%s.byteLength);", CamelCase(fieldName), fieldName, ref)
		default:
			w.Linef("    const [%sState, err%sState] = rt.listCountState(%s.length);", CamelCase(fieldName), fieldName, ref)
		}
		w.Linef("    if (err%sState !== null) return err%sState;", fieldName, fieldName)
	}
	w.Line("    const headerStates = [")
	for _, field := range st.Fields {
		fieldName := PascalCase(field.Name)
		ref := "s." + CamelCase(field.Name)
		switch {
		case field.Type.Name == "bool":
			w.Linef("        %s ? 1 : 0,", ref)
		case g.tagWidth(field.Type) == 1:
			w.Linef("        %s ? 1 : 0,", g.nonDefaultExpr(field.Type, ref))
		default:
			w.Linef("        %sState,", CamelCase(fieldName))
		}
	}
	w.Line("    ];")
	w.Linef("    const [header, errHeader] = rt.writeHeader(%s, headerStates);", headerWidthsName)
	w.Line("    if (errHeader !== undefined) { buf.rewindWrite(startOffset); return new Error(`set header: ${errHeader.message}`); }")
	w.Line("    const errHeaderWrite = buf.write(header);")
	w.Line("    if (errHeaderWrite !== null) { buf.rewindWrite(startOffset); return errHeaderWrite; }")
	for _, field := range st.Fields {
		fieldName := PascalCase(field.Name)
		ref := "s." + CamelCase(field.Name)
		switch {
		case field.Type.Name == "bool":
		case field.Type.Kind == TplKindBase && !field.Type.IsList && field.Type.Name == "text":
			w.Linef("    { const err = rt.setText(buf, %sState, %s); if (err !== null) { buf.rewindWrite(startOffset); return new Error(`set %s %s: ${err.message}`); } }", CamelCase(fieldName), ref, name, fieldName)
		case field.Type.Kind == TplKindBase && !field.Type.IsList && field.Type.Name == "bin":
			w.Linef("    { const err = rt.setBin(buf, %sState, %s); if (err !== null) { buf.rewindWrite(startOffset); return new Error(`set %s %s: ${err.message}`); } }", CamelCase(fieldName), ref, name, fieldName)
		case field.Type.Kind == TplKindBase && field.Type.IsList:
			if field.Type.Name == "bool" || field.Type.Name == "text" || field.Type.Name == "bin" {
				w.Linef("    { const err = %s; if (err !== null) { buf.rewindWrite(startOffset); return new Error(`set %s %s: ${err.message}`); } }", g.listSetExpr(field.Type, ref, CamelCase(fieldName)+"State"), name, fieldName)
			} else {
				w.Linef("    { const err = %s; if (err !== undefined) { buf.rewindWrite(startOffset); return new Error(`set %s %s: ${err.message}`); } }", g.listSetExpr(field.Type, ref, CamelCase(fieldName)+"State"), name, fieldName)
			}
		case field.Type.Kind == TplKindStruct && field.Type.IsList:
			w.Linef("    { const err = _.set%sListBody(buf, %sState, %s); if (err !== undefined) { buf.rewindWrite(startOffset); return new Error(`set %s %s: ${err.message}`); } }", PascalCase(field.Type.Name), CamelCase(fieldName), ref, name, fieldName)
		case field.Type.Kind == TplKindEnum && field.Type.IsList:
			w.Linef("    { const err = _.set%sListBody(buf, %sState, %s); if (err !== undefined) { buf.rewindWrite(startOffset); return new Error(`set %s %s: ${err.message}`); } }", PascalCase(field.Type.Name), CamelCase(fieldName), ref, name, fieldName)
		case field.Type.Kind == TplKindStruct:
			w.Linef("    if (!_.isZero%s(%s)) {", PascalCase(field.Type.Name), ref)
			w.Linef("        const err = _.set%s(buf, %s!);", PascalCase(field.Type.Name), ref)
			w.Linef("        if (err !== undefined) { buf.rewindWrite(startOffset); return new Error(`set %s %s: ${err.message}`); }", name, fieldName)
			w.Line("    }")
		case field.Type.Kind == TplKindEnum:
			w.Linef("    if (!_.IsDefault%s(%s as any)) {", PascalCase(field.Type.Name), ref)
			w.Linef("        const err = rt.setU8(buf, _.Normalize%s(%s as any) as any);", PascalCase(field.Type.Name), ref)
			w.Linef("        if (err !== null) { buf.rewindWrite(startOffset); return new Error(`set %s %s: ${err.message}`); }", name, fieldName)
			w.Line("    }")
		default:
			w.Linef("    if (%s) {", g.nonDefaultExpr(field.Type, ref))
			w.Linef("        const err = rt.%s(buf, %s as any);", g.primitiveSetter(field.Type.Name), ref)
			w.Linef("        if (err !== null) { buf.rewindWrite(startOffset); return new Error(`set %s %s: ${err.message}`); }", name, fieldName)
			w.Line("    }")
		}
	}
	w.Line("    return undefined;")
	w.Line("};")
	w.Blank()
	w.Linef("export const read%s = (buf: rt.Buffer): [%s, rt.Err] => get%s(buf);", name, name, name)
	w.Blank()
	w.Linef("export const eq%s = (a: %s | null | undefined, b: %s | null | undefined): boolean => {", name, name, name)
	w.Linef("    if (isZero%s(a as any) && isZero%s(b as any)) return true;", name, name)
	w.Line("    if (a === null || a === undefined || b === null || b === undefined) return false;")
	for _, field := range st.Fields {
		fieldName := CamelCase(field.Name)
		w.Linef("    if (!%s) return false;", g.eqExpr(field.Type, "a."+fieldName, "b."+fieldName))
	}
	w.Line("    return true;")
	w.Line("};")
	w.Blank()
	g.renderStructListBodyHelpers(&w, st)
	return w.String()
}

func (g *TsGenerator) renderStructListBodyHelpers(w *sourceWriter, st TplStruct) {
	name := PascalCase(st.Name)
	w.Linef("export const get%sListBody = (buf: rt.Buffer, state: number): [%s[], rt.Err] => {", name, name)
	w.Linef("    const [list, err] = rt.getDefaultList<%s>(", name)
	w.Line("        buf,")
	w.Line("        state,")
	w.Linef("        () => new%s(),", name)
	w.Linef("        (buf) => read%s(buf),", name)
	w.Line("    );")
	w.Line("    return [list, err];")
	w.Line("};")
	w.Blank()
	w.Linef("export const set%sListBody = (buf: rt.Buffer, state: number, v: %s[]): rt.Err => {", name, name)
	w.Linef("    return rt.setDefaultList<%s>(", name)
	w.Line("        buf,")
	w.Line("        state,")
	w.Line("        v,")
	w.Linef("        (item) => isZero%s(item),", name)
	w.Linef("        (buf, item) => set%s(buf, item),", name)
	w.Line("    );")
	w.Line("}")
}

func (g *TsGenerator) renderIndexFile(files []string) string {
	var w sourceWriter
	w.Line("export * from \"./type\"")
	for _, file := range files {
		w.Linef("export * from \"./%s\"", file)
	}
	return w.String()
}

func (g *TsGenerator) structHeaderWidthsName(st TplStruct) string {
	return CamelCase(st.Name) + "HeaderWidths"
}

func (g *TsGenerator) structHeaderWidthsLiteral(st TplStruct) string {
	parts := make([]string, 0, len(st.Fields))
	for _, field := range st.Fields {
		parts = append(parts, fmt.Sprintf("%d", g.tagWidth(field.Type)))
	}
	return "[" + strings.Join(parts, ", ") + "] as const"
}

func (g *TsGenerator) renderRPCFile(apis []TplApi) string {
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
	w.Line("export type RpcStatus = RpcErrCode | number;")
	w.Blank()
	w.Line("export interface RpcConfig {")
	w.Line("    host: string;")
	w.Line("    headers?: Record<string, string>;")
	w.Line("    timeout?: number;")
	w.Line("    retries?: number;")
	w.Line("    maxRespBytes?: number;")
	w.Line("}")
	w.Blank()
	w.Line("export interface RpcClient {")
	w.Line("    setHeader: (key: string, value: string) => void;")
	w.Line("    getHeader: (key: string) => string | undefined;")
	w.Line("    removeHeader: (key: string) => void;")
	w.Line("    setAuthorization: (token: string) => void;")
	w.Line("    getAuthorization: () => string | undefined;")
	w.Line("    removeAuthorization: () => void;")
	w.Line("    isAuthorized: () => boolean;")
	for _, api := range apis {
		methodName := CamelCase(api.Name)
		promiseType := strings.ReplaceAll(g.rpcPromiseType(api.Result), "RpcErrCode", "RpcStatus")
		w.Linef("    %s: (%s) => Promise<%s>;", methodName, g.rpcArgList(api.Args), promiseType)
	}
	w.Line("}")
	w.Blank()
	w.Line("type _RpcClientState = RpcClient & {")
	w.Line("    config: RpcConfig;")
	w.Line("    headers: Record<string, string>;")
	w.Line("    timeout: number;")
	w.Line("    retries: number;")
	w.Line("    maxRespBytes: number;")
	w.Line("    _fetch: (path: string, body: Uint8Array) => Promise<[Uint8Array | null, RpcStatus]>;")
	w.Line("};")
	w.Line("type RpcClientCtor = new (config: RpcConfig) => RpcClient;")
	w.Blank()
	w.Line("function _RpcClientCtor(this: _RpcClientState, config: RpcConfig): void {")
	w.Line("    this.config = { ...config, host: config.host.replace(/\\/+$/, \"\") };")
	w.Line("    this.headers = config.headers ? { ...config.headers } : {};")
	w.Line("    const cfgTimeout = config.timeout;")
	w.Line("    this.timeout = cfgTimeout !== undefined && Number.isFinite(cfgTimeout) && cfgTimeout >= 0 ? Math.min(Math.floor(cfgTimeout), maxTimeoutMs) : 5000;")
	w.Line("    const cfgRetries = config.retries;")
	w.Line("    this.retries = cfgRetries !== undefined && Number.isFinite(cfgRetries) && cfgRetries >= 0 ? Math.floor(cfgRetries) : 3;")
	w.Line("    const cfgMaxRespBytes = config.maxRespBytes;")
	w.Line("    this.maxRespBytes = cfgMaxRespBytes !== undefined && Number.isFinite(cfgMaxRespBytes) && cfgMaxRespBytes > 0 ? Math.min(Math.floor(cfgMaxRespBytes), maxSafeRespBytes) : defaultMaxRespBytes;")
	w.Line("}")
	w.Line("const _rpcClientProto = _RpcClientCtor.prototype as _RpcClientState;")
	w.Blank()
	w.Line("_rpcClientProto.setHeader = function(this: _RpcClientState, key: string, value: string): void { this.headers[key] = value; };")
	w.Line("_rpcClientProto.getHeader = function(this: _RpcClientState, key: string): string | undefined { return this.headers[key]; };")
	w.Line("_rpcClientProto.removeHeader = function(this: _RpcClientState, key: string): void { delete this.headers[key]; };")
	w.Blank()
	w.Line("_rpcClientProto.setAuthorization = function(this: _RpcClientState, token: string): void { this.setHeader(\"Authorization\", `Bearer ${token}`); };")
	w.Line("_rpcClientProto.getAuthorization = function(this: _RpcClientState): string | undefined { return this.getHeader(\"Authorization\"); };")
	w.Line("_rpcClientProto.removeAuthorization = function(this: _RpcClientState): void { this.removeHeader(\"Authorization\"); };")
	w.Line("_rpcClientProto.isAuthorized = function(this: _RpcClientState): boolean { return !!this.getAuthorization(); };")
	w.Blank()
	w.Line("async function readResponseBytes(res: Response, maxRespBytes: number): Promise<[Uint8Array | null, RpcErrCode | null]> {")
	w.Line("    const contentLength = res.headers.get(\"content-length\");")
	w.Line("    if (contentLength !== null) {")
	w.Line("        const size = Number(contentLength);")
	w.Line("        if (Number.isFinite(size) && size > maxRespBytes) return [null, RpcErrCode.RespErr];")
	w.Line("    }")
	w.Blank()
	w.Line("    if (res.body !== null) {")
	w.Line("        const reader = res.body.getReader();")
	w.Line("        const chunks: Uint8Array[] = [];")
	w.Line("        let total = 0;")
	w.Line("        while (true) {")
	w.Line("            const { done, value } = await reader.read();")
	w.Line("            if (done) break;")
	w.Line("            if (value === undefined) continue;")
	w.Line("            total += value.byteLength;")
	w.Line("            if (total > maxRespBytes) {")
	w.Line("                try {")
	w.Line("                    await reader.cancel();")
	w.Line("                } catch {")
	w.Line("                }")
	w.Line("                return [null, RpcErrCode.RespErr];")
	w.Line("            }")
	w.Line("            chunks.push(value);")
	w.Line("        }")
	w.Blank()
	w.Line("        const bytes = new Uint8Array(total);")
	w.Line("        let offset = 0;")
	w.Line("        for (const chunk of chunks) {")
	w.Line("            bytes.set(chunk, offset);")
	w.Line("            offset += chunk.byteLength;")
	w.Line("        }")
	w.Line("        return [bytes, null];")
	w.Line("    }")
	w.Blank()
	w.Line("    const bytes = new Uint8Array(await res.arrayBuffer());")
	w.Line("    if (bytes.byteLength > maxRespBytes) return [null, RpcErrCode.RespErr];")
	w.Line("    return [bytes, null];")
	w.Line("}")
	w.Blank()
	w.Line("_rpcClientProto._fetch = async function(this: _RpcClientState, path: string, body: Uint8Array): Promise<[Uint8Array | null, RpcStatus]> {")
	w.Line("    let lastStatus = RpcErrCode.NoConn;")
	w.Line("    for (let i = 0; i <= this.retries; i++) {")
	w.Line("        if (i > 0) await new Promise((res) => setTimeout(res, i * 1000));")
	w.Line("        const controller = new AbortController();")
	w.Line("        let timeoutId: ReturnType<typeof setTimeout> | null = null;")
	w.Line("        if (this.timeout > 0) timeoutId = setTimeout(() => controller.abort(), this.timeout);")
	w.Line("        try {")
	w.Line("            const res = await fetch(`${this.config.host}/${path}`, {")
	w.Line("                method: \"POST\",")
	w.Line("                headers: { \"Content-Type\": \"application/octet-stream\", ...this.headers },")
	w.Line("                body: body as any,")
	w.Line("                signal: controller.signal,")
	w.Line("            });")
	w.Line("            if (res.ok) {")
	w.Line("                const [bytes, readErr] = await readResponseBytes(res, this.maxRespBytes);")
	w.Line("                if (readErr !== null) return [null, readErr];")
	w.Line("                return [bytes, RpcErrCode.Ok];")
	w.Line("            }")
	w.Line("            lastStatus = res.status;")
	w.Line("            if (res.status === 408 && i < this.retries) continue;")
	w.Line("            return [null, res.status];")
	w.Line("        } catch (e: any) {")
	w.Line("            lastStatus = e && e.name === \"AbortError\" ? RpcErrCode.Timeout : RpcErrCode.NoConn;")
	w.Line("            if (i < this.retries) continue;")
	w.Line("        } finally {")
	w.Line("            if (timeoutId !== null) clearTimeout(timeoutId);")
	w.Line("        }")
	w.Line("    }")
	w.Line("    return [null, lastStatus];")
	w.Line("};")
	for _, api := range apis {
		methodName := CamelCase(api.Name)
		defaultVal := g.rpcDefaultValue(api.Result)
		promiseType := strings.ReplaceAll(g.rpcPromiseType(api.Result), "RpcErrCode", "RpcStatus")
		w.Blank()
		w.WriteLineComment("// ", api.Note)
		w.Linef("_rpcClientProto.%s = async function(this: _RpcClientState, %s): Promise<%s> {", methodName, g.rpcArgList(api.Args), promiseType)
		w.Line("    const buf = new _.Buffer();")
		for _, arg := range api.Args {
			g.tsDirectWrite(&w, arg.Type, CamelCase(arg.Name), "buf", g.rpcReqErrReturn(api.Result, defaultVal))
		}
		w.Linef("    const [bytes, status] = await this._fetch(\"%s\", buf.bytes);", api.Name)
		w.Linef("    if (status !== RpcErrCode.Ok || bytes === null) return %s;", g.rpcStatusReturn(api.Result, defaultVal))
		if api.Result.Name == "nil" {
			w.Line("    if (bytes.byteLength !== 0) return RpcErrCode.RespErr;")
			w.Line("    return RpcErrCode.Ok;")
			w.Line("};")
			continue
		}
		w.Line("    const respBuf = new _.Buffer(bytes);")
		w.Linef("    let result = %s as any;", defaultVal)
		g.tsDirectRead(&w, api.Result, "result", "respBuf", fmt.Sprintf("[%s, RpcErrCode.RespErr]", defaultVal))
		w.Linef("    if (respBuf.len !== 0) return [%s, RpcErrCode.RespErr];", defaultVal)
		w.Line("    return [result as any, RpcErrCode.Ok];")
		w.Line("};")
	}
	w.Blank()
	w.Line("export const RpcClient = _RpcClientCtor as unknown as RpcClientCtor;")
	return w.String()
}

func (g *TsGenerator) renderSmokeTest(apis []TplApi) string {
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

func (g *TsGenerator) renderEnumSmokeTest(enums []TplEnum) string {
	var w sourceWriter
	w.Line("import { describe, expect, test } from \"bun:test\";")
	w.Blank()
	w.Line("import * as _ from \"./_\";")
	w.Blank()
	w.Line("describe(\"enum smoke\", () => {")
	for _, enum := range enums {
		if len(enum.Children) == 0 {
			continue
		}
		first := PascalCase(enum.Children[0].Name)
		last := PascalCase(enum.Children[len(enum.Children)-1].Name)
		enumName := PascalCase(enum.Name)
		w.Linef("    test(\"%s validator accepts generated values\", () => {", enumName)
		w.Linef("        expect(_.Is%s(_.%s.%s)).toBe(true);", enumName, enumName, first)
		w.Linef("        expect(_.Is%s(_.%s.%s)).toBe(true);", enumName, enumName, last)
		w.Linef("        expect(_.Is%s(255 as any)).toBe(false);", enumName)
		w.Line("    });")
	}
	w.Line("});")
	return w.String()
}

func (g *TsGenerator) fieldType(t TplType) string {
	base := g.getTsType(TplType{Name: t.Name, Kind: t.Kind})
	switch t.Kind {
	case TplKindEnum, TplKindStruct:
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

func (g *TsGenerator) defaultValue(t TplType) string {
	if t.IsList {
		return "[]"
	}
	switch t.Kind {
	case TplKindBase:
		return g.getTsValue(t.Name)
	case TplKindEnum:
		return fmt.Sprintf("_.Default%s()", PascalCase(t.Name))
	default:
		return "null"
	}
}

func (g *TsGenerator) headerBits(st TplStruct) int {
	total := 0
	for _, field := range st.Fields {
		total += g.tagWidth(field.Type)
	}
	return total
}

func (g *TsGenerator) tagWidth(t TplType) int {
	if t.Name == "bool" && !t.IsList {
		return 1
	}
	if t.Kind == TplKindBase && (t.Name == "text" || t.Name == "bin" || t.IsList) {
		return 2
	}
	if t.IsList {
		return 2
	}
	return 1
}

func (g *TsGenerator) nonDefaultExpr(t TplType, ref string) string {
	switch {
	case t.IsList:
		return fmt.Sprintf("%s.length !== 0", ref)
	case t.Kind == TplKindStruct:
		return fmt.Sprintf("!_.isZero%s(%s)", PascalCase(t.Name), ref)
	case t.Kind == TplKindEnum:
		return fmt.Sprintf("!_.IsDefault%s(%s as any)", PascalCase(t.Name), ref)
	case t.Name == "text":
		return fmt.Sprintf("%s !== \"\"", ref)
	case t.Name == "bin":
		return fmt.Sprintf("%s.byteLength !== 0", ref)
	case t.Name == "bool":
		return ref
	case t.Name == "i64", t.Name == "u64":
		return fmt.Sprintf("%s !== 0n", ref)
	default:
		return fmt.Sprintf("%s !== 0", ref)
	}
}

func (g *TsGenerator) primitiveGetter(name string) string {
	switch name {
	case "i8":
		return "getI8"
	case "u8":
		return "getU8"
	case "i16":
		return "getI16"
	case "u16":
		return "getU16"
	case "i32":
		return "getI32"
	case "u32":
		return "getU32"
	case "i64":
		return "getI64"
	case "u64":
		return "getU64"
	case "f32":
		return "getF32"
	case "f64":
		return "getF64"
	default:
		return ""
	}
}

func (g *TsGenerator) primitiveSetter(name string) string {
	switch name {
	case "i8":
		return "setI8"
	case "u8":
		return "setU8"
	case "i16":
		return "setI16"
	case "u16":
		return "setU16"
	case "i32":
		return "setI32"
	case "u32":
		return "setU32"
	case "i64":
		return "setI64"
	case "u64":
		return "setU64"
	case "f32":
		return "setF32"
	case "f64":
		return "setF64"
	default:
		return ""
	}
}

func (g *TsGenerator) primitiveDefault(name string) string {
	switch name {
	case "i64", "u64":
		return "0n"
	default:
		return g.getTsValue(name)
	}
}

func (g *TsGenerator) primitiveEq(name string) string {
	switch name {
	case "i8":
		return "rt.eqI8"
	case "u8":
		return "rt.eqU8"
	case "i16":
		return "rt.eqI16"
	case "u16":
		return "rt.eqU16"
	case "i32":
		return "rt.eqI32"
	case "u32":
		return "rt.eqU32"
	case "i64":
		return "rt.eqI64"
	case "u64":
		return "rt.eqU64"
	case "f32":
		return "rt.eqF32"
	case "f64":
		return "rt.eqF64"
	case "bool":
		return "rt.eqBool"
	case "text":
		return "rt.eqText"
	case "bin":
		return "rt.eqBin"
	default:
		return ""
	}
}

func (g *TsGenerator) listGetExpr(t TplType, stateVar string) string {
	switch t.Name {
	case "bool":
		return fmt.Sprintf("rt.getBoolList(buf, %s)", stateVar)
	case "text":
		return fmt.Sprintf("rt.getTextList(buf, %s)", stateVar)
	case "bin":
		return fmt.Sprintf("rt.getBinList(buf, %s)", stateVar)
	default:
		if t.Kind == TplKindBase {
			return fmt.Sprintf("rt.getZeroList<%s>(buf, %s, %s, rt.%s)", g.getTsType(TplType{Name: t.Name, Kind: t.Kind}), stateVar, g.primitiveDefault(t.Name), g.primitiveGetter(t.Name))
		}
		return fmt.Sprintf("rt.getDefaultList<%s>(buf, %s, %s, %s)", g.getTsType(TplType{Name: t.Name, Kind: t.Kind}), stateVar, g.bitmapDefaultFactory(t), g.bitmapGetter(t))
	}
}

func (g *TsGenerator) listSetExpr(t TplType, ref, stateVar string) string {
	switch t.Name {
	case "bool":
		return fmt.Sprintf("rt.setBoolList(buf, %s, %s)", stateVar, ref)
	case "text":
		return fmt.Sprintf("rt.setTextList(buf, %s, %s)", stateVar, ref)
	case "bin":
		return fmt.Sprintf("rt.setBinList(buf, %s, %s)", stateVar, ref)
	default:
		if t.Kind == TplKindBase {
			return fmt.Sprintf("rt.setZeroList<%s>(buf, %s, %s, %s, rt.%s)", g.getTsType(TplType{Name: t.Name, Kind: t.Kind}), stateVar, ref, g.primitiveDefault(t.Name), g.primitiveSetter(t.Name))
		}
		return fmt.Sprintf("rt.setDefaultList<%s>(buf, %s, %s, %s, %s)", g.getTsType(TplType{Name: t.Name, Kind: t.Kind}), stateVar, ref, g.bitmapIsDefault(t), g.bitmapSetter(t))
	}
}

func (g *TsGenerator) bitmapDefaultFactory(t TplType) string {
	switch t.Kind {
	case TplKindEnum:
		return fmt.Sprintf("() => _.Default%s()", PascalCase(t.Name))
	case TplKindStruct:
		return fmt.Sprintf("() => _.new%s()", PascalCase(t.Name))
	default:
		return fmt.Sprintf("() => %s", g.primitiveDefault(t.Name))
	}
}

func (g *TsGenerator) bitmapIsDefault(t TplType) string {
	switch t.Kind {
	case TplKindEnum:
		return fmt.Sprintf("(item) => _.IsDefault%s(item as any)", PascalCase(t.Name))
	case TplKindStruct:
		return fmt.Sprintf("(item) => _.isZero%s(item)", PascalCase(t.Name))
	default:
		if t.Name == "i64" || t.Name == "u64" {
			return "(item) => item === 0n"
		}
		return "(item) => item === 0"
	}
}

func (g *TsGenerator) bitmapGetter(t TplType) string {
	switch t.Kind {
	case TplKindEnum:
		return fmt.Sprintf("(buf) => { const [value, err] = rt.getU8(buf); if (err !== null) return [_.Default%s(), rt.errU(err)]; const item = value as _.%s; if (!_.Is%s(item)) return [_.Default%s(), new Error(`非法枚举值: ${item}`)]; return [item, undefined]; }", PascalCase(t.Name), PascalCase(t.Name), PascalCase(t.Name), PascalCase(t.Name))
	case TplKindStruct:
		return fmt.Sprintf("(buf) => _.read%s(buf)", PascalCase(t.Name))
	default:
		return fmt.Sprintf("(buf) => rt.resultU(...rt.%s(buf))", g.primitiveGetter(t.Name))
	}
}

func (g *TsGenerator) bitmapSetter(t TplType) string {
	switch t.Kind {
	case TplKindEnum:
		return fmt.Sprintf("(buf, item) => { if (!_.Is%s(item as any)) return new Error(`非法枚举值: ${item}`); return rt.setU8(buf, item as any); }", PascalCase(t.Name))
	case TplKindStruct:
		return fmt.Sprintf("(buf, item) => _.set%s(buf, item)", PascalCase(t.Name))
	default:
		return fmt.Sprintf("(buf, item) => rt.errU(rt.%s(buf, item as any))", g.primitiveSetter(t.Name))
	}
}

func (g *TsGenerator) eqExpr(t TplType, left, right string) string {
	switch {
	case t.Kind == TplKindStruct && t.IsList:
		return fmt.Sprintf("rt.eqList(%s, %s, _.eq%s)", left, right, PascalCase(t.Name))
	case t.Kind == TplKindStruct:
		return fmt.Sprintf("_.eq%s(%s, %s)", PascalCase(t.Name), left, right)
	case t.Kind == TplKindEnum && t.IsList:
		return fmt.Sprintf("_.eq%sList(%s as any, %s as any)", PascalCase(t.Name), left, right)
	case t.Kind == TplKindEnum:
		return fmt.Sprintf("_.eq%sValue(%s as any, %s as any)", PascalCase(t.Name), left, right)
	case t.Kind == TplKindBase && t.IsList && t.Name == "bin":
		return fmt.Sprintf("rt.eqBinList(%s, %s)", left, right)
	case t.Kind == TplKindBase && t.IsList:
		return fmt.Sprintf("rt.eqList(%s, %s, %s)", left, right, g.primitiveEq(t.Name))
	default:
		return fmt.Sprintf("%s(%s, %s)", g.primitiveEq(t.Name), left, right)
	}
}

func (g *TsGenerator) rpcArgList(args []TplApiArg) string {
	parts := make([]string, 0, len(args))
	for _, arg := range args {
		prefix := ""
		if arg.Type.Kind != TplKindBase {
			prefix = "_."
		}
		parts = append(parts, fmt.Sprintf("%s: %s%s", CamelCase(arg.Name), prefix, g.getTsLogicType(arg.Type)))
	}
	return strings.Join(parts, ", ")
}

func (g *TsGenerator) rpcPromiseType(result TplType) string {
	if result.Name == "nil" {
		return "RpcErrCode"
	}
	prefix := ""
	if result.Kind != TplKindBase {
		prefix = "_."
	}
	return fmt.Sprintf("[%s%s, RpcErrCode]", prefix, g.getTsLogicType(result))
}

func (g *TsGenerator) rpcDefaultValue(result TplType) string {
	if result.IsList {
		return "[]"
	}
	switch result.Kind {
	case TplKindBase:
		return g.getTsValue(result.Name)
	case TplKindEnum:
		return fmt.Sprintf("_.Default%s()", PascalCase(result.Name))
	case TplKindStruct:
		return fmt.Sprintf("_.new%s()", PascalCase(result.Name))
	default:
		return "null"
	}
}

func (g *TsGenerator) rpcReqErrReturn(result TplType, defaultVal string) string {
	if result.Name == "nil" {
		return "RpcErrCode.ReqErr"
	}
	return fmt.Sprintf("[%s, RpcErrCode.ReqErr]", defaultVal)
}

func (g *TsGenerator) rpcStatusReturn(result TplType, defaultVal string) string {
	if result.Name == "nil" {
		return "status"
	}
	return fmt.Sprintf("[%s, status]", defaultVal)
}

func (g *TsGenerator) apiReqStruct(api TplApi) (TplStruct, bool) {
	if len(api.Args) == 0 {
		return TplStruct{}, false
	}
	fields := make([]TplStructField, 0, len(api.Args))
	for _, arg := range api.Args {
		fields = append(fields, TplStructField{Name: arg.Name, Type: arg.Type})
	}
	return TplStruct{Name: "api_" + SnakeCase(api.Name) + "_req", Fields: fields}, true
}

func (g *TsGenerator) apiRespStruct(api TplApi) (TplStruct, bool) {
	if api.Result.Name == "nil" {
		return TplStruct{}, false
	}
	return TplStruct{
		Name:   "api_" + SnakeCase(api.Name) + "_resp",
		Fields: []TplStructField{{Name: "result", Type: api.Result}},
	}, true
}

func (g *TsGenerator) apiReqTypeName(api TplApi) (string, bool) {
	req, ok := g.apiReqStruct(api)
	if !ok {
		return "", false
	}
	return PascalCase(req.Name), true
}

func (g *TsGenerator) apiRespTypeName(api TplApi) (string, bool) {
	resp, ok := g.apiRespStruct(api)
	if !ok {
		return "", false
	}
	return PascalCase(resp.Name), true
}

func (g *TsGenerator) tsDirectRead(w *sourceWriter, t TplType, target, bufVar, errReturn string) {
	w.Line("        {")
	switch {
	case t.Name == "bool" && !t.IsList:
		w.Linef("            const [value, err] = _.getU8(%s);", bufVar)
		w.Linef("            if (err !== null) return %s;", errReturn)
		w.Linef("            if (value > 1) return %s;", errReturn)
		w.Linef("            %s = (value !== 0) as any;", target)
	case t.Kind == TplKindBase && !t.IsList && t.Name == "text":
		w.Linef("            const [state, errState] = _.getU8(%s);", bufVar)
		w.Linef("            if (errState !== null) return %s;", errReturn)
		w.Linef("            const [value, err] = _.getText(%s, state);", bufVar)
		w.Linef("            if (err !== null) return %s;", errReturn)
		w.Linef("            %s = value as any;", target)
	case t.Kind == TplKindBase && !t.IsList && t.Name == "bin":
		w.Linef("            const [state, errState] = _.getU8(%s);", bufVar)
		w.Linef("            if (errState !== null) return %s;", errReturn)
		w.Linef("            const [value, err] = _.getBin(%s, state);", bufVar)
		w.Linef("            if (err !== null) return %s;", errReturn)
		w.Linef("            %s = value as any;", target)
	case t.IsList:
		w.Linef("            const [state, errState] = _.getU8(%s);", bufVar)
		w.Linef("            if (errState !== null) return %s;", errReturn)
		w.Linef("            const [value, err] = %s;", g.tsDirectListGetExpr(t, "state", bufVar))
		if t.Name == "bool" || t.Name == "text" || t.Name == "bin" {
			w.Linef("            if (err !== null) return %s;", errReturn)
		} else {
			w.Linef("            if (err !== undefined) return %s;", errReturn)
		}
		w.Linef("            %s = value as any;", target)
	case t.Kind == TplKindStruct:
		w.Linef("            const [value, err] = _.read%s(%s);", PascalCase(t.Name), bufVar)
		w.Linef("            if (err !== undefined) return %s;", errReturn)
		w.Linef("            %s = value as any;", target)
	case t.Kind == TplKindEnum:
		typeName := PascalCase(t.Name)
		w.Linef("            const [value, err] = _.getU8(%s);", bufVar)
		w.Linef("            if (err !== null) return %s;", errReturn)
		w.Linef("            const item = value as _.%s;", typeName)
		w.Linef("            if (!_.Is%s(item)) return %s;", typeName, errReturn)
		w.Linef("            %s = item as any;", target)
	default:
		w.Linef("            const [value, err] = _.%s(%s);", g.primitiveGetter(t.Name), bufVar)
		w.Linef("            if (err !== null) return %s;", errReturn)
		w.Linef("            %s = value as any;", target)
	}
	w.Line("        }")
}

func (g *TsGenerator) tsDirectWrite(w *sourceWriter, t TplType, ref, bufVar, errReturn string) {
	w.Line("        {")
	switch {
	case t.Name == "bool" && !t.IsList:
		w.Linef("            const encodeErr = _.setU8(%s, %s ? 1 : 0);", bufVar, ref)
		w.Linef("            if (encodeErr !== null) return %s;", errReturn)
	case t.Kind == TplKindBase && !t.IsList && t.Name == "text":
		w.Linef("            const [state, errState] = _.textState(%s);", ref)
		w.Linef("            if (errState !== null) return %s;", errReturn)
		w.Linef("            const err0 = _.setU8(%s, state);", bufVar)
		w.Linef("            if (err0 !== null) return %s;", errReturn)
		w.Linef("            const err1 = _.setText(%s, state, %s);", bufVar, ref)
		w.Linef("            if (err1 !== null) return %s;", errReturn)
	case t.Kind == TplKindBase && !t.IsList && t.Name == "bin":
		w.Linef("            const [state, errState] = _.binState(%s.byteLength);", ref)
		w.Linef("            if (errState !== null) return %s;", errReturn)
		w.Linef("            const err0 = _.setU8(%s, state);", bufVar)
		w.Linef("            if (err0 !== null) return %s;", errReturn)
		w.Linef("            const err1 = _.setBin(%s, state, %s);", bufVar, ref)
		w.Linef("            if (err1 !== null) return %s;", errReturn)
	case t.IsList:
		w.Linef("            const [state, errState] = _.listCountState(%s.length);", ref)
		w.Linef("            if (errState !== null) return %s;", errReturn)
		w.Linef("            const err0 = _.setU8(%s, state);", bufVar)
		w.Linef("            if (err0 !== null) return %s;", errReturn)
		w.Linef("            const err1 = %s;", g.tsDirectListSetExpr(t, ref, "state", bufVar))
		if t.Name == "bool" || t.Name == "text" || t.Name == "bin" {
			w.Linef("            if (err1 !== null) return %s;", errReturn)
		} else {
			w.Linef("            if (err1 !== undefined) return %s;", errReturn)
		}
	case t.Kind == TplKindStruct:
		w.Linef("            const encodeErr = _.set%s(%s, %s);", PascalCase(t.Name), bufVar, ref)
		w.Linef("            if (encodeErr !== undefined) return %s;", errReturn)
	case t.Kind == TplKindEnum:
		typeName := PascalCase(t.Name)
		w.Linef("            if (!_.IsAssignable%s(%s as any)) return %s;", typeName, ref, errReturn)
		w.Linef("            const encodeErr = _.setU8(%s, _.Normalize%s(%s as any) as any);", bufVar, typeName, ref)
		w.Linef("            if (encodeErr !== null) return %s;", errReturn)
	default:
		w.Linef("            const encodeErr = _.%s(%s, %s as any);", g.primitiveSetter(t.Name), bufVar, ref)
		w.Linef("            if (encodeErr !== null) return %s;", errReturn)
	}
	w.Line("        }")
}

func (g *TsGenerator) tsDirectListGetExpr(t TplType, stateVar, bufVar string) string {
	switch {
	case t.Name == "bool":
		return fmt.Sprintf("_.getBoolList(%s, %s)", bufVar, stateVar)
	case t.Name == "text":
		return fmt.Sprintf("_.getTextList(%s, %s)", bufVar, stateVar)
	case t.Name == "bin":
		return fmt.Sprintf("_.getBinList(%s, %s)", bufVar, stateVar)
	case t.Kind == TplKindBase:
		return fmt.Sprintf("_.getDefaultList<%s>(%s, %s, %s, %s)", g.getTsType(TplType{Name: t.Name, Kind: t.Kind}), bufVar, stateVar, g.bitmapDefaultFactory(t), g.bitmapGetter(t))
	case t.Kind == TplKindEnum:
		return fmt.Sprintf("_.get%sListBody(%s, %s)", PascalCase(t.Name), bufVar, stateVar)
	default:
		return fmt.Sprintf("_.get%sListBody(%s, %s)", PascalCase(t.Name), bufVar, stateVar)
	}
}

func (g *TsGenerator) tsDirectListSetExpr(t TplType, ref, stateVar, bufVar string) string {
	switch {
	case t.Name == "bool":
		return fmt.Sprintf("_.setBoolList(%s, %s, %s)", bufVar, stateVar, ref)
	case t.Name == "text":
		return fmt.Sprintf("_.setTextList(%s, %s, %s)", bufVar, stateVar, ref)
	case t.Name == "bin":
		return fmt.Sprintf("_.setBinList(%s, %s, %s)", bufVar, stateVar, ref)
	case t.Kind == TplKindBase:
		return fmt.Sprintf("_.setDefaultList<%s>(%s, %s, %s, %s, %s)", g.getTsType(TplType{Name: t.Name, Kind: t.Kind}), bufVar, stateVar, ref, g.bitmapIsDefault(t), g.bitmapSetter(t))
	case t.Kind == TplKindEnum:
		return fmt.Sprintf("_.set%sListBody(%s, %s, %s)", PascalCase(t.Name), bufVar, stateVar, ref)
	default:
		return fmt.Sprintf("_.set%sListBody(%s, %s, %s)", PascalCase(t.Name), bufVar, stateVar, ref)
	}
}

func (g *TsGenerator) readErrSuffix(t TplType) string {
	if t.Name == "bool" && !t.IsList {
		return "State"
	}
	if g.tagWidth(t) == 1 {
		return "Present"
	}
	return "State"
}
