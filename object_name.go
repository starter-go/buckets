package buckets

// PathString 是一个表示路径的字符串 (以 '/' 分隔路径中的各个元素)
type PathString string

////////////////////////////////////////////////////////////////////////////////

// ObjectName 表示对象的索引名称 (类似于一个'路径')
type ObjectName PathString

func (name ObjectName) String() string {
	return string(name)
}
