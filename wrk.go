package main

import "flag"

var (
	versionFlag = flag.Bool("version", false, "Print wrk version")
	helpFlag    = flag.Bool("help", false, "Print Help")
	http2       = flag.Bool("http2", true, "Use HTTP/2")
	goroutines  = flag.Int("c", 10, "Number of goroutines to use (concurrent connections)")
	duration    = flag.Int("d", 10, "Duration of test in second")
	cpus        = flag.Int("cpus", 0, "Numbers of cpus, i.e. GOMAXPROCS. 0 = system default")
	method      = flag.String("M", "GET", "HTTP method")
	host        = flag.String("host", "", "Host header")
	header      = flag.String("H", "", "Header to add for each request (You can define -H multiply)")
	body        = flag.String("b", "", "request body string or @filename")
	cert        = flag.String("cert", "", "CA certificate file to verify peer against (SSL/TLS)")
	key         = flag.String("key", "", "Private key file name(SSL/TLS)")
	caCert      = flag.String("ca", "", "CA file to verify peer against(SSL/TLS)")
)

func main() {
	flag.Parse()
}
