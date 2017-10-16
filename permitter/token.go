package permitter

import (
	"time"

	"github.com/SimonRichardson/resilience/bucket"
)

type token struct {
	bucket bucket.Bucket
}

func New(n int64, d time.Duration) Permitter {
	return &token{
		bucket: bucket.NewProvisionBucket(n, d),
	}
}

func (t *token) Run(fn func() error) error {
	if value := t.bucket.Take(1); value == 1 {
		defer t.bucket.Put(1)

		return fn()
	}
	return nil
}
