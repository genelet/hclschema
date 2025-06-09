package jsm07

import (
	"encoding/json"
	"fmt"

	"github.com/genelet/hcllight/light"
)

// UnmarshalHCL unmarshals HCL data into a Schema object.
func (self *Schema) UnmarshalHCL(data []byte, labels ...string) error {
	s, err := ParseSchema(data)
	if err != nil {
		return err
	}
	*self = *s
	return nil
}

// parseSchema parses a HCL string representing a JSON schema and returns a Schema object.
func ParseSchema(data []byte) (*Schema, error) {
	body, err := light.ParseBody(data)
	if err != nil {
		return nil, err
	}
	return parseSchemaFromBody(body)
}

func parseSchemaFromBody(body *light.Body) (*Schema, error) {
	if body == nil {
		return nil, nil
	}

	var typ *StringOrStringArray
	if attr, ok := body.Attributes["type"]; ok {
		typ = exprToStringOrStringArray(attr.Expr)
	}

	schema := &Schema{}
	for k, v := range body.Attributes {
		expr := v.Expr
		var err error
		switch k {
		case "_id":
			schema.ID = light.TextValueExprToString(expr)
		case "_schema":
			schema.Schema = light.TextValueExprToString(expr)
		case "_ref":
			schema.Ref = light.TextValueExprToString(expr)
		case "_comment":
			schema.Comment = light.TextValueExprToString(expr)
		case "title":
			schema.Title = light.TextValueExprToString(expr)
		case "description":
			schema.Description = light.TextValueExprToString(expr)
		case "format":
			schema.Format = light.TextValueExprToString(expr)
		case "contentMediaType":
			schema.ContentMediaType = light.TextValueExprToString(expr)
		case "contentEncoding":
			schema.ContentEncoding = light.TextValueExprToString(expr)
		case "pattern":
			schema.Pattern = light.TextValueExprToString(expr)

		case "maxLength":
			schema.MaxLength = light.LiteralValueExprToInt64(expr)
		case "minLength":
			schema.MinLength = light.LiteralValueExprToInt64(expr)
		case "maxItems":
			schema.MaxItems = light.LiteralValueExprToInt64(expr)
		case "minItems":
			schema.MinItems = light.LiteralValueExprToInt64(expr)
		case "maxProperties":
			schema.MaxProperties = light.LiteralValueExprToInt64(expr)
		case "minProperties":
			schema.MinProperties = light.LiteralValueExprToInt64(expr)

		case "readOnly":
			schema.ReadOnly = light.LiteralValueExprToBoolean(expr)
		case "writeOnly":
			schema.WriteOnly = light.LiteralValueExprToBoolean(expr)
		case "uniqueItems":
			schema.UniqueItems = light.LiteralValueExprToBoolean(expr)

		case "const":
			schema.Const = exprToJSONRaw(expr)
		case "default":
			schema.Default = exprToJSONRaw(expr)
		case "example":
			schema.Examples = exprToJSONRaw(expr)

		case "type":
			schema.Type = typ
		case "multipleOf":
			schema.MultipleOf = exprToIntegerOrFloat(expr, typ)
		case "maximum":
			schema.Maximum = exprToIntegerOrFloat(expr, typ)
		case "exclusiveMaximum":
			schema.ExclusiveMaximum = exprToIntegerOrFloat(expr, typ)
		case "minimum":
			schema.Minimum = exprToIntegerOrFloat(expr, typ)
		case "exclusiveMinimum":
			schema.ExclusiveMinimum = exprToIntegerOrFloat(expr, typ)

		case "additionalItems":
			schema.AdditionalItems, err = expressionToCombined(expr)
		case "propertyNames":
			schema.PropertyNames, err = expressionToCombined(expr)
		case "additionalProperties":
			schema.AdditionalProperties, err = expressionToCombined(expr)
		case "contains":
			schema.Contains, err = expressionToCombined(expr)
		case "if":
			schema.If, err = expressionToCombined(expr)
		case "then":
			schema.Then, err = expressionToCombined(expr)
		case "else":
			schema.Else, err = expressionToCombined(expr)
		case "not":
			schema.Not, err = expressionToCombined(expr)

		case "required":
			schema.Required = light.TupleConsExprToStringArray(expr)
		case "enum":
			schema.Enumeration, err = tupleConsExprToEnum(expr)

		default:
			// ignore
		}
		if err != nil {
			return nil, err
		}

	}

	var combs []*Combined
	props := make(map[string]*Combined)
	var nullProps int

	for _, block := range body.Blocks {
		switch block.Type {
		case "items", "additionalProperties", "additionalItems", "propertyNames", "contains", "if", "then", "else", "not":
			c, err := newCombinedFromBody(block.Bdy)
			if err != nil {
				return nil, err
			}
			switch block.Type {
			case "items":
				combs = append(combs, c)
			case "additionalProperties":
				schema.AdditionalProperties = c
			case "additionalItems":
				schema.AdditionalItems = c
			case "propertyNames":
				schema.PropertyNames = c
			case "contains":
				schema.Contains = c
			case "if":
				schema.If = c
			case "then":
				schema.Then = c
			case "else":
				schema.Else = c
			case "not":
				schema.Not = c
			}

		case "allOf", "anyOf", "oneOf":
			combined, err := newCombinedFromBody(block.Bdy)
			if err != nil {
				return nil, err
			}
			switch block.Type {
			case "allOf":
				schema.AllOf = append(schema.AllOf, combined)
			case "anyOf":
				schema.AnyOf = append(schema.AnyOf, combined)
			case "oneOf":
				schema.OneOf = append(schema.OneOf, combined)
			}

		case "dependencies":
			combinedOrStringArray, err := bodyCombinedOrStringArray(block.Bdy)
			if err != nil {
				return nil, err
			}
			if schema.Dependencies == nil {
				schema.Dependencies = make(map[string]*CombinedOrStringArray)
			}
			schema.Dependencies[block.Labels[0]] = combinedOrStringArray

		case "definitions", "properties", "patternProperties":
			combined, err := newCombinedFromBody(block.Bdy)
			if err != nil {
				return nil, err
			}
			switch block.Type {
			case "definitions":
				if schema.Definitions == nil {
					schema.Definitions = make(map[string]*Combined)
				}
				schema.Definitions[block.Labels[0]] = combined
			case "patternProperties":
				if schema.PatternProperties == nil {
					schema.PatternProperties = make(map[string]*Combined)
				}
				schema.PatternProperties[block.Labels[0]] = combined
			case "properties":
				if len(block.Labels) == 0 {
					nullProps++
				} else {
					props[block.Labels[0]] = combined
				}
			}

		default:
			// ignore
		}
	}

	if len(combs) > 1 {
		schema.Items = NewCombinedOrCombinedArrayWithCombinedArray(combs)
	} else if len(combs) == 1 {
		schema.Items = NewCombinedOrCombinedArrayWithCombined(combs[0])
	}

	if len(props) > 0 {
		schema.Properties = props
	} else if nullProps > 0 {
		schema.Properties = map[string]*Combined{}
	}

	return schema, nil
}

