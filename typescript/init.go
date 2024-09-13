package typescript

import (
	"fmt"
	"path/filepath"

	v3 "github.com/pb33f/libopenapi/datamodel/high/v3"

	"columba-livia/common"
	c "columba-livia/content"
)

type Project struct {
	Project *string `yaml:"project"` // 项目地址，如果为空则使用当前工作目录

	// 相对于项目地址的目录
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
	renderPackage := func(
		path *string,
		fileP func(),
		fileRenderMap map[string]render,
	) (err error) {
		if path == nil {
			return nil
		}

		packagePath := filepath.Join(*p.Project, *path)
		fileMap := make(map[string]string)

		for name, fileRender := range fileRenderMap {
			file = &_file{
				importMap: make(map[string]struct{}),
			}

			fileP()

			content := fileRender()

			if file.needModels && p.Models != nil {
				file.importMap[fmt.Sprintf(
					`import * as %s from "%s";`,
					modelsPackageName, *p.Models,
				)] = struct{}{}
			}

			if file.needMessage && p.Message != nil {
				file.importMap[fmt.Sprintf(
					`import * as %s from "%s";`,
					messagePackageName, *p.Message,
				)] = struct{}{}
			}

			// 整理成文件
			header := c.C(`// 由 columba-livia 生成`)

			imports := imports()

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
			"index.ts": models(doc.Components.Schemas),
		},
	)
	if err != nil {
		return nil, err
	}

	//																			message

	fileRenderMap := make(map[string]render)
	for _, tag := range doc.Tags {
		fileRenderMap[fmt.Sprintf("%s.ts", tag.Name)] = messageAPI(tag, doc.Paths.PathItems)
	}
	fileRenderMap["index.ts"] = messageIndex(doc.Tags)

	err = renderPackage(
		p.Message,
		func() {},
		fileRenderMap,
	)
	if err != nil {
		return nil, err
	}

	//																			client

	fileRenderMap = make(map[string]render)
	for _, tag := range doc.Tags {
		fileRenderMap[fmt.Sprintf("%s.ts", tag.Name)] = clientAPI(tag, doc.Paths.PathItems)
	}
	fileRenderMap["index.ts"] = clientIndex(doc.Tags)

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
