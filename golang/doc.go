package golang

import (
	c "columba-livia/content"
	"fmt"
	"strings"

	"github.com/iancoleman/strcase"
)

func doc(
	name string,
	doc string,
) c.C {
	if doc == "" {
		return ""
	}

	doc = strings.TrimSpace(doc)

	lines := make([]string, 0)
	if name != "" {
		lines = append(lines, fmt.Sprintf("// %s //", name))
	}

	for _, line := range strings.Split(
		strings.TrimSpace(doc),
		"\n",
	) {
		lines = append(lines, "// "+line)
	}
	return c.C(strings.Join(lines, "\n"))
}

func ExportName(
	name string,
) string {
	return strcase.ToCamel(name)
}
