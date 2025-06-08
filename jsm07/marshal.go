package jsm07

import (
	"encoding/json"
	"strings"

	"github.com/genelet/determined/dethcl"
	"github.com/genelet/hcllight/light"
)

func assignRaw(attrs map[string]*light.Attribute, key string, val *json.RawMessage) bool {
	if val == nil {
		return false
	}
	attrs[key] = &light.Attribute{
		Name: key,
		Expr: light.StringToTextValueExpr(string(*val)),
	}
	return true
}

func assignInt64(attrs map[string]*light.Attribute, key string, val *int64) bool {
	if val == nil {
		return false
	}
	attrs[key] = &light.Attribute{
		Name: key,
		Expr: light.Int64ToLiteralValueExpr(*val),
	}
	return true
}

func assignBool(attrs map[string]*light.Attribute, key string, val *bool) bool {
	if val == nil {
		return false
	}
	attrs[key] = &light.Attribute{
		Name: key,
		Expr: light.BooleanToLiteralValueExpr(*val),
	}
	return true
}

func assignIntegerOrFloat(attrs map[string]*light.Attribute, key string, val *IntegerOrFloat) bool {
	if val == nil {
		return false
	}
	if val.Float != nil {
		attrs[key] = &light.Attribute{
			Name: key,
			Expr: light.Float64ToLiteralValueExpr(*val.Float),
		}
	} else {
		attrs[key] = &light.Attribute{
			Name: key,
			Expr: light.Int64ToLiteralValueExpr(*val.Integer),
		}
	}
	return true
}

func assignString(attrs map[string]*light.Attribute, key string, val *string) bool {
	if val == nil {
		return false
	}
	if key[0] == '$' {
		key = `_` + key[1:]
	}
	attrs[key] = &light.Attribute{
		Name: key,
		Expr: light.StringToTextValueExpr(*val),
	}
	return true
}

func assignCombined(attrs map[string]*light.Attribute, key string, val *Combined) bool {
	if val == nil || val.Boolean == nil {
		return false
	}

	attrs[key] = &light.Attribute{
		Name: key,
		Expr: light.BooleanToLiteralValueExpr(*val.Boolean),
	}
	return true
}

func assignEnum(attrs map[string]*light.Attribute, key string, val []SchemaEnumValue) bool {
	if val == nil {
		return false
	}

	var exprs []*light.Expression
	for _, v := range val {
		if v.String != nil {
			exprs = append(exprs, light.StringToTextValueExpr(*v.String))
		} else if v.Bool != nil {
			exprs = append(exprs, light.BooleanToLiteralValueExpr(*v.Bool))
		} else if v.Number != nil {
			if v.Number.Float != nil {
				exprs = append(exprs, light.Float64ToLiteralValueExpr(*v.Number.Float))
			} else if v.Number.Integer != nil {
				exprs = append(exprs, light.Int64ToLiteralValueExpr(*v.Number.Integer))
			}
		} else if v.Null != nil && *v.Null {
			exprs = append(exprs, light.StringToTextValueExpr("null"))
		}
	}

	attrs[key] = &light.Attribute{
		Name: key,
		Expr: &light.Expression{
			ExpressionClause: &light.Expression_Tcexpr{
				Tcexpr: &light.TupleConsExpr{
					Exprs: exprs,
				},
			},
		},
	}

	return true
}

