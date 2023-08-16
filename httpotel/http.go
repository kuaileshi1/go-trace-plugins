package httpotel

import (
	"context"
	"io"
	"net/http"
	"time"

	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
)

type ClientOption struct {
	Transport   *http.Transport
	ContentType string
	Headers     map[string]string
}

type ClientOptionFunc func(*ClientOption)

// WithTransport
// @Description: 设置client transport
// @Auth shigx 2023-08-11 09:26:47
// @param transport
// @return ClientOptionFunc
func WithTransport(transport *http.Transport) ClientOptionFunc {
	return func(option *ClientOption) {
		option.Transport = transport
	}
}

// WithContentType
// @Description: 设置content-type
// @Auth shigx 2023-08-11 09:26:57
// @param contentType
// @return ClientOptionFunc
func WithContentType(contentType string) ClientOptionFunc {
	return func(option *ClientOption) {
		option.ContentType = contentType
	}
}

// WithHeaders
// @Description: 设置headers
// @Auth shigx 2023-08-11 15:42:20
// @param headers
// @return ClientOptionFunc
func WithHeaders(headers map[string]string) ClientOptionFunc {
	return func(option *ClientOption) {
		option.Headers = headers
	}
}

// CallApi
// @Description: http请求
// @Auth shigx 2023-08-11 09:27:14
// @param ctx 上下文
// @param url 请求地址
// @param method 请求方法
// @param reqBody 请求参数
// @param timeOut 超时时间
// @param option 请求配置
// @return []byte 响应结果
// @return error 错误信息
func CallApi(ctx context.Context, url string, method string, reqBody io.Reader, timeOut time.Duration, option ...ClientOptionFunc) ([]byte, error) {
	clientOption := &ClientOption{}
	for _, o := range option {
		o(clientOption)
	}

	client := http.Client{Timeout: timeOut * time.Second, Transport: otelhttp.NewTransport(http.DefaultTransport)}
	if clientOption.Transport != nil {
		client.Transport = otelhttp.NewTransport(clientOption.Transport)
	}
	req, err := http.NewRequestWithContext(ctx, method, url, reqBody)
	if err != nil {
		return nil, err
	}
	if clientOption.ContentType != "" {
		req.Header.Set("Content-Type", clientOption.ContentType)
	}
	// 设置headers
	if len(clientOption.Headers) > 0 {
		for k, v := range clientOption.Headers {
			req.Header.Set(k, v)
		}
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	resBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return resBody, nil
}
