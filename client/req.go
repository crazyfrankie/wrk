package client

import (
	"bytes"
	"crypto/tls"
	"crypto/x509"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"sync/atomic"
	"time"

	histogram "github.com/HdrHistogram/hdrhistogram-go"
)

type Request struct {
	config *Config
	Stats  chan *RequestStats
}

func NewRequest(cfg *Config) *Request {
	stats := make(chan *RequestStats, cfg.Goroutines)

	return &Request{Stats: stats, config: cfg}
}

func newClient(req *Request) (*http.Client, error) {
	tlsConfig := &tls.Config{
		InsecureSkipVerify: req.config.SkipVerify,
	}

	if req.config.ClientCert != "" || req.config.ClientKey != "" || req.config.CACert != "" {
		if req.config.ClientCert == "" || req.config.ClientKey == "" {
			return nil, fmt.Errorf("both client certificate and key must be provided")
		}

		cert, err := tls.LoadX509KeyPair(req.config.ClientCert, req.config.ClientKey)
		if err != nil {
			return nil, fmt.Errorf("failed to load client certificate: %v", err)
		}
		tlsConfig.Certificates = []tls.Certificate{cert}

		if req.config.CACert != "" {
			caCert, err := os.ReadFile(req.config.CACert)
			if err != nil {
				return nil, fmt.Errorf("failed to read CA certificate: %v", err)
			}
			caCertPool := x509.NewCertPool()
			if !caCertPool.AppendCertsFromPEM(caCert) {
				return nil, fmt.Errorf("failed to parse CA certificate")
			}
			tlsConfig.RootCAs = caCertPool
		}
	}

	transport := &http.Transport{
		TLSClientConfig:       tlsConfig,
		ResponseHeaderTimeout: time.Duration(req.config.Timeout) * time.Millisecond,
		ForceAttemptHTTP2:     req.config.HTTP2,
	}

	//if req.config.HTTP2 && tlsConfig.NextProtos == nil {
	//	tlsConfig.NextProtos = []string{"h2", "http/1.1"}
	//}

	client := &http.Client{
		Transport: transport,
	}

	return client, nil
}

func (r *Request) Run() {
	res := &RequestStats{ErrMap: make(map[string]int), Histogram: histogram.New(1, int64(r.config.Duration*1000000), 4)}
	start := time.Now()

	client, err := newClient(r)
	if err != nil {
		log.Fatal(err)
	}

	for time.Since(start).Seconds() <= float64(r.config.Duration) && atomic.LoadInt32(&r.config.Interrupted) == 0 {
		respSize, reqDur, err := r.Do(client)
		if err != nil {
			res.ErrMap[unwrap(err).Error()] += 1
			res.NumErrs++
		} else if respSize > 0 {
			res.TotalRespSize += int64(respSize)
			res.TotalDuration += reqDur
			res.NumRequests++
			res.Histogram.RecordValue(reqDur.Microseconds())
		} else {
			res.NumErrs++
		}
	}

	r.Stats <- res
}

// Do send a http request
func (r *Request) Do(client *http.Client) (int, time.Duration, error) {
	var respSize int
	var duration time.Duration

	url := escapeUrlStr(r.config.URL)

	var buf io.Reader
	if len(r.config.Body) > 0 {
		buf = bytes.NewBufferString(r.config.Body)
	}

	req, err := http.NewRequest(r.config.Method, url, buf)
	if err != nil {
		return 0, 0, err
	}

	for k, v := range r.config.Header {
		req.Header.Add(k, v)
	}
	req.Header.Add("User-Agent", USERAGENT)

	if r.config.Host != "" {
		req.Host = r.config.Host
	}
	start := time.Now()

	resp, err := client.Do(req)
	if err != nil {
		return 0, 0, err
	}

	if resp == nil {
		return 0, 0, errors.New("empty response")
	}
	defer func() {
		if resp != nil && resp.Body != nil {
			resp.Body.Close()
		}
	}()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return 0, 0, err
	}
	if resp.StatusCode/100 == 2 {
		duration = time.Since(start)
		respSize = len(body) + int(httpHeadersSize(resp.Header))
	} else if resp.StatusCode == http.StatusMovedPermanently || resp.StatusCode == http.StatusPermanentRedirect {
		duration = time.Since(start)
		respSize = int(resp.ContentLength) + int(httpHeadersSize(resp.Header))
	} else {
		return 0, 0, errors.New(fmt.Sprint("received status code ", resp.StatusCode))
	}

	return respSize, duration, nil
}

func (r *Request) Stop() {
	atomic.StoreInt32(&r.config.Interrupted, 1)
}
