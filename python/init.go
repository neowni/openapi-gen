package python

import (
	"fmt"
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
				importMap:  make(map[string]struct{}),
				additional: make([]c.C, 0),
			}

			fileP()

			content := fileRender()

			if file.needModels && p.Models != nil {
				file.importMap[fmt.Sprintf(
					"import %s as %s",
					strings.ReplaceAll(*p.Models, "/", "."),
					modelsPackageName,
				)] = struct{}{}
			}

			if file.needMessage && p.Message != nil {
				file.importMap[fmt.Sprintf(
					"import %s as %s",
					strings.ReplaceAll(*p.Message, "/", "."),
					messagePackageName,
				)] = struct{}{}
			}

			// 整理成文件
			header := c.C(`""" 由 columba-livia 生成 """`)

			imports := imports()

			additional := c.List(2, file.additional...)

			fileMap[name] = c.List(
				2,
				c.List(1,
					header,
					imports,
				),
				additional,
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
			"__init__.py": models(doc.Components.Schemas),
		},
	)
	if err != nil {
		return nil, err
	}

	//																			message

	fileRenderMap := make(map[string]render)
	for _, tag := range doc.Tags {
		fileRenderMap[fmt.Sprintf("%s.py", tag.Name)] = messageApi(tag, doc.Paths.PathItems)
	}
	fileRenderMap["__init__.py"] = messageInit(doc.Tags)

	err = renderPackage(
		p.Message,
		func() {},
		fileRenderMap,
	)
	if err != nil {
		return nil, err
	}

	// 																			渲染 server

	fileRenderMap = make(map[string]render)
	for _, tag := range doc.Tags {
		fileRenderMap[fmt.Sprintf("%s.py", tag.Name)] = serverApi(tag, doc.Paths.PathItems)
	}
	fileRenderMap["__init__.py"] = serverInit(doc.Tags)

	err = renderPackage(
		p.Server,
		func() {},
		fileRenderMap,
	)
	if err != nil {
		return nil, err
	}

	return ignoreList, nil
}
