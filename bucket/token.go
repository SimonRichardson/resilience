package bucket

import "sync/atomic"

type tokenBucket struct {
	tokens, capacity int64
}

// NewTokenBucket creates a new Bucket with a max capacity
func NewTokenBucket(capacity int64) Bucket {
	return &tokenBucket{
		tokens:   capacity,
		capacity: capacity,
	}
}

func (b *tokenBucket) Take(n int64) (taken int64) {
TAKE:
	if tokens := atomic.LoadInt64(&b.tokens); tokens == 0 {
		return 0
	} else if n <= tokens {
		if !atomic.CompareAndSwapInt64(&b.tokens, tokens, tokens-n) {
			goto TAKE
		}
		return n
	} else {
		if !atomic.CompareAndSwapInt64(&b.tokens, tokens, 0) {
			goto TAKE
		}
		return tokens
	}
}

func (b *tokenBucket) Put(n int64) (added int64) {
PUT:
	if tokens := atomic.LoadInt64(&b.tokens); tokens == b.capacity {
		return 0
	} else if left := b.capacity - tokens; n <= left {
		if !atomic.CompareAndSwapInt64(&b.tokens, tokens, tokens+n) {
			goto PUT
		}
		return n
	} else {
		if !atomic.CompareAndSwapInt64(&b.tokens, tokens, b.capacity) {
			goto PUT
		}
		return left
	}
}

func (b *tokenBucket) Close() {}
