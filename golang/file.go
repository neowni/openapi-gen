package golang

import (
	"fmt"
	"sort"
	"strings"

	c "columba-livia/content"
)

// 使用 models/message 包时的包名，无论原有的包名是什么
const (
	modelsPackageName  = "models"
	messagePackageName = "message"
)

// 当前正在渲染的文件上下文
var file *_file = nil

type _file struct {
	// 需要导入的 package
	importMap map[string]string

	// 用于 models 与其他包关联
	isModels   bool
	needModels bool // 是否需要导入 models

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

func imports(projectName string) c.C {
	standard := make([]string, 0)
	thirdParty := make([]string, 0)
	internal := make([]string, 0)

	for path, name := range file.importMap {
		line := fmt.Sprintf(`"%s"`, path)
		if name != "" {
			line = name + " " + line
		}

		switch true {
		case strings.HasPrefix(path, "github.com/"):
			thirdParty = append(thirdParty, line)
		case strings.HasPrefix(path, projectName+"/"):
			internal = append(internal, line)
		default:
			standard = append(standard, line)
		}
	}

	sort.Strings(standard)
	sort.Strings(thirdParty)
	sort.Strings(internal)

	// 整理导入
	importList := standard
	if len(thirdParty) != 0 {
		if len(importList) != 0 {
			importList = append(importList, "")
		}
		importList = append(importList, thirdParty...)
	}
	if len(internal) != 0 {
		if len(importList) != 0 {
			importList = append(importList, "")
		}
		importList = append(importList, internal...)
	}

	content := c.C(strings.Join(importList, "\n"))
	if content == "" {
		return ""
	}

	return c.JoinSpace(
		"import",
		c.BodyC(content.IndentTab(1)),
	)
}
