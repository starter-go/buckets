package localfiles

import (
	"fmt"

	"github.com/starter-go/afs"
)

type innerTempFileCtrl struct {
	target afs.Path
	tmp    afs.Path
}

func (inst *innerTempFileCtrl) init(target afs.Path) error {

	fs := target.GetFS()
	dir := target.GetParent()
	name := target.GetName()

	if !dir.Exists() {
		opt := inst.innerMakeOptionsForDir()
		dir.Mkdirs(opt)
	}

	tmp, err := fs.CreateTempFile(name, ".tmp~", dir)
	if err != nil {
		return err
	}

	inst.tmp = tmp
	inst.target = target

	return nil
}

func (inst *innerTempFileCtrl) innerMakeOptionsForDir() *afs.Options {

	om := new(afs.OptionsMaker)
	opt := new(afs.Options)

	om.WriteOnly()
	om.SetMode(7, 5, 5)
	*opt = om.Options()

	return opt
}

func (inst *innerTempFileCtrl) prepareToWrite() (afs.Path, *afs.Options, error) {

	file := inst.tmp
	om := new(afs.OptionsMaker)
	opt := new(afs.Options)

	om.Create()
	om.WriteOnly()
	om.SetMode(6, 4, 4)
	*opt = om.Options()

	return file, opt, nil
}

func (inst *innerTempFileCtrl) commit() error {

	src := inst.tmp
	dst := inst.target
	if dst == nil || src == nil {
		return fmt.Errorf("src|dst file is nil")
	}

	if dst.IsFile() {
		dst.Delete()
	}

	om := new(afs.OptionsMaker)
	om.Create()
	om.WriteOnly()
	om.SetMode(6, 4, 4)
	opt := om.Options()

	return src.MoveTo(dst, &opt)
}

func (inst *innerTempFileCtrl) finish() error {

	tmp := inst.tmp

	// clean

	if tmp != nil {
		if tmp.IsFile() {
			tmp.Delete()
		}
	}

	return nil
}
