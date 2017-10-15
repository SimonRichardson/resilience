package bucket

import (
	"math"
	"time"
)

type provision struct {
	bucket Bucket
	freq   time.Duration
	inc    int64
	stop   chan chan struct{}
}

// NewProvisionBucket auto provisions a bucket at a given frequency rate
func NewProvisionBucket(capacity int64, freq time.Duration) Bucket {
	p := &provision{
		bucket: NewTokenBucket(capacity),
		freq:   freq,
		stop:   make(chan chan struct{}),
	}

	if freq < 0 {
		return p
	} else if evenFreq := time.Duration(1e9 / capacity); freq < evenFreq {
		freq = evenFreq
	}

	p.freq = freq
	p.inc = int64(math.Floor(.5 + (float64(capacity) * freq.Seconds())))

	go p.run()

	return p
}

func (p *provision) Take(n int64) int64 {
	return p.bucket.Take(n)
}

func (p *provision) Put(n int64) int64 {
	return p.bucket.Put(n)
}

func (p *provision) Close() {
	c := make(chan struct{})
	p.stop <- c
	<-c
}

func (p *provision) run() {
	ticker := time.NewTicker(p.freq)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			p.Put(p.inc)
		case q := <-p.stop:
			close(q)
			return
		}
	}
}
