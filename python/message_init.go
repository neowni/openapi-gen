package python

import (
	c "columba-livia/content"
	"fmt"

	"github.com/pb33f/libopenapi/datamodel/high/base"
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
