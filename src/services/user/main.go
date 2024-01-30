package main

import (
	"context"
	"github.com/sirupsen/logrus"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"google.golang.org/grpc"
	"net"
	"tiktok/src/constant/config"
	"tiktok/src/extra/profiling"
	"tiktok/src/extra/tracing"
	"tiktok/src/rpc/health"
	"tiktok/src/rpc/user"
	healthImpl "tiktok/src/services/health"
	"tiktok/src/utils/consul"
	"tiktok/src/utils/logging"
)

func main() {
	tp, err := tracing.SetTraceProvider(config.UserRpcServerName)

	if err != nil {
		logging.Logger.WithFields(logrus.Fields{
			"err": err,
		}).Panicf("Error to set the trace")
	}
	defer func() {
		if err := tp.Shutdown(context.Background()); err != nil {
			logging.Logger.WithFields(logrus.Fields{
				"err": err,
			}).Errorf("Error to set the trace")
		}
	}()

	// Configure Pyroscope
	profiling.InitPyroscope("TikTok.UserService")

	s := grpc.NewServer(
		grpc.UnaryInterceptor(otelgrpc.UnaryServerInterceptor()),
	)

	log := logging.LogService(config.UserRpcServerName)
	lis, err := net.Listen("tcp", config.UserRpcServerPort)

	if err != nil {
		log.Panicf("Rpc %s listen happens error: %v", config.UserRpcServerName, err)
	}

	var srv UserServiceImpl
	var probe healthImpl.ProbeImpl
	user.RegisterUserServiceServer(s, srv)
	health.RegisterHealthServer(s, &probe)
	if err := consul.RegisterConsul(config.UserRpcServerName, config.UserRpcServerPort); err != nil {
		log.Panicf("Rpc %s register consul hanpens error for: %v", config.UserRpcServerName, err)
	}
	srv.New()
	log.Infof("Rpc %s is running at %s now", config.UserRpcServerName, config.UserRpcServerPort)
	if err := s.Serve(lis); err != nil {
		log.Panicf("Rpc %s listen hanpens error for: %v", config.UserRpcServerName, err)
	}
}
