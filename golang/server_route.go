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
				return c.F("{{.field}} *{{.type}}").Format(map[string]any{
					"field": publicName(item.Name),
					"type":  item.Name,
				})
			})...,
		).IndentTab(1))

		routeInit := c.List(0,
			c.ForList(tags, func(item *base.Tag) c.C {
				return c.F("s.{{.field}} = &{{.type}}{engine}").Format(map[string]any{
					"field": publicName(item.Name),
					"type":  item.Name,
				})
			})...,
		).IndentTab(1)

		return c.F(`
type Server = struct {{.body}}

func New(engine gin.IRoutes) *Server {
	s := new(Server)
{{.init}}
	return s
}
`).Format(map[string]any{
			"body": routeStructBody,
			"init": routeInit,
		})
	}
}
