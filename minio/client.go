package minio

import (
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"gitlab.tessan.com/data-center/tessan-erp-common/config"
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
