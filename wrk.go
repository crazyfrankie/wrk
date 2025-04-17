package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"runtime"
	"strings"
	"syscall"
	"time"

	histogram "github.com/HdrHistogram/hdrhistogram-go"

	"github.com/crazyfrankie/wrk/client"
)

const (
	Version = "0.1.1"
)

var (
	versionFlag bool
	helpFlag    bool
	http2       bool
	skipVerify  bool
	goroutines  int
	duration    int
	cpus        int
	timeout     int
	method      string
	host        string
	body        string
	cert        string
	key         string
	caCert      string
	headerFlags HeaderList
)

func initFlag() {
	flag.BoolVar(&versionFlag, "v", false, "Print wrk version")
	flag.BoolVar(&helpFlag, "help", false, "Print Help")
	flag.BoolVar(&http2, "http2", true, "Use HTTP/2")
	flag.IntVar(&goroutines, "c", 10, "Number of goroutines to use (concurrent connections)")
	flag.IntVar(&duration, "d", 10, "Duration of test in second")
	flag.IntVar(&cpus, "cpus", 0, "Numbers of cpus, i.e. GOMAXPROCS. 0 = system default")
	flag.IntVar(&timeout, "T", 1000, "Socket/request timeout in ms")
	flag.StringVar(&method, "M", "GET", "HTTP method")
	flag.StringVar(&host, "host", "", "Host header")
	flag.Var(&headerFlags, "H", "Header to add for each request (You can define -H multiply)")
	flag.StringVar(&body, "b", "", "request body string or @filename")
	flag.StringVar(&cert, "cert", "", "CA certificate file to verify peer against (SSL/TLS)")
	flag.StringVar(&key, "key", "", "Private key file name(SSL/TLS)")
	flag.StringVar(&caCert, "ca", "", "CA file to verify peer against(SSL/TLS)")
}

func main() {
	initFlag()

	flag.Parse()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT)

	header := make(map[string]string)
	for _, hdr := range headerFlags {
		hp := strings.SplitN(hdr, ":", 2)
		header[hp[0]] = hp[1]
	}

	url := flag.Arg(0)

	if versionFlag {
		fmt.Println("Version", Version)
		return
	} else if helpFlag || len(url) == 0 {
		printDefaults()
		return
	}

	if cpus > 0 {
		runtime.GOMAXPROCS(cpus)
	}

	fmt.Printf("Running %vs test @ %v\n  %v goroutines running concurrently\n", duration, url, goroutines)

	// if user input a file for client body
	if len(body) > 0 && body[0] == '@' {
		filename := body[1:]
		data, err := os.ReadFile(filename)
		if err != nil {
			fmt.Println(fmt.Errorf("could not read file %q: %v", filename, err))
			os.Exit(1)
		}
		body = string(data)
	}

	req := client.NewRequest(&client.Config{
		Duration:   duration,
		Goroutines: goroutines,
		URL:        url,
		Body:       body,
		Method:     method,
		Host:       host,
		Header:     header,
		Timeout:    timeout,
		SkipVerify: skipVerify,
		ClientCert: cert,
		ClientKey:  key,
		CACert:     caCert,
		HTTP2:      http2,
	})

	start := time.Now()

	for i := 0; i < goroutines; i++ {
		go req.Run()
	}

	responders := 0
	stats := client.RequestStats{ErrMap: make(map[string]int), Histogram: histogram.New(1, int64(duration*1000000), 4)}

	for responders < goroutines {
		select {
		case <-quit:
			req.Stop()
			fmt.Printf("stopping...\n")
		case res := <-req.Stats:
			stats.NumErrs += res.NumErrs
			stats.NumRequests += res.NumRequests
			stats.TotalRespSize += res.TotalRespSize
			stats.TotalDuration += res.TotalDuration
			responders++
			for k, v := range res.ErrMap {
				stats.ErrMap[k] = v
			}
			stats.Histogram.Merge(res.Histogram)
		}
	}

	duration := time.Now().Sub(start)

	if stats.NumRequests == 0 {
		fmt.Println("Error: No statistics collected / no requests found")
		fmt.Printf("Number of Errors:\t%v\n", stats.NumErrs)
		if stats.NumErrs > 0 {
			fmt.Printf("Error Counts:\t\t%v\n", convertErr(stats.ErrMap))
		}
		return
	}

	avgGoroutineDur := stats.TotalDuration / time.Duration(responders)

	reqRate := float64(stats.NumRequests) / avgGoroutineDur.Seconds()
	bytesRate := float64(stats.TotalRespSize) / avgGoroutineDur.Seconds()

	overallReqRate := float64(stats.NumRequests) / duration.Seconds()
	overallBytesRate := float64(stats.TotalRespSize) / duration.Seconds()

	fmt.Printf("%v request in %v, %v read\n", stats.NumRequests, avgGoroutineDur, Bytes{Size: float64(stats.TotalRespSize)})
	fmt.Printf("Request/sec:\t\t%.2f\nTransfer/sec:\t\t%v\n", reqRate, Bytes{Size: bytesRate})
	fmt.Printf("Overall Requests/sec:\t%.2f\nOverall Transfer/sec:\t%v\n", overallReqRate, Bytes{Size: overallBytesRate})
	fmt.Printf("Fastest Request:\t%v\n", convertDuration(stats.Histogram.Min()))
	fmt.Printf("Avg Req Time:\t\t%v\n", convertDuration(int64(stats.Histogram.Mean())))
	fmt.Printf("Slowest Request:\t%v\n", convertDuration(stats.Histogram.Max()))
	fmt.Printf("Number of Errors:\t%v\n", stats.NumErrs)
	if stats.NumErrs > 0 {
		fmt.Printf("Error Counts:\t\t%v\n", convertErr(stats.ErrMap))
	}
	fmt.Printf("10%%:\t\t\t%v\n", convertDuration(stats.Histogram.ValueAtPercentile(.10)))
	fmt.Printf("50%%:\t\t\t%v\n", convertDuration(stats.Histogram.ValueAtPercentile(.50)))
	fmt.Printf("75%%:\t\t\t%v\n", convertDuration(stats.Histogram.ValueAtPercentile(.75)))
	fmt.Printf("99%%:\t\t\t%v\n", convertDuration(stats.Histogram.ValueAtPercentile(.99)))
	fmt.Printf("99.9%%:\t\t\t%v\n", convertDuration(stats.Histogram.ValueAtPercentile(.999)))
	fmt.Printf("99.9999%%:\t\t%v\n", convertDuration(stats.Histogram.ValueAtPercentile(.999999)))
	fmt.Printf("99.99999%%:\t\t%v\n", convertDuration(stats.Histogram.ValueAtPercentile(.9999999)))
	fmt.Printf("stddev:\t\t\t%v\n", convertDuration(int64(stats.Histogram.StdDev())))
}

// printDefaults is used when users enter -help or do not enter the url correctly
func printDefaults() {
	fmt.Println("Usage: wrk-go <option> <url>")
	fmt.Println("Options:")
	flag.VisitAll(func(f *flag.Flag) {
		fmt.Println("\t-"+f.Name, "\t", f.Usage, "(Default ", f.DefValue+")")
	})
}

func convertErr(err map[string]int) string {
	s := make([]string, 0, len(err))
	for k, v := range err {
		s = append(s, fmt.Sprint(k, "=", v))
	}
	return strings.Join(s, ",")
}

func convertDuration(usecs int64) time.Duration {
	return time.Duration(usecs * 1000)
}
