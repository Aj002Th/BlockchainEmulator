package metrics

import "encoding/json"

// 这是Desc的组成元素。
type DescElem struct {
	Name string      `json:"name"`
	Desc string      `json:"desc"`
	Val  interface{} `json:"val"`
}

type Desc = []DescElem

// // 为GetDesc提供便利。反正遵循这个框架。Go没有type alias所以只能这样
// type Desc struct {
// 	val []DescElem
// }

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
func NewDescBuilder() *DescBuilder {
	return &DescBuilder{}
}

// AddElem 向描述中添加元素
func (builder *DescBuilder) AddElem(name string, desc string, val interface{}) *DescBuilder {
	builder.desc = append(builder.desc, DescElem{Name: name, Desc: desc, Val: val})
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
