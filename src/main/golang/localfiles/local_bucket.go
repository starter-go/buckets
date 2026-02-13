package localfiles

import (
	"context"
	"fmt"
	"io"
	"strings"

	"github.com/starter-go/afs"
	"github.com/starter-go/buckets"
)

type layout struct {
	configFile      afs.Path // @{workspace}/.bucket/config
	dotBucketDir    afs.Path // @{workspace}/.bucket
	workspaceFolder afs.Path // @{workspace}
}

////////////////////////////////////////////////////////////////////////////////

type innerLocalBucket struct {
	config1 buckets.Configuration
	oo      buckets.OpenOptions
	layout  layout
	context context.Context
}

func (inst *innerLocalBucket) _impl() buckets.Bucket {
	return inst
}

// func (inst *innerLocalBucket) getPathOf(o *buckets.Object) (string, error) {
// 	o2, err := inst.toInnerObjectHolder(o)
// 	if err != nil {
// 		return "", err
// 	}
// 	path := o2.dataFile.GetPath()
// 	if path == "" {
// 		return "", fmt.Errorf("bad path")
// 	}
// 	return path, nil
// }

func (inst *innerLocalBucket) isChildOf(child, parent afs.Path) bool {
	sep := child.GetFS().Separator()
	path1 := parent.GetPath()
	path2 := child.GetPath()
	if !strings.HasSuffix(path1, sep) {
		path1 = path1 + sep
	}
	return strings.HasPrefix(path2, path1)
}

func (inst *innerLocalBucket) computeObjectPath(name buckets.ObjectName) afs.Path {
	dot := inst.layout.dotBucketDir
	wks := inst.layout.workspaceFolder
	child := wks.GetChild(name.String())
	if !inst.isChildOf(child, dot) && inst.isChildOf(child, wks) {
		return child
	}
	panic("bad path of object: " + name.String())
}

func (inst *innerLocalBucket) toInnerObjectHolder(o1 *buckets.Object) (*innerObjectHolder, error) {

	if o1 == nil {
		return nil, fmt.Errorf("object is nil")
	}

	path := inst.computeObjectPath(o1.Name)

	h := new(innerObjectHolder)
	h.dataFile = path.GetChild("object.data")
	h.metaFile = path.GetChild("object.meta")
	h.object = *o1
	h.object.Bucket = inst

	return h, nil
}

func (inst *innerLocalBucket) SetContext(ctx context.Context) buckets.Bucket {

	if ctx == nil {
		return inst
	}

	inst.context = ctx
	return inst
}

func (inst *innerLocalBucket) GetContext() context.Context {

	ctx := inst.context

	if ctx == nil {
		ctx = context.Background()
		inst.context = ctx
	}

	return ctx
}

func (inst *innerLocalBucket) GetObject(name buckets.ObjectName) *buckets.Object {
	o := new(buckets.Object)
	o.Bucket = inst
	o.Name = name
	o.Size = -1
	return o
}

func (inst *innerLocalBucket) Fetch(o1 *buckets.Object) (*buckets.Object, error) {

	h, err := inst.toInnerObjectHolder(o1)
	if err != nil {
		return nil, err
	}

	data, err := h.openDataReader()
	if err != nil {
		return nil, err
	}

	o2 := inst.GetObject(o1.Name)
	err = h.readMeta(o2)
	if err != nil {
		return nil, err
	}

	o2.Data = data
	return o2, nil
}

func (inst *innerLocalBucket) Put(o1 *buckets.Object) (*buckets.Object, error) {

	h, err := inst.toInnerObjectHolder(o1)
	if err != nil {
		return nil, err
	}

	o2 := inst.GetObject(h.object.Name)

	wtr, err := h.openDataWriter()
	if err != nil {
		return nil, err
	}
	defer wtr.Close()

	var total int64 = 0
	rdr := o1.Data
	if rdr != nil {
		cb, err := io.Copy(wtr, rdr)
		if err == io.EOF {
			total = cb
		} else if err != nil {
			return nil, err
		}
	}

	err = h.computeMeta(o2)
	if err != nil {
		return nil, err
	}
	o2.Size = total

	err = h.writeMeta(o2)
	if err != nil {
		return nil, err
	}

	// make result
	err = h.readMeta(o2)
	return o2, err
}

func (inst *innerLocalBucket) GetMeta(o1 *buckets.Object) (*buckets.Object, error) {

	h, err := inst.toInnerObjectHolder(o1)
	if err != nil {
		return nil, err
	}

	name := h.object.Name
	o2 := inst.GetObject(name)
	err = h.readMeta(o2)
	return o2, err
}

func (inst *innerLocalBucket) Exists(o1 *buckets.Object) (bool, error) {
	h, err := inst.toInnerObjectHolder(o1)
	if err != nil {
		return false, err
	}
	ex1 := h.metaFile.IsFile()
	ex2 := h.dataFile.IsFile()
	return (ex1 && ex2), nil
}
