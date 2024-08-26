package common

import "github.com/pb33f/libopenapi/orderedmap"

// 目前支持的 json 类型
const (
	TypeObject  = "object"
	TypeArray   = "array"
	TypeString  = "string"
	TypeBoolean = "boolean"
	TypeInteger = "integer"
	TypeNumber  = "number"
)

func Range[K comparable, V any](
	pairMap *orderedmap.Map[K, V],
) (pairList []orderedmap.Pair[K, V]) {
	if pairMap == nil {
		return make([]orderedmap.Pair[K, V], 0)
	}

	pairList = make([]orderedmap.Pair[K, V], 0, pairMap.Len())

	item := pairMap.First()
	for item != nil {
		pairList = append(pairList, item)
		item = item.Next()
	}

	return pairList
}

func Keys[K comparable, V any](
	pairMap *orderedmap.Map[K, V],
) (pairList []K) {
	if pairMap == nil {
		return make([]K, 0)
	}

	pairList = make([]K, 0, pairMap.Len())

	item := pairMap.First()
	for item != nil {
		pairList = append(pairList, item.Key())
		item = item.Next()
	}

	return pairList
}
