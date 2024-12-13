package entitygen

import (
	"fmt"

	j "github.com/dave/jennifer/jen"
)

func (g *generator) entityFile(entity *entitySpec) *genFile {
	// Create a Go file with header
	file := j.NewFile(g.pkgName)
	file.HeaderComment(headerComment)

	g.entityStruct(entity, file)
	g.entityRepositoryStruct(entity, file)
	g.entityRepositoryMethods(entity, file)

	return &genFile{FileName: entity.FileName, File: file}
}

func (g *generator) entityStruct(entity *entitySpec, file *j.File) {
	// Comment
	if entity.GoDesc != "" {
		file.Comment(fmt.Sprintf("%s - %s", entity.GoName, entity.GoDesc))
	}

	// Fields
	file.Type().Id(entity.GoName).Struct(g.fields(entity)...).Line()

	// Methods
	file.
		Comment("TrackChanges internally stores entity actual state to track changes for the Update operation.").Line().
		Func().Params(j.Id("e").Op("*").Id(entity.GoName)).Id("TrackChanges").Params().
		Block(
			j.Id("clone").Clone().Op(":=").Op("*").Id("e"),
			j.Id("e").Dot("original").Op("=").Op("&").Id("clone"),
		)
	file.
		Comment("resetChanges after the Update operation.").Line().
		Func().Params(j.Id("e").Op("*").Id(entity.GoName)).Id("resetChanges").Params().
		Block(
			j.Id("e").Dot("original").Op("=").Nil(),
		)
}

func (g *generator) fields(entity *entitySpec) []j.Code {
	var fields []j.Code

	// Self reference, to store original state for updates
	originalComment := "original field stores snapshot of the entity state for the Update operation, see TrackChanges method."
	fields = append(fields, j.Comment(originalComment).Line().Id("original").Op("*").Id(entity.GoName))

	for _, field := range entity.Fields {
		fields = append(fields, g.field(field))
	}

	return fields
}

func (*generator) field(field *fieldSpec) *j.Statement {
	code := j.Empty()

	if field.GoDesc != "" {
		code = code.Comment(fmt.Sprintf("%s - %s", field.GoName, field.GoDesc))
	}

	return code.Id(field.GoName).Add(field.GoType).Tag(field.Tags)
}

func (*generator) entityRepositoryStruct(entity *entitySpec, file *j.File) {
	// Struct
	file.Type().Id(entity.RepositoryName).Struct(clientParam.Clone())

	// Constructor
	file.Func().Id("new" + entity.RepositoryName).Params(clientParam.Clone()).Op("*").Id(entity.RepositoryName).Block(
		j.Return(j.Op("&").Id(entity.RepositoryName).Values(j.Dict{clientVar.Clone(): clientVar.Clone()})),
	)
}

func (g *generator) entityRepositoryMethods(entity *entitySpec, file *j.File) {
	g.createMethod(entity, file)
	g.updateMethod(entity, file)
	g.deleteMethod(entity, file)
	g.allMethod(entity, file)
	g.byIDMethod(entity, file)
	g.lookupHelperMethods(entity, file)
	g.createPayloadMethod(entity, file)
	g.updatePayloadMethod(entity, file)
}

func (g *generator) createMethod(entity *entitySpec, file *j.File) {
	resultType := j.Id(entity.GoName)
	file.
		Comment("Create entity. After successful operation, the new primary ID will be set to the original entity.").Line().
		Add(g.entityRepoMethod(entity, "Create")).
		Params(entityVar.Clone().Op("*").Add(resultType.Clone())).
		// Returns
		Op("*").Qual(webAPIPkg, "APIRequest").Types(resultType.Clone()).
		// Body
		Block(
			j.List(payloadVar.Clone(), j.Id("err")).Op(":=").Id("r").Dot("createPayload").Call(entityVar.Clone()),
			j.If(j.Id("err").Op("!=").Nil()).Block(
				j.Return(
					j.Qual(webAPIPkg, "NewAPIRequestError").Types(resultType.Clone()).Call(
						entityVar.Clone(), j.Id("err"),
					),
				),
			),
			j.Line(),
			pathVar.Clone().Op(":=").Add(g.collectionPath(entity)),
			resultVar.Clone().Op(":=").Add(entityVar.Clone()),
			httpReqVar.Clone().Op(":=").Add(g.httpRequest(resultVar.Clone(), "Post", payloadVar.Clone())),
			httpReqVar.Clone().Dot("Header").Call(j.Lit("Prefer"), j.Lit("return=representation")),
			httpReqVar.Clone().Dot("ExpectStatus").Call(j.Qual("net/http", "StatusCreated")),
			j.Return(
				j.Qual(webAPIPkg, "NewAPIRequest").Call(resultVar.Clone(), httpReqVar.Clone()),
			),
		).
		Line()
}

