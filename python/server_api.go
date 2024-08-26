package python

import (
	"columba-livia/common"
	c "columba-livia/content"
	"strings"

	"github.com/pb33f/libopenapi/datamodel/high/base"
	v3 "github.com/pb33f/libopenapi/datamodel/high/v3"
	"github.com/pb33f/libopenapi/orderedmap"
)

func serverPath(path string) string {
	return common.PathReg.ReplaceAllString(path, "<$1>")
}

func serverApi(
	tag *base.Tag,
	pathItems *orderedmap.Map[string, *v3.PathItem],
) (render render) {
	return func() c.C {
		file.needMessage = true

		file.importMap["from flask import Flask as _Flask"] = struct{}{}
		file.importMap["from flask import request as _request"] = struct{}{}
		file.importMap["import typing as _typing"] = struct{}{}

		list := make([]c.C, 0)
		for _, operation := range common.TagOperationList(tag.Name, pathItems) {
			list = append(list, serverApiOperation(operation))
		}

		return c.List(1,
			// 类名
			c.List(0,
				c.C(`class %s:`).Format(tag.Name),
				c.C(`"""%s"""`).Format(tag.Description).Indent(4),
			),
			// init 函数
			c.List(0,
				c.C("def __init__(self, app: _Flask):"),
				c.C("self.__app = app").Indent(4),
			).Indent(4),
			// 方法
			c.List(1, list...).Indent(4),
		)
	}
}

// api 注册函数
func serverApiOperation(
	op *common.Operation,
) c.C {
	uriExist := len(op.URI) != 0
	qryExist := len(op.Qry) != 0
	reqExist := common.ReqSchemaProxy(op.Req).ContentType != common.ContentEmpty

	reqType := common.ReqSchemaProxy(op.Req)
	rspType := common.RspSchemaProxy(op.Rsp)

	//                                                                 handler 声明

	handlerDecl := c.C("handler: _typing.Callable%s,").Format(c.BodyS(
		c.List(0,
			// 参数
			c.C("%s,").Format(c.BodyS(
				c.List(0,
					c.If(uriExist, c.C("_message.%s.%s.uri,").Format(op.Tag, op.ID)),
					c.If(qryExist, c.C("_message.%s.%s.qry,").Format(op.Tag, op.ID)),
					c.If(reqExist, c.C("_message.%s.%s.req,").Format(op.Tag, op.ID)),
				).Indent(4),
			)),
			// 返回值
			c.C("_typing.Awaitable[_typing.Tuple%s],").Format(c.BodyS(
				c.List(0,
					c.ForList(rspType, func(item common.ContentSchema) c.C {
						return c.C("_typing.Optional[_message.%s.%s.rsp%s],").Format(op.Tag, op.ID, item.RspCode)
					})...,
				).Indent(4),
			)),
		).Indent(4),
	))

	//                                                                  注册内容

	registerName := c.C(`
@self.__app.%s("%s")
async def %s(%s):
`).TrimSpace().Format(
		strings.ToLower(op.Method), serverPath(op.Path),
		op.ID, c.If(uriExist, "**path"),
	)

	handlerURI := c.If(uriExist, c.C("uri = _message.%s.%s.uri(**path)").Format(op.Tag, op.ID))

	handlerQry := c.If(qryExist, c.C(`
args: dict = _request.args.to_dict()
qry = _message.%s.%s.qry(**args)
`).
		TrimSpace().Format(op.Tag, op.ID),
	)

	handlerReq := c.C("")
	switch reqType.ContentType {
	case common.ContentJson:
		if common.SchemaType(reqType.SchemaProxy.Schema()) == common.TypeArray {
			handlerReq = c.C(`
reqType = _typing.get_args(_message.%s.%s.req)[0]
req = [reqType(**i) for i in _request.get_json()]
`).TrimSpace().Format(op.Tag, op.ID)
		} else {
			handlerReq = c.C("req = _message.%s.%s.req(**_request.get_json())").Format(op.Tag, op.ID)
		}

	case common.ContentText:
		handlerReq = c.C(`req = _request.get_data().decode("utf-8")`)
	}

	handlerHandleReturn := c.JoinSpace(
		c.ForList(rspType, func(item common.ContentSchema) c.C {
			return c.C("rsp%s,").Format(item.RspCode)
		})...,
	)

	handlerHandle := c.C("%s await handler(%s)").Format(
		c.If(handlerHandleReturn != "", c.C("(%s) =").Format(handlerHandleReturn)),
		c.Join(", ",
			c.If(uriExist, "uri"),
			c.If(qryExist, "qry"),
			c.If(reqExist, "req"),
		),
	)

	handlerRsp := make([]c.C, 0)
	for _, rsp := range rspType {
		r := c.C(`""`)
		if rsp.ContentType == common.ContentJson {
			if common.SchemaType(rsp.SchemaProxy.Schema()) == common.TypeArray {
				r = c.C("[i.model_dump_json(exclude_none=True) for i in rsp%s]").Format(rsp.RspCode)
			} else {
				r = c.C("rsp%s.model_dump_json(exclude_none=True)").Format(rsp.RspCode)
			}
		}
		if rsp.ContentType == common.ContentText {
			r = c.C("rsp%s").Format(rsp.RspCode)
		}

		handlerRsp = append(handlerRsp,
			c.C(`
if rsp%s is not None:
    return %s, %s
`).
				TrimSpace().Format(
				rsp.RspCode,
				r, rsp.RspCode,
			),
		)
	}

	handlerRsp = append(handlerRsp, `return "", 500`)

	return c.List(1,
		c.List(0,
			c.C("def %s%s:").Format(
				op.ID,
				c.BodyC(c.List(0,
					"self,",
					handlerDecl,
				).Indent(4)),
			),
			c.C(`"""%s"""`).Format(op.Description).Indent(4),
		),
		c.List(0,
			registerName,
			c.List(1,
				handlerURI,
				handlerQry,
				handlerReq,
				handlerHandle,
				c.List(1, handlerRsp...),
			).Indent(4),
		).Indent(4),
	)
}
