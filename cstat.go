package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"sort"
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

func poll(host string, opts options, result_chan chan [2]string) {
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

	result_chan <- [2]string{host, string(json)}
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

	result_chan := make(chan [2]string)
	results := make(map[string]string)
	var keys []string
	hosts := strings.Split(opts.hosts, ",")

	for {
		fmt.Println()

		for _, host := range hosts {
			go poll(host, opts, result_chan)
		}

		for i := 0; i < len(hosts); i++ {
			result := <-result_chan
			results[result[0]] = result[1]
			keys = append(keys, result[0])
		}

		sort.Strings(keys)

		for _, k := range keys {
			fmt.Printf("%s: %s\n", k, results[k])
		}

		time.Sleep(1 * time.Second)
		keys = []string{}
	}
}