func (g *generator) updateMethod(entity *entitySpec, file *j.File) {
	resultType := j.Id(entity.GoName)
	file.
		Comment("Update entity. A diff of modifications is generated and saved via API.").Line().
		Comment("Before making changes, it is necessary to call the TrackChanges entity method.").Line().
		Add(g.entityRepoMethod(entity, "Update")).
		Params(entityVar.Clone().Op("*").Add(resultType.Clone())).
		// Returns
		Op("*").Qual(webAPIPkg, "APIRequest").Types(resultType.Clone()).
		// Body
		Block(
			j.Line(),
			j.If(j.Id("entity").Dot("original").Op("==").Nil()).Block(
				j.Id("err").Op(":=").Qual("errors", "New").
					Call(j.Lit(`changes are not tracked: use "TrackChanges" entity method to track changes and allow "Update" operation`)),
				j.Return(
					j.Qual(webAPIPkg, "NewAPIRequestError").Types(resultType.Clone()).Call(
						entityVar.Clone(), j.Id("err"),
					),
				),
			),
			j.Line(),
			j.List(payloadVar.Clone(), j.Id("err")).Op(":=").
				Id("r").Dot("updatePayload").Call(entityVar.Clone().Dot("original"), entityVar.Clone()),
			j.If(j.Id("err").Op("!=").Nil()).Block(
				j.Return(
					j.Qual(webAPIPkg, "NewAPIRequestError").Types(resultType.Clone()).Call(
						entityVar.Clone(), j.Id("err"),
					),
				),
			),
			j.Line(),
			idVar.Clone().Op(":=").Id("entity").Dot(entity.PrimaryIDField.GoName),
			pathVar.Clone().Op(":=").Add(g.entityPathByID(entity)),
			resultVar.Clone().Op(":=").Add(entityVar.Clone()),
			httpReqVar.Clone().Op(":=").Add(g.httpRequest(resultVar.Clone(), "Patch", payloadVar)),
			httpReqVar.Clone().Dot("Header").Call(j.Lit("If-Match"), j.Lit("*")).Comment("prevent create if entity not exists"),
			httpReqVar.Clone().Dot("Header").Call(j.Lit("Prefer"), j.Lit("return=representation")),
			httpReqVar.Clone().Dot("OnSuccess").Call(
				j.Func().
					Params(
						j.Id("ctx").Qual("context", "Context"),
						j.Id("c").Op("*").Add(resultType.Clone()),
					).
					Id("error").Block(
					j.Id("c").Dot("resetChanges").Call(),
					j.Return(j.Nil()),
				),
			),
			j.Return(
				j.Qual(webAPIPkg, "NewAPIRequest").Call(resultVar.Clone(), httpReqVar.Clone()),
			),
		).
		Line()
}

func (g *generator) deleteMethod(entity *entitySpec, file *j.File) {
	file.
		Comment("Delete entity by the ID.").Line().
		Add(g.entityRepoMethod(entity, "Delete")).
		Params(idVar.Clone().String()).
		// Returns
		Op("*").Qual(webAPIPkg, "APIRequest").Types(j.Qual(webAPIPkg, "NoResult")).
		// Body
		Block(
			pathVar.Clone().Op(":=").Add(g.entityPathByID(entity)),
			httpReqVar.Clone().Op(":=").Add(g.httpRequest(j.Op("&").Qual(webAPIPkg, "NoResult").Block(), "Delete", j.Nil())),
			httpReqVar.Clone().Dot("ExpectStatus").Call(j.Qual("net/http", "StatusNoContent")),
			j.Return(
				j.Qual(webAPIPkg, "NewAPIRequest").Call(j.Op("&").Qual(webAPIPkg, "NoResult").Block(), httpReqVar.Clone()),
			),
		).
		Line()
}

