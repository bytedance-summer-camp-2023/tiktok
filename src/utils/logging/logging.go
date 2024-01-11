package logging

import (
	"github.com/opentracing/opentracing-go"
	log "github.com/sirupsen/logrus"
	"github.com/uber/jaeger-client-go"
	"os"
	"tiktok/src/constant/config"
)

func init() {
	log.SetOutput(os.Stdout)
	log.SetFormatter(&log.JSONFormatter{
		PrettyPrint: true,
	})
	switch config.EnvCfg.LoggerLevel {
	case "DEBUG":
		log.SetLevel(log.DebugLevel)
	case "INFO":
		log.SetLevel(log.InfoLevel)
	case "WARN", "WARNING":
		log.SetLevel(log.WarnLevel)
	case "ERROR":
		log.SetLevel(log.ErrorLevel)
	case "FATAL":
		log.SetLevel(log.FatalLevel)
	}
}

var Logger = log.WithFields(log.Fields{
	"Tied": config.EnvCfg.TiedLogging,
})

func LogMethod(name string) *log.Entry {
	return Logger.WithFields(log.Fields{
		"Action": name,
	})
}

func LogService(name string) *log.Entry {
	return Logger.WithFields(log.Fields{
		"Service": name,
	})
}

func GetSpanLogger(span opentracing.Span, method string) *log.Entry {
	return log.WithFields(log.Fields{
		"operation": method,
		"trace_id":  span.Context().(jaeger.SpanContext).TraceID().String(),
		"span_id":   span.Context().(jaeger.SpanContext).SpanID().String(),
	})
}
