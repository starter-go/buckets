package buckets

import (
	"context"
	"io"
)

// // ObjectMeta 结构包含对象的元数据信息
// type ObjectMeta struct {
// }

// Object  代表一个具体的对象
type Object struct {

	// 操作该对象的上下文
	Context context.Context

	// 对象的名称 (相当于路径)
	Name ObjectName

	// 返回包含该对象的 bucket
	Bucket Bucket

	// 对象的大小 (负数表示空值)
	Size int64

	// 对象的 mime-type
	Type ContentType

	// 校验和
	Sum SUM

	// 对象是否存在
	Existed bool

	// 对象的数据来源
	Data io.ReadCloser
}
