package buckets

type Driver interface {
	GetLoader() Loader

	GetRegistration() *DriverRegistration

	Accept(cfg *Configuration) bool
}

type DriverManager interface {

	// 根据配置查找对应的驱动
	FindDriver(cfg *Configuration) (Driver, error)
}

type DriverRegistry interface {
	ListDriverRegistrations() []*DriverRegistration
}

type DriverRegistration struct {
	Name     string
	Priority int
	Enabled  bool
	Driver   Driver
}
