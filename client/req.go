package client

import "net/http"

type Request struct {
	config *Config
	stats  chan *requestStats
}

func NewRequest(cfg *Config) *Request {
	stats := make(chan *requestStats, cfg.Goroutines)

	return &Request{stats: stats, config: cfg}
}

func newClient(req *Request) (*http.Client, error) {
	return nil, nil
}

func (r *Request) Run() {

}
