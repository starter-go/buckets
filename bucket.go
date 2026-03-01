package buckets

import (
	"context"
	"time"
)

type Bucket interface {
	SetContext(ctx context.Context) Bucket

	GetContext() context.Context

	GetObject(name ObjectName) *Object

	Fetch(o *Object) (*Object, error)

	Put(o *Object) (*Object, error)

	GetMeta(o *Object) (*Object, error)

	Delete(o *Object) error

	Exists(o *Object) (bool, error)
}

type OpenOptions struct {
	Context context.Context
	Flag    int
	Timeout time.Duration
}

type Loader interface {
	Open(cfg *Configuration, options *OpenOptions) (Bucket, error)
}