func exprToJSONRaw(expr *light.Expression) *json.RawMessage {
	switch expr.ExpressionClause.(type) {
	case *light.Expression_Lvexpr:
		i := light.LiteralValueExprToInterface(expr)
		raw := json.RawMessage([]byte(fmt.Sprintf("%v", i)))
		return &raw
	case *light.Expression_Texpr:
		str := *light.TextValueExprToString(expr)
		raw := json.RawMessage([]byte(`"` + str + `"`))
		return &raw
	default:
	}
	return nil
}

func exprToIntegerOrFloat(expr *light.Expression, typ *StringOrStringArray) *IntegerOrFloat {
	if typ != nil && typ.String != nil {
		if *typ.String == "integer" {
			i64 := light.LiteralValueExprToInt64(expr)
			return NewIntegerOrFloatWithInteger(*i64)
		} else {
			i64 := light.LiteralValueExprToInt64(expr)
			if i64 != nil {
				return NewIntegerOrFloatWithInteger(*i64)
			}
			f64 := light.LiteralValueExprToFloat64(expr)
			return NewIntegerOrFloatWithFloat(*f64)
		}
	}
	return nil
}

func newCombinedFromBody(body *light.Body) (*Combined, error) {
	if body == nil {
		return nil, nil
	}

	schema, err := parseSchemaFromBody(body)
	if err != nil {
		return nil, err
	}
	return NewCombinedWithSchema(schema), nil
}

func exprToStringOrStringArray(expr *light.Expression) *StringOrStringArray {
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

func tupleConsExprToEnum(expr *light.Expression) ([]SchemaEnumValue, error) {
	t := expr.GetTcexpr()
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

func bodyCombinedOrStringArray(b *light.Body) (*CombinedOrStringArray, error) {
	if b == nil {
		return nil, nil
	}

	for _, v := range b.Attributes {
		// label is "" if body is attribute only
		if v.Expr.GetTcexpr() != nil {
			stringArray := light.TupleConsExprToStringArray(v.Expr)
			return NewCombinedOrStringArrayWithStringArray(stringArray), nil
		} else {
			c, err := expressionToCombined(v.Expr)
			if err != nil {
				return nil, err
			}
			return NewCombinedOrStringArrayWithCombined(c), nil
		}
	}

	s, err := parseSchemaFromBody(b)
	if err != nil {
		return nil, err
	}
	return NewCombinedOrStringArrayWithCombined(NewCombinedWithSchema(s)), nil
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
	// case *light.Expression_Fcexpr:
	//   return fcexprToSchema(expr.GetFcexpr())
	case *light.Expression_Ocexpr:
		body := expr.GetOcexpr().ToBody()
		return parseSchemaFromBody(body)
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

func expressionToReference(expr *light.Expression) (string, error) {
	// in case there is only one level of reference which is parsed as lvexpr
	if x := expr.GetLvexpr(); x != nil {
		return "#/" + x.Val.GetStringValue(), nil
	} else if x := light.TraversalToString(expr); x != nil {
		return "#/" + *x, nil
	}
	return "", fmt.Errorf("1 invalid expression: %#v", expr)
}
