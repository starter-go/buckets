package buckets

// Configuration 结构包含了访问 bucket 所需的配置信息
type Configuration struct {
	Name            string // the name of bucket
	URL             string // url to the bucket
	Driver          string // the name of driver
	AccessKeyID     string
	AccessKeySecret string

	MaxObjectSize int // 指定存储桶支持的最大对象大小 (0表示加载默认值; 小于0表示无限制)

	Location Location //resolved location
}
