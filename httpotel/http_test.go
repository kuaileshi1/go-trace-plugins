package httpotel

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"
)

func TestCallApi(t *testing.T) {
	// GET 请求
	res, err := CallApi(context.Background(), "http://localhost", "GET", nil, 10)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(string(res))

	// POST 请求
	reqBody := map[string]interface{}{
		"id":   1,
		"name": "test",
	}
	postData, err := json.Marshal(reqBody)
	if err != nil {
		t.Fatal(err)
	}
	res, err = CallApi(context.Background(), "http://localhost", "POST", bytes.NewReader(postData), 10, WithContentType("application/json"))
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(string(res))

	// 自定义Transport
	transport := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	res, err = CallApi(context.Background(), "http://localhost", "GET", nil, 10, WithTransport(transport))
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(string(res))
}
