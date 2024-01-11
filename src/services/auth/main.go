package main

import (
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"io"
	"net"
	"tiktok/src/constant/config"
	"tiktok/src/gateway/middleware"
	"tiktok/src/rpc/auth"
	"tiktok/src/rpc/health"
	healthImpl "tiktok/src/services/health"
	"tiktok/src/utils/consul"
	"tiktok/src/utils/interceptor"
	"tiktok/src/utils/logging"
)

func main() {
	tracer, closer, err := middleware.NewTracer(config.AuthRpcServerName)
	if err != nil {
		logging.Logger.WithFields(logrus.Fields{
			"err": err,
		}).Errorf("Can not init Jaeger")
		return
	}
	defer func(closer io.Closer) {
		err := closer.Close()
		if err != nil {
			logging.Logger.WithFields(logrus.Fields{
				"err": err,
			}).Errorf("Error when close closer")
		}
	}(closer)
	s := grpc.NewServer(
		grpc.UnaryInterceptor(interceptor.OpentracingServerInterceptor(tracer)),
	)
	log := logging.LogService(config.AuthRpcServerName)
	lis, err := net.Listen("tcp", config.AuthRpcServerPort)

	if err != nil {
		log.Panicf("Rpc %s listen happens error: %v", config.AuthRpcServerName, err)
	}

	var srv AuthServiceImpl
	var probe healthImpl.ProbeImpl
	auth.RegisterAuthServiceServer(s, srv)
	health.RegisterHealthServer(s, &probe)
	if err := consul.RegisterConsul(config.AuthRpcServerName, config.AuthRpcServerPort); err != nil {
		log.Panicf("Rpc %s register consul hanpens error for: %v", config.AuthRpcServerName, err)
	}
	log.Infof("Rpc %s is running at %s now", config.AuthRpcServerName, config.AuthRpcServerPort)
	if err := s.Serve(lis); err != nil {
		log.Panicf("Rpc %s listen hanpens error for: %v", config.AuthRpcServerName, err)
	}
}
