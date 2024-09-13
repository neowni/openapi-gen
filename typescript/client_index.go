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

			filedList = append(filedList, c.C(`%s: %s;`).Format(
				name, name,
			))
			initList = append(initList, c.C(`this.%s = new %s(instance);`).Format(
				name, name,
			))
		}

		return c.C(`
export default class {
%s
  constructor(instance: AxiosInstance) {
%s
  }
}
`).TrimSpace().Format(
			c.List(0, filedList...).IndentSpace(2),
			c.List(0, initList...).IndentSpace(4),
		)
	}
}
