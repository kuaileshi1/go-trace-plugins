package tracer

import (
	"context"
	"go.opentelemetry.io/otel"
	"testing"
)

func TestTracerProvider(t *testing.T) {
	// 设置jaeger provider
	c := NewConf("MacBook-Pro.local", "wp-api", "dev", "http://localhost:14268/api/traces")
	tp, err := c.TracerProvider()

	// 设置 zipkin provider和采样率
	c = NewConf("MacBook-Pro.local", "wp-api", "dev", "http://localhost:9411")
	tp, err = c.TracerProvider(c.WithProvider("zipkin"), c.WithSampling(0.8))

	// 设置otlp-http provider
	c = NewConf("MacBook-Pro.local", "wp-api", "dev", "localhost:4318")
	tp, err = c.TracerProvider(c.WithProvider("otlp-http"), c.WithSampling(0.5))

	if err != nil {
		t.Fatal(err)
	}
	if tp == nil {
		t.Fatal("tp is nil")
	}
	otel.SetTracerProvider(tp)

	defer func() {
		if err = tp.Shutdown(context.Background()); err != nil {
			t.Fatal(err)
		}
	}()
}
