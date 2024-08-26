package typescript

import (
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
	importList := make([]c.C, 0, len(file.importMap))

	if len(file.importMap) != 0 {
		importList = append(importList, "")
		for line := range file.importMap {
			importList = append(importList, c.C(line))
		}
	}

	return c.List(0, importList...)
}
