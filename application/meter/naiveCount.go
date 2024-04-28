package meter

var txCount uint64
var BlockCount uint64

// 数数，很简单的。在节点上计算就好了。
func onCommited(data struct { // 我发现这种匿名结构体可以拿来做解耦。类似鸭子类型的interface？有效避免循环引用。
	tx    int
	block int
}) {

}
