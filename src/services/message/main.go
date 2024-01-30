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
	"tiktok/src/rpc/chat"
	"tiktok/src/rpc/health"
	healthImpl "tiktok/src/services/health"
	"tiktok/src/utils/consul"
	"tiktok/src/utils/logging"
)

func main() {
	tp, err := tracing.SetTraceProvider(config.MessageRpcServerName)

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
	profiling.InitPyroscope("TikTok.ChatService")

	s := grpc.NewServer(
		grpc.UnaryInterceptor(otelgrpc.UnaryServerInterceptor()),
	)

	log := logging.LogService(config.MessageRpcServerName)

	lis, err := net.Listen("tcp", config.MessageRpcServerPort)

	if err != nil {
		log.Panicf("Rpc %s listen happens error: %v", config.MessageRpcServerName, err)
	}

	var srv MessageServiceImpl
	var probe healthImpl.ProbeImpl

	chat.RegisterChatServiceServer(s, srv)

	health.RegisterHealthServer(s, &probe)

	if err := consul.RegisterConsul(config.MessageRpcServerName, config.MessageRpcServerPort); err != nil {
		log.Panicf("Rpc %s register consul happens error for: %v", config.MessageRpcServerName, err)
	}
	srv.New()
	log.Infof("Rpc %s is running at %s now", config.MessageRpcServerName, config.MessageRpcServerPort)
	if err := s.Serve(lis); err != nil {
		log.Panicf("Rpc %s listen happens error for: %v", config.MessageRpcServerName, err)
	}
}
