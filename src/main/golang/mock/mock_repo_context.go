package mock

import "github.com/starter-go/buckets"

type myRepoContext struct {
	table map[buckets.ObjectName]*innerMockObjectHolder
}

func (inst *myRepoContext) init() {
	inst.table = make(map[buckets.ObjectName]*innerMockObjectHolder)
}
