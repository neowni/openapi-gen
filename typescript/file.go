package typescript

import (
	"sort"
	"strings"

	c "columba-livia/content"
)

const (
	modelsPackageName  = "models" // 使用 models 包时的包名，无论原有的包名是什么
	messagePackageName = "message"
)

// 当前正在渲染的文件上下文
var file *_file = nil

type _file struct {
	// 需要导入的 package
	importMap map[string]struct{}

	// 用于 models 与其他包关联
	isModels   bool
	needModels bool

	// 用于 message 与其他包关联
	needMessage bool
}

func (b *_file) modelsNamespace() string {
	if b.isModels {
		return ""
	}

	b.needModels = true
	return modelsPackageName + "."
}

type render = func() c.C

func imports() c.C {
	importLine := make([]string, 0, len(file.importMap))
	for line := range file.importMap {
		importLine = append(importLine, line)
	}

	sort.Slice(importLine, func(i, j int) bool {
		is, _ := strings.CutPrefix(importLine[i], "import ")
		is, _ = strings.CutPrefix(is, "{ ")
		is = strings.ToLower(is)

		js, _ := strings.CutPrefix(importLine[j], "import ")
		js, _ = strings.CutPrefix(js, "{ ")
		js = strings.ToLower(js)

		return is < js
	})

	importLineList := make([]c.C, 0, len(file.importMap))

	if len(importLine) != 0 {
		importLineList = append(importLineList, "")
		for _, line := range importLine {
			importLineList = append(importLineList, c.C(line))
		}
	}

	return c.List(0, importLineList...)
}
