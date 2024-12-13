package entitygen

import (
	"context"
	"errors"
	"fmt"
	"sort"
	"strings"

	"github.com/michaljurecko/ch-demo/internal/pkg/common/dataverse/metadata"

	j "github.com/dave/jennifer/jen"
)

// loader composes registry from API responses.
type loader struct {
	*registry
}

// registry of temporary entity definitions.
type registry struct {
	// Entities maps LogicalName to entity
	Entities map[string]*entitySpec
}

// entitySpec is a temporary entity definition, an intermediate step between API and Go file.
type entitySpec struct {
	LogicalName      string
	EntitySetName    string
	GoName           string
	RepositoryName   string
	FileName         string
	GoDesc           string
	PrimaryIDField   *fieldSpec
	PrimaryNameField *fieldSpec
	Fields           []*fieldSpec
}

// fieldSpec is a temporary entity definition, an intermediate step between API and Go file.
type fieldSpec struct {
	SchemaName  string
	LogicalName string
	IsPrimaryID bool
	GoName      string
	APIType     string
	GoType      *j.Statement
	GoDesc      string
	Tags        map[string]string
}

func newRegistry(ctx context.Context, api *metadata.API) (*registry, error) {
	r := &registry{Entities: make(map[string]*entitySpec)}
	l := &loader{registry: r}
	if err := l.loadEntities(ctx, api); err != nil {
		return nil, err
	}
	return r, nil
}

// loadEntities metadata from API to temporary registry.
func (l *loader) loadEntities(ctx context.Context, api *metadata.API) error {
	var errs []error

	// Load entities metadata
	entities, err := api.CustomEntityDefinitions().Do(ctx)
	if err != nil {
		return fmt.Errorf("cannot load entities metadata: %w", err)
	}

	// Create entities specifications
	for _, entity := range entities.Definitions {
		if err := l.addEntity(entity); err != nil {
			errs = append(errs, err)
		}
	}

	// Load attributes metadata
	for _, entity := range l.Entities {
		if err := l.loadAttributes(ctx, api, entity); err != nil {
			errs = append(errs, fmt.Errorf("cannot load attributes for entity %s: %w", entity.LogicalName, err))
		}
	}

	return errors.Join(errs...)
}

// loadAttributes metadata from API to temporary registry.
func (l *loader) loadAttributes(ctx context.Context, api *metadata.API, entity *entitySpec) error {
	attrs, err := api.CustomEntityAttributes(entity.LogicalName).Do(ctx)
	if err != nil {
		return fmt.Errorf(`cannot load attributes for entity "%s": %w`, entity.LogicalName, err)
	}

	sort.SliceStable(attrs.Attributes, func(i, j int) bool {
		return attrs.Attributes[i].ColumnNumber < attrs.Attributes[j].ColumnNumber
	})

	var errs []error
	for _, attr := range attrs.Attributes {
		if err := l.addField(entity, attr); err != nil {
			errs = append(errs, err)
		}
	}

	return errors.Join(errs...)
}

func (l *loader) addEntity(entity metadata.Entity) error {
	goName, err := entityGoName(entity)
	if err != nil {
		return err
	}

	fileName, err := entityFileName(entity)
	if err != nil {
		return err
	}

	spec := &entitySpec{
		LogicalName:    entity.LogicalName,
		EntitySetName:  entity.EntitySetName,
		GoName:         goName,
		RepositoryName: goName + "Repository",
		FileName:       fileName,
		GoDesc:         entityDesc(entity),
	}

	l.Entities[spec.LogicalName] = spec

	return nil
}

func (l *loader) addField(entity *entitySpec, attr metadata.Attribute) error {
	goName, err := attributeGoName(entity, attr, entity.GoName)
	if err != nil {
		return err
	}

	goType, err := l.resolveGoAttrType(entity, attr)
	if err != nil {
		return err
	}

	jsonTag := attr.LogicalName
	if attr.AttributeType == lookupFieldType {
		jsonTag = fmt.Sprintf("_%s_value", jsonTag)
	}
	if attr.IsPrimaryID {
		jsonTag += ",omitempty" // don't send empty ID on create
	}

	spec := &fieldSpec{
		SchemaName:  attr.SchemaName,
		LogicalName: attr.LogicalName,
		GoName:      goName,
		IsPrimaryID: attr.IsPrimaryID,
		APIType:     attr.AttributeType,
		GoType:      goType,
		GoDesc:      attributeDesc(attr),
		Tags:        map[string]string{"json": jsonTag},
	}

	if attr.IsPrimaryID {
		entity.PrimaryIDField = spec
	}

	if attr.IsPrimaryName && attr.AttributeType == "String" {
		entity.PrimaryNameField = spec
	}

	entity.Fields = append(entity.Fields, spec)
	return nil
}

func (l *loader) resolveGoAttrType(entity *entitySpec, attr metadata.Attribute) (*j.Statement, error) {
	switch attr.AttributeType {
	case "Lookup":
		return l.resolveGoAttrLookupType(entity, attr)
	case "Uniqueidentifier":
		return j.String(), nil
	case "String":
		return j.String(), nil
	case "Integer":
		return j.Int(), nil
	case "DateTime":
		return j.Qual("time", "Time"), nil
	default:
		return nil, fmt.Errorf(`not implemented attribute "%s" type "%s"`, attr.LogicalName, attr.AttributeType)
	}
}

func (l *loader) resolveGoAttrLookupType(entity *entitySpec, attr metadata.Attribute) (*j.Statement, error) {
	// Map referenced types to matching Go structs
	var refStructs []j.Code
	for _, refType := range attr.Targets {
		refEntity, ok := l.Entities[refType]
		if !ok {
			return nil, fmt.Errorf(
				`lookup attribute "%s.%s" refers to the type "%s", but it is not known`,
				entity.LogicalName, attr.LogicalName, refType,
			)
		}
		refStructs = append(refStructs, j.Id(refEntity.GoName))
	}

	// Complete the type
	switch len(refStructs) {
	case 0:
		return nil, fmt.Errorf(`empty "targets"" for lookup attribute "%s.%s"`, entity.LogicalName, attr.LogicalName)
	case 1:
		return j.Qual(webAPIPkg, "Lookup").Types(refStructs[0]), nil
	default:
		return nil, fmt.Errorf(
			`lookup attribute "%s.%s" refers to multiple types "%s", not implemented`,
			entity.LogicalName, attr.LogicalName, strings.Join(attr.Targets, `", "`),
		)
	}
}
