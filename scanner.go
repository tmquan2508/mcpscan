package main

import (
	"fmt"
	"io"
	"os"
	"sync"
	"time"
)

type ScanResult struct {
	Port       int
	JSONString string
	Err        error
}

func RunScan(config ScannerConfig, domain string, writer io.Writer) error {
	fmt.Printf("[*] Starting host scan: %s (from port %d to %d)\n", domain, config.StartPort, config.EndPort)
	if !config.DebugMode {
		fmt.Printf("[*] Using %d concurrent scan threads.\n", config.NumWorkers)
	}

	portsToScan := make(chan int, config.NumWorkers)
	scanResults := make(chan ScanResult, config.NumWorkers)
	rateLimiter := time.NewTicker(time.Second / time.Duration(config.ScansPerSecond)).C

	var workersWg, loggerWg sync.WaitGroup

	loggerWg.Add(1)
	go logAndWriteResults(config, domain, scanResults, writer, &loggerWg)

	for i := 1; i <= config.NumWorkers; i++ {
		workersWg.Add(1)
		go worker(domain, config, portsToScan, scanResults, &workersWg, rateLimiter)
	}

	for port := config.StartPort; port <= config.EndPort; port++ {
		portsToScan <- port
	}
	close(portsToScan)

	workersWg.Wait()
	close(scanResults)
	loggerWg.Wait()

	return nil
}

func worker(host string, config ScannerConfig, ports <-chan int, results chan<- ScanResult, wg *sync.WaitGroup, rateLimiter <-chan time.Time) {
	defer wg.Done()
	for port := range ports {
		<-rateLimiter
		jsonString, err := getPingResult(host, uint16(port), config.ScanTimeout, config.DebugMode)
		results <- ScanResult{Port: port, JSONString: jsonString, Err: err}
	}
}

func logAndWriteResults(config ScannerConfig, host string, results <-chan ScanResult, writer io.Writer, wg *sync.WaitGroup) {
	defer wg.Done()

	if config.DebugMode {
		if config.DebugMode {
			fmt.Println("[!] DEBUG MODE IS ENABLED.")
		}
		for result := range results {
			address := fmt.Sprintf("%s:%d", host, result.Port)
			if result.Err != nil {
				fmt.Fprintf(os.Stderr, "%s -> No server found.\n", address)
			} else {
				fmt.Printf("%s -> Server found!\n", address)
				fmt.Fprintf(writer, "%s - %s\n", address, result.JSONString)
			}
		}
		return
	}

	startTime := time.Now()
	totalPorts := config.EndPort - config.StartPort + 1
	scannedCount := 0
	foundCount := 0

	for result := range results {
		scannedCount++

		if result.Err == nil {
			foundCount++
			address := fmt.Sprintf("%s:%d", host, result.Port)
			fmt.Fprintf(writer, "%s - %s\n", address, result.JSONString)
		}

		percentage := (float64(scannedCount) / float64(totalPorts)) * 100
		elapsedTime := time.Since(startTime).Seconds()

		fmt.Fprintf(os.Stderr,
			"\r[*] Scanning %s: %d/%d (%.2f%%) | Time: %.1fs | Found: %d ",
			host,
			scannedCount,
			totalPorts,
			percentage,
			elapsedTime,
			foundCount,
		)
	}
	fmt.Fprint(os.Stderr, "\n")
}