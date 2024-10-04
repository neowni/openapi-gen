package python

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
			list = append(list, typeDecl(
				pair.Key(),
				pair.Value(),
			))
		}

		return c.List(2, list...)
	}

	return render
}

//                                                                              类型声明

func typeDecl(
	name string,
	schemaProxy *base.SchemaProxy,
) c.C {
	jsonType := common.SchemaType(schemaProxy.Schema())

	// 非引用
	if common.SchemaRef(schemaProxy) == "" {
		// object
		// 直接定义 class
		if jsonType == common.TypeObject {
			return object(name, schemaProxy)
		}

		// enum
		if len(schemaProxy.Schema().Enum) != 0 {
			if jsonType != common.TypeString {
				panic(jsonType)
			}

			return enum(name, schemaProxy)
		}
	}

	// 其他类型
	return c.F("{{.name}} = {{.type}}").Format(map[string]any{
		"name": name,
		"type": typeName(schemaProxy),
	})
}

//                                                                              类型

// 提供类型名
// 如果需要进行额外的定义，例如 不使用引用的 object[] ，则在文件处进行额外的私有定义
func typeName(
	schemaProxy *base.SchemaProxy,
) c.C {
	// 引用
	ref := common.SchemaRef(schemaProxy)
	if ref != "" {
		return c.C(file.modelsNamespace() + ref)
	}

	jsonType := common.SchemaType(schemaProxy.Schema())

	// 基础类型
	pythonType := strictTypeMap(jsonType)
	if pythonType != "" {
		return c.C(pythonType)
	}

	// object 类型
	if jsonType == common.TypeObject {
		// 定义一个私有类用于提供名称
		privateName := privateName()
		file.additional = append(file.additional, object(privateName, schemaProxy))

		return c.C(privateName)
	}

	// array 类型
	if jsonType == common.TypeArray {
		itemTypeName := typeName(common.SchemaItems(schemaProxy.Schema()))

		file.importMap["import typing as _typing"] = struct{}{}
		return c.F("_typing.List[{{.}}]").Format(itemTypeName)
	}

	panic("")
}

func baseTypeName(
	schemaProxy *base.SchemaProxy,
) c.C {
	jsonType := common.SchemaType(schemaProxy.Schema())

	// 基础类型
	pythonType := typeMap(jsonType)
	if pythonType != "" {
		return c.C(pythonType)
	}

	panic("")
}

// 																				enum

func enum(
	name string,
	schemaProxy *base.SchemaProxy,
) c.C {
	file.importMap["from enum import Enum as _Enum"] = struct{}{}

	return c.List(0,
		c.List(0,
			c.F(`class {{.}}(_Enum):`).Format(name),
			doc(schemaProxy.Schema().Description).IndentSpace(4),
		),
		c.List(0, c.ForList(
			schemaProxy.Schema().Enum,
			func(item *yaml.Node) c.C {
				value := item.Value

				return c.F(`{{.}} = "{{.}}"`).Format(value)
			},
		)...).IndentSpace(4),
	)
}

// 																				object - class类定义

func object(
	name string,
	schemaProxy *base.SchemaProxy,
) c.C {
	fieldList := make([]c.C, 0)

	requiredNameList := schemaProxy.Schema().Required
	for _, pair := range common.Range(schemaProxy.Schema().Properties) {
		name := pair.Key()
		typeName := typeName(pair.Value())

		required := slices.Contains(requiredNameList, name)
		if !required {
			file.importMap["import typing as _typing"] = struct{}{}
			typeName = c.F("_typing.Optional[{{.}}] = None").Format(typeName)
		}

		// 字段
		fieldList = append(fieldList,
			c.List(0,
				// 文档
				comment(pair.Value().Schema().Description),
				// 声明
				c.F("{{.name}}: {{.type}}").Format(map[string]any{
					"name": name,
					"type": typeName,
				}),
			),
		)
	}

	return c.List(1,
		c.List(0,
			c.F(`class {{.}}(_pydantic.BaseModel):`).Format(name),
			doc(schemaProxy.Schema().Description).IndentSpace(4),
		),
		c.List(1, fieldList...).IndentSpace(4),
	)
}

//																				到 python 的类型转换

func typeMap(jsonType string) (pythonType string) {
	switch jsonType {
	case common.TypeBoolean:
		return "bool"

	case common.TypeInteger:
		return "int"

	case common.TypeNumber:
		return "float"

	case common.TypeString:
		return "str"

	default:
		return ""
	}
}

func strictTypeMap(jsonType string) (pythonType string) {
	switch jsonType {
	case common.TypeBoolean:
		file.importMap["import pydantic as _pydantic"] = struct{}{}
		return "_pydantic.StrictBool"

	case common.TypeInteger:
		file.importMap["import pydantic as _pydantic"] = struct{}{}
		return "_pydantic.StrictInt"

	case common.TypeNumber:
		file.importMap["import pydantic as _pydantic"] = struct{}{}
		return "_pydantic.StrictFloat"

	case common.TypeString:
		file.importMap["import pydantic as _pydantic"] = struct{}{}
		return "_pydantic.StrictStr"

	default:
		return ""
	}
}
