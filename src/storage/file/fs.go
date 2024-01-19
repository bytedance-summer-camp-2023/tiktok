package file

import (
	"context"
	"io"
	"net/url"
	"os"
	"path"
	"tiktok/src/constant/config"
	"tiktok/src/extra/tracing"
	"tiktok/src/utils/logging"

	"github.com/sirupsen/logrus"
)

type FSStorage struct {
}

func (f FSStorage) Upload(ctx context.Context, fileName string, content io.Reader) (output *PutObjectOutput, err error) {
	ctx, span := tracing.Tracer.Start(ctx, "FSStorage-Upload")
	defer span.End()
	logger := logging.LogService("FSStorage.Upload").WithContext(ctx)
	logger = logger.WithFields(logrus.Fields{
		"file_name": fileName,
	})
	logger.Debugf("Process start")

	all, err := io.ReadAll(content)
	if err != nil {
		logger.WithFields(logrus.Fields{
			"err": err,
		}).Debug("Failed reading content")
		return nil, err
	}

	filePath := path.Join(config.EnvCfg.FileSystemStartPath, fileName)
	dir := path.Dir(filePath)
	err = os.MkdirAll(dir, os.FileMode(0755))
	if err != nil {
		logger.WithFields(logrus.Fields{
			"err": err,
		}).Debug("Failed creating directory before writing file")
		return nil, err
	}

	err = os.WriteFile(filePath, all, os.FileMode(0755))
	if err != nil {
		logger.WithFields(logrus.Fields{
			"err": err,
		}).Debug("Failed writing content to file")
		return nil, err
	}

	return &PutObjectOutput{}, nil
}

func (f FSStorage) GetLink(ctx context.Context, fileName string) (string, error) {
	_, span := tracing.Tracer.Start(ctx, "FSStorage-GetLink")
	defer span.End()
	return url.JoinPath(config.EnvCfg.FileSystemBaseUrl, fileName)
}
