package common

import (
	"github.com/pb33f/libopenapi/datamodel/high/base"
)

// SchemaType //
// openapi 3.1 schema 为 string / []string
// 仅支持单一类型
func SchemaType(
	schema *base.Schema,
) (t string) {
	if schema == nil {
		return ""
	}
	schemaType := schema.Type

	switch len(schemaType) {
	case 0:
		return ""
	case 1:
		return schemaType[0]
	default:
		panic("schema type len > 1")
	}
}

// SchemaItems //
// schema item 可能为 SchemaProxy / bool
func SchemaItems(
	schema *base.Schema,
) *base.SchemaProxy {
	items := schema.Items
	if items.IsA() {
		return items.A
	}

	panic("schema items is bool")
}

func SchemaRef(
	schemaProxy *base.SchemaProxy,
) string {
	if schemaProxy.IsReference() {
		return refComponentsSchemas(schemaProxy.GetReference())
	}

	return ""

}
