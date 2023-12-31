package minio

import (
	"fmt"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

// no using sync.Once
// as currently the usage of Init
// is controlled in the main entrypoint
var client *minio.Client

func Init(endpoint, accessKey, secretKey string) error {
	var err error

	client, err = minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKey, secretKey, ""),
		Secure: false,
	})

	if err != nil {
		return fmt.Errorf("fail to init the minio client")
	}

	return nil
}
