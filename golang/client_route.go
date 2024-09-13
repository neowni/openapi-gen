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
				return c.C("%s *%s").Format(
					ExportName(item.Name),
					item.Name,
				)
			})...,
		).IndentTab(1))

		routeInit := c.List(0,
			c.ForList(tags, func(item *base.Tag) c.C {
				return c.C("c.%s = &%s{client}").Format(
					ExportName(item.Name),
					item.Name,
				)
			})...,
		).IndentTab(1)

		return c.C(`
type Client = struct%s

func New(client *resty.Client) *Client {
	c := new(Client)
%s
	return c
}
`).TrimSpace().Format(
			routeStructBody,
			routeInit,
		)
	}
}
