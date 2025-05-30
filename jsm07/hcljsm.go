package jsm07

import (
	jsonschema "github.com/genelet/hclschema/jsonschema07"
)

func NewSchemaFromJSM(s *jsonschema.Schema) *Schema {
	if s == nil {
		return nil
	}

	if isFull(s) {
		return schemaFullToHcl(s)
	}

	if s.Ref != nil {
		return &Schema{
			Reference: referenceToHcl(s),
		}
	}

	common := commonToHcl(s)

	switch *s.Type.String {
	case "boolean":
		return &Schema{
			Common: common,
		}
	case "number", "integer":
		return &Schema{
			Common:       common,
			SchemaNumber: numberToHcl(s),
		}
	case "string":
		return &Schema{
			Common:       common,
			SchemaString: stringToHcl(s),
		}
	case "array":
		return &Schema{
			Common:      common,
			SchemaArray: arrayToHcl(s),
		}
	case "object":
		if isMap(s) && !isObject(s) {
			typ := "map"
			common.Type.String = &typ
			return &Schema{
				Common:    common,
				SchemaMap: mapToHcl(s),
			}
		}

		return &Schema{
			Common:       common,
			SchemaObject: objectToHcl(s),
		}
	default:
	}

	return schemaFullToHcl(s)
}

func (s *Schema) ToJSM() *jsonschema.Schema {
	if s == nil {
		return nil
	}
	if s.isFull {
		return schemaFullToJSM(s)
	}

	if s.Reference != nil {
		return referenceToJSM(s.Reference)
	}

	schema := commonToJSM(s.Common)
	if s.SchemaString != nil {
		return stringToJSM(schema, s.SchemaString)
	}
	if s.SchemaNumber != nil {
		return numberToJSM(schema, s.SchemaNumber)
	}
	if s.SchemaArray != nil {
		return arrayToJSM(schema, s.SchemaArray)
	}
	if s.SchemaObject != nil {
		return objectToJSM(schema, s.SchemaObject)
	}
	if s.SchemaMap != nil {
		return mapToJSM(schema, s.SchemaMap)
	} else if s.Common != nil {
		// boolean
		return schema
	}

	return nil
}

func namedSchemaArrayToMap(s *[]*jsonschema.NamedSchema) map[string]*Schema {
	if s == nil {
		return nil
	}
	m := make(map[string]*Schema)
	for _, v := range *s {
		m[v.Name] = NewSchemaFromJSM(v.Value)
	}
	return m
}

func mapToNamedSchemaArray(s map[string]*Schema) *[]*jsonschema.NamedSchema {
	if s == nil {
		return nil
	}
	var arr []*jsonschema.NamedSchema
	for k, v := range s {
		arr = append(arr, &jsonschema.NamedSchema{
			Name:  k,
			Value: v.ToJSM(),
		})
	}
	return &arr
}

func namedSchemaOrStringArrayArrayToMap(s *[]*jsonschema.NamedSchemaOrStringArray) map[string]*SchemaOrStringArray {
	if s == nil {
		return nil
	}
	m := make(map[string]*SchemaOrStringArray)
	for _, v := range *s {
		if v.Value.Schema != nil {
			m[v.Name] = &SchemaOrStringArray{
				Schema: NewSchemaFromJSM(v.Value.Schema),
			}
		} else {
			var arr []string
			arr = append(arr, *v.Value.StringArray...)
			m[v.Name] = &SchemaOrStringArray{
				StringArray: arr,
			}
		}
	}
	return m
}

func mapToNamedSchemaOrStringArrayArray(s map[string]*SchemaOrStringArray) *[]*jsonschema.NamedSchemaOrStringArray {
	if s == nil {
		return nil
	}
	var arr []*jsonschema.NamedSchemaOrStringArray
	for k, v := range s {
		if v.Schema != nil {
			arr = append(arr, &jsonschema.NamedSchemaOrStringArray{
				Name: k,
				Value: &jsonschema.SchemaOrStringArray{
					Schema: v.Schema.ToJSM(),
				},
			})
		} else {
			var sa []string
			sa = append(sa, v.StringArray...)
			arr = append(arr, &jsonschema.NamedSchemaOrStringArray{
				Name: k,
				Value: &jsonschema.SchemaOrStringArray{
					StringArray: &sa,
				},
			})
		}
	}
	return &arr
}

func sliceToHcl(allof *[]*jsonschema.Schema) []*Schema {
	if allof == nil {
		return nil
	}
	var arr []*Schema
	for _, v := range *allof {
		arr = append(arr, NewSchemaFromJSM(v))
	}
	return arr
}

func sliceToJSM(allof []*Schema) *[]*jsonschema.Schema {
	if allof == nil {
		return nil
	}
	var arr []*jsonschema.Schema
	for _, v := range allof {
		arr = append(arr, v.ToJSM())
	}
	return &arr
}

