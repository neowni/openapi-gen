package typescript

import (
	"fmt"

	"github.com/pb33f/libopenapi/datamodel/high/base"

	c "columba-livia/content"
)

func messageIndex(
	tags []*base.Tag,
) (render render) {
	return func() c.C {

		for _, tag := range tags {
			file.importMap[fmt.Sprintf(
				`export * as %s from "./%s";`,
				tag.Name, tag.Name,
			)] = struct{}{}
		}

		return ""
	}
}
