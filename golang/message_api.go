package golang

import (
	"github.com/pb33f/libopenapi/datamodel/high/base"
	v3 "github.com/pb33f/libopenapi/datamodel/high/v3"
	"github.com/pb33f/libopenapi/orderedmap"

	"columba-livia/common"
	c "columba-livia/content"
)

func message(
	tag *base.Tag,
	pathItems *orderedmap.Map[string, *v3.PathItem],
) (render render) {
	render = func() c.C {
		list := make([]c.C, 0)

		for _, op := range common.TagOperationList(tag.Name, pathItems) {
			list = append(list,
				messageURI(op),
				messageQry(op),
				messageReq(op),
				messageRsp(op),
			)
		}

		return c.List(1, list...)
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

	list := make([]c.C, 0)

	for _, par := range op.URI {
		name := par.Name
		type_ := baseType(par.Schema)

		list = append(list,
			c.List(0,
				// 文档
				doc("", par.Description),
				// 声明
				c.F("{{.field}} {{.type}} `uri:\"{{.name}}\" binding:\"required\"`").Format(map[string]any{
					"field": publicName(name),
					"type":  type_,
					"name":  name,
				}),
			),
		)
	}

	return c.F(`type {{.name}} = struct {{.body}}`).Format(map[string]any{
		"name": publicName(op.Tag + publicName(op.ID) + "Uri"),
		"body": c.BodyF(
			c.List(1, list...).IndentTab(1),
		),
	})
}

// 																				qry

func messageQry(
	op *common.Operation,
) c.C {
	if len(op.Qry) == 0 {
		return ""
	}

	list := make([]c.C, 0)

	for _, par := range op.Qry {
		name := par.Name
		type_ := baseType(par.Schema)

		requiredTag := ""
		if common.ParRequired(par) {
			requiredTag = ` binding:"required"`
		} else {
			type_ = "*" + type_

		}

		list = append(list,
			c.List(0,
				// 文档
				doc("", par.Description),
				// 声明
				c.F("{{.field}} {{.type}} `form:\"{{.name}}\"{{.tag}}`").Format(map[string]any{
					"field": publicName(name),
					"name":  name,
					"type":  type_,
					"tag":   requiredTag,
				}),
			),
		)
	}

	return c.F(`type {{.name}} = struct {{.body}}`).Format(map[string]any{
		"name": publicName(op.Tag + publicName(op.ID) + "Qry"),
		"body": c.BodyF(c.List(1, list...).IndentTab(1)),
	})
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
		publicName(op.Tag+publicName(op.ID)+"Req"),
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
		name := publicName(op.Tag + publicName(op.ID) + "Rsp" + rspSchemaProxy.RspCode)

		if rspSchemaProxy.ContentType == common.ContentEmpty {
			list = append(list,
				c.F(`type {{.}} = struct {}`).Format(name),
			)
		} else {
			list = append(list,
				typeDecl(
					name,
					rspSchemaProxy.SchemaProxy,
				),
			)
		}
	}

	return c.List(1, list...)
}
