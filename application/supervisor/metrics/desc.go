package metrics

import (
	"encoding/json"
	"math"
	"strconv"
)

// 这是Desc的组成元素。
type DescElem struct {
	Name string      `json:"name"`
	Desc string      `json:"desc"`
	Val  interface{} `json:"val"`
}

type DescHead struct {
	Name string `json:"name"`
	Desc string `json:"desc"`
}

type Desc struct {
	Name  string     `json:"name"`
	Desc  string     `json:"desc"`
	Elems []DescElem `json:"elems"`
}

func DescDump(d Desc) []byte {
	bs, err := json.Marshal(d)
	if err != nil {
		panic(err)
	}
	return bs
}

// DescBuilder 构建器结构体。方便构造。
type DescBuilder struct {
	desc Desc
}

// NewDescBuilder 创建一个新的 DescBuilder 实例
func NewDescBuilder(name string, desc string) *DescBuilder {
	return &DescBuilder{desc: Desc{Name: name, Desc: desc}}
}

func checkAndConvertFloat(value interface{}) interface{} {
	// 判断类型是否为 float, float32, float64
	switch value.(type) {
	case float32, float64:
		// 检查值是否为 NaN
		if math.IsNaN(value.(float64)) {
			// 如果是 NaN，将其转换为字符串
			return strconv.FormatFloat(value.(float64), 'f', -1, 64)
		}
	}

	// 如果不是 float 类型或者不是 NaN，则原样返回
	return value
}

// AddElem 向描述中添加元素
func (builder *DescBuilder) AddElem(name string, desc string, val interface{}) *DescBuilder {
	builder.desc.Elems = append(builder.desc.Elems, DescElem{Name: name, Desc: desc, Val: checkAndConvertFloat(val)})
	return builder
}

// GetDesc 获取构建好的描述
func (builder *DescBuilder) GetDesc() Desc {
	return builder.desc
}

// -----------------------------------------------------

func DescMarshal(d Desc) []byte {
	a, b := json.Marshal(d)
	if b != nil {
		panic("PrintDesc")
	}
	return a
}

func DescPrintJson(d Desc) string {
	a := DescMarshal(d)
	return string(a)
}
