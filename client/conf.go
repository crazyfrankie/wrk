package client

import (
	"time"

	histogram "github.com/HdrHistogram/hdrhistogram-go"
)

const (
	USERAGENT = "wkr-go"
)

type Config struct {
	Duration    int
	Goroutines  int
	Timeout     int
	URL         string
	Body        string
	Host        string
	Method      string
	Header      map[string]string
	SkipVerify  bool
	HTTP2       bool
	Interrupted int32
	ClientCert  string
	ClientKey   string
	CACert      string
}

type requestStats struct {
	TotalRespSize int64
	TotalDuration time.Duration
	NumRequests   int
	NumErrs       int
	ErrMap        map[string]int
	Histogram     *histogram.Histogram
}
