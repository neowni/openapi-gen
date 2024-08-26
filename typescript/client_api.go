package typescript

import (
	"github.com/pb33f/libopenapi/datamodel/high/base"
	v3 "github.com/pb33f/libopenapi/datamodel/high/v3"
	"github.com/pb33f/libopenapi/orderedmap"

	"columba-livia/common"
	c "columba-livia/content"
)

func clientPath(path string) string {
	return common.PathReg.ReplaceAllString(path, "${uri.$1}")
}

func clientApi(
	tag *base.Tag,
	pathItems *orderedmap.Map[string, *v3.PathItem],
) (render render) {
	return func() c.C {
		file.needMessage = true

		file.importMap[`import { AxiosInstance } from "axios";`] = struct{}{}

		list := make([]c.C, 0)
		for _, op := range common.TagOperationList(tag.Name, pathItems) {
			list = append(list, clientApiOperation(op))
		}

		return c.C(`
export default class {
  private instance: AxiosInstance;
  constructor(instance: AxiosInstance) {
    this.instance = instance;
  }

%s
}
`).TrimSpace().Format(c.List(1, list...).Indent(2))
	}
}

func clientApiOperation(
	op *common.Operation,
) c.C {
	uriExist := len(op.URI) != 0
	qryExist := len(op.Qry) != 0
	reqExist := common.ReqSchemaProxy(op.Req).ContentType != common.ContentEmpty

	rspType := common.RspSchemaProxy(op.Rsp)

	// 函数名称/参数
	funcNameArg := c.C("async %s%s:").Format(
		op.ID,
		c.BodyC(c.List(0,
			c.If(uriExist, c.C("uri: message.%s.%sURI,").Format(op.Tag, op.ID)),
			c.If(qryExist, c.C("qry: message.%s.%sQry,").Format(op.Tag, op.ID)),
			c.If(reqExist, c.C("req: message.%s.%sReq,").Format(op.Tag, op.ID)),
		).Indent(2)),
	)

	// 函数返回值
	funcReturn := c.C("Promise<%s>").Format(c.BodyF(
		c.List(0, c.ForList(
			rspType,
			func(item common.ContentSchema) c.C {
				return c.C("_%s?: message.%s.%sRsp%s;").Format(
					item.RspCode,
					op.Tag, op.ID,
					item.RspCode,
				)
			},
		)...).Indent(2),
	))

	// 函数主体发起请求
	funcBodyHandle := c.C(`const rsp = await this.instance.request(%s)`).TrimSpace().Format(
		c.BodyF(c.List(0,
			c.C("url: `%s`,").Format(clientPath(op.Path)),
			c.C("method: \"%s\",").Format(op.Method),
			c.If(qryExist, c.C("params: qry,")),
			c.If(reqExist, c.C("data: req,")),
		).Indent(2)),
	)

	// 函数主体处理响应
	funcBodyReturn := c.C(`
switch (rsp.status) {
%s
  default:
    return {};
}
	`).TrimSpace().Format(c.List(0,
		c.ForList(rspType, func(item common.ContentSchema) c.C {
			return c.C(`
case %s:
  return { _%s: rsp.data };
`).TrimSpace().Format(item.RspCode, item.RspCode)
		})...,
	).Indent(2))

	// 函数主体
	funcBody := c.BodyF(c.List(1,
		funcBodyHandle,
		funcBodyReturn,
	).Indent(2))

	return c.List(0,
		doc(op.Description),
		c.JoinSpace(
			funcNameArg, funcReturn, funcBody,
		),
	)
}
