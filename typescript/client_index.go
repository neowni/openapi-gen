package typescript

import (
	"fmt"

	"github.com/pb33f/libopenapi/datamodel/high/base"

	c "columba-livia/content"
)

func clientIndex(
	tags []*base.Tag,
) (render render) {
	return func() c.C {
		file.importMap[`import { AxiosInstance } from "axios";`] = struct{}{}

		filedList := make([]c.C, 0)
		initList := make([]c.C, 0)

		for _, tag := range tags {
			name := tag.Name

			file.importMap[fmt.Sprintf(`import %s from "./%s";`, name, name)] = struct{}{}

			filedList = append(filedList, c.F(`{{.}}: {{.}};`).Format(name))
			initList = append(initList, c.F(`this.{{.}} = new {{.}}(instance);`).Format(name))
		}

		return c.F(`
export default class {
{{.body}}
  constructor(instance: AxiosInstance) {
{{.init}}
  }
}
`).Format(map[string]any{
			"body": c.List(0, filedList...).IndentSpace(2),
			"init": c.List(0, initList...).IndentSpace(4),
		})
	}
}
