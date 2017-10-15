package bucket

// Bucket defines a generic lock-free implementation of a Token Bucket.
type Bucket interface {

	// Take attempts to take n tokens out of the bucket.
	Take(int64) int64

	// Put attempts to add n tokens to the bucket.
	Put(int64) int64

	// Close stops the filling of a given bucket if it was started.
	Close()
}
