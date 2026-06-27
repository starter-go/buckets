package buckets

import (
	"context"
	"hash"
	"time"
)

type BucketFileAPI interface {
	Bucket() Bucket

	FetchFile(o *ObjectFile) (*ObjectFile, error)

	PutFile(o *ObjectFile) (*ObjectFile, error)
}

type BucketNativeSumAPI interface {
	Bucket() Bucket

	Algorithm() CheckSumAlgorithm

	Hash() hash.Hash
}

type Bucket interface {
	SetContext(ctx context.Context) Bucket

	GetContext() context.Context

	GetObject(name ObjectName) *Object

	Fetch(o *Object) (*Object, error)

	Put(o *Object) (*Object, error)

	GetMeta(o *Object) (*Object, error)

	SetMeta(o *Object) (*Object, error)

	Delete(o *Object) error

	Exists(o *Object) (bool, error)

	ForFiles() BucketFileAPI

	ForSum() BucketNativeSumAPI
}

type OpenOptions struct {
	Context context.Context
	Flag    int
	Timeout time.Duration
}

type Loader interface {
	Open(cfg *Configuration, options *OpenOptions) (Bucket, error)
}
