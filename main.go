package main

import (
	"flag"
	"fmt"
	"os"
)

func main() {

	hostFlag := flag.String("host", "", "The domain name to analyze")
	flag.Parse()

	var host string

	if *hostFlag != "" {
		host = *hostFlag
	}
	if host == "" && len(flag.Args()) > 0 {
		host = flag.Args()[0]
	}

	if host == "" {
		fmt.Println("Error: You must provide a host.")
		fmt.Println("Usage:")
		fmt.Println("  - Via argument: go run main.go google.com")
		fmt.Println("  - Via flag:      go run main.go -host= google.com")
		fmt.Println("  - Via flag (with space instead of equals): go run main.go -host google.com")
		os.Exit(1)
	}

	analyze(host)

}
