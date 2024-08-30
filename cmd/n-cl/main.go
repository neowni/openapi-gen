package main

import (
	"os"
	"path/filepath"
	"slices"
	"strings"

	"github.com/alecthomas/kong"
	"github.com/pb33f/libopenapi"
	"gopkg.in/yaml.v3"

	"columba-livia/common"
	"columba-livia/golang"
	"columba-livia/python"
	"columba-livia/typescript"
)

func main() {
	ctx := kong.Parse(&CLI{})

	err := ctx.Run()
	if err != nil {
		panic(err)
	}
}

type CLI struct {
	WorkDir    string `arg:"" help:"工作目录" type:"path" default:"-"`
	ConfigPath string `arg:"" help:"配置文件路径" type:"string" default:"./.n-cl.yaml"`
}

func (f *CLI) Run() (err error) {
	if f.WorkDir != "-" {
		err = os.Chdir(f.WorkDir)
		if err != nil {
			return err
		}
	} else {
		f.WorkDir = "."
	}

	// 读取配置
	config := new(config)
	file, err := os.Open(
		filepath.Join(f.WorkDir, f.ConfigPath),
	)
	if err != nil {
		return err
	}
	err = yaml.NewDecoder(file).Decode(config)
	if err != nil {
		return err
	}

	if config.OpenAPI == "" {
		config.OpenAPI = "openapi.yaml"
	}

	// 解析 openapi
	content, err := os.ReadFile(config.OpenAPI)
	if err != nil {
		return err
	}
	document, err := libopenapi.NewDocument(content)
	if err != nil {
		return err
	}

	v3Model, errors := document.BuildV3Model()
	if len(errors) != 0 {
		return errors[0]
	}

	// 生成
	doc := common.Tidy(v3Model.Model)

	pathList := make([]string, 0)

	// golang
	ignoreList, err := config.Golang.Render(doc)
	if err != nil {
		return err
	}
	pathList = append(pathList, ignoreList...)

	// python
	ignoreList, err = config.Python.Render(doc)
	if err != nil {
		return err
	}
	pathList = append(pathList, ignoreList...)

	// typescript
	ignoreList, err = config.TypeScript.Render(doc)
	if err != nil {
		return err
	}
	pathList = append(pathList, ignoreList...)

	pathList = slices.DeleteFunc(pathList, func(path string) bool {
		return path == "/"
	})

	// 修改 gitignore
	content, err = os.ReadFile(".gitignore")
	if err != nil {
		if os.IsNotExist(err) {
			content = []byte("")
		} else {
			return err
		}
	}

	lineList := strings.Split(string(content), "\n")
	startIndex := slices.Index(lineList, gitignoreStart)
	endIndex := slices.Index(lineList, gitignoreEnd)

	if startIndex != -1 {
		lineList = slices.Delete(lineList, startIndex, endIndex+1)
	}
	lineList = []string{
		strings.TrimSpace(strings.Join(lineList, "\n")),
		"",
	}

	lineList = append(lineList, gitignoreStart)
	lineList = append(lineList, pathList...)
	lineList = append(lineList, gitignoreEnd, "")

	err = os.WriteFile(".gitignore", []byte(strings.Join(lineList, "\n")), 0o644)
	if err != nil {
		return err
	}

	return nil
}

type config struct {
	OpenAPI string `yaml:"openapi"`

	Golang     *golang.Project     `yaml:"golang"`
	Python     *python.Project     `yaml:"python"`
	TypeScript *typescript.Project `yaml:"typescript"`
}

const (
	gitignoreStart = "# columba-livia v"
	gitignoreEnd   = "# columba-livia ^"
)
