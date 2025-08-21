package buckets

import "context"

type BucketHolder struct {
	context context.Context
	name    string // name of the bucket
	service Service
	bucket  Bucket
	lazy    bool // is-lazy-load
}

func (inst *BucketHolder) SetName(name string) *BucketHolder {
	inst.name = name
	return inst
}

func (inst *BucketHolder) SetService(s Service) *BucketHolder {
	inst.service = s
	return inst
}

func (inst *BucketHolder) SetContext(ctx context.Context) *BucketHolder {
	inst.context = ctx
	return inst
}

func (inst *BucketHolder) SetLazy(lazy bool) *BucketHolder {
	inst.lazy = lazy
	return inst
}

func (inst *BucketHolder) Init() error {
	if inst.lazy {
		return nil
	}
	_, err := inst.GetBucket()
	return err
}

func (inst *BucketHolder) loadBucket() (Bucket, error) {
	ctx := inst.context
	name := inst.name
	return inst.service.GetBucket(ctx, name)
}

func (inst *BucketHolder) GetBucket() (Bucket, error) {
	b := inst.bucket
	if b == nil {
		b2, err := inst.loadBucket()
		if err != nil {
			return nil, err
		}
		b = b2
		inst.bucket = b2
	}
	return b, nil
}
