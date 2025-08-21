package core

import (
	"encoding/json"
	"fmt"
	"sort"
	"sync"

	"github.com/starter-go/buckets"
)

type BucketDriverManagerImpl struct {

	//starter:component

	_as func(buckets.DriverManager) //starter:as("#")

	RawDriverList []buckets.DriverRegistry //starter:inject(".")

	cache *innerBucketDriversCache
	mutex sync.Mutex
}

func (inst *BucketDriverManagerImpl) _impl() buckets.DriverManager {
	return inst
}

func (inst *BucketDriverManagerImpl) FindDriver(cfg *buckets.Configuration) (buckets.Driver, error) {

	inst.mutex.Lock()
	defer inst.mutex.Unlock()

	cache := inst.getCache()
	reg, err := cache.find(cfg)
	if err != nil {
		return nil, err
	}
	return reg.Driver, err
}

func (inst *BucketDriverManagerImpl) getCache() *innerBucketDriversCache {
	c := inst.cache
	if c == nil {
		c = inst.loadCache()
		inst.cache = c
	}
	return c
}

func (inst *BucketDriverManagerImpl) loadCache() *innerBucketDriversCache {

	cache := new(innerBucketDriversCache)
	cache.init()
	src := inst.RawDriverList

	for _, r1 := range src {
		tmp := r1.ListDriverRegistrations()
		for _, r2 := range tmp {
			cache.add(r2)
		}
	}

	cache.sort()
	cache.loadTable()

	return cache
}

////////////////////////////////////////////////////////////////////////////////

type innerBucketDriversCache struct {
	cachedDriverList  []*buckets.DriverRegistration
	cachedDriverTable map[string]*buckets.DriverRegistration
}

func (inst *innerBucketDriversCache) init() {
	inst.cachedDriverList = make([]*buckets.DriverRegistration, 0)
	inst.cachedDriverTable = make(map[string]*buckets.DriverRegistration)
}

func (inst *innerBucketDriversCache) isDriverReady(item *buckets.DriverRegistration) bool {

	if item == nil {
		return false
	}

	if !item.Enabled {
		return false
	}

	if item.Driver == nil {
		return false
	}

	return true
}

func (inst *innerBucketDriversCache) add(item *buckets.DriverRegistration) {
	if !inst.isDriverReady(item) {
		return
	}
	inst.cachedDriverList = append(inst.cachedDriverList, item)
}

func (inst *innerBucketDriversCache) loadTable() {
	src := inst.cachedDriverList
	dst := inst.cachedDriverTable
	for _, item := range src {
		key := item.Name
		older := dst[key]
		if older == nil {
			dst[key] = item
		}
	}
}

func (inst *innerBucketDriversCache) sort() {
	sort.Sort(inst)
}
func (inst *innerBucketDriversCache) Len() int {
	return len(inst.cachedDriverList)
}
func (inst *innerBucketDriversCache) Less(i1, i2 int) bool {
	list := inst.cachedDriverList
	o1 := list[i1]
	o2 := list[i2]
	return (o1.Priority > o2.Priority)
}
func (inst *innerBucketDriversCache) Swap(i1, i2 int) {
	list := inst.cachedDriverList
	list[i1], list[i2] = list[i2], list[i1]
}

func (inst *innerBucketDriversCache) find(cfg *buckets.Configuration) (*buckets.DriverRegistration, error) {
	item, err := inst.find0(cfg)
	if err != nil {
		return nil, err
	}
	result := new(buckets.DriverRegistration)
	*result = *item
	return result, nil
}

func (inst *innerBucketDriversCache) find0(cfg *buckets.Configuration) (*buckets.DriverRegistration, error) {
	item, err := inst.find1(cfg)
	if err == nil && item != nil {
		return item, nil
	}
	return inst.find2(cfg)
}

func (inst *innerBucketDriversCache) find1(cfg *buckets.Configuration) (*buckets.DriverRegistration, error) {
	// by name
	name := cfg.Driver
	tab := inst.cachedDriverTable
	item := tab[name]
	if item == nil {
		return nil, fmt.Errorf("no bucket-driver with name: " + name)
	}
	return item, nil
}

func (inst *innerBucketDriversCache) find2(cfg *buckets.Configuration) (*buckets.DriverRegistration, error) {
	// by accept

	all := inst.cachedDriverList
	for _, item := range all {
		if item.Driver.Accept(cfg) {
			return item, nil
		}
	}

	c2 := new(buckets.Configuration)
	*c2 = *cfg
	c2.AccessKeyID = "***"
	c2.AccessKeySecret = "***"
	js, err := json.Marshal(c2)
	if err != nil {
		return nil, err
	}
	jstr := string(js)
	return nil, fmt.Errorf("no bucket driver accept the config: " + jstr)
}

////////////////////////////////////////////////////////////////////////////////
