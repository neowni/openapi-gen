package common

import (
	"fmt"
	"net/http"
	"regexp"
	"strings"

	"github.com/pb33f/libopenapi/datamodel/high/base"
	v3 "github.com/pb33f/libopenapi/datamodel/high/v3"
	"github.com/pb33f/libopenapi/orderedmap"
)

// Tidy //
// 整理 doc
//
//  1. 整理 path 中的 tag
//     如果 Operation 的 tag 为空则添加 tag `default`
//     如果 Operation 的 tag 不在 tags 中定义，则创建定义
//
//  2. Operation
//     - 必须要有 operation id
//     - request/response content 只能够有单一类型
//     - responses，存在 default 则使用 code 200，不存在 2xx 返回值则使用 code 204
func Tidy(
	doc v3.Document,
) v3.Document {
	// 日志
	Log(">>> 预处理 openapi.yaml 文件")

	tags := make(map[string]struct{}, len(doc.Tags))
	for _, tag := range doc.Tags {
		tags[tag.Name] = struct{}{}
	}

	for i := doc.Paths.PathItems.First(); i != nil; i = i.Next() {
		pathItem := i.Value()

		path := i.Key()

		for method, op := range map[string]*v3.Operation{
			http.MethodGet:     pathItem.Get,
			http.MethodPut:     pathItem.Put,
			http.MethodPost:    pathItem.Post,
			http.MethodDelete:  pathItem.Delete,
			http.MethodOptions: pathItem.Options,
			http.MethodHead:    pathItem.Head,
			http.MethodPatch:   pathItem.Patch,
			http.MethodTrace:   pathItem.Trace,
		} {
			if op == nil {
				continue
			}

			// 日志
			Log("%-50s %s", path, method)

			// 需要含有 OperationId
			if op.OperationId == "" {
				panic(fmt.Sprintf("%s must has operationId", i.Key()))
			}

			// 如果不存在 tag ，则默认使用 default tag
			if len(op.Tags) == 0 {
				op.Tags = []string{"default"}
			}

			for _, tag := range op.Tags {
				_, exist := tags[tag]
				if exist {
					continue
				}

				tags[tag] = struct{}{}
				doc.Tags = append(doc.Tags, &base.Tag{
					Name: tag,
				})
			}

			// req

			// content 只能有单一类型
			if op.RequestBody != nil && op.RequestBody.Content.Len() > 1 {
				panic("request content type")
			}

			// rsp

			if op.Responses == nil {
				op.Responses = new(v3.Responses)
				op.Responses.Codes = orderedmap.New[string, *v3.Response]()
			}

			// 默认返回内容使用 200
			if op.Responses.Default != nil {
				op.Responses.Codes.Set("200", op.Responses.Default)
				err := op.Responses.Codes.MoveToFront("200")
				if err != nil {
					panic(err)
				}
			}

			// 不存在 2xx 返回则使用 204
			has200 := false
			for rsp := op.Responses.Codes.First(); rsp != nil; rsp = rsp.Next() {
				if strings.HasPrefix(rsp.Key(), "2") {
					has200 = true
				}

				if rsp.Value().Content == nil {
					continue
				}

				// content 只能有单一类型
				if rsp.Value().Content.Len() > 1 {
					panic("response content type")
				}
			}
			if !has200 {
				op.Responses.Codes.Set("204", &v3.Response{})
			}
		}
	}

	return doc
}

var PathReg = regexp.MustCompile(`\{([^}]*)\}`)

func ref(prefix string) (f func(string) string) {
	return func(s string) string {
		return strings.TrimPrefix(s, prefix)
	}
}

var refComponentsSchemas = ref("#/components/schemas/")
