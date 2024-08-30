package golang

import (
	"columba-livia/common"
	c "columba-livia/content"

	"github.com/pb33f/libopenapi/datamodel/high/base"
	v3 "github.com/pb33f/libopenapi/datamodel/high/v3"
	"github.com/pb33f/libopenapi/orderedmap"
)

func serverPath(path string) string {
	return common.PathReg.ReplaceAllString(path, ":$1")
}

func serverApi(
	tag *base.Tag,
	pathItems *orderedmap.Map[string, *v3.PathItem],
) (render render) {
	return func() c.C {
		file.needMessage = true

		file.importMap["context"] = ""
		file.importMap["github.com/gin-gonic/gin"] = ""

		list := make([]c.C, 0)
		list = append(list, serverApiStruct(tag.Name))

		for _, op := range common.TagOperationList(tag.Name, pathItems) {
			list = append(list, serverApiOperation(op))
		}

		return c.List(1, list...)
	}
}

// api 结构体
func serverApiStruct(
	name string,
) c.C {
	return c.C(`
type %s struct {
	engine gin.IRoutes
}
`).
		TrimSpace().
		Format(name)
}

// api 注册函数
func serverApiOperation(
	op *common.Operation,
) c.C {
	uriExist := len(op.URI) != 0
	qryExist := len(op.Qry) != 0
	reqExist := common.ReqSchemaProxy(op.Req).ContentType != common.ContentEmpty

	reqType := common.ReqSchemaProxy(op.Req).ContentType
	rspType := common.RspSchemaProxy(op.Rsp)

	// 																handler 声明

	handlerDecl := c.C(`
handler func%s %s,
`).
		TrimSpace().
		Format(
			// 参数
			c.BodyC(c.List(0,
				c.C("ctx context.Context,"),
				c.If(uriExist, c.C("uri *message.%s,").Format(ExportName(op.ID+"URI"))),
				c.If(qryExist, c.C("qry *message.%s,").Format(ExportName(op.ID+"Qry"))),
				c.If(reqExist, c.C("req *message.%s,").Format(ExportName(op.ID+"Req"))),
			).Indent(4)),
			// 返回值
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

	// 																注册函数名称

	registerName := c.C("func (tag *%s) %s%s").Format(
		op.Tag, ExportName(op.ID),
		c.BodyC(handlerDecl.Indent(4)),
	)

	// 																  handle函数

	handlerURI := c.If(uriExist,
		c.C(`
// 解析 uri
uri := new(message.%s)
if convert.BindURI(ctx, uri) {
	return
}
`).TrimSpace().Format(ExportName(op.ID+"URI")),
	)

	handlerQry := c.If(qryExist,
		c.C(`
// 解析 qry
qry := new(message.%s)
if convert.BindQry(ctx, qry) {
	return
}
`).TrimSpace().Format(ExportName(op.ID+"Qry")),
	)

	handlerReq := c.If(reqExist,
		c.C(`
// 解析 req
req := new(message.%s)
if convert.BindReq%s(ctx, req) {
	return
}
`).TrimSpace().Format(
			ExportName(op.ID+"Req"),
			ExportName(string(reqType)),
		),
	)

	// 处理请求与异常
	handlerHandle := c.C(`
// 处理请求
%s := handler(%s)

// 返回异常
if convert.ResponseError(ctx, err) {
	return
}
`).TrimSpace().Format(
		// 返回值
		c.Join(", ", append(
			c.ForList(op.Rsp, func(item orderedmap.Pair[string, *v3.Response]) c.C {
				return c.C("rsp%s").Format(item.Key())
			}),
			"err",
		)...),
		// 参数
		c.Join(", ",
			"ctx",
			c.If(uriExist, "uri"),
			c.If(qryExist, "qry"),
			c.If(reqExist, "req"),
		),
	)

	// 响应请求
	handlerReturn := c.C(`
// 返回响应
switch true {
%s
}
`).
		TrimSpace().Format(
		c.List(0,
			c.ForList(rspType, func(item common.ContentSchema) c.C {
				return c.C(`
case rsp%s != nil:
	convert.Response%s(ctx, %s, rsp%s)
`,
				).Format(item.RspCode, ExportName(string(item.ContentType)), item.RspCode, item.RspCode)
			})...,
		),
	)

	handlerFunc := c.C("func(ctx *gin.Context) %s,").Format(
		c.BodyF(
			c.List(1,
				handlerURI,
				handlerQry,
				handlerReq,
				handlerHandle,
				handlerReturn,
			).Indent(4),
		),
	)

	// 																注册函数主体

	registerBody := c.BodyF(
		c.C("tag.engine.Handle%s").Format(c.BodyC(
			c.List(0,
				c.C(`"%s",`).Format(op.Method),
				c.C(`"%s",`).Format(serverPath(op.Path)),
				handlerFunc,
			).Indent(4),
		)).Indent(4),
	)

	return c.JoinSpace(registerName, registerBody)
}
