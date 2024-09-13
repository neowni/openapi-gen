package typescript

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
			list = append(list, c.List(1,
				doc(pair.Value().Schema().Description),
				c.C("export %s").Format(typeDecl(
					pair.Key(),
					pair.Value(),
				)),
			))
		}

		return c.List(1, list...)
	}

	return render
}

// 																				typeDecl 类型声明

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

		return c.C(`enum %s %s`).Format(
			name,
			c.BodyF(
				c.List(0, c.ForList(
					schemaProxy.Schema().Enum,
					func(item *yaml.Node) c.C {
						value := item.Value

						return c.C(`%s = "%s",`).Format(
							value, value,
						)
					},
				)...).IndentSpace(2),
			),
		)
	}

	// 其他类型
	return c.C("type %s = %s;").Format(name, type_(schemaProxy))
}

//                                                                              类型

func type_(
	schemaProxy *base.SchemaProxy,
) c.C {
	// 引用
	if schemaProxy.IsReference() {
		return c.C(file.modelsNamespace() + common.SchemaRef(schemaProxy))
	}

	// 基础类型
	jsonType := common.SchemaType(schemaProxy.Schema())
	typescriptType := typeMap(jsonType)
	if typescriptType != "" {
		return c.C(typescriptType)
	}

	// 数组类型
	if jsonType == common.TypeArray {
		return c.C("%s[]").Format(type_(common.SchemaItems(schemaProxy.Schema())))
	}

	// object 类型
	if jsonType == common.TypeObject {
		return object(schemaProxy)
	}

	// 类型不支持
	panic("")
}

func baseTypeName(
	schemaProxy *base.SchemaProxy,
) c.C {
	jsonType := common.SchemaType(schemaProxy.Schema())

	typescriptType := typeMap(jsonType)
	if typescriptType != "" {
		return c.C(typescriptType)
	}

	// 类型不支持
	panic("")
}

func object(
	schemaProxy *base.SchemaProxy,
) c.C {
	fieldList := make([]c.C, 0)

	requiredNameList := schemaProxy.Schema().Required

	// 添加 field
	for _, pair := range common.Range(schemaProxy.Schema().Properties) {
		name := pair.Key()
		type_ := type_(pair.Value())

		required := ""
		if !slices.Contains(requiredNameList, name) {
			required = "?"
		}

		fieldList = append(fieldList, c.List(1,
			doc(pair.Value().Schema().Description),
			c.C("%s%s: %s;").Format(
				name, required, type_,
			),
		))
	}

	return c.BodyF(
		c.List(1, fieldList...).IndentSpace(2),
	)
}

//																				到 typescript 的类型转换

func typeMap(jsonType string) (typescript string) {
	switch jsonType {
	case common.TypeBoolean:
		return "boolean"

	case common.TypeInteger:
		return "number"

	case common.TypeNumber:
		return "number"

	case common.TypeString:
		return "string"

	default:
		return ""
	}
}
