package entitygen

import (
	j "github.com/dave/jennifer/jen"
)

func (g *generator) commonFile(entities []*entitySpec) *genFile {
	// Create a Go file with header
	file := j.NewFile(g.pkgName)
	file.HeaderComment(headerComment)

	g.commonRepository(entities, file)
	g.lookupFullIDOrNilFunc(entities, file)

	return &genFile{FileName: "zz_common.go", File: file}
}

func (*generator) commonRepository(entities []*entitySpec, file *j.File) {
	// Repository contains one field for each entity repository, generate fields and their initializations
	var fields []j.Code
	var fieldsInit []j.Code

	// Client field
	fields = append(fields, j.Id("client").Op("*").Qual(webAPIPkg, "Client"))
	fieldsInit = append(fieldsInit, j.Line().Id("client").Op(": ").Id("client"))

	// Fields for entity repositories
	for _, entity := range entities {
		fieldName := j.Id(lowerFirst(entity.RepositoryName))
		fields = append(fields, fieldName.Clone().Op("*").Id(entity.RepositoryName))
		fieldsInit = append(fieldsInit, j.Line().Add(fieldName.Clone()).Op(": ").Id("new"+entity.RepositoryName).Call(clientVar.Clone()))
	}

	file.Type().Id("Repository").Struct(fields...).Line()

	// Constructor
	file.Func().Id("NewRepository").Params(clientParam).Op("*").Id("Repository").Block(
		j.Return(j.Op("&").Id("Repository").Values(fieldsInit...)),
	).Line()

	// NewChangeSet helper
	file.Func().Params(j.Id("r").Op("*").Id("Repository")).Id("NewChangeSet").Params().Op("*").Qual(webAPIPkg, "ChangeSet").Block(
		j.Return(j.Qual(webAPIPkg, "NewChangeSet").Call(j.Id("r").Dot("client"))),
	).Line()

	// Getters
	for _, entity := range entities {
		file.Func().Params(j.Id("r").Op("*").Id("Repository")).Id(entity.GoName).Params().Op("*").Id(entity.RepositoryName).Block(
			j.Return(j.Id("r").Dot(lowerFirst(entity.RepositoryName))),
		).Line()
	}
}

func (*generator) lookupFullIDOrNilFunc(entities []*entitySpec, file *j.File) {
	setNameVar := j.Id("set")

	// Create switch case for each entity struct
	var cases []j.Code
	for _, entity := range entities {
		if entity.PrimaryIDField == nil {
			continue
		}
		cases = append(cases,
			j.Case(j.Qual(webAPIPkg, "Lookup").Types(j.Id(entity.GoName))).Block(
				setNameVar.Clone().Op("=").Lit("/"+entity.EntitySetName),
			),
		)
	}

	// Default error
	cases = append(cases,
		j.Default().Block(
			j.Return(j.Nil(), j.Qual("fmt", "Errorf").Call(j.Lit("unexpected entity type %T"), j.Id("v"))),
		),
	)

	file.
		Func().Id("lookupFullIDOrNil").
		Params(j.Id("v").Qual(webAPIPkg, "LookupInterface")).
		Parens(j.List(j.Op("*").String(), j.Id("error"))).
		Block(
			j.If(j.Id("ref").Op(":=").Id("v").Dot("ContentIDRef").Call().Op(";").Id("ref").Op("!=").Nil()).Block(
				j.Return(j.Id("ref"), j.Nil()),
			),
			j.Line(),
			j.Var().Add(setNameVar.Clone()).String(),
			j.Switch(j.Id("v").Op(".").Parens(j.Id("type"))).Block(cases...),
			j.Line(),
			idVar.Clone().Op(":=").Add(j.Id("v").Clone().Dot("ID").Call()),
			j.If(idVar.Clone().Op("==").Lit("")).Block(
				j.Return(j.Nil(), j.Nil()),
			),
			j.Line(),
			j.Id("fullID").Op(":=").Add(setNameVar.Clone()).
				Op("+").Lit("(").Op("+").
				Qual(webAPIPkg, "ID").Call(idVar.Clone()).
				Op("+").Lit(")"),
			j.Return(j.Op("&").Id("fullID"), j.Nil()),
		)
}
