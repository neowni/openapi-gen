package c

import (
	"strings"
)

type C string

func (c C) String() string {
	return string(c)
}

// TrimSpace //
// 去除 空格
func (c C) TrimSpace() C {
	return C(strings.TrimSpace(string(c)))
}

// IndentSpace //
// 缩进
func (c C) IndentSpace(indent int) C {
	indentSpace := strings.Repeat(" ", indent)

	lines := strings.Split(string(c), "\n")
	for index, line := range lines {
		if line == "" {
			continue
		}
		lines[index] = indentSpace + line
	}

	return C(strings.Join(lines, "\n"))
}

func (c C) IndentTab(indent int) C {
	indentSpace := strings.Repeat("	", indent)

	lines := strings.Split(string(c), "\n")
	for index, line := range lines {
		if line == "" {
			continue
		}
		lines[index] = indentSpace + line
	}

	return C(strings.Join(lines, "\n"))
}

// BodyF //
// 花括号
func BodyF(c C) C {
	if c == "" {
		return C("{}")
	}

	return C("{\n" + c + "\n}")
}

// BodyS //
// 方括号
func BodyS(c C) C {
	if c == "" {
		return C("[]")
	}

	return C("[\n" + c + "\n]")
}

// BodyC //
// 圆括号
func BodyC(c C) C {
	if c == "" {
		return C("()")
	}

	return C("(\n" + c + "\n)")
}

func List(gap int, list ...C) C {
	if len(list) == 0 {
		return ""
	}

	lines := make([]string, 0)
	for _, item := range list {
		line := item.String()
		if line == "" {
			continue
		}
		lines = append(lines, line)
	}

	gapLine := strings.Repeat("\n", gap+1)
	return C(strings.Join(lines, gapLine))
}

func Flat(lists ...[]C) []C {
	list := make([]C, 0)

	for _, l := range lists {
		list = append(list, l...)
	}

	return list
}

func JoinSpace(list ...C) C {
	return Join(" ", list...)
}

func Join(c string, list ...C) C {
	if len(list) == 0 {
		return ""
	}

	lines := make([]string, 0)
	for _, item := range list {
		line := item.String()
		if line == "" {
			continue
		}
		lines = append(lines, line)
	}

	return C(strings.Join(lines, c))
}

func If(condition bool, c C) C {
	if condition {
		return c
	} else {
		return ""
	}
}

func ForList[T any](
	items []T,
	f func(item T) C,
) []C {
	list := make([]C, len(items))

	for _, item := range items {
		list = append(list, f(item))
	}

	return list
}

func ForMap[K comparable, V any](
	items map[K]V,
	f func(key K, value V) C,
) []C {
	list := make([]C, len(items))

	for key, value := range items {
		list = append(list, f(key, value))
	}

	return list
}
