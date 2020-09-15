package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"net"

	"github.com/fatih/color"
)

var (
	verbose = false
)

func main() {
	flag.Usage = func() {
		_, _ = fmt.Fprintf(flag.CommandLine.Output(), "Prints AWS SDK metrics. Enable them with 'export AWS_CSM_ENABLED=true'\n")
		flag.PrintDefaults()
	}
	flag.BoolVar(&verbose, "v", false, "verbose mode, prints raw messages")
	flag.Parse()

	s, err := net.ResolveUDPAddr("udp4", ":31000")
	if err != nil {
		panic(err)
	}

	connection, err := net.ListenUDP("udp4", s)
	if err != nil {
		panic(err)
	}
	defer connection.Close()

	for {
		res, err := Read(connection)
		if err != nil {
			_ = fmt.Errorf("%v: %w", "unknown response", err)
			continue
		}
		if res.Type == "ApiCallAttempt" { // omit Verdicts
			msg := fmt.Sprintf("%v %-50.50v %v\n", res.HTTPStatusCode, fmt.Sprintf("%s:%s", res.Service, res.API), res.Fqdn)
			if res.HTTPStatusCode >= 400 {
				color.Red(msg)
			} else {
				color.Green(msg)
			}
		}
	}
}

type Response struct {
	// Attempt
	Version        int    `json:"Version"`
	ClientID       string `json:"ClientId"`
	Type           string `json:"Type"`
	Service        string `json:"Service"`
	API            string `json:"Api"`
	Timestamp      int64  `json:"Timestamp"`
	AttemptLatency int    `json:"AttemptLatency"`
	Fqdn           string `json:"Fqdn"`
	UserAgent      string `json:"UserAgent"`
	AccessKey      string `json:"AccessKey"`
	Region         string `json:"Region"`
	SessionToken   string `json:"SessionToken"`
	HTTPStatusCode int    `json:"HttpStatusCode"`
	XAmznRequestID string `json:"XAmznRequestId"`

	// Verdict
	AttemptCount        int `json:"AttemptCount"`
	FinalHTTPStatusCode int `json:"FinalHttpStatusCode"`
	Latency             int `json:"Latency"`
	MaxRetriesExceeded  int `json:"MaxRetriesExceeded"`
}

func Read(conn *net.UDPConn) (*Response, error) {
	b := make([]byte, 1024)
	oob := make([]byte, 40)

	//nolint:dogsled
	_, _, _, _, err := conn.ReadMsgUDP(b, oob)
	if err != nil {
		return nil, err
	}

	// Remove NULL characters
	b = bytes.Trim(b, "\x00")

	if verbose {
		fmt.Println(string(b))
	}

	ret := &Response{}
	if err := json.Unmarshal(b, ret); err != nil {
		return nil, err
	}

	return ret, err
}
