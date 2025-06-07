package jsonschema07

import (
	"fmt"

	"github.com/genelet/hcllight/light"
	"gopkg.in/yaml.v3"
)

// ParseSchema parses a HCL string representing a JSON schema and returns a Schema object.
func ParseSchema(data string) (*Schema, error) {
	s, err := light.ParseBody([]byte(data))
	if err != nil {
		return nil, err
	}
	return NewSchemaFromBody(s)
}

func NewSchemaFromBody(body *light.Body) (*Schema, error) {
	if body == nil {
		return nil, nil
	}

	var typ *StringOrStringArray
	if attr, ok := body.Attributes["type"]; ok {
		typ = expressionToStringOrStringArray(attr.Expr)
	}
	integerOrFloat := func(expr *light.Expression) *IntegerOrFloat {
		if typ != nil && typ.String != nil {
			if *typ.String == "integer" {
				i64 := light.LiteralValueExprToInt64(expr)
				return NewIntegerOrFloatWithInteger(*i64)
			} else {
				f64 := light.LiteralValueExprToFloat64(expr)
				return NewIntegerOrFloatWithFloat(*f64)
			}
		}
		return nil
	}

	schema := &Schema{}
	for k, v := range body.Attributes {
		expr := v.Expr
		var err error
		switch k {
		case "$id":
			schema.ID = light.TextValueExprToString(expr)
		case "$schema":
			schema.Schema = light.TextValueExprToString(expr)
		case "$ref":
			schema.Ref = light.TextValueExprToString(expr)
		case "$comment":
			schema.Comment = light.TextValueExprToString(expr)
		case "title":
			schema.Title = light.TextValueExprToString(expr)
		case "description":
			schema.Description = light.TextValueExprToString(expr)

		case "default":
			schema.Default = expressionToYamlNode(expr)
		case "readOnly":
			schema.ReadOnly = light.LiteralValueExprToBoolean(expr)
		case "writeOnly":
			schema.WriteOnly = light.LiteralValueExprToBoolean(expr)
		case "example":
			schema.Examples, err = tupleConsExprToSlice(expr.GetTcexpr())
		case "multipleOf":
			schema.MultipleOf = integerOrFloat(expr)
		case "maximum":
			schema.Maximum = integerOrFloat(expr)
		case "exclusiveMaximum":
			schema.ExclusiveMaximum = integerOrFloat(expr)
		case "minimum":
			schema.Minimum = integerOrFloat(expr)
		case "exclusiveMinimum":
			schema.ExclusiveMinimum = integerOrFloat(expr)

		case "maxLength":
			schema.MaxLength = light.LiteralValueExprToInt64(expr)
		case "minLength":
			schema.MinLength = light.LiteralValueExprToInt64(expr)
		case "pattern":
			schema.Pattern = light.TextValueExprToString(expr)

		case "additionalItems":
			schema.AdditionalItems, err = expressionToCombined(expr)
		case "items":
			schema.Items, err = expressionToCombinedOrCombinedArray(expr)
		case "maxItems":
			schema.MaxItems = light.LiteralValueExprToInt64(expr)
		case "minItems":
			schema.MinItems = light.LiteralValueExprToInt64(expr)
		case "uniqueItems":
			schema.UniqueItems = light.LiteralValueExprToBoolean(expr)

		case "contains":
			schema.Contains, err = expressionToCombined(expr)
		case "maxProperties":
			schema.MaxProperties = light.LiteralValueExprToInt64(expr)
		case "minProperties":
			schema.MinProperties = light.LiteralValueExprToInt64(expr)
		case "required":
			schema.Required = light.TupleConsExprToStringArray(expr)
		case "additionalProperties":
			schema.AdditionalProperties, err = expressionToCombined(expr)
		case "definitions":
			schema.Definitions, err = objectConsExprToMapCombined(expr.GetOcexpr())
		case "properties":
			schema.Properties, err = objectConsExprToMapCombined(expr.GetOcexpr())
		case "patternProperties":
			schema.PatternProperties, err = objectConsExprToMapCombined(expr.GetOcexpr())
		case "dependencies":
			//stringArray := light.TupleConsExprToStringArray(v.Expr)
			//value := NewCombinedOrStringArrayWithStringArray(stringArray)
			//schema.Dependencies = append(schema.Dependencies, &NamedCombinedOrStringArray{Name: k, Value: value})
		case "propertiesNames":
			schema.PropertyNames, err = expressionToCombined(expr)
		case "const":
			schema.Const = expressionToYamlNode(expr)
		case "enum":
			schema.Enumeration, err = tupleConsExprToEnum(expr.GetTcexpr())
		case "format":
			schema.Format = light.TextValueExprToString(expr)
		case "type":
			schema.Type = typ
		case "contentEncoding":
			schema.ContentEncoding = light.TextValueExprToString(expr)
		case "contentMediaType":
			schema.ContentMediaType = light.TextValueExprToString(expr)

		case "if":
			schema.If, err = expressionToCombined(expr)
		case "then":
			schema.Then, err = expressionToCombined(expr)
		case "else":
			schema.Else, err = expressionToCombined(expr)
		case "allOf":
			schema.AllOf, err = tupleConsExprToSlice(expr.GetTcexpr())
		case "anyOf":
			schema.AnyOf, err = tupleConsExprToSlice(expr.GetTcexpr())
		case "oneOf":
			schema.OneOf, err = tupleConsExprToSlice(expr.GetTcexpr())
		case "not":
			schema.Not, err = expressionToCombined(expr)

		default:
			err = fmt.Errorf("unknown attribute in schema: %s", k)
		}
		if err != nil {
			return nil, err
		}
	}

	for _, block := range body.Blocks {
		var err error
		switch block.Type {
		case "dependencies":
			combinedOrStringArray, err := bodyCombinedOrStringArray(block.Labels[0], block.Bdy)
			if err != nil {
				return nil, err
			}
			schema.Dependencies = append(schema.Dependencies, combinedOrStringArray)
		case "definitions":
			schema.Definitions, err = bodyToMapCombined(block.Bdy)
		case "properties":
			schema.Properties, err = bodyToMapCombined(block.Bdy)
		case "patternProperties":
			schema.PatternProperties, err = bodyToMapCombined(block.Bdy)
		case "examples", "allOf", "anyOf", "oneOf":
		//schema.Examples, err = bodyToSliceSchema(block.Bdy)
		case "items":
			// only for combined; for array of combined, it would be in attribute "items"
			s, err := newCombinedFromBody(block.Bdy)
			if err != nil {
				return nil, err
			}
			schema.Items = NewCombinedOrCombinedArrayWithCombined(s)
		case "additionalItems":
			schema.AdditionalItems, err = newCombinedFromBody(block.Bdy)
		case "contains":
			schema.Contains, err = newCombinedFromBody(block.Bdy)
		case "additionalProperties":
			schema.AdditionalProperties, err = newCombinedFromBody(block.Bdy)
		case "propertyNames":
			schema.PropertyNames, err = newCombinedFromBody(block.Bdy)
		case "if":
			schema.If, err = newCombinedFromBody(block.Bdy)
		case "then":
			schema.Then, err = newCombinedFromBody(block.Bdy)
		case "else":
			schema.Else, err = newCombinedFromBody(block.Bdy)
		case "not":
			schema.Not, err = newCombinedFromBody(block.Bdy)
		default:
			// ignore
		}
		if err != nil {
			return nil, err
		}
	}

	return schema, nil
}

