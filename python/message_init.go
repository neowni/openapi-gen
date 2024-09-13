package python

import (
	"fmt"

	"github.com/pb33f/libopenapi/datamodel/high/base"

	c "columba-livia/content"
)

func messageInit(
	tags []*base.Tag,
) (render render) {
	return func() c.C {

		for _, tag := range tags {
			file.importMap[fmt.Sprintf(
				"from . import %s",
				tag.Name,
			)] = struct{}{}
		}

		return ""
	}
}
