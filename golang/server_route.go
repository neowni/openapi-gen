package golang

import (
	"github.com/pb33f/libopenapi/datamodel/high/base"

	c "columba-livia/content"
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
		).IndentTab(1))

		routeInit := c.List(0,
			c.ForList(tags, func(item *base.Tag) c.C {
				return c.C("s.%s = &%s{engine}").Format(
					ExportName(item.Name),
					item.Name,
				)
			})...,
		).IndentTab(1)

		return c.C(`
type Server = struct %s

func New(engine gin.IRoutes) *Server {
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
