package buckets

import "context"

type Service interface {

	// get a bucket with name
	GetBucket(ctx context.Context, name string) (Bucket, error)
}