func newCombinedFromBody(body *light.Body) (*Combined, error) {
	if body == nil {
		return nil, nil
	}

	schema, err := NewSchemaFromBody(body)
	if err != nil {
		return nil, err
	}
	return NewCombinedWithSchema(schema), nil
}

func expressionToYamlNode(expr *light.Expression) *yaml.Node {
	switch expr.ExpressionClause.(type) {
	case *light.Expression_Lvexpr:
		return &yaml.Node{
			Kind:  yaml.ScalarNode,
			Value: fmt.Sprintf("%v", light.LiteralValueExprToInterface(expr)),
		}
	case *light.Expression_Texpr:
		return &yaml.Node{
			Kind:  yaml.ScalarNode,
			Value: *light.TextValueExprToString(expr),
		}
	default:
	}
	return nil
}

func expressionToStringOrStringArray(expr *light.Expression) *StringOrStringArray {
	if expr == nil {
		return nil
	}

	switch expr.ExpressionClause.(type) {
	case *light.Expression_Texpr:
		return &StringOrStringArray{
			String: light.TextValueExprToString(expr),
		}
	default:
	}
	x := light.TupleConsExprToStringArray(expr)
	return &StringOrStringArray{
		StringArray: &x,
	}
}

func objectConsExprToMapCombined(o *light.ObjectConsExpr) ([]*NamedCombined, error) {
	if o == nil {
		return nil, nil
	}
	var m []*NamedCombined
	for _, item := range o.Items {
		k := light.KeyValueExprToString(item.KeyExpr)
		v, err := expressionToCombined(item.ValueExpr)
		if err != nil {
			return nil, err
		}
		m = append(m, &NamedCombined{Name: *k, Value: v})
	}
	return m, nil
}

func tupleConsExprToSlice(t *light.TupleConsExpr) ([]*Combined, error) {
	if t == nil {
		return nil, nil
	}
	exprs := t.Exprs
	if len(exprs) == 0 {
		return nil, nil
	}

	var items []*Combined
	for _, expr := range exprs {
		s, err := expressionToCombined(expr)
		if err != nil {
			return nil, err
		}
		items = append(items, s)
	}

	return items, nil
}

func expressionToCombinedOrCombinedArray(expr *light.Expression) (*CombinedOrCombinedArray, error) {
	if expr.GetTcexpr() != nil {
		items, err := tupleConsExprToSlice(expr.GetTcexpr())
		if err != nil {
			return nil, err
		}
		return NewCombinedOrCombinedArrayWithCombinedArray(items), nil
	} else {
		s, err := expressionToCombined(expr)
		if err != nil {
			return nil, err
		}
		return NewCombinedOrCombinedArrayWithCombined(s), nil
	}
}

