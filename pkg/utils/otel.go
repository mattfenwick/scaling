package utils

import (
	"net/http"
	"time"

	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
)

// TODO is this necessary?
//func OtelHttpClient() *http.Client {
//	return &http.Client{Transport: OtelTransport()}
//}

func OtelTransport() *otelhttp.Transport {
	return otelhttp.NewTransport(transport())
}

func transport() *http.Transport {
	return &http.Transport{
		Proxy:               http.ProxyFromEnvironment,
		MaxIdleConns:        100,
		MaxIdleConnsPerHost: 100,
		IdleConnTimeout:     90 * time.Second,
	}
}
