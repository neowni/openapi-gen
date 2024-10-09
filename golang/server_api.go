package golang

import (
	"github.com/pb33f/libopenapi/datamodel/high/base"
	v3 "github.com/pb33f/libopenapi/datamodel/high/v3"
	"github.com/pb33f/libopenapi/orderedmap"

	"columba-livia/common"
	c "columba-livia/content"
)

func serverPath(path string) string {
	return common.PathReg.ReplaceAllString(path, ":$1")
}

func serverAPI(
	tag *base.Tag,
	pathItems *orderedmap.Map[string, *v3.PathItem],
) (render render) {
	return func() c.C {
		file.needMessage = true

		file.importMap["context"] = ""
		file.importMap["github.com/gin-gonic/gin"] = ""

		list := make([]c.C, 0)
		list = append(list, serverAPIStruct(tag.Name))

		for _, op := range common.TagOperationList(tag.Name, pathItems) {
			list = append(list, serverAPIOperation(op))
		}

		return c.List(1, list...)
	}
}

// api 结构体
func serverAPIStruct(
	name string,
) c.C {
	return c.F(`
type {{.}} struct {
	engine gin.IRoutes
}
`).
		Format(name)
}

// api 注册函数
func serverAPIOperation(
	op *common.Operation,
) c.C {
	uriExist := len(op.URI) != 0
	qryExist := len(op.Qry) != 0
	reqExist := common.ReqSchemaProxy(op.Req).ContentType != common.ContentEmpty

	reqType := common.ReqSchemaProxy(op.Req).ContentType
	rspType := common.RspSchemaProxy(op.Rsp)

	// 																handler 声明

	handlerDecl := c.F(`handler func{{.args}} {{.return}},`).
		Format(map[string]any{
			"args": c.BodyC(c.List(0,
				c.C("ctx context.Context,"),
				c.If(uriExist, c.F("uri *message.{{.}},").Format(publicName(op.Tag+publicName(op.ID)+"Uri"))),
				c.If(qryExist, c.F("qry *message.{{.}},").Format(publicName(op.Tag+publicName(op.ID)+"Qry"))),
				c.If(reqExist, c.F("req *message.{{.}},").Format(publicName(op.Tag+publicName(op.ID)+"Req"))),
			).IndentTab(1)),

			"return": c.BodyC(c.List(0, c.Flat(
				c.ForList(op.Rsp, func(item orderedmap.Pair[string, *v3.Response]) c.C {
					return c.F("rsp{{.code}} *message.{{.name}},").
						Format(map[string]any{
							"code": item.Key(),
							"name": publicName(op.Tag + publicName(op.ID) + "Rsp" + item.Key()),
						})
				}),
				[]c.C{c.C("err error,")},
			)...).IndentTab(1)),
		})

	// 																注册函数名称

	registerName := c.F("func (tag *{{.tag}}) {{.op}}{{.decl}}").Format(map[string]any{
		"tag": op.Tag,
		"op":  publicName(op.ID),
		// 参数：提供 handler 函数
		"decl": c.BodyC(handlerDecl.IndentTab(1)),
	})

	// 																  handle函数

	handlerURI := c.If(uriExist,
		c.F(`
// 解析 uri
uri := new(message.{{.}})
if convert.BindURI(ctx, uri) {
	return
}
`).Format(publicName(op.Tag+publicName(op.ID)+"Uri")))

	handlerQry := c.If(qryExist,
		c.F(`
// 解析 qry
qry := new(message.{{.}})
if convert.BindQry(ctx, qry) {
	return
}
`).Format(publicName(op.Tag+publicName(op.ID)+"Qry")))

	handlerReq := c.If(reqExist,
		c.F(`
// 解析 req
req := new(message.{{.name}})
if convert.BindReq{{.type}}(ctx, req) {
	return
}
`).Format(map[string]any{
			"name": publicName(op.Tag + publicName(op.ID) + "Req"),
			"type": publicName(string(reqType)),
		}),
	)

	// 处理请求与异常
	handlerHandle := c.F(`
// 处理请求
{{.return}} := handler({{.args}})

// 返回异常
if convert.ResponseError(ctx, err) {
	return
}
`).Format(map[string]any{
		"return": c.Join(", ",
			c.Flat(
				c.ForList(op.Rsp, func(item orderedmap.Pair[string, *v3.Response]) c.C {
					return c.F("rsp{{.}}").Format(item.Key())
				}),
				[]c.C{"err"},
			)...,
		),
		"args": c.Join(", ",
			"ctx",
			c.If(uriExist, "uri"),
			c.If(qryExist, "qry"),
			c.If(reqExist, "req"),
		),
	})

	// 响应请求
	handlerReturnCases := c.List(0,
		c.ForList(rspType, func(item common.ContentSchema) c.C {
			return c.F(`
case rsp{{.code}} != nil:
	convert.Response{{.type}}(ctx, {{.code}}, rsp{{.code}})
`,
			).Format(map[string]any{
				"code": item.RspCode,
				"type": publicName(string(item.ContentType)),
			})
		})...,
	)

	handlerReturn := c.F(`
// 返回响应
switch true {
{{.}}
}
`).Format(handlerReturnCases)

	handlerFunc := c.F("func(ctx *gin.Context) {{.}},").Format(
		c.BodyF(c.List(1,
			handlerURI,
			handlerQry,
			handlerReq,
			handlerHandle,
			handlerReturn,
		).IndentTab(1)),
	)

	// 																注册函数主体

	registerBody := c.BodyF(
		c.F("tag.engine.Handle{{.}}").Format(c.BodyC(
			c.List(0,
				c.F(`"{{.}}",`).Format(op.Method),
				c.F(`"{{.}}",`).Format(serverPath(op.Path)),
				handlerFunc,
			).IndentTab(1),
		)).IndentTab(1),
	)

	return c.JoinSpace(registerName, registerBody)
}
