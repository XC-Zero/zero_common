package minio

import (
	"github.com/XC-Zero/zero_common/config"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

func InitMinio(config config.MinioConfig) (*minio.Client, error) {
	client, err := minio.New(config.EndPoint, &minio.Options{
		Creds:        credentials.NewStaticV4(config.AccessKey, config.SecretKey, ""),
		Secure:       false,
		Transport:    nil,
		Region:       "",
		BucketLookup: 0,
		CustomMD5:    nil,
		CustomSHA256: nil,
	})

	//client.GetBucket
	return client, err
}
