package python

import (
	"fmt"

	"github.com/pb33f/libopenapi/datamodel/high/base"

	c "columba-livia/content"
)

func serverInit(
	tags []*base.Tag,
) (render render) {
	return func() c.C {

		file.importMap["from flask import Flask as _Flask"] = struct{}{}

		fieldList := make([]c.C, 0)

		for _, tag := range tags {
			// 导入
			file.importMap[fmt.Sprintf(
				"from .%s import %s as _%s",
				tag.Name, tag.Name, tag.Name,
			)] = struct{}{}

			// 字段内容
			fieldList = append(fieldList,
				c.C("self.%s = _%s(app)").Format(tag.Name, tag.Name),
			)
		}

		return c.C(`
class Server:
    def __init__(self, app: _Flask):
%s
		`).TrimSpace().Format(c.List(0, fieldList...).IndentSpace(8))
	}
}
