package mock

import "github.com/starter-go/buckets"

type innerMockObjectHolder struct {
	repo   *myRepoContext
	bucket buckets.Bucket
	object buckets.Object
	data   []byte
}

// func (inst *innerMockObjectHolder) _impl() buckets.Object {
// 	return inst
// }

// func (inst *innerMockObjectHolder) Bucket() buckets.Bucket {
// 	return inst.bucket
// }

// func (inst *innerMockObjectHolder) Name() buckets.ObjectName {
// 	return inst.name
// }
