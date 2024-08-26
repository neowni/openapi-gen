package golang

import (
	"columba-livia/common"
	c "columba-livia/content"

	"github.com/pb33f/libopenapi/datamodel/high/base"
	v3 "github.com/pb33f/libopenapi/datamodel/high/v3"
	"github.com/pb33f/libopenapi/orderedmap"
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
				c.JoinSpace(
					c.C("%s %s").Format(
						ExportName(name),
						type_,
					),
					c.C("`uri:\"%s\" binding:\"required\"`").Format(
						name,
					),
				),
			),
		)
	}

	return c.C(`type %s = struct %s`).Format(
		ExportName(op.ID+"URI"),
		c.BodyF(
			c.List(1, list...).Indent(4),
		),
	)
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
				c.JoinSpace(
					c.C("%s %s").Format(
						ExportName(name),
						type_,
					),
					c.C("`form:\"%s\"%s`").Format(
						name,
						requiredTag,
					),
				),
			),
		)
	}

	return c.C(`type %s = struct %s`).Format(
		ExportName(op.ID+"Qry"),
		c.BodyF(
			c.List(1, list...).Indent(4),
		),
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
		ExportName(op.ID+"Req"),
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
		name := op.ID + "Rsp" + rspSchemaProxy.RspCode

		if rspSchemaProxy.ContentType == common.ContentEmpty {
			list = append(list,
				c.C(`type %s = struct {}`).Format(
					ExportName(name),
				),
			)
		} else {
			list = append(list,
				typeDecl(
					ExportName(name),
					rspSchemaProxy.SchemaProxy,
				),
			)
		}
	}

	return c.List(1, list...)
}
