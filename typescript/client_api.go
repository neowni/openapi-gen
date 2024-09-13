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

		return c.C(`
export default class {
  private instance: AxiosInstance;
  constructor(instance: AxiosInstance) {
    this.instance = instance;
  }

%s
}
`).TrimSpace().Format(c.List(1, list...).IndentSpace(2))
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
	funcNameArg := c.C("async %s%s:").Format(
		op.ID,
		c.BodyC(c.List(0,
			c.If(uriExist, c.C("uri: message.%s.%sURI,").Format(op.Tag, op.ID)),
			c.If(qryExist, c.C("qry: message.%s.%sQry,").Format(op.Tag, op.ID)),
			c.If(reqExist, c.C("req: message.%s.%sReq,").Format(op.Tag, op.ID)),
		).IndentSpace(2)),
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
		)...).IndentSpace(2),
	))

	// 函数主体发起请求

	// 对于没有参数的路径，使用 " 而非 ` 包含其中字符
	uri := c.C("url: `%s`,").Format(clientPath(op.Path))
	if !strings.Contains(uri.String(), "$") {
		uri = c.C("url: \"%s\",").Format(clientPath(op.Path))
	}
	funcBodyHandle := c.C(`const rsp = await this.instance.request(%s);`).TrimSpace().Format(
		c.BodyF(c.List(0,
			uri,
			c.C("method: \"%s\",").Format(op.Method),
			c.If(qryExist, c.C("params: qry,")),
			c.If(reqExist, c.C("data: req,")),
		).IndentSpace(2)),
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
	).IndentSpace(2))

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