func referenceToHcl(s *jsonschema.Schema) *Reference {
	if s == nil || !isReference(s) {
		return nil
	}
	return &Reference{
		Ref: s.Ref,
	}
}

func referenceToJSM(r *Reference) *jsonschema.Schema {
	if r == nil {
		return nil
	}
	return &jsonschema.Schema{
		Ref: r.Ref,
	}
}

func commonToHcl(s *jsonschema.Schema) *Common {
	if s == nil || !isCommon(s) {
		return nil
	}
	common := &Common{
		Type:    s.Type,
		Format:  s.Format,
		Default: s.Default,
	}
	if s.Enumeration != nil {
		common.Enumeration = *s.Enumeration
	}
	return common
}

func commonToJSM(c *Common) *jsonschema.Schema {
	if c == nil {
		return nil
	}

	jsm := &jsonschema.Schema{}
	jsm.Type = c.Type
	jsm.Format = c.Format
	jsm.Default = c.Default
	if c.Enumeration != nil {
		jsm.Enumeration = &c.Enumeration
	}
	return jsm
}

func stringToHcl(s *jsonschema.Schema) *SchemaString {
	if s == nil || !isString(s) {
		return nil
	}

	return &SchemaString{
		MaxLength: s.MaxLength,
		MinLength: s.MinLength,
		Pattern:   s.Pattern,
	}
}

func stringToJSM(jsm *jsonschema.Schema, s *SchemaString) *jsonschema.Schema {
	if s == nil {
		return jsm
	}
	if jsm == nil {
		jsm = &jsonschema.Schema{}
	}

	jsm.MaxLength = s.MaxLength
	jsm.MinLength = s.MinLength
	jsm.Pattern = s.Pattern
	return jsm
}

func numberToHcl(s *jsonschema.Schema) *SchemaNumber {
	if s == nil || !isNumber(s) {
		return nil
	}

	return &SchemaNumber{
		MultipleOf:       s.MultipleOf,
		Maximum:          s.Maximum,
		ExclusiveMaximum: s.ExclusiveMaximum,
		Minimum:          s.Minimum,
		ExclusiveMinimum: s.ExclusiveMinimum,
	}
}

func numberToJSM(jsm *jsonschema.Schema, n *SchemaNumber) *jsonschema.Schema {
	if n == nil {
		return jsm
	}
	if jsm == nil {
		jsm = &jsonschema.Schema{}
	}

	jsm.MultipleOf = n.MultipleOf
	jsm.Maximum = n.Maximum
	jsm.ExclusiveMaximum = n.ExclusiveMaximum
	jsm.Minimum = n.Minimum
	jsm.ExclusiveMinimum = n.ExclusiveMinimum

	return jsm
}
func arrayToHcl(s *jsonschema.Schema) *SchemaArray {
	if s == nil || !isArray(s) {
		return nil
	}

	items := new(SchemaOrSchemaArray)
	if s.Items != nil {
		if s.Items.Schema != nil {
			items.Schema = NewSchemaFromJSM(s.Items.Schema)
		} else {
			items.SchemaArray = sliceToHcl(s.Items.SchemaArray)
		}
	}

	return &SchemaArray{
		Items:       items,
		MaxItems:    s.MaxItems,
		MinItems:    s.MinItems,
		UniqueItems: s.UniqueItems,
	}
}

func arrayToJSM(jsm *jsonschema.Schema, a *SchemaArray) *jsonschema.Schema {
	if a == nil {
		return jsm
	}
	if jsm == nil {
		jsm = &jsonschema.Schema{}
	}

	if a.Items != nil {
		if jsm.Items == nil {
			jsm.Items = &jsonschema.SchemaOrSchemaArray{}
		}
		if a.Items.Schema != nil {
			jsm.Items.Schema = a.Items.Schema.ToJSM()
		} else {
			jsm.Items.SchemaArray = sliceToJSM(a.Items.SchemaArray)
		}
	}

	jsm.MaxItems = a.MaxItems
	jsm.MinItems = a.MinItems
	jsm.UniqueItems = a.UniqueItems
	return jsm
}

func mapToHcl(s *jsonschema.Schema) *SchemaMap {
	if s == nil || !isMap(s) {
		return nil
	}

	return &SchemaMap{
		AdditionalProperties: &SchemaOrBoolean{
			Schema:  NewSchemaFromJSM(s.AdditionalProperties.Schema),
			Boolean: s.AdditionalProperties.Boolean,
		},
	}
}

func mapToJSM(jsm *jsonschema.Schema, m *SchemaMap) *jsonschema.Schema {
	if m == nil {
		return jsm
	}

	if jsm == nil {
		jsm = &jsonschema.Schema{}
	}
	if jsm.AdditionalProperties == nil {
		jsm.AdditionalProperties = &jsonschema.SchemaOrBoolean{}
	}
	if m.AdditionalProperties.Schema != nil {
		jsm.AdditionalProperties.Schema = m.AdditionalProperties.Schema.ToJSM()
	} else {
		jsm.AdditionalProperties.Boolean = m.AdditionalProperties.Boolean
	}

	return jsm
}