func (g *generator) allMethod(entity *entitySpec, file *j.File) {
	resultType := j.Qual(webAPIPkg, "Collection").Types(j.Id(entity.GoName))
	file.Add(g.entityRepoMethod(entity, "All")).
		Params().
		// Returns
		Op("*").Qual(webAPIPkg, "APIRequest").Types(resultType.Clone()).
		// Body
		Block(
			pathVar.Clone().Op(":=").Add(g.collectionPath(entity)),
			resultVar.Clone().Op(":=").Op("&").Add(resultType.Clone()).Block(),
			httpReqVar.Clone().Op(":=").Add(g.httpRequest(resultVar.Clone(), "Get", j.Nil())),
			j.Return(
				j.Qual(webAPIPkg, "NewAPIRequest").Call(resultVar.Clone(), httpReqVar.Clone()),
			),
		).
		Line()
}

func (g *generator) byIDMethod(entity *entitySpec, file *j.File) {
	if entity.PrimaryIDField == nil {
		return
	}

	resultType := j.Id(entity.GoName)
	file.Add(g.entityRepoMethod(entity, "ByID")).
		Params(idVar.Clone().String()).
		// Returns
		Op("*").Qual(webAPIPkg, "APIRequest").Types(resultType.Clone()).
		// Body
		Block(
			pathVar.Clone().Op(":=").Add(g.entityPathByID(entity)),
			resultVar.Clone().Op(":=").Op("&").Add(resultType.Clone()).Block(),
			httpReqVar.Clone().Op(":=").Add(g.httpRequest(resultVar.Clone(), "Get", j.Nil())),
			j.Return(
				j.Qual(webAPIPkg, "NewAPIRequest").Call(resultVar.Clone(), httpReqVar.Clone()),
			),
		).
		Line()
}

func (g *generator) lookupHelperMethods(entity *entitySpec, file *j.File) {
	for _, field := range entity.Fields {
		if field.APIType == lookupFieldType {
			g.lookupHelperMethodsForField(entity, field, file)
		}
	}
}

func (g *generator) lookupHelperMethodsForField(entity *entitySpec, field *fieldSpec, file *j.File) {
	resultType := j.Qual(webAPIPkg, "Collection").Types(j.Id(entity.GoName))
	file.Add(g.entityRepoMethod(entity, "By"+field.GoName)).
		Params(idVar.Clone().String()).
		// Returns
		Op("*").Qual(webAPIPkg, "APIRequest").Types(resultType.Clone()).
		// Body
		Block(
			pathVar.Clone().Op(":=").Add(g.collectionPath(entity)),
			resultVar.Clone().Op(":=").Op("&").Add(resultType.Clone()).Block(),
			httpReqVar.Clone().Op(":=").Add(g.httpRequest(resultVar.Clone(), "Get", j.Nil())),
			httpReqVar.Clone().Dot("Filter").
				Call(j.
					Lit(field.LogicalName+" eq '").
					Op("+").
					Add(j.Qual(webAPIPkg, "ID").Call(idVar.Clone())).
					Op("+").
					Lit("'"),
				),
			j.Return(
				j.Qual(webAPIPkg, "NewAPIRequest").Call(resultVar.Clone(), httpReqVar.Clone()),
			),
		).
		Line()
}

func (g *generator) createPayloadMethod(entity *entitySpec, file *j.File) {
	// Init payload variable, create a map
	mapType := j.Map(j.String()).Any()
	blocks := []j.Code{
		payloadVar.Clone().Op(":=").Add(mapType).Clone().Block(),
	}

	// Generate diff for each field
	for _, field := range entity.Fields {
		if field.IsPrimaryID {
			continue
		}

		if fieldBlocks := g.payloadFieldSetter(field, entityVar); len(fieldBlocks) > 1 {
			blocks = append(blocks, j.Block(fieldBlocks...))
		} else {
			blocks = append(blocks, fieldBlocks...)
		}
	}

	// Return
	blocks = append(blocks, j.Return(payloadVar.Clone(), j.Nil()))

	// Compose method
	file.
		Comment("createPayload method generates payload for the Create operation.").Line().
		Add(g.entityRepoMethod(entity, "createPayload")).
		Params(entityVar.Clone().Op("*").Id(entity.GoName)).
		// Returns
		Parens(j.List(mapType.Clone(), j.Id("error"))).
		// Body
		Block(blocks...).
		Line()
}

