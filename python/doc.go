package python

import (
	c "columba-livia/content"
	"strings"
)

func doc(
	doc string,
) c.C {
	if doc == "" {
		return ""
	}

	doc = strings.TrimSpace(doc)

	lines := make([]string, 0)

	for _, line := range strings.Split(
		strings.TrimSpace(doc),
		"\n",
	) {
		lines = append(lines, "# "+line)
	}
	return c.C(strings.Join(lines, "\n"))
}