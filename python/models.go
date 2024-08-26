package python

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

	// object / 非引用
	// 直接定义 class
	if !schemaProxy.IsReference() && jsonType == common.TypeObject {
		return object(name, schemaProxy)
	}

	// 其他类型
	return c.C("%s = %s").Format(name, typeName(
		schemaProxy,
	))
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
		return c.C("_typing.List[%s]").Format(itemTypeName)
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
			typeName = c.C("_typing.Optional[%s] = None").Format(typeName)
		}

		// 字段
		fieldList = append(fieldList,
			c.List(0,
				// 文档
				doc(pair.Value().Schema().Description),
				// 声明
				c.C("%s: %s").Format(name, typeName),
			),
		)
	}

	return c.List(1,
		c.List(0,
			c.C(`class %s(_pydantic.BaseModel):`).Format(name),
			c.C(`"""%s"""`).Format(schemaProxy.Schema().Description).Indent(4),
		),
		c.List(1, fieldList...).Indent(4),
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
