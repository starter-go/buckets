package localfiles

import (
	"fmt"
	"strings"

	"github.com/starter-go/afs"
	"github.com/starter-go/buckets"
)

const theDriverName = "file"

type Driver struct {

	//starter:component

	_as func(buckets.DriverRegistry) //starter:as(".")

	Enabled  bool //starter:inject("${buckets-driver.file.enabled}")
	Priority int  //starter:inject("${buckets-driver.file.priority}")

	AFS afs.FS //starter:inject("#")

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
	r := new(buckets.DriverRegistration)
	r.Name = theDriverName
	r.Driver = inst
	r.Enabled = inst.Enabled
	r.Priority = inst.Priority
	return r
}

func (inst *Driver) Accept(cfg *buckets.Configuration) bool {
	if cfg == nil {
		return false
	}
	return (cfg.Driver == theDriverName)
}

func (inst *Driver) locateWorkspaceFolder(cfg *buckets.Configuration, options *buckets.OpenOptions) (afs.Path, error) {

	if cfg == nil {
		return nil, fmt.Errorf("configuration is nil")
	}

	l, err := buckets.ParseLocation(cfg.URL)
	if err != nil {
		return nil, err
	}

	cfg.Location = *l
	scheme := l.Scheme
	path := l.Path

	if scheme == "file" && strings.HasPrefix(path, "/") {
		p := inst.AFS.NewPath(path)
		return p, nil
	}

	return nil, fmt.Errorf("bad location of local-file-system bucket: " + cfg.URL)
}

func (inst *Driver) checkLayout(layout *layout) error {
	config := layout.configFile
	ok := config.IsFile()
	if !ok {
		path := config.GetPath()
		return fmt.Errorf("bad layout of local-file-system bucket folder, @" + path)
	}
	return nil
}

func (inst *Driver) Open(cfg *buckets.Configuration, options *buckets.OpenOptions) (buckets.Bucket, error) {

	workspace, err := inst.locateWorkspaceFolder(cfg, options)
	if err != nil {
		return nil, err
	}

	b := new(innerLocalBucket)

	if cfg != nil {
		b.config1 = *cfg
	}

	if options != nil {
		b.oo = *options
	}

	b.layout.workspaceFolder = workspace
	b.layout.dotBucketDir = workspace.GetChild(".bucket")
	b.layout.configFile = workspace.GetChild(".bucket/config")

	err = inst.checkLayout(&b.layout)
	if err != nil {
		return nil, err
	}

	return b, nil
}
