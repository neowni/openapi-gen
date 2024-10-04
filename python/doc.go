package python

import (
	"strings"

	c "columba-livia/content"
)

func comment(
	comment string,
) c.C {
	if comment == "" {
		return ""
	}

	comment = strings.TrimSpace(comment)

	lines := make([]string, 0)

	for _, line := range strings.Split(
		strings.TrimSpace(comment),
		"\n",
	) {
		lines = append(lines, "# "+line)
	}
	return c.C(strings.Join(lines, "\n"))
}

func doc(
	doc string,
) c.C {
	if doc == "" {
		doc = "-"
	}
	return c.F(`"""{{.}}"""`).Format(doc)
}
