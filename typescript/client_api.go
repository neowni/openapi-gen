package typescript

import (
	"strings"

	"github.com/pb33f/libopenapi/datamodel/high/base"
	v3 "github.com/pb33f/libopenapi/datamodel/high/v3"
	"github.com/pb33f/libopenapi/orderedmap"

	"columba-livia/common"
	c "columba-livia/content"
)

func clientPath(path string) string {
	return common.PathReg.ReplaceAllString(path, "${uri.$1}")
}

func clientAPI(
	tag *base.Tag,
	pathItems *orderedmap.Map[string, *v3.PathItem],
) (render render) {
	return func() c.C {
		file.needMessage = true

		file.importMap[`import { AxiosInstance } from "axios";`] = struct{}{}

		list := make([]c.C, 0)
		for _, op := range common.TagOperationList(tag.Name, pathItems) {
			list = append(list, clientAPIOperation(op))
		}

		return c.F(`
export default class {
  private instance: AxiosInstance;
  constructor(instance: AxiosInstance) {
    this.instance = instance;
  }

{{.}}
}
`).Format(c.List(1, list...).IndentSpace(2))
	}
}

func clientAPIOperation(
	op *common.Operation,
) c.C {
	uriExist := len(op.URI) != 0
	qryExist := len(op.Qry) != 0
	reqExist := common.ReqSchemaProxy(op.Req).ContentType != common.ContentEmpty

	rspType := common.RspSchemaProxy(op.Rsp)

	// 函数名称/参数
	funcNameArg := c.F("async {{.op}}{{.body}}:").Format(map[string]any{
		"op": op.ID,
		"body": c.BodyC(c.List(0,
			c.If(uriExist, c.F("uri: message.{{.tag}}.{{.op}}URI,").Format(map[string]any{"tag": op.Tag, "op": op.ID})),
			c.If(qryExist, c.F("qry: message.{{.tag}}.{{.op}}Qry,").Format(map[string]any{"tag": op.Tag, "op": op.ID})),
			c.If(reqExist, c.F("req: message.{{.tag}}.{{.op}}Req,").Format(map[string]any{"tag": op.Tag, "op": op.ID})),
		).IndentSpace(2)),
	})

	// 函数返回值
	funcReturn := c.F("Promise<{{.}}>").Format(c.BodyF(
		c.List(0, c.ForList(
			rspType,
			func(item common.ContentSchema) c.C {
				return c.F("_{{.code}}?: message.{{.tag}}.{{.op}}Rsp{{.code}};").Format(map[string]any{
					"tag":  op.Tag,
					"op":   op.ID,
					"code": item.RspCode,
				})
			},
		)...).IndentSpace(2),
	))

	// 函数主体发起请求

	// 对于没有参数的路径，使用 " 而非 ` 包含其中字符
	uri := c.F("url: `{{.}}`,").Format(clientPath(op.Path))
	if !strings.Contains(uri.String(), "$") {
		uri = c.F("url: \"{{.}}\",").Format(clientPath(op.Path))
	}
	funcBodyHandle := c.F(`const rsp = await this.instance.request({{.}});`).Format(
		c.BodyF(c.List(0,
			uri,
			c.F("method: \"{{.}}\",").Format(op.Method),
			c.If(qryExist, c.C("params: qry,")),
			c.If(reqExist, c.C("data: req,")),
		).IndentSpace(2)),
	)

	// 函数主体处理响应
	funcBodyReturnCases := c.List(0,
		c.ForList(rspType, func(item common.ContentSchema) c.C {
			return c.F(`
case {{.}}:
  return { _{{.}}: rsp.data };
`).Format(item.RspCode)
		})...,
	)

	funcBodyReturn := c.F(`
switch (rsp.status) {
{{.}}
  default:
    return {};
}
	`).Format(funcBodyReturnCases.IndentSpace(2))

	// 函数主体
	funcBody := c.BodyF(c.List(1,
		funcBodyHandle,
		funcBodyReturn,
	).IndentSpace(2))

	return c.List(0,
		doc(op.Description),
		c.JoinSpace(
			funcNameArg, funcReturn, funcBody,
		),
	)
}
