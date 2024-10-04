package python

import (
	"strings"

	"github.com/pb33f/libopenapi/datamodel/high/base"
	v3 "github.com/pb33f/libopenapi/datamodel/high/v3"
	"github.com/pb33f/libopenapi/orderedmap"

	"columba-livia/common"
	c "columba-livia/content"
)

func serverPath(path string) string {
	return common.PathReg.ReplaceAllString(path, "<$1>")
}

func serverAPI(
	tag *base.Tag,
	pathItems *orderedmap.Map[string, *v3.PathItem],
) (render render) {
	return func() c.C {
		file.needMessage = true

		file.importMap["from flask import Flask as _Flask"] = struct{}{}
		file.importMap["from flask import request as _request"] = struct{}{}
		file.importMap["import typing as _typing"] = struct{}{}
		file.importMap["from . import _convert"] = struct{}{}

		list := make([]c.C, 0)
		for _, operation := range common.TagOperationList(tag.Name, pathItems) {
			list = append(list, serverAPIOperation(operation))
		}

		return c.List(1,
			// 类名
			c.List(0,
				c.F(`class {{.}}:`).Format(tag.Name),
				doc(tag.Description).IndentSpace(4),
			),
			// init 函数
			c.List(0,
				c.C("def __init__(self, app: _Flask):"),
				c.C("self.__app = app").IndentSpace(4),
			).IndentSpace(4),
			// 方法
			c.List(1, list...).IndentSpace(4),
		)
	}
}

// api 注册函数
func serverAPIOperation(
	op *common.Operation,
) c.C {
	uriExist := len(op.URI) != 0
	qryExist := len(op.Qry) != 0
	reqExist := common.ReqSchemaProxy(op.Req).ContentType != common.ContentEmpty

	reqType := common.ReqSchemaProxy(op.Req)
	rspType := common.RspSchemaProxy(op.Rsp)

	//                                                                 handler 声明

	handlerDecl := c.F("handler: _typing.Callable{{.}},").Format(c.BodyS(
		c.List(0,
			// 参数
			c.F("{{.}},").Format(c.BodyS(
				c.List(0,
					c.If(uriExist, c.F("_message.{{.tag}}.{{.op}}.uri,").Format(map[string]any{"tag": op.Tag, "op": op.ID})),
					c.If(qryExist, c.F("_message.{{.tag}}.{{.op}}.qry,").Format(map[string]any{"tag": op.Tag, "op": op.ID})),
					c.If(reqExist, c.F("_message.{{.tag}}.{{.op}}.req,").Format(map[string]any{"tag": op.Tag, "op": op.ID})),
				).IndentSpace(4),
			)),
			// 返回值
			c.F("_typing.Awaitable[_typing.Tuple{{.}}],").Format(c.BodyS(
				c.List(0,
					c.ForList(rspType, func(item common.ContentSchema) c.C {
						return c.F("_typing.Optional[_message.{{.tag}}.{{.op}}.rsp{{.code}}],").Format(map[string]any{
							"tag":  op.Tag,
							"op":   op.ID,
							"code": item.RspCode,
						})
					})...,
				).IndentSpace(4),
			)),
		).IndentSpace(4),
	))

	//                                                                  注册内容

	registerName := c.F(`
@self.__app.{{.method}}("{{.path}}")
async def {{.op}}({{.args}}):
`).
		Format(map[string]any{
			"method": strings.ToLower(op.Method),
			"path":   serverPath(op.Path),
			"op":     op.ID,
			"args":   c.If(uriExist, "**path"),
		})

	handlerURI := c.If(uriExist, c.F("uri = _message.{{.tag}}.{{.op}}.uri(**path)").
		Format(map[string]any{"tag": op.Tag, "op": op.ID}))

	handlerQry := c.If(qryExist, c.F(`
args: dict = _request.args.to_dict()
qry = _message.{{.tag}}.{{.op}}.qry(**args)
`).
		Format(map[string]any{"tag": op.Tag, "op": op.ID}),
	)

	handlerReq := c.C("")
	switch reqType.ContentType {
	case common.ContentJSON:
		if common.SchemaType(reqType.SchemaProxy.Schema()) == common.TypeArray {
			handlerReq = c.F(`
reqType = _typing.get_args(_message.{{.tag}}.{{.op}}.req)[0]
req = [reqType(**i) for i in _request.get_json()]
`).Format(map[string]any{"tag": op.Tag, "op": op.ID})
		} else {
			handlerReq = c.F("req = _message.{{.tag}}.{{.op}}.req(**_request.get_json())").Format(map[string]any{"tag": op.Tag, "op": op.ID})
		}

	case common.ContentText:
		handlerReq = c.C(`req = _request.get_data().decode("utf-8")`)
	}

	handlerHandleReturn := c.JoinSpace(
		c.ForList(rspType, func(item common.ContentSchema) c.C {
			return c.F("rsp{{.}},").Format(item.RspCode)
		})...,
	)

	handlerHandle := c.F("{{.return}} await handler({{.args}})").Format(map[string]any{
		"return": c.If(handlerHandleReturn != "", c.F("({{.}}) =").Format(handlerHandleReturn)),
		"args": c.Join(", ",
			c.If(uriExist, "uri"),
			c.If(qryExist, "qry"),
			c.If(reqExist, "req"),
		),
	})

	handlerRsp := make([]c.C, 0)
	for _, rsp := range rspType {
		handlerRsp = append(handlerRsp,
			c.F(`
if rsp{{.}} is not None:
    return _convert.rsp(rsp{{.}}), {{.}}
`).Format(rsp.RspCode),
		)
	}

	handlerRsp = append(handlerRsp, `return "", 500`)

	return c.List(1,
		c.List(0,
			c.F("def {{.op}}{{.args}}:").Format(map[string]any{
				"op": op.ID,
				"args": c.BodyC(c.List(0,
					"self,",
					handlerDecl,
				).IndentSpace(4)),
			}),
			doc(op.Description).IndentSpace(4),
		),
		c.List(0,
			registerName,
			c.List(1,
				handlerURI,
				handlerQry,
				handlerReq,
				handlerHandle,
				c.List(1, handlerRsp...),
			).IndentSpace(4),
		).IndentSpace(4),
	)
}
