package grpc

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/keepalive"
	"tiktok/src/constant/config"
	"tiktok/src/utils/logging"
	"time"
)

func Connect(serviceName string) (conn *grpc.ClientConn) {
	kacp := keepalive.ClientParameters{
		Time:                10 * time.Second, // send pings every 10 seconds if there is no activity
		Timeout:             time.Second,      // wait 1 second for ping ack before considering the connection dead
		PermitWithoutStream: true,             // send pings even without active streams
	}

	conn, err := grpc.Dial(
		fmt.Sprintf("consul://%s/%s?wait=15s", config.EnvCfg.ConsulAddr, config.EnvCfg.ConsulAnonymityPrefix+serviceName),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy": "round_robin"}`),
		grpc.WithUnaryInterceptor(otelgrpc.UnaryClientInterceptor()),
		grpc.WithKeepaliveParams(kacp),
	)

	logging.Logger.Debugf("connect")

	if err != nil {
		logging.Logger.WithFields(logrus.Fields{
			"service": config.EnvCfg.ConsulAnonymityPrefix + serviceName,
			"err":     err,
		}).Errorf("Cannot connect to %v service", config.EnvCfg.ConsulAnonymityPrefix+serviceName)
	}
	return
}
