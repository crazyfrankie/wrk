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

	"github.com/crazyfrankie/wrk/client"
)

const (
	Version = "0.1.0"
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
	flag.BoolVar(&versionFlag, "version", false, "Print wrk version")
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

	// TODO
}

// printDefaults is used when users enter -help or do not enter the url correctly
func printDefaults() {
	fmt.Println("Usage: wrk-go <option> <url>")
	fmt.Println("Options:")
	flag.VisitAll(func(f *flag.Flag) {
		fmt.Println("\t-"+f.Name, "\t", f.Usage, "(Default ", f.DefValue+")")
	})
}