func (g *generator) updatePayloadMethod(entity *entitySpec, file *j.File) {
	// Init payload variable, create a map
	mapType := j.Map(j.String()).Any()
	blocks := []j.Code{
		payloadVar.Clone().Op(":=").Add(mapType).Clone().Block(),
	}

	// Generate diff for each field
	for _, field := range entity.Fields {
		if field.IsPrimaryID {
			continue
		}
		blocks = append(blocks, g.compareFieldValues(field).Block(g.payloadFieldSetter(field, modifiedVar)...))
	}

	// Return
	blocks = append(blocks, j.Return(payloadVar.Clone(), j.Nil()))

	// Compose method
	file.
		Comment("updatePayload method generates diff of original and modified entity state for the Update operation.").Line().
		Add(g.entityRepoMethod(entity, "updatePayload")).
		Params(originalVar.Clone(), modifiedVar.Clone().Op("*").Id(entity.GoName)).
		// Returns
		Parens(j.List(mapType.Clone(), j.Id("error"))).
		// Body
		Block(blocks...).
		Line()
}

func (*generator) payloadFieldSetter(field *fieldSpec, entityValueVar *j.Statement) (codes []j.Code) {
	value := entityValueVar.Clone().Dot(field.GoName)
	switch field.APIType {
	case lookupFieldType:
		fullIDVar := j.Id("idOrNil")
		codes = append(codes,
			j.List(fullIDVar.Clone(), j.Id("err")).Op(":=").Id("lookupFullIDOrNil").Call(value.Clone()),
			j.If(j.Id("err").Op("!=").Nil()).Block(
				j.Return(j.Nil(), j.Id("err")),
			),
			payloadVar.Clone().
				Index(j.Lit(field.SchemaName+"@odata.bind")).Op("=").
				Add(fullIDVar.Clone()),
		)
	default:
		codes = append(codes,
			payloadVar.Clone().
				Index(j.Lit(field.LogicalName)).Op("=").
				Add(value.Clone()),
		)
	}

	return codes
}

func (*generator) compareFieldValues(field *fieldSpec) *j.Statement {
	var originalValue, modifiedValue *j.Statement
	switch field.APIType {
	case lookupFieldType:
		originalValue = originalVar.Clone().Dot(field.GoName).Clone().Dot("ID").Call()
		modifiedValue = modifiedVar.Clone().Dot(field.GoName).Clone().Dot("ID").Call()
	default:
		originalValue = originalVar.Clone().Dot(field.GoName)
		modifiedValue = modifiedVar.Clone().Dot(field.GoName)
	}

	return j.If(originalValue.Clone().Op("!=").Add(modifiedValue.Clone()))
}

func (*generator) entityRepoMethod(entity *entitySpec, method string) *j.Statement {
	return j.Func().Params(j.Id("r").Op("*").Id(entity.RepositoryName)).Id(method)
}

func (g *generator) entityPathByID(entity *entitySpec) *j.Statement {
	return g.collectionPath(entity).
		Op("+").
		Lit("(").
		Op("+").
		Qual(webAPIPkg, "ID").Call(idVar.Clone()).
		Op("+").
		Lit(")")
}

func (*generator) collectionPath(entity *entitySpec) *j.Statement {
	return j.Lit(entity.EntitySetName)
}

func (*generator) httpRequest(result *j.Statement, method string, body *j.Statement) *j.Statement {
	return j.Qual(webAPIPkg, "NewHTTPRequest").Call(
		result,
		j.Id("r.client"),
		j.Qual("net/http", "Method"+method),
		j.Id("path"),
		body,
	)
}
