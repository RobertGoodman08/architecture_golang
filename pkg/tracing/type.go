package tracing

import (
	"fmt"
	"io"

	opentracing "github.com/opentracing/opentracing-go"
	"github.com/spf13/viper"
	"github.com/uber/jaeger-client-go"
	"github.com/uber/jaeger-client-go/config"

	"architecture_go/pkg/type/context"
	log "architecture_go/pkg/type/logger"
)

func init() {
	viper.SetConfigName(".env")
	viper.SetConfigType("dotenv")
	viper.AddConfigPath(".")
	viper.AutomaticEnv()

	viper.SetDefault("JAEGER_AGENT_HOST", "127.0.0.1")
	viper.SetDefault("JAEGER_AGENT_PORT", 6831)
}

func New(ctx context.Context) (io.Closer, error) {

	cfg := &config.Configuration{
		ServiceName: viper.GetString("SERVICE_NAME"),
		RPCMetrics:  true,
		Sampler: &config.SamplerConfig{
			Type:  "const",
			Param: 1,
		},
		Reporter: &config.ReporterConfig{
			LogSpans:           false,
			LocalAgentHostPort: fmt.Sprintf("%s:%d", viper.GetString("JAEGER_AGENT_HOST"), viper.GetUint32("JAEGER_AGENT_PORT")),
		},
	}

	tracer, closer, err := cfg.NewTracer(config.Logger(jaeger.StdLogger))
	if err != nil {
		return nil, log.ErrorWithContext(ctx, err)
	}

	opentracing.SetGlobalTracer(tracer)

	return closer, nil
}
