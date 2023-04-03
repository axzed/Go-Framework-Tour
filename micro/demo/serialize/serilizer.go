package serialize

type Serializer interface {
	// Code 协议里面对应字段的值
	Code() byte
	// Encode 编码
	Encode(val any) ([]byte, error)
	// Decode 解码
	// 如果你要求实现者处理 data==nil 或者 data.len == 0 的场景
	// 要在这里说明清楚
	Decode(data []byte, val any) error
}
