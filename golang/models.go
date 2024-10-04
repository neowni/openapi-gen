package golang

import (
	"slices"

	"github.com/pb33f/libopenapi/datamodel/high/base"
	"github.com/pb33f/libopenapi/orderedmap"
	"gopkg.in/yaml.v3"

	"columba-livia/common"
	c "columba-livia/content"
)

func models(
	schemas *orderedmap.Map[string, *base.SchemaProxy],
) (render render) {
	render = func() c.C {
		list := make([]c.C, 0)

		// 渲染 类型声明
		for _, pair := range common.Range(schemas) {
			typeDecl := typeDecl(pair.Key(), pair.Value())
			list = append(list, typeDecl)
		}

		return c.List(1, list...)
	}
	return render
}

// 																				json 类型

func typeMap(jsonType string) (type_ string) {
	switch jsonType {
	case common.TypeBoolean:
		return "bool"

	case common.TypeInteger:
		return "int"

	case common.TypeNumber:
		return "float64"

	case common.TypeString:
		return "string"

	default:
		return ""
	}
}

//																				类型声明

func typeDecl(
	name string,
	schemaProxy *base.SchemaProxy,
) c.C {
	// 枚举类型
	if common.SchemaRef(schemaProxy) == "" &&
		len(schemaProxy.Schema().Enum) != 0 {

		jsonType := common.SchemaType(schemaProxy.Schema())
		if jsonType != common.TypeString {
			panic(jsonType)
		}

		return c.List(1,
			c.F("type {{.}} = string").Format(publicName(name)),

			c.F(`const {{.}}`).Format(c.BodyC(
				c.List(0, c.ForList(
					schemaProxy.Schema().Enum,
					func(item *yaml.Node) c.C {
						value := item.Value

						return c.F(`{{.const}} {{.type}} = "{{.value}}"`).Format(map[string]any{
							"const": publicName(name) + publicName(value),
							"type":  publicName(name),
							"value": value,
						})
					},
				)...).IndentTab(1),
			)),
		)
	}

	// 其他类型
	return c.F("type {{.name}} = {{.type}}").Format(map[string]any{
		"name": publicName(name),
		"type": type_(schemaProxy),
	})
}

//																				类型

func type_(
	schemaProxy *base.SchemaProxy,
) c.C {
	// 引用类型
	ref := common.SchemaRef(schemaProxy)
	if ref != "" {
		return c.C(file.modelsNamespace() + publicName(ref))
	}

	jsonType := common.SchemaType(schemaProxy.Schema())

	// 基础类型
	golangType := typeMap(jsonType)
	if golangType != "" {
		return c.C(golangType)
	}

	// object
	if jsonType == common.TypeObject {
		return object(schemaProxy)
	}

	// array
	if jsonType == common.TypeArray {
		return c.F("[]{{.}}").Format(
			type_(common.SchemaItems(schemaProxy.Schema())),
		)
	}

	panic(jsonType)
}

func baseType(
	schemaProxy *base.SchemaProxy,
) c.C {
	jsonType := common.SchemaType(schemaProxy.Schema())

	golangType := typeMap(jsonType)
	if golangType != "" {
		return c.C(golangType)
	}

	// 类型不支持
	panic("")
}

//																				object

func object(
	schemaProxy *base.SchemaProxy,
) c.C {
	requiredNameList := schemaProxy.Schema().Required

	// 整理 字段
	fieldList := make([]c.C, 0)
	for _, pair := range common.Range(schemaProxy.Schema().Properties) {
		name := pair.Key()
		type_ := type_(pair.Value())

		// 必要字段 信息
		requiredTag := ""
		omitempty := ""

		required := slices.Contains(requiredNameList, name)
		if required {
			requiredTag = ` binding:"required"`
		} else {
			type_ = "*" + type_
			omitempty = ",omitempty"
		}

		fieldList = append(
			fieldList,
			c.List(0,
				// 文档
				doc("", pair.Value().Schema().Description),

				// 声明
				c.F("{{.field}} {{.type}} `json:\"{{.name}}{{.omitempty}}\"{{.requiredTag}}`").Format(map[string]any{
					"field":       publicName(name),
					"type":        type_,
					"name":        name,
					"omitempty":   omitempty,
					"requiredTag": requiredTag,
				}),
			),
		)
	}

	return c.F("struct {{.}}").Format(
		c.BodyF(
			c.List(1, fieldList...).IndentTab(1),
		),
	)
}
