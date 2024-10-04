package golang

import (
	"strings"

	"github.com/pb33f/libopenapi/datamodel/high/base"
	v3 "github.com/pb33f/libopenapi/datamodel/high/v3"
	"github.com/pb33f/libopenapi/orderedmap"

	"columba-livia/common"
	c "columba-livia/content"
)

func clientPath(path string) string {
	return common.PathReg.ReplaceAllString(path, "{$1}")
}

func clientAPI(
	tag *base.Tag,
	pathItems *orderedmap.Map[string, *v3.PathItem],
) (render render) {
	return func() c.C {
		file.needMessage = true

		file.importMap["github.com/go-resty/resty/v2"] = ""

		list := make([]c.C, 0)
		list = append(list, clientAPIStruct(tag.Name))

		for _, op := range common.TagOperationList(tag.Name, pathItems) {
			list = append(list, clientAPIOperation(op))
		}

		return c.List(1, list...)
	}
}

func clientAPIStruct(
	name string,
) c.C {
	return c.F(`
type {{.}} struct {
	client *resty.Client
}
`).
		Format(privateName(name))
}

func clientAPIOperation(
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

	funcName := c.F("func (tag *{{.tag}}) {{.name}}{{.args}} {{.return}}").Format(map[string]any{
		"tag":  op.Tag,
		"name": publicName(op.ID),

		"args": c.BodyC(c.List(0,
			c.If(uriExist, c.F("uri *message.{{.}},").Format(publicName(op.ID+"Uri"))),
			c.If(qryExist, c.F("qry *message.{{.}},").Format(publicName(op.ID+"Qry"))),
			c.If(reqExist, c.F("req *message.{{.}},").Format(publicName(op.ID+"Req"))),
		).IndentTab(1)),

		"return": c.BodyC(c.List(0, c.Flat(
			c.ForList(op.Rsp, func(item orderedmap.Pair[string, *v3.Response]) c.C {
				return c.F("rsp{{.code}} *message.{{.name}},").
					Format(map[string]any{
						"code": item.Key(),
						"name": publicName(op.ID + "Rsp" + item.Key()),
					})
			}),
			[]c.C{c.C("err error,")},
		)...).IndentTab(1)),
	})

	//																	函数主体

	funcBodyNew := c.C(`
// 新建请求
r := tag.client.R()
`).TrimSpace()

	funcBodyURI := c.If(uriExist, c.List(0,
		c.ForMap(
			uriType, func(name string, type_ string) c.C {
				return c.F(`convert.uri{{.type}}(r, "{{.name}}", uri.{{.field}})`).Format(map[string]any{
					"type":  publicName(type_),
					"name":  name,
					"field": publicName(name),
				})
			},
		)...,
	))

	funcBodyQry := c.If(qryExist, c.List(0, c.ForMap(
		qryType, func(name string, type_ string) c.C {
			return c.F(`convert.qry{{.type}}(r, "{{.name}}", qry.{{.field}})`).Format(map[string]any{
				"type":  publicName(type_),
				"name":  name,
				"field": publicName(name),
			})
		},
	)...))

	funcBodyReq := c.If(reqExist,
		c.F(`r.SetBody({{.}}req)`).Format(
			// 对于 text 类型，需要将 *string 转换为 string
			c.If(reqType == common.ContentText, "*"),
		),
	)

	// 返回语句内容
	funcReturn := c.Join(", ",
		c.Flat(
			c.ForList(op.Rsp, func(item orderedmap.Pair[string, *v3.Response]) c.C {
				return c.F("rsp{{.}}").Format(item.Key())
			}),
			[]c.C{c.C("err")},
		)...,
	)

	funcBodyHandle := c.F(`
resp, err := r.{{.method}}("{{.path}}")
if err != nil {
	return {{.return}}
}
`).Format(map[string]any{
		"method": publicName(strings.ToLower(op.Method)),
		"path":   clientPath(op.Path),
		"return": funcReturn,
	})

	funcBodyRspCases := c.List(0, c.ForList(
		rspType, func(item common.ContentSchema) c.C {
			return c.F(`
case {{.code}}:
	rsp{{.code}} = new(message.{{.name}})
	err = convert.Response{{.type}}(resp, rsp{{.code}})
	return {{.return}}
`).
				Format(map[string]any{
					"code":   item.RspCode,
					"name":   publicName(op.ID + "Rsp" + item.RspCode),
					"type":   publicName(string(item.ContentType)),
					"return": funcReturn,
				})
		},
	)...)

	funcBodyRsp := c.F(`
switch resp.StatusCode() {
{{.cases}}
default:
	err = convert.ResponseError(resp)
	return {{.return}}
}
`).
		Format(map[string]any{
			"cases":  funcBodyRspCases,
			"return": funcReturn,
		})

	return c.List(0,
		doc(publicName(op.ID), op.Description),
		c.JoinSpace(funcName, c.BodyF(
			c.List(1,
				funcBodyNew,
				funcBodyURI,
				funcBodyQry,
				funcBodyReq,
				funcBodyHandle,
				funcBodyRsp,
			).IndentTab(1),
		)),
	)
}
