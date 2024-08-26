package golang

import (
	c "columba-livia/content"

	"github.com/pb33f/libopenapi/datamodel/high/base"
)

func serverRoute(
	tags []*base.Tag,
) (render render) {
	return func() c.C {
		file.importMap["github.com/gin-gonic/gin"] = ""

		routeStructBody := c.BodyF(c.List(0,
			c.ForList(tags, func(item *base.Tag) c.C {
				return c.C("%s *%s").Format(
					ExportName(item.Name),
					item.Name,
				)
			})...,
		).Indent(4))

		routeInit := c.List(0,
			c.ForList(tags, func(item *base.Tag) c.C {
				return c.C("s.%s = &%s{engine}").Format(
					ExportName(item.Name),
					item.Name,
				)
			})...,
		).Indent(4)

		return c.C(`
type Server = struct%s

func New(engine *gin.Engine) *Server {
	s := new(Server)
%s
	return s
}
`).TrimSpace().Format(
			routeStructBody,
			routeInit,
		)
	}
}
