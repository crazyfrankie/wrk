package client

import (
	"errors"
	"net/http"
	"net/url"
	"strings"
	"time"

	histogram "github.com/HdrHistogram/hdrhistogram-go"
)

const (
	USERAGENT = "wrk-go"
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

type RequestStats struct {
	TotalRespSize int64
	TotalDuration time.Duration
	NumRequests   int
	NumErrs       int
	ErrMap        map[string]int
	Histogram     *histogram.Histogram
}

func unwrap(err error) error {
	for errors.Unwrap(err) != nil {
		err = errors.Unwrap(err)
	}
	return err
}

func escapeUrlStr(in string) string {
	qm := strings.Index(in, "?")
	if qm != -1 {
		qry := in[qm+1:]
		qrys := strings.Split(qry, "&")
		var query string
		var qEscaped string
		first := true
		for _, q := range qrys {
			qSplit := strings.Split(q, "=")
			if len(qSplit) == 2 {
				qEscaped = qSplit[0] + "=" + url.QueryEscape(qSplit[1])
			} else {
				qEscaped = qSplit[0]
			}
			if first {
				first = false
			} else {
				query += "&"
			}
			query += qEscaped

		}
		return in[:qm] + "?" + query
	} else {
		return in
	}
}

func httpHeadersSize(headers http.Header) (result int64) {
	result = 0

	for k, v := range headers {
		result += int64(len(k) + len(": \r\n"))
		for _, s := range v {
			result += int64(len(s))
		}
	}

	result += int64(len("\r\n"))

	return result
}
