package python

import (
	"github.com/pb33f/libopenapi/datamodel/high/base"
	v3 "github.com/pb33f/libopenapi/datamodel/high/v3"
	"github.com/pb33f/libopenapi/orderedmap"

	"columba-livia/common"
	c "columba-livia/content"
)

func messageAPI(
	tag *base.Tag,
	pathItems *orderedmap.Map[string, *v3.PathItem],
) (render render) {
	render = func() c.C {
		list := make([]c.C, 2)

		for _, op := range common.TagOperationList(tag.Name, pathItems) {
			list = append(list,
				c.List(1,
					// 使用 class 用作 命名空间
					c.List(0,
						c.C("class %s:").Format(op.ID),
						c.C(`"""%s"""`).Format(op.Description).IndentSpace(4),
					),
					// message 定义
					c.List(1,
						messageURI(op),
						messageQry(op),
						messageReq(op),
						messageRsp(op),
					).IndentSpace(4),
				),
			)
		}

		return c.List(2, list...)
	}
	return render
}

// 																				uri

func messageURI(
	op *common.Operation,
) c.C {
	if len(op.URI) == 0 {
		return ""
	}

	fieldList := make([]c.C, 0)

	for _, par := range op.URI {
		name := par.Name
		typeName := baseTypeName(par.Schema)

		fieldList = append(fieldList,
			c.List(0,
				// 文档
				doc(par.Description),
				// 声明
				c.C("%s: %s").Format(name, typeName),
			),
		)
	}

	return c.List(1,
		c.C(`class uri(_pydantic.BaseModel):`),
		c.List(1, fieldList...).IndentSpace(4),
	)
}

// 																				qry

func messageQry(
	op *common.Operation,
) c.C {
	if len(op.Qry) == 0 {
		return ""
	}

	fieldList := make([]c.C, 0)

	for _, par := range op.Qry {
		name := par.Name
		typeName := baseTypeName(par.Schema)

		if !common.ParRequired(par) {
			file.importMap["import typing as _typing"] = struct{}{}
			typeName = c.C("_typing.Optional[%s] = None").Format(typeName)
		}

		// 字段
		fieldList = append(fieldList,
			c.List(0,
				// 文档
				doc(par.Description),
				// 声明
				c.C("%s: %s").Format(name, typeName),
			),
		)
	}

	return c.List(1,
		c.C(`class qry(_pydantic.BaseModel):`),
		c.List(1, fieldList...).IndentSpace(4),
	)
}

// 																				req

func messageReq(
	op *common.Operation,
) c.C {
	reqSchemaProxy := common.ReqSchemaProxy(op.Req)

	// 空类型
	if reqSchemaProxy.ContentType == common.ContentEmpty {
		return ""
	}

	return typeDecl(
		"req",
		reqSchemaProxy.SchemaProxy,
	)
}

// 																				rsp

func messageRsp(
	op *common.Operation,
) c.C {
	list := make([]c.C, 0)

	rspSchemaProxyList := common.RspSchemaProxy(op.Rsp)
	for _, rspSchemaProxy := range rspSchemaProxyList {
		name := "rsp" + rspSchemaProxy.RspCode

		if rspSchemaProxy.ContentType == common.ContentEmpty {
			// 空类型
			list = append(list,
				c.C("%s = str").Format(name),
			)
		} else {
			// 引用
			refName := common.SchemaRef(rspSchemaProxy.SchemaProxy)

			// 基础类型
			jsonType := common.SchemaType(rspSchemaProxy.SchemaProxy.Schema())
			pythonType := typeMap(jsonType)

			if refName == "" && pythonType != "" {
				// 非引用 且 为基础类型
				return c.C("%s = %s").Format(name, pythonType)
			} else {
				// 其他类型
				list = append(list,
					typeDecl(
						name,
						rspSchemaProxy.SchemaProxy,
					),
				)
			}
		}
	}

	return c.List(1, list...)
}
