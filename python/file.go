package python

import (
	"fmt"
	"slices"
	"sort"
	"strings"

	c "columba-livia/content"
)

const (
	modelsPackageName  = "_models" // 使用 models 包时的包名，无论原有的包名是什么
	messagePackageName = "_message"
)

// 																				工具函数

var privateNum = 0

func privateName() string {
	privateNum += 1
	return fmt.Sprintf("_private%d", privateNum)
}

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

	// 文件需要额外定义的私有内容
	additional []c.C
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
	importList := make([]string, 0, len(file.importMap))

	// 默认导入 annotations
	importList = append(importList, "from __future__ import annotations as _annotations")

	// 整理 import
	standard := make([]string, 0)
	thirdParty := make([]string, 0)
	internal := make([]string, 0)

	for line := range file.importMap {
		switch importPackageName(line) {
		case "typing", "enum", "collections.abc":
			standard = append(standard, line)
		case "pydantic", "flask":
			thirdParty = append(thirdParty, line)
		default:
			internal = append(internal, line)
		}
	}

	sort.Strings(standard)
	sort.Strings(thirdParty)
	slices.SortFunc(internal, func(a string, b string) int {
		aName := importPackageName(a)
		bName := importPackageName(b)

		r := 0

		if aName < bName {
			r += 2
		} else {
			r -= 2
		}

		if a > b {
			r += 1
		} else {
			r -= 1
		}

		return r
	})

	if len(standard) != 0 {
		importList = append(importList, "")
		importList = append(importList, standard...)
	}

	if len(thirdParty) != 0 {
		importList = append(importList, "")
		importList = append(importList, thirdParty...)
	}

	if len(internal) != 0 {
		importList = append(importList, "")
		importList = append(importList, internal...)
	}

	return c.C(strings.Join(importList, "\n"))
}

func importPackageName(line string) string {
	name := ""

	// 查找包名
	fields := strings.Fields(line)
	index := slices.Index(fields, "from")
	if index == -1 {
		index = slices.Index(fields, "import")
	}

	if index >= 0 {
		name = fields[index+1]
	}

	return name
}
