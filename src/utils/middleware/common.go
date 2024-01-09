package middleware

import (
	"context"
	"github.com/cloudwego/kitex/pkg/endpoint"
	"github.com/cloudwego/kitex/pkg/rpcinfo"
	"go.uber.org/zap"
	z "tiktok/utils/zap"
)

var (
	_      endpoint.Middleware = CommonMiddleware
	logger *zap.SugaredLogger
)

func init() {
	logger = z.InitLogger()
	defer logger.Sync()
}

func CommonMiddleware(next endpoint.Endpoint) endpoint.Endpoint {
	return func(ctx context.Context, req, resp interface{}) (err error) {
		ri := rpcinfo.GetRPCInfo(ctx)
		// get real request
		logger.Debugf("real request: %+v", req)
		// get remote service information
		logger.Debugf("remote service name: %s, remote method: %s", ri.To().ServiceName(), ri.To().Method())
		if err := next(ctx, req, resp); err != nil {
			return err
		}
		// get real response
		logger.Infof("real response: %+v", resp)
		return nil
	}
}
