package localfiles

import (
	"fmt"
	"io"

	"github.com/starter-go/afs"
)

type innerDataWriter struct {
	dataFile afs.Path
	opened   bool
	err      error
	writer   io.Writer
	closer   io.Closer
}

func (inst *innerDataWriter) as() io.WriteCloser {
	return inst
}

func (inst *innerDataWriter) open() error {

	if inst.opened {
		return inst.err
	}

	file := inst.dataFile
	om := afs.OptionsMaker{}

	// make parent dir
	om.SetFlags().Create()
	om.SetPermissions().SetMode(7, 5, 5)
	opt1 := om.Options()
	file.MakeParents(&opt1)

	// write file
	om.SetFlags().Create().Truncate().WriteOnly()
	om.SetPermissions().SetMode(6, 4, 4)
	opt2 := om.Options()
	wtr, err := file.GetIO().OpenWriter(&opt2)

	inst.closer = wtr
	inst.writer = wtr
	inst.err = err
	inst.opened = true
	return err
}

func (inst *innerDataWriter) Close() error {
	cl := inst.closer
	inst.closer = nil
	if cl == nil {
		return nil
	}
	return cl.Close()
}

func (inst *innerDataWriter) Write(p []byte) (int, error) {
	w := inst.writer
	if w == nil {
		err := inst.open()
		if err != nil {
			return 0, err
		}
		w = inst.writer
		if w == nil {
			return 0, fmt.Errorf("no writer opened")
		}
	}
	return w.Write(p)
}