func tupleConsExprToEnum(t *light.TupleConsExpr) ([]SchemaEnumValue, error) {
	if t == nil {
		return nil, nil
	}
	exprs := t.Exprs
	if len(exprs) == 0 {
		return nil, nil
	}
	var enums []SchemaEnumValue
	for _, expr := range exprs {
		switch expr.ExpressionClause.(type) {
		case *light.Expression_Texpr:
			enums = append(enums, SchemaEnumValue{String: light.TextValueExprToString(expr)})
		case *light.Expression_Lvexpr:
			enums = append(enums, SchemaEnumValue{Bool: light.LiteralValueExprToBoolean(expr)})
		default:
		}
	}
	return enums, nil
}

func bodyToMapCombined(b *light.Body) ([]*NamedCombined, error) {
	if b == nil {
		return nil, nil
	}
	var m []*NamedCombined
	for k, v := range b.Attributes {
		s, err := expressionToCombined(v.Expr)
		if err != nil {
			return nil, err
		}
		m = append(m, &NamedCombined{Name: k, Value: s})
	}
	for _, block := range b.Blocks {
		s, err := NewSchemaFromBody(block.Bdy)
		if err != nil {
			return nil, err
		}
		combined := NewCombinedWithSchema(s)
		m = append(m, &NamedCombined{Name: block.Type, Value: combined})
	}

	return m, nil
}

func bodyCombinedOrStringArray(label string, b *light.Body) (*NamedCombinedOrStringArray, error) {
	if b == nil {
		return nil, nil
	}

	for k, v := range b.Attributes {
		// label is "" if body is attribute only
		if v.Expr.GetTcexpr() != nil {
			stringArray := light.TupleConsExprToStringArray(v.Expr)
			combinedOrStringArray := NewCombinedOrStringArrayWithStringArray(stringArray)
			return &NamedCombinedOrStringArray{Name: k, Value: combinedOrStringArray}, nil
		} else {
			s, err := expressionToCombined(v.Expr)
			if err != nil {
				return nil, err
			}
			return &NamedCombinedOrStringArray{Name: k, Value: NewCombinedOrStringArrayWithCombined(s)}, nil
		}
	}

	s, err := NewSchemaFromBody(b)
	if err != nil {
		return nil, err
	}
	v := NewCombinedOrStringArrayWithCombined(NewCombinedWithSchema(s))
	return &NamedCombinedOrStringArray{Name: label, Value: v}, nil
}

func bodyToMapCombinedOrStringArray(b *light.Body) ([]*NamedCombinedOrStringArray, error) {
	if b == nil {
		return nil, nil
	}
	var m []*NamedCombinedOrStringArray
	for k, v := range b.Attributes {
		if v.Expr.GetTcexpr() != nil {
			stringArray := light.TupleConsExprToStringArray(v.Expr)
			combinedOrStringArray := NewCombinedOrStringArrayWithStringArray(stringArray)
			m = append(m, &NamedCombinedOrStringArray{Name: k, Value: combinedOrStringArray})
		} else {
			s, err := expressionToCombined(v.Expr)
			if err != nil {
				return nil, err
			}
			m = append(m, &NamedCombinedOrStringArray{Name: k, Value: NewCombinedOrStringArrayWithCombined(s)})
		}
	}
	for _, block := range b.Blocks {
		s, err := NewSchemaFromBody(block.Bdy)
		if err != nil {
			return nil, err
		}
		v := NewCombinedOrStringArrayWithCombined(NewCombinedWithSchema(s))
		m = append(m, &NamedCombinedOrStringArray{Name: block.Type, Value: v})
	}
	return m, nil
}

func expressionToSchema(expr *light.Expression) (*Schema, error) {
	if expr == nil {
		return nil, nil
	}

	switch expr.ExpressionClause.(type) {
	case *light.Expression_Stexpr:
		ref, err := expressionToReference(expr)
		if err != nil {
			return nil, err
		}
		return &Schema{
			Ref: &ref,
		}, nil
	//case *light.Expression_Fcexpr:
	//	return fcexprToSchema(expr.GetFcexpr())
	case *light.Expression_Ocexpr:
		body := expr.GetOcexpr().ToBody()
		return NewSchemaFromBody(body)
	default:
	}

	return nil, fmt.Errorf("not supported expression: %#v", expr)
}

func expressionToCombined(expr *light.Expression) (*Combined, error) {
	if expr == nil {
		return nil, nil
	}
	schemaOrBoolean := &Combined{}
	switch expr.ExpressionClause.(type) {
	case *light.Expression_Lvexpr:
		v := light.LiteralValueExprToBoolean(expr)
		schemaOrBoolean.Boolean = v
	default:
		schema, err := expressionToSchema(expr)
		if err != nil {
			return nil, err
		}
		schemaOrBoolean.Schema = schema
	}
	return schemaOrBoolean, nil
}
