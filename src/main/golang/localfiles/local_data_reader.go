package localfiles

import (
	"fmt"
	"io"

	"github.com/starter-go/afs"
)

type innerDataReader struct {
	dataFile afs.Path
	opened   bool
	err      error
	reader   io.Reader
	closer   io.Closer
}

func (inst *innerDataReader) as() io.ReadCloser {
	return inst
}

func (inst *innerDataReader) open() error {

	if inst.opened {
		return inst.err
	}

	file := inst.dataFile
	om := afs.OptionsMaker{}

	// read file
	om.SetFlags().ReadOnly()
	om.SetPermissions()
	opt2 := om.Options()
	rdr, err := file.GetIO().OpenReader(&opt2)

	inst.closer = rdr
	inst.reader = rdr
	inst.err = err
	inst.opened = true
	return err
}

func (inst *innerDataReader) Close() error {
	cl := inst.closer
	inst.closer = nil
	if cl == nil {
		return nil
	}
	return cl.Close()
}

func (inst *innerDataReader) Read(p []byte) (int, error) {
	r := inst.reader
	if r == nil {
		err := inst.open()
		if err != nil {
			return 0, err
		}
		r = inst.reader
		if r == nil {
			return 0, fmt.Errorf("no reader opened")
		}
	}
	return r.Read(p)
}
