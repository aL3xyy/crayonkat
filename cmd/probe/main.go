package main

import (
	"fmt"
	"net/http"
	"os"
)

func main() {

	// Makes sure if the user typed: probe check <url>
	if len(os.Args) != 3 || os.Args[1] != "check" {
		fmt.Fprintln(os.Stderr, "usage: probe check <url>")
		os.Exit(1)
	}

	url := os.Args[2]

	// visit the website and checks the response
	resp, err := http.Get(url)
	if err != nil {
		fmt.Fprint(os.Stderr, "DOWN", url, err)
		os.Exit(1)
	}

	defer resp.Body.Close()

	fmt.Println("UP", url, resp.StatusCode)
}
