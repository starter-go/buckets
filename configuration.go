package buckets

// Configuration 结构包含了访问 bucket 所需的配置信息
type Configuration struct {
	Name            string // the name of bucket
	URL             string // url to the bucket
	Driver          string // the name of driver
	AccessKeyID     string
	AccessKeySecret string

	Location Location //resolved location
}