func (self *Schema) MarshalHCL() ([]byte, error) {
	attrs := map[string]*light.Attribute{}

	trimmed := *self

	if trimmed.Type != nil {
		if trimmed.Type.String != nil {
			attrs["type"] = &light.Attribute{
				Name: "type",
				Expr: light.StringToTextValueExpr(*trimmed.Type.String),
			}
		} else {
			attrs["type"] = &light.Attribute{
				Name: "type",
				Expr: light.StringArrayToTupleConsEpr(*trimmed.Type.StringArray),
			}
		}
		trimmed.Type = nil
	}

	if assignString(attrs, "$id", trimmed.ID) {
		trimmed.ID = nil
	}
	if assignString(attrs, "$ref", trimmed.Ref) {
		trimmed.Ref = nil
	}
	if assignString(attrs, "$schema", trimmed.Schema) {
		trimmed.Schema = nil
	}
	if assignString(attrs, "$comment", trimmed.Comment) {
		trimmed.Comment = nil
	}
	if assignString(attrs, "format", trimmed.Format) {
		trimmed.Format = nil
	}
	if assignString(attrs, "contentMediaType", trimmed.ContentMediaType) {
		trimmed.ContentMediaType = nil
	}
	if assignString(attrs, "contentEncoding", trimmed.ContentEncoding) {
		trimmed.ContentEncoding = nil
	}
	if assignString(attrs, "title", trimmed.Title) {
		trimmed.Title = nil
	}
	if assignString(attrs, "description", trimmed.Description) {
		trimmed.Description = nil
	}
	if assignString(attrs, "pattern", trimmed.Pattern) {
		trimmed.Pattern = nil
	}

	if assignInt64(attrs, "maxLength", trimmed.MaxLength) {
		trimmed.MaxLength = nil
	}
	if assignInt64(attrs, "minLength", trimmed.MinLength) {
		trimmed.MinLength = nil
	}
	if assignInt64(attrs, "maxItems", trimmed.MaxItems) {
		trimmed.MaxItems = nil
	}
	if assignInt64(attrs, "minItems", trimmed.MinItems) {
		trimmed.MinItems = nil
	}
	if assignInt64(attrs, "maxProperties", trimmed.MaxProperties) {
		trimmed.MaxProperties = nil
	}
	if assignInt64(attrs, "minProperties", trimmed.MinProperties) {
		trimmed.MinProperties = nil
	}

	if assignBool(attrs, "readOnly", trimmed.ReadOnly) {
		trimmed.ReadOnly = nil
	}
	if assignBool(attrs, "writeOnly", trimmed.WriteOnly) {
		trimmed.WriteOnly = nil
	}
	if assignBool(attrs, "uniqueItems", trimmed.UniqueItems) {
		trimmed.UniqueItems = nil
	}

	if assignRaw(attrs, "const", trimmed.Const) {
		trimmed.Const = nil
	}
	if assignRaw(attrs, "default", trimmed.Default) {
		trimmed.Default = nil
	}
	if assignRaw(attrs, "examples", trimmed.Examples) {
		trimmed.Examples = nil
	}

	if assignIntegerOrFloat(attrs, "multipleOf", trimmed.MultipleOf) {
		trimmed.MultipleOf = nil
	}
	if assignIntegerOrFloat(attrs, "maximum", trimmed.Maximum) {
		trimmed.Maximum = nil
	}
	if assignIntegerOrFloat(attrs, "exclusiveMaximum", trimmed.ExclusiveMaximum) {
		trimmed.ExclusiveMaximum = nil
	}
	if assignIntegerOrFloat(attrs, "minimum", trimmed.Minimum) {
		trimmed.Minimum = nil
	}
	if assignIntegerOrFloat(attrs, "exclusiveMinimum", trimmed.ExclusiveMinimum) {
		trimmed.ExclusiveMinimum = nil
	}

	if assignCombined(attrs, "additionalProperties", trimmed.AdditionalProperties) {
		trimmed.AdditionalProperties = nil
	}
	if assignCombined(attrs, "additionalItems", trimmed.AdditionalItems) {
		trimmed.AdditionalItems = nil
	}
	if assignCombined(attrs, "propertyNames", trimmed.PropertyNames) {
		trimmed.PropertyNames = nil
	}
	if assignCombined(attrs, "contains", trimmed.Contains) {
		trimmed.Contains = nil
	}
	if assignCombined(attrs, "if", trimmed.If) {
		trimmed.If = nil
	}
	if assignCombined(attrs, "then", trimmed.Then) {
		trimmed.Then = nil
	}
	if assignCombined(attrs, "else", trimmed.Else) {
		trimmed.Else = nil
	}
	if assignCombined(attrs, "not", trimmed.Not) {
		trimmed.Not = nil
	}

	if assignEnum(attrs, "enum", trimmed.Enumeration) {
		trimmed.Enumeration = nil
	}

	bs, err := dethcl.Marshal(trimmed)
	if err != nil {
		return nil, err
	}
	if len(attrs) == 0 {
		return bs, nil
	}

	body := &light.Body{
		Attributes: attrs,
	}
	data, err := body.MarshalHCL()
	if err != nil {
		return nil, err
	}

	if len(bs) == 0 {
		str := "  " + strings.TrimSpace(string(data))
		return []byte(str), nil
	}

	str := "  " + strings.TrimSpace(string(append(data, bs...)))
	return []byte(str), nil
}
