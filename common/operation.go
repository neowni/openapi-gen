package common

import (
	"fmt"
	"net/http"
	"slices"
	"strings"

	"github.com/pb33f/libopenapi/datamodel/high/base"
	v3 "github.com/pb33f/libopenapi/datamodel/high/v3"
	"github.com/pb33f/libopenapi/orderedmap"
)

// 整理 属于 tag 的所有 operation
func TagOperationList(
	tag string,
	pathItems *orderedmap.Map[string, *v3.PathItem],
) (messageList []*Operation) {
	for _, pathPair := range Range(pathItems) {
		path := pathPair.Key()
		pathItem := pathPair.Value()

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

			if !slices.Contains(op.Tags, tag) {
				continue
			}

			uri := make([]*v3.Parameter, 0)
			qry := make([]*v3.Parameter, 0)

			for _, parameter := range slices.Concat(
				pathItem.Parameters,
				op.Parameters,
			) {
				switch parameter.In {
				case "path":
					uri = append(uri, parameter)
				case "query":
					qry = append(qry, parameter)
				}
			}

			req := op.RequestBody
			rsp := make([]orderedmap.Pair[string, *v3.Response], 0)

			for _, responsePair := range Range(op.Responses.Codes) {
				if !strings.HasPrefix(responsePair.Key(), "2") {
					continue
				}

				rsp = append(rsp, responsePair)
			}

			operation := &Operation{
				Tag:         tag,
				ID:          op.OperationId,
				Description: op.Description,

				Path:   path,
				Method: method,

				URI: uri,
				Qry: qry,
				Req: req,
				Rsp: rsp,
			}

			messageList = append(messageList, operation)
		}
	}

	return messageList
}

// Operation 详细信息整理
type Operation struct {
	Tag         string
	ID          string
	Description string

	Path   string
	Method string

	URI []*v3.Parameter
	Qry []*v3.Parameter
	Req *v3.RequestBody
	Rsp []orderedmap.Pair[string, *v3.Response] // 只包含 2xx 的响应码
}

// 																				content 类型

type ContentType string

const (
	ContentText  = ContentType("text")
	ContentJson  = ContentType("json")
	ContentEmpty = ContentType("empty")
)

type ContentSchema struct {
	ContentType ContentType
	SchemaProxy *base.SchemaProxy

	RspCode string // 如果是 rsp 的 content，则表示其返回状态码
}

func contentType(
	contentsMap *orderedmap.Map[string, *v3.MediaType],
) (
	contentSchema ContentSchema,
) {
	// 空类型
	if contentsMap == nil || contentsMap.Len() == 0 {
		return ContentSchema{
			ContentType: ContentEmpty,
			SchemaProxy: nil,
		}
	}

	switch t := contentsMap.First().Key(); t {
	case "text/plain":
		return ContentSchema{
			ContentType: ContentText,
			SchemaProxy: contentsMap.First().Value().Schema,
		}
	case "application/json":
		return ContentSchema{
			ContentType: ContentJson,
			SchemaProxy: contentsMap.First().Value().Schema,
		}
	default:
		panic(fmt.Sprintf("media type %s not support", t))
	}
}

//																				uri / qry

func ParSchemaProxy(par *v3.Parameter) (
	schemaProxy *base.SchemaProxy,
) {
	return par.Schema
}

func ParRequired(par *v3.Parameter) (required bool) {
	return par.Required != nil && *par.Required
}

// 																				req

func ReqSchemaProxy(req *v3.RequestBody) (
	contentSchema ContentSchema,
) {
	if req == nil {
		return ContentSchema{
			ContentType: ContentEmpty,
			SchemaProxy: nil,
		}
	} else {
		return contentType(req.Content)
	}
}

// 																				rsp

func RspSchemaProxy(rspList []orderedmap.Pair[string, *v3.Response]) (
	contentSchemaList []ContentSchema,
) {
	contentSchemaList = make([]ContentSchema, 0, len(rspList))

	for _, rsp := range rspList {
		contentType := contentType(rsp.Value().Content)
		contentType.RspCode = rsp.Key()

		contentSchemaList = append(contentSchemaList, contentType)
	}

	return contentSchemaList
}
