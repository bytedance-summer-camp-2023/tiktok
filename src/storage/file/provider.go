package file

import (
	"context"
	"io"
	"tiktok/src/constant/config"
)

var client storageProvider

type storageProvider interface {
	Upload(ctx context.Context, fileName string, content io.Reader) (*PutObjectOutput, error)
	GetLink(ctx context.Context, fileName string) (string, error)
	GetLocalPath(ctx context.Context, fileName string) string
}

type PutObjectOutput struct{}

func init() {
	switch config.EnvCfg.StorageType { // Append more type here to provide more file action ability
	case "fs":
		client = FSStorage{}
	}
}

func Upload(ctx context.Context, fileName string, content io.Reader) (*PutObjectOutput, error) {
	return client.Upload(ctx, fileName, content)
}

func GetLocalPath(ctx context.Context, fileName string) string {
	return client.GetLocalPath(ctx, fileName)
}

func GetLink(ctx context.Context, fileName string) (link string, err error) {
	return client.GetLink(ctx, fileName)
}
