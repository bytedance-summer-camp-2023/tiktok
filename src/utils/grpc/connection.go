package grpc

import (
	"fmt"
	capi "github.com/hashicorp/consul/api"
	"github.com/sirupsen/logrus"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/keepalive"
	"tiktok/src/utils/consul"
	"tiktok/src/utils/logging"
	"time"
)

func Connect(serviceName string) (conn *grpc.ClientConn) {
	service, err := consul.ResolveService(serviceName)
	if err != nil {
		logging.Logger.WithFields(logrus.Fields{
			"service": serviceName,
			"err":     err,
		}).Errorf("Cannot find %v rpc server", serviceName)
	}

	logging.Logger.Debugf("Found service %v in %v:%v", service.ServiceName, service.Address, service.ServicePort)

	conn, err = connect(service)
	if err != nil {
		logging.Logger.WithFields(logrus.Fields{
			"service": service.ServiceName,
			"err":     err,
		}).Errorf("Cannot connect to %v rpc server in %v:%v", service.ServiceName, service.Address, service.ServicePort)
	}
	return
}

func connect(service *capi.CatalogService) (conn *grpc.ClientConn, err error) {
	var kacp = keepalive.ClientParameters{
		Time:                10 * time.Second, // send pings every 10 seconds if there is no activity
		Timeout:             time.Second,      // wait 1 second for ping ack before considering the connection dead
		PermitWithoutStream: false,            // send pings even without active streams
	}

	conn, err = grpc.Dial(fmt.Sprintf("%v:%v", service.Address, service.ServicePort),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithKeepaliveParams(kacp),
		grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy": "round_robin"}`),
		grpc.WithUnaryInterceptor(otelgrpc.UnaryClientInterceptor()))
	return
}
