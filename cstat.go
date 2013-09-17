package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"time"
)

type options struct {
	hosts    string
	username string
	password string
	port     uint
	metric   string
}

func poll(opts options) {
	for _, host := range strings.Split(opts.hosts, ",") {
		host_url := fmt.Sprintf("http://%s:%s@%s:%d", opts.username, opts.password, host, opts.port)
		payload, err := json.Marshal(map[string]string{"Name": opts.metric})
		resp, err := http.Post(host_url, "application/json", strings.NewReader(string(payload)))

		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		defer resp.Body.Close()

		json, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		fmt.Printf("%s: %s\n", host, string(json))
	}
}

func main() {
	var opts options

	flag.StringVar(&opts.hosts, "hosts", "", "comma-separated host list")
	flag.StringVar(&opts.username, "username", "cirgonus", "username to use for authentication")
	flag.StringVar(&opts.password, "password", "cirgonus", "password to use for authentication")
	flag.UintVar(&opts.port, "port", 8000, "port to use for connection")
	flag.StringVar(&opts.metric, "metric", "", "metric to fetch for polling")
	flag.Parse()

	if opts.hosts == "" {
		fmt.Println("Please provide a list of hosts to monitor")
		os.Exit(1)
	}

	if opts.metric == "" {
		fmt.Println("Please provide a metric to monitor")
		os.Exit(1)
	}

	if opts.username == "" || opts.password == "" {
		fmt.Println("Please supply a username and password for authentication to the cirgonus agent(s)")
		os.Exit(1)
	}

	for {
		fmt.Println()
		poll(opts)
		time.Sleep(1 * time.Second)
	}
}
