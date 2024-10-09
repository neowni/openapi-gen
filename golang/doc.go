package golang

import (
	"fmt"
	"strings"
	"unicode"

	"columba-livia/common"
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

func publicName(
	name string,
) (publicName string) {
	name = common.NameCamelCase(name)
	runes := []rune(name)
	runes[0] = unicode.ToUpper(runes[0])
	publicName = string(runes)

	// 特定首字母缩写
	publicName = strings.ReplaceAll(publicName, "Id", "ID")
	publicName = strings.ReplaceAll(publicName, "Ssh", "SSH")
	publicName = strings.ReplaceAll(publicName, "Api", "API")
	publicName = strings.ReplaceAll(publicName, "Url", "URL")
	publicName = strings.ReplaceAll(publicName, "Uri", "URI")
	publicName = strings.ReplaceAll(publicName, "Dns", "DNS")
	publicName = strings.ReplaceAll(publicName, "Uuid", "UUID")
	publicName = strings.ReplaceAll(publicName, "Json", "JSON")

	return publicName
}

func privateName(
	name string,
) (privateName string) {
	name = common.NameCamelCase(name)
	runes := []rune(name)
	runes[0] = unicode.ToLower(runes[0])
	privateName = string(runes)

	return privateName
}
