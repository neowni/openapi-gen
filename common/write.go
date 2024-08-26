package common

import (
	"os"
	"path/filepath"
	"strings"
)

func P[T any](v T) *T {
	return &v
}

// 																				写入文件

func WriteDir(path string, fileMap map[string]string) (err error) {
	err = os.MkdirAll(path, 0o755)
	if err != nil {
		return err
	}

	for file, content := range fileMap {
		// 渲染内容
		lines := strings.Split(content, "\n")
		for i, line := range lines {
			lines[i] = strings.TrimRight(line, " ")
		}
		if lines[len(lines)-1] != "" {
			lines = append(lines, "")
		}

		// 创建目录
		err = os.MkdirAll(filepath.Join(path, filepath.Dir(file)), 0o755)
		if err != nil {
			return err
		}

		// 写入文件
		err = os.WriteFile(
			filepath.Join(path, file),
			[]byte(strings.Join(lines, "\n")),
			0o644,
		)
		if err != nil {
			return err
		}
	}

	// 删除其他非 raw 开头的文件
	entries, err := os.ReadDir(path)
	if err != nil {
		return err
	}
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		_, exist := fileMap[entry.Name()]
		if exist {
			continue
		}

		// 排除的文件
		if strings.HasPrefix(entry.Name(), "raw") {
			continue
		}

		err = os.Remove(
			filepath.Join(path, entry.Name()),
		)
		if err != nil {
			return err
		}
	}

	return nil
}
