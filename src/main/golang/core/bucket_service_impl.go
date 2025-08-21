package core

import (
	"context"
	"sync"
	"time"

	"github.com/starter-go/application"
	"github.com/starter-go/buckets"
)

////////////////////////////////////////////////////////////////////////////////

type BucketServiceImpl struct {

	//starter:component

	_as func(buckets.Service) //starter:as("#")

	Drivers buckets.DriverManager //starter:inject("#")
	AC      application.Context   //starter:inject("context")

	cache map[string]*innerBucketCacheItem
	mutex sync.Mutex
}

func (inst *BucketServiceImpl) _impl() buckets.Service {
	return inst
}

func (inst *BucketServiceImpl) GetBucket(ctx context.Context, name string) (buckets.Bucket, error) {

	inst.mutex.Lock()
	defer inst.mutex.Unlock()

	cache := inst.getCache()
	item := cache[name]
	if item != nil {
		if item.bucket != nil {
			return item.bucket, nil
		}
	}

	// try load
	b, err := inst.loadBucket(ctx, name)
	if err != nil {
		return nil, err
	}

	// put into cache
	item = new(innerBucketCacheItem)
	item.bucket = b
	item.name = name
	cache[name] = item

	return item.bucket, nil
}

func (inst *BucketServiceImpl) getCache() map[string]*innerBucketCacheItem {
	c := inst.cache
	if c == nil {
		c = make(map[string]*innerBucketCacheItem)
		inst.cache = c
	}
	return c
}

func (inst *BucketServiceImpl) loadBucket(ctx context.Context, name string) (buckets.Bucket, error) {

	if ctx == nil {
		ctx = inst.AC
	}

	cfg, err := inst.loadConfig(ctx, name)
	if err != nil {
		return nil, err
	}

	driver, err := inst.Drivers.FindDriver(cfg)
	if err != nil {
		return nil, err
	}

	opt := &buckets.OpenOptions{
		Timeout: time.Second * 30,
		Context: ctx,
	}

	return driver.GetLoader().Open(cfg, opt)
}

func (inst *BucketServiceImpl) loadConfig(ctx context.Context, name string) (*buckets.Configuration, error) {

	if name == "" {
		name = "default"
	}

	prefix := "bucket." + name + "."
	getter := inst.AC.GetProperties().Getter().Required()
	cfg := new(buckets.Configuration)

	cfg.Name = getter.GetString(prefix + "name")
	cfg.URL = getter.GetString(prefix + "url")
	cfg.Driver = getter.GetString(prefix + "driver")
	cfg.AccessKeyID = getter.GetString(prefix + "access-key-id")
	cfg.AccessKeySecret = getter.GetString(prefix + "access-key-secret")

	err := getter.Error()
	if err != nil {
		return nil, err
	}
	return cfg, nil
}

////////////////////////////////////////////////////////////////////////////////

type innerBucketCacheItem struct {
	name   string
	bucket buckets.Bucket
}

////////////////////////////////////////////////////////////////////////////////