func objectToHcl(s *jsonschema.Schema) *SchemaObject {
	if s == nil || !isObject(s) {
		return nil
	}

	object := &SchemaObject{
		MaxProperties: s.MaxProperties,
		MinProperties: s.MinProperties,
		Properties:    namedSchemaArrayToMap(s.Properties),
	}
	if s.Required != nil {
		object.Required = *s.Required
	}
	return object
}

func objectToJSM(jsm *jsonschema.Schema, o *SchemaObject) *jsonschema.Schema {
	if o == nil {
		return jsm
	}

	if jsm == nil {
		jsm = &jsonschema.Schema{}
	}

	jsm.Properties = mapToNamedSchemaArray(o.Properties)
	jsm.MaxProperties = o.MaxProperties
	jsm.MinProperties = o.MinProperties
	if o.Required != nil {
		jsm.Required = &o.Required
	}
	return jsm
}

func schemaFullToHcl(s *jsonschema.Schema) *Schema {
	full := &SchemaFull{
		Schema:            s.Schema,
		ID:                s.ID,
		ReadOnly:          s.ReadOnly,
		WriteOnly:         s.WriteOnly,
		PatternProperties: namedSchemaArrayToMap(s.PatternProperties),
		Dependencies:      namedSchemaOrStringArrayArrayToMap(s.Dependencies),

		Reference:    referenceToHcl(s),
		Common:       commonToHcl(s),
		SchemaNumber: numberToHcl(s),
		SchemaString: stringToHcl(s),
		SchemaArray:  arrayToHcl(s),
		SchemaMap:    mapToHcl(s),
		SchemaObject: objectToHcl(s),

		AllOf:       sliceToHcl(s.AllOf),
		AnyOf:       sliceToHcl(s.AnyOf),
		OneOf:       sliceToHcl(s.OneOf),
		Not:         NewSchemaFromJSM(s.Not),
		Definitions: namedSchemaArrayToMap(s.Definitions),

		Title:       s.Title,
		Description: s.Description,
	}
	if s.AdditionalItems != nil {
		if s.AdditionalItems.Schema != nil {
			full.AdditionalItems = &SchemaOrBoolean{
				Schema: NewSchemaFromJSM(s.AdditionalItems.Schema),
			}
		} else {
			full.AdditionalItems = &SchemaOrBoolean{
				Boolean: s.AdditionalItems.Boolean,
			}
		}
	}
	return &Schema{
		SchemaFull: full,
		isFull:     true,
	}
}

func schemaFullToJSM(s *Schema) *jsonschema.Schema {
	if s == nil || !s.isFull {
		return nil
	}
	full := s.SchemaFull
	jsm := commonToJSM(full.Common)
	if jsm == nil {
		jsm = &jsonschema.Schema{}
	}
	jsm.Schema = full.Schema
	jsm.ID = full.ID
	if full.Reference != nil {
		jsm.Ref = full.Reference.Ref
	}
	jsm.ReadOnly = full.ReadOnly
	jsm.WriteOnly = full.WriteOnly
	if full.SchemaString != nil {
		jsm = stringToJSM(jsm, full.SchemaString)
	}
	if full.SchemaNumber != nil {
		jsm = numberToJSM(jsm, full.SchemaNumber)
	}
	if full.SchemaArray != nil {
		jsm = arrayToJSM(jsm, full.SchemaArray)
	}
	if full.SchemaObject != nil {
		jsm = objectToJSM(jsm, full.SchemaObject)
	}
	if full.SchemaMap != nil {
		jsm = mapToJSM(jsm, full.SchemaMap)
	}

	if full.AdditionalItems != nil {
		if full.AdditionalItems.Schema != nil {
			jsm.AdditionalItems.Schema = (full.AdditionalItems.Schema).ToJSM()
		} else {
			jsm.AdditionalItems.Boolean = full.AdditionalItems.Boolean
		}
	}
	jsm.PatternProperties = mapToNamedSchemaArray(full.PatternProperties)
	jsm.Dependencies = mapToNamedSchemaOrStringArrayArray(full.Dependencies)

	jsm.AllOf = sliceToJSM(full.AllOf)
	jsm.AnyOf = sliceToJSM(full.AnyOf)
	jsm.OneOf = sliceToJSM(full.OneOf)
	if s.Not != nil {
		jsm.Not = s.Not.ToJSM()
	}
	jsm.Definitions = mapToNamedSchemaArray(full.Definitions)

	jsm.Title = full.Title
	jsm.Description = full.Description

	return jsm
}
