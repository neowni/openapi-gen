package golang

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strings"

	v3 "github.com/pb33f/libopenapi/datamodel/high/v3"

	"columba-livia/common"
	c "columba-livia/content"
)

type Project struct {
	Project *string `yaml:"project"` // 项目地址，如果为空则使用当前工作目录

	// 相对于项目地址的目录
	Server  *string `yaml:"server"`
	Client  *string `yaml:"client"`
	Models  *string `yaml:"models"`
	Message *string `yaml:"message"`
}

func (p *Project) Render(
	doc v3.Document,
) (ignoreList []string, err error) {
	if p == nil {
		return nil, nil
	}

	if p.Project == nil {
		p.Project = common.P(".")
	}

	// 																			前期处理

	// 获取项目名
	fileMod, err := os.ReadFile(filepath.Join(*p.Project, "go.mod"))
	if err != nil {
		panic(err)
	}
	fileModItems := strings.Fields(
		strings.SplitN(string(fileMod), "\n", 2)[0],
	)
	if len(fileModItems) != 2 || fileModItems[0] != "module" {
		panic("go.mod format")
	}

	projectName := fileModItems[1]

	// 整理 models/message 导入时的路径
	modelsImportPath := path.Join(projectName, modelsPackageName)
	if p.Models != nil {
		modelsImportPath = path.Join(projectName, *p.Models)
	}
	messageImportPath := path.Join(projectName, messagePackageName)
	if p.Message != nil {
		messageImportPath = path.Join(projectName, *p.Message)
	}

	renderPackage := func(
		path *string,
		pFile func(), // 预处理
		fileRenderMap map[string]render,
	) (err error) {
		if path == nil {
			return nil
		}

		packagePath := filepath.Join(*p.Project, *path)
		packageName := filepath.Base(*path)
		fileMap := make(map[string]string)

		for name, render := range fileRenderMap {
			// 文件上下文
			file = &_file{
				importMap: make(map[string]string),
			}

			// 预处理文件
			pFile()

			// 渲染
			content := render()

			if file.needModels && p.Models != nil {
				file.importMap[modelsImportPath] = modelsPackageName
			}

			if file.needMessage && p.Message != nil {
				file.importMap[messageImportPath] = messagePackageName
			}

			// 整理成文件
			header := c.C("// 由 columba-livia 生成\npackage %s").Format(packageName)

			imports := imports(projectName)

			fileMap[name] = c.List(
				1,
				header,
				imports,
				content,
			).String()

			// 解除上下文
			file = nil
		}
		err = common.WriteDir(packagePath, fileMap)
		if err != nil {
			return err
		}

		ignoreList = append(ignoreList, packagePath)

		return nil
	}

	//																			models

	err = renderPackage(
		p.Models,
		func() {
			file.isModels = true
		},
		map[string]render{
			"models.go": models(doc.Components.Schemas),
		},
	)
	if err != nil {
		return nil, err
	}

	//																			message

	fileRenderMap := make(map[string]render)
	for _, tag := range doc.Tags {
		fileRenderMap[fmt.Sprintf("%s.go", tag.Name)] = message(tag, doc.Paths.PathItems)
	}

	err = renderPackage(
		p.Message,
		func() {},
		fileRenderMap,
	)
	if err != nil {
		return nil, err
	}

	//																			server

	fileRenderMap = make(map[string]render)
	for _, tag := range doc.Tags {
		fileRenderMap[fmt.Sprintf("%s.go", tag.Name)] = serverAPI(tag, doc.Paths.PathItems)
	}
	fileRenderMap["convert.go"] = serverConvert()
	fileRenderMap["route.go"] = serverRoute(doc.Tags)

	err = renderPackage(
		p.Server,
		func() {},
		fileRenderMap,
	)
	if err != nil {
		return nil, err
	}

	//																			client

	fileRenderMap = make(map[string]render)
	for _, tag := range doc.Tags {
		fileRenderMap[fmt.Sprintf("%s.go", tag.Name)] = clientAPI(tag, doc.Paths.PathItems)
	}
	fileRenderMap["convert.go"] = clientConvert()
	fileRenderMap["route.go"] = clientRoute(doc.Tags)

	err = renderPackage(
		p.Client,
		func() {},
		fileRenderMap,
	)
	if err != nil {
		return nil, err
	}

	return ignoreList, nil
}
