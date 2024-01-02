package storage

import "io"

type Store interface {
	Put(bucket, key string, file io.Reader) error
	Get(bucket, key string) (io.Reader, error)
}
