package python

import (
	c "columba-livia/content"
)

func serverConvert() (render render) {
	return func() c.C {
		file.importMap["import typing as _typing"] = struct{}{}
		file.importMap["import pydantic as _pydantic"] = struct{}{}
		file.importMap["from collections.abc import Sequence as _Sequence"] = struct{}{}

		return c.C(serverConvertContent).TrimSpace()
	}
}

const serverConvertContent = `
bodyType = _typing.Union[
    _pydantic.StrictStr,
    _pydantic.BaseModel,
    _typing.List["bodyType"],
    _typing.Sequence["bodyType"],
]


def rsp(body: bodyType):
    """转换 rsp body 格式"""
    if isinstance(body, str):
        return body

    if isinstance(body, _pydantic.BaseModel):
        return body.model_dump_json(exclude_none=True)

    if isinstance(body, list):
        return [rsp(i) for i in body]

    if isinstance(body, _Sequence):
        return [rsp(i) for i in body]

    return None
`
