package typescript

import (
	"columba-livia/common"
	c "columba-livia/content"
	"fmt"

	"github.com/pb33f/libopenapi/datamodel/high/base"
	v3 "github.com/pb33f/libopenapi/datamodel/high/v3"
	"github.com/pb33f/libopenapi/orderedmap"
)

func messageApi(
	tag *base.Tag,
	pathItems *orderedmap.Map[string, *v3.PathItem],
) (render render) {
	render = func() c.C {
		list := make([]c.C, 2)

		for _, op := range common.TagOperationList(tag.Name, pathItems) {
			list = append(list,
				messageURI(op),
				messageQry(op),
				messageReq(op),
				messageRsp(op))
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

	fieldList := make([]c.C, 0)
	for _, parameter := range op.URI {
		fieldList = append(fieldList, c.List(0,
			doc(parameter.Description),
			c.C("%s: %s;").Format(
				parameter.Name,
				baseTypeName(parameter.Schema),
			),
		))
	}

	return c.C("export type %sURI = %s;").Format(
		op.ID,
		c.BodyF(c.List(1, fieldList...).Indent(2)),
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
	for _, parameter := range op.Qry {
		required := ""
		if !common.ParRequired(parameter) {
			required = "?"
		}

		fieldList = append(fieldList, c.List(0,
			doc(parameter.Description),
			c.C("%s%s: %s;").Format(
				parameter.Name,
				required,
				baseTypeName(parameter.Schema),
			),
		))
	}

	return c.C("export type %sQry = %s;").Format(
		op.ID,
		c.BodyF(c.List(1, fieldList...).Indent(2)),
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

	return c.C("export %s").Format(
		typeDecl(
			fmt.Sprintf("%sReq", op.ID),
			op.Req.Content.First().Value().Schema,
		),
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
				c.C(`export type %s = string`).Format(
					name,
				),
			)
		} else {
			list = append(list,
				c.C("export %s").Format(typeDecl(
					name,
					rspSchemaProxy.SchemaProxy,
				)),
			)
		}
	}

	return c.List(1,
		list...,
	)
}
