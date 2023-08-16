package tracer

import (
	"context"
	"fmt"
	"go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/exporters/zipkin"
	"go.opentelemetry.io/otel/sdk/resource"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.16.0"
)

const (
	ProviderJaeger   = "jaeger"
	ProviderZipkin   = "zipkin"
	ProviderOtlpHttp = "otlp-http"
)

// Conf
// @Description: 配置
type Conf struct {
	Id          string // 服务ID
	Name        string // 服务名称
	Environment string // 服务环境
	Endpoint    string // 服务地址
}

type Option struct {
	provider string  // provider类型
	sampling float64 // 采样率
}

type OptionFunc func(*Option)

// NewConf
// @Description: 实例化配置
// @Auth shigx 2023-08-10 16:11:00
// @param id 服务ID
// @param name 服务名称
// @param env 环境标识
// @param endpoint 服务地址
// @return *Conf
func NewConf(id, name, env, endpoint string) *Conf {
	return &Conf{
		Id:          id,
		Name:        name,
		Environment: env,
		Endpoint:    endpoint,
	}
}

// WithProvider
// @Description: 设置provider类型
// @Auth shigx 2023-08-10 16:25:37
// @param provider provider类型 jaeger zipkin ...
// @return OptionFunc
func (c *Conf) WithProvider(provider string) OptionFunc {
	return func(o *Option) {
		o.provider = provider
	}
}

// WithSampling
// @Description: 设置采样率
// @Auth shigx 2023-08-10 16:42:08
// @param sampling 采样率 0.0-1.0
// @return OptionFunc
func (c *Conf) WithSampling(sampling float64) OptionFunc {
	return func(o *Option) {
		o.sampling = sampling
	}
}

// TracerProvider
// @Description: 创建TracerProvider
// @Auth shigx 2023-08-10 16:48:07
// @param option
// @return *tracesdk.TracerProvider
// @return error
func (c *Conf) TracerProvider(option ...OptionFunc) (*tracesdk.TracerProvider, error) {
	o := &Option{
		provider: ProviderJaeger,
		sampling: 1.0,
	}
	for _, opt := range option {
		opt(o)
	}

	exp, err := c.createProvider(o.provider)
	if err != nil {
		return nil, err
	}

	tp := tracesdk.NewTracerProvider(
		tracesdk.WithSampler(tracesdk.ParentBased(tracesdk.TraceIDRatioBased(o.sampling))), // 基于父span的采样率设置
		tracesdk.WithBatcher(exp), // 始终确保在生产中批量处理
		tracesdk.WithResource( // 在资源中记录有关此应用程序的信息
			resource.NewWithAttributes(
				semconv.SchemaURL,
				semconv.ServiceNameKey.String(c.Name),
				semconv.DeploymentEnvironmentKey.String(c.Environment),
				semconv.ContainerIDKey.String(c.Id),
			)),
	)

	return tp, nil
}

// createProvider
// @Description: 创建provider
// @Auth shigx 2023-08-10 16:47:32
// @param provider
// @return tracesdk.SpanExporter
// @return error
func (c *Conf) createProvider(provider string) (tracesdk.SpanExporter, error) {
	switch provider {
	case ProviderJaeger:
		return jaeger.New(
			jaeger.WithCollectorEndpoint(jaeger.WithEndpoint(c.Endpoint)),
		)
	case ProviderZipkin:
		return zipkin.New(c.Endpoint)
	case ProviderOtlpHttp:
		ctx := context.Background()

		traceClientHttp := otlptracehttp.NewClient(
			otlptracehttp.WithEndpoint(c.Endpoint),
			otlptracehttp.WithInsecure())
		otlptracehttp.WithCompression(1)

		return otlptrace.New(ctx, traceClientHttp)
	default:
		return nil, fmt.Errorf("unknown exporter: %s", provider)
	}
}
