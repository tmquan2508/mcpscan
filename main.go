package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"time"
)

func main() {
	config, hostInput, outputFilePath, err := ParseFlags(os.Args[1:])
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n\n", err)
		PrintUsage()
		os.Exit(1)
	}

	if config.ShowHelp {
		PrintUsage()
		os.Exit(0)
	}

	hosts, err := ResolveHosts(hostInput)
	if err != nil {
		log.Fatalf("Error resolving hosts: %v", err)
	}

	file, err := os.Create(outputFilePath)
	if err != nil {
		log.Fatalf("Error creating output file: %v", err)
	}
	defer file.Close()
	writer := bufio.NewWriter(file)
	defer writer.Flush()

	totalStartTime := time.Now()
	fmt.Printf("[*] Starting scan for %d host(s). Results will be saved to %s\n", len(hosts), outputFilePath)

	for i, host := range hosts {
		fmt.Printf("\n--- Scanning host %d/%d: %s ---\n", i+1, len(hosts), host)
		err := RunScan(config, host, writer)
		if err != nil {
			log.Printf("Error scanning host %s: %v", host, err)
		}
	}

	totalTime := time.Since(totalStartTime).Seconds()
	fmt.Printf("\n\n[*] ENTIRE SCAN PROCESS COMPLETE! (Total time: %.2fs)\n", totalTime)
	fmt.Printf("[+] All found servers saved to file: %s\n", outputFilePath)
}