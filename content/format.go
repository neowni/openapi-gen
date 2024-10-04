package c

import (
	"bytes"
	"strings"
	"text/template"
)

var formatterMap = make(map[string]*template.Template)

type F string

func (f F) Format(value any) C {
	var err error

	// 获取 formatter，避免重复创建
	formatter, exist := formatterMap[string(f)]
	if !exist {
		formatter, err = template.New("").Parse(
			strings.TrimSpace(string(f)),
		)
		if err != nil {
			panic(err)
		}

		formatterMap[string(f)] = formatter
	}

	// 渲染结果
	buffer := new(bytes.Buffer)

	err = formatter.Execute(buffer, value)
	if err != nil {
		panic(err)
	}

	return C(buffer.String())
}
