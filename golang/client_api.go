package golang

import (
	"columba-livia/common"
	c "columba-livia/content"
	"strings"

	"github.com/pb33f/libopenapi/datamodel/high/base"
	v3 "github.com/pb33f/libopenapi/datamodel/high/v3"
	"github.com/pb33f/libopenapi/orderedmap"
)

func clientPath(path string) string {
	return common.PathReg.ReplaceAllString(path, "{$1}")
}

func clientApi(
	tag *base.Tag,
	pathItems *orderedmap.Map[string, *v3.PathItem],
) (render render) {
	return func() c.C {
		file.needMessage = true

		file.importMap["github.com/go-resty/resty/v2"] = ""

		list := make([]c.C, 0)
		list = append(list, clientApiStruct(tag.Name))

		for _, op := range common.TagOperationList(tag.Name, pathItems) {
			list = append(list, clientApiOperation(op))
		}

		return c.List(1, list...)
	}
}

func clientApiStruct(
	name string,
) c.C {
	return c.C(`
type %s struct {
	client *resty.Client
}
`).
		TrimSpace().
		Format(name)
}

func clientApiOperation(
	op *common.Operation,
) c.C {
	uriExist := len(op.URI) != 0
	qryExist := len(op.Qry) != 0
	reqExist := common.ReqSchemaProxy(op.Req).ContentType != common.ContentEmpty

	uriType := make(map[string]string)
	for _, u := range op.URI {
		uriType[u.Name] = typeMap(common.SchemaType(u.Schema.Schema()))
	}

	qryType := make(map[string]string)
	for _, q := range op.Qry {
		type_ := typeMap(common.SchemaType(q.Schema.Schema()))
		if common.ParRequired(q) {
			type_ += "R"
		}
		qryType[q.Name] = type_
	}

	reqType := common.ReqSchemaProxy(op.Req).ContentType
	rspType := common.RspSchemaProxy(op.Rsp)

	// 																	函数名称

	funcName := c.C("func (tag *%s) %s%s %s").Format(
		op.Tag, ExportName(op.ID),
		// 函数参数
		c.BodyC(c.List(0,
			c.If(uriExist, c.C("uri *message.%s,").Format(ExportName(op.ID+"URI"))),
			c.If(qryExist, c.C("qry *message.%s,").Format(ExportName(op.ID+"Qry"))),
			c.If(reqExist, c.C("req *message.%s,").Format(ExportName(op.ID+"Req"))),
		).Indent(4)),
		// 函数返回值
		c.BodyC(c.List(0,
			append(
				c.ForList(op.Rsp, func(item orderedmap.Pair[string, *v3.Response]) c.C {
					return c.C("rsp%s *message.%s,").Format(
						item.Key(), ExportName(op.ID+"Rsp"+item.Key()),
					)
				}),
				c.C("err error,"),
			)...,
		).Indent(4)),
	)

	//																	函数主体

	funcBodyNew := c.C(`
// 新建请求
r := tag.client.R()
`).TrimSpace()

	funcBodyURI := c.If(uriExist, c.List(0,
		c.ForMap(
			uriType, func(key string, value string) c.C {
				return c.C(`convert.uri%s(r, "%s", uri.%s)`).Format(
					ExportName(value), key, ExportName(key),
				)
			},
		)...,
	))

	funcBodyQry := c.If(qryExist, c.List(0,
		c.ForMap(
			qryType, func(key string, value string) c.C {
				return c.C(`convert.qry%s(r, "%s", qry.%s)`).Format(
					ExportName(value), key, ExportName(key),
				)
			},
		)...,
	))

	funcBodyReq := c.If(reqExist,
		c.C(`r.SetBody(%sreq)`).Format(
			// 对于 text 类型，需要将 *string 转换为 string
			c.If(reqType == common.ContentText, "*"),
		),
	)

	// 返回语句内容
	funcReturn := c.Join(", ", append(
		c.ForList(op.Rsp, func(item orderedmap.Pair[string, *v3.Response]) c.C {
			return c.C("rsp%s").Format(item.Key())
		}),
		c.C("err"),
	)...)

	funcBodyHandle := c.C(`
resp, err := r.%s("%s")
if err != nil {
	return %s
}
`).TrimSpace().Format(
		ExportName(strings.ToLower(op.Method)),
		clientPath(op.Path),
		funcReturn,
	)

	funcBodyRsp := c.C(`
switch resp.StatusCode() {
%s
default:
	err = convert.ResponseError(resp)
	return %s
}
`).TrimSpace().Format(
		c.List(0,
			c.ForList(rspType, func(item common.ContentSchema) c.C {
				return c.C(`
case %s:
	rsp%s = new(message.%s)
	err = convert.Response%s(resp, rsp%s)
	return %s
`).TrimSpace().Format(
					item.RspCode,
					item.RspCode, ExportName(op.ID+"Rsp"+item.RspCode),
					ExportName(string(item.ContentType)), item.RspCode,
					funcReturn,
				)
			})...,
		),
		// default返回
		funcReturn,
	)

	return c.List(0,
		doc(ExportName(op.ID), op.Description),
		c.JoinSpace(funcName, c.BodyF(
			c.List(1,
				funcBodyNew,
				funcBodyURI,
				funcBodyQry,
				funcBodyReq,
				funcBodyHandle,
				funcBodyRsp,
			).Indent(4),
		)),
	)
}
