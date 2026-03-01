package mock

import (
	"bytes"
	"context"
	"crypto/sha256"
	"fmt"
	"io"

	"github.com/starter-go/base/lang"
	"github.com/starter-go/buckets"
)

type innerMockBucket struct {
	config  buckets.Configuration
	oo      buckets.OpenOptions
	repo    *myRepoContext
	context context.Context
}

// Delete implements buckets.Bucket.
func (inst *innerMockBucket) Delete(o1 *buckets.Object) error {

	name := o1.Name
	h := inst.repo.table[name]
	inst.repo.table[name] = nil

	if h == nil {
		return fmt.Errorf("no object with name: %v", name)
	}

	return nil
}

// GetContext implements buckets.Bucket.
func (inst *innerMockBucket) GetContext() context.Context {

	ctx := inst.context
	if ctx == nil {
		ctx = context.Background()
		inst.context = ctx
	}
	return ctx
}

// SetContext implements buckets.Bucket.
func (inst *innerMockBucket) SetContext(ctx context.Context) buckets.Bucket {

	if ctx != nil {
		inst.context = ctx
	}

	return inst
}

func (inst *innerMockBucket) _impl() buckets.Bucket {
	return inst
}

func (inst *innerMockBucket) init() {
	repo := new(myRepoContext)
	repo.init()
	inst.repo = repo
}

func (inst *innerMockBucket) GetObject(name buckets.ObjectName) *buckets.Object {
	obj := new(buckets.Object)
	obj.Bucket = inst
	obj.Name = name
	return obj
}

func (inst *innerMockBucket) createHolderForObject(o *buckets.Object) (*innerMockObjectHolder, error) {

	if o == nil {
		return nil, fmt.Errorf("object is nil")
	}

	holder := new(innerMockObjectHolder)
	holder.object = *o
	holder.data = make([]byte, 0)
	holder.bucket = inst
	holder.repo = inst.repo

	// data
	src := o.Data
	if src != nil {
		defer src.Close()
		data, _ := io.ReadAll(src)
		size := len(data)
		holder.data = data
		holder.object.Size = int64(size)
	}

	// check-sum
	sum := sha256.Sum256(holder.data)
	holder.object.Sum.Algorithm = buckets.AlgorithmSHA256
	holder.object.Sum.Value = lang.HexFromBytes(sum[:])
	holder.object.Data = nil

	return holder, nil
}

func (inst *innerMockBucket) Fetch(o1 *buckets.Object) (*buckets.Object, error) {

	name := o1.Name
	h := inst.repo.table[name]
	if h == nil {
		return nil, fmt.Errorf("no object with name: %v", name)
	}
	o2 := new(buckets.Object)
	*o2 = h.object

	// prepare data
	data1 := bytes.NewReader(h.data)
	data2 := io.NopCloser(data1)
	o2.Data = data2

	return o2, nil
}

func (inst *innerMockBucket) Put(o1 *buckets.Object) (*buckets.Object, error) {

	h, err := inst.createHolderForObject(o1)
	if err != nil {
		return nil, err
	}
	name := o1.Name
	inst.repo.table[name] = h
	o2 := new(buckets.Object)
	*o2 = h.object
	return o2, nil
}

func (inst *innerMockBucket) GetMeta(want *buckets.Object) (*buckets.Object, error) {

	name := want.Name
	holder := inst.repo.table[name]
	if holder == nil {
		return nil, fmt.Errorf("no object with name: %v", name)
	}
	have := new(buckets.Object)
	*have = holder.object
	return have, nil
}

func (inst *innerMockBucket) Exists(want *buckets.Object) (bool, error) {

	name := want.Name
	holder := inst.repo.table[name]

	if holder == nil {
		return false, nil
	}

	if holder.data == nil {
		return false, nil
	}

	return true, nil
}
