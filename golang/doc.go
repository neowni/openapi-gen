package golang

import (
	"fmt"
	"strings"
	"unicode"

	c "columba-livia/content"
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

// ExportName //
// 转换为标准标志符名称
func ExportName(
	name string,
) (exportName string) {
	runes := []rune(name)
	runes[0] = unicode.ToTitle(runes[0])
	exportName = string(runes)

	// 特定首字母缩写
	exportName = strings.ReplaceAll(exportName, "Id", "ID")
	exportName = strings.ReplaceAll(exportName, "Api", "API")
	exportName = strings.ReplaceAll(exportName, "Url", "URL")
	exportName = strings.ReplaceAll(exportName, "Uri", "URI")
	exportName = strings.ReplaceAll(exportName, "Dns", "DNS")
	exportName = strings.ReplaceAll(exportName, "Uuid", "UUID")
	exportName = strings.ReplaceAll(exportName, "Json", "JSON")

	return exportName
}
