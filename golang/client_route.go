package golang

import (
	"github.com/pb33f/libopenapi/datamodel/high/base"

	c "columba-livia/content"
)

func clientRoute(
	tags []*base.Tag,
) (render render) {
	return func() c.C {
		file.importMap["github.com/go-resty/resty/v2"] = ""

		routeStructBody := c.BodyF(c.List(0,
			c.ForList(tags, func(item *base.Tag) c.C {
				return c.F("{{.field}} *{{.type}}").Format(map[string]any{
					"field": publicName(item.Name),
					"type":  privateName(item.Name),
				})
			})...,
		).IndentTab(1))

		routeInit := c.List(0,
			c.ForList(tags, func(item *base.Tag) c.C {
				return c.F("c.{{.field}} = &{{.type}}{client}").Format(map[string]any{
					"field": publicName(item.Name),
					"type":  privateName(item.Name),
				})
			})...,
		).IndentTab(1)

		return c.F(`
type Client = struct {{.struct}}

func New(client *resty.Client) *Client {
	c := new(Client)
{{.init}}
	return c
}
`).
			Format(map[string]any{
				"struct": routeStructBody,
				"init":   routeInit,
			})
	}
}
