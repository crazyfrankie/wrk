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
369784 request in 10.023184433s, 38.44MB read
Request/sec:            36892.87
Transfer/sec:           3.84MB
Overall Requests/sec:   31587.22
Overall Transfer/sec:   3.28MB
Fastest Request:        90µs
Avg Req Time:           55.511ms
Slowest Request:        931.615ms
Number of Errors:       0
10%:                    133µs
50%:                    181µs
75%:                    201µs
99%:                    218µs
99.9%:                  218µs
99.9999%:               218µs
99.99999%:              218µs
stddev:                 52.955ms
```