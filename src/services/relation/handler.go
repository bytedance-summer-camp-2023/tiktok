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
	"tiktok/src/rpc/relation"
	healthImpl "tiktok/src/services/health"
	"tiktok/src/utils/consul"
	"tiktok/src/utils/logging"
)

func main() {
	tp, err := tracing.SetTraceProvider(config.RelationRpcServerName)

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
	profiling.InitPyroscope("GuGoTik.RelationService")

	s := grpc.NewServer(
		grpc.UnaryInterceptor(otelgrpc.UnaryServerInterceptor()),
	)

	log := logging.LogService(config.RelationRpcServerName)
	lis, err := net.Listen("tcp", config.RelationRpcServerPort)

	if err != nil {
		log.Panicf("Rpc %s listen happens error: %v", config.RelationRpcServerName, err)
	}

	var srv RelationServiceImpl
	var probe healthImpl.ProbeImpl
	relation.RegisterRelationServiceServer(s, srv)
	health.RegisterHealthServer(s, &probe)
	if err := consul.RegisterConsul(config.RelationRpcServerName, config.RelationRpcServerPort); err != nil {
		log.Panicf("Rpc %s register consul happens error for: %v", config.RelationRpcServerName, err)
	}
	log.Infof("Rpc %s is running at %s now", config.RelationRpcServerName, config.RelationRpcServerPort)
	if err := s.Serve(lis); err != nil {
		log.Panicf("Rpc %s listen happens error for: %v", config.RelationRpcServerName, err)
	}
}