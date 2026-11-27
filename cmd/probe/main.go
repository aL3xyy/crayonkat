package main

import (
	"fmt"
	"net/http"
	"os"
	"time"
)

func main() {
	// Make sure the user typed: probe check <url>
	if len(os.Args) != 3 || os.Args[1] != "check" {
		fmt.Fprintln(os.Stderr, "usage: probe check <url>")
		os.Exit(1)
	}

	url := os.Args[2]

	// Give up after 10s so a server that never answers can't hang the probe
	client := &http.Client{Timeout: 10 * time.Second}

	// Send the request
	resp, err := client.Get(url)
	if err != nil {
		fmt.Fprintln(os.Stderr, "DOWN", err, url)
		os.Exit(1)
	}

	defer resp.Body.Close()

	// 200-399 means OK; 4xx/5xx means the site is broken
	if resp.StatusCode < 200 || resp.StatusCode >= 400 {
		fmt.Fprintln(os.Stderr, "DOWN", resp.StatusCode, url)
		os.Exit(1)
	}

	fmt.Println("UP", resp.StatusCode, url)
}
