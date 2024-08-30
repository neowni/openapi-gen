package golang

import (
	"columba-livia/common"
	c "columba-livia/content"
	"slices"

	"github.com/pb33f/libopenapi/datamodel/high/base"
	"github.com/pb33f/libopenapi/orderedmap"
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
	return c.C("type %s = %s").Format(
		ExportName(name),
		type_(schemaProxy),
	)
}

//																				类型

func type_(
	schemaProxy *base.SchemaProxy,
) c.C {
	// 引用类型
	ref := common.SchemaRef(schemaProxy)
	if ref != "" {
		return c.C(file.modelsNamespace() + ExportName(ref))
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
		return c.C("[]%s").Format(
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
				c.List(-1,
					c.C("%s %s").Format(
						ExportName(name),
						type_,
					),
					c.C(" `json:\"%s%s\"%s`").Format(
						name, omitempty,
						requiredTag,
					),
				),
			),
		)
	}

	return c.C("struct %s").Format(
		c.BodyF(
			c.List(1, fieldList...).Indent(4),
		),
	)
}
