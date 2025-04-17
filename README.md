# wrk-go

wrk-go is an HTTP benchmark testing tool developed based on the Go language. In concurrent scenarios, it uses goroutines for asynchronous I/O and can achieve a good concurrent volume in a short time

# Install
```
go install github.com/crazyfrankie/wrk@latest
```

Command line parameters(wrk -help)
```
Usage: wrk-go <option> <url>
Options:
        -H       Header to add for each request (You can define -H multiply) (Default  )
        -M       HTTP method (Default  GET)
        -T       Socket/request timeout in ms (Default  1000)
        -b       request body string or @filename (Default  )
        -c       Number of goroutines to use (concurrent connections) (Default  10)
        -ca      CA file to verify peer against(SSL/TLS) (Default  )
        -cert    CA certificate file to verify peer against (SSL/TLS) (Default  )
        -cpus    Numbers of cpus, i.e. GOMAXPROCS. 0 = system default (Default  0)
        -d       Duration of test in second (Default  10)
        -help    Print Help (Default  false)
        -host    Host header (Default  )
        -http2   Use HTTP/2 (Default  true)
        -key     Private key file name(SSL/TLS) (Default  )
        -v       Print wrk version (Default  false)
```

# Usage
```
wrk -c 2048 -d 10 http://localhost:8082/hello
```

This runs a benchmark for 10 seconds, using 2048 go routines (connections)

The Output will like this:
```
Running 10s test @ http://localhost:8082/hello
  2048 goroutines running concurrently
428253 request in 42733.48199446375,44.52MB read
Request/sec:            42733.48
Transfer/sec:           4.44MB
Overall Requests/sec:   37116.50
Overall Transfer/sec:   3.86MB
Fastest Request:        81µs
Avg Req Time:           47.924ms
Slowest Request:        520.127ms
Number of Errors:       0
10%:                    129µs
50%:                    169µs
75%:                    188µs
99%:                    203µs
99.9%:                  204µs
99.9999%:               204µs
99.99999%:              204µs
stddev:                 43.175ms
```