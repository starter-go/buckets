package mock

import "github.com/starter-go/buckets"

const theDriverName = "mock"

type Driver struct {

	//starter:component

	_as func(buckets.DriverRegistry) //starter:as('.')

	Enabled  bool //starter:inject("${buckets-driver.mock.enabled}")
	Priority int  //starter:inject("${buckets-driver.mock.priority}")

}

func (inst *Driver) _impl() buckets.DriverRegistry {
	return inst
}

func (inst *Driver) ListDriverRegistrations() []*buckets.DriverRegistration {
	r1 := inst.GetRegistration()
	return []*buckets.DriverRegistration{r1}
}

func (inst *Driver) GetLoader() buckets.Loader {
	return inst
}

func (inst *Driver) GetRegistration() *buckets.DriverRegistration {
	r1 := &buckets.DriverRegistration{
		Name:     theDriverName,
		Enabled:  inst.Enabled,
		Priority: inst.Priority,
		Driver:   inst,
	}
	return r1
}

func (inst *Driver) Accept(cfg *buckets.Configuration) bool {
	if cfg == nil {
		return false
	}
	return (cfg.Driver == theDriverName)
}

func (inst *Driver) Open(cfg *buckets.Configuration, options *buckets.OpenOptions) (buckets.Bucket, error) {

	b := new(innerMockBucket)

	if options != nil {
		b.oo = *options
	}

	if cfg != nil {
		b.config = *cfg
	}

	b.init()
	return b, nil
}
