package main

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

type ScannerConfig struct {
	NumWorkers      int
	ScansPerSecond  int
	ScanTimeout     time.Duration
	StartPort       int
	EndPort         int
	DebugMode       bool
	ShowHelp        bool
}

func findFlagValue(args []string, long, short string) (string, bool) {
	longFlag := "--" + long
	shortFlag := "-" + short
	for i, arg := range args {
		if arg == longFlag || (short != "" && arg == shortFlag) {
			if i+1 < len(args) && !strings.HasPrefix(args[i+1], "-") {
				return args[i+1], true
			}
		}
	}
	return "", false
}

func hasBoolFlag(args []string, long, short string) bool {
	longFlag := "--" + long
	shortFlag := "-" + short
	for _, arg := range args {
		if arg == longFlag || (short != "" && arg == shortFlag) {
			return true
		}
	}
	return false
}

func ParseFlags(args []string) (ScannerConfig, string, string, error) {
	config := ScannerConfig{
		NumWorkers:      150,
		ScansPerSecond:  200,
		ScanTimeout:     5 * time.Second,
		StartPort:       25000,
		EndPort:         30000,
		DebugMode:       false,
		ShowHelp:        false,
	}

	if hasBoolFlag(args, "help", "?") {
		config.ShowHelp = true
		return config, "", "", nil
	}
	
	hostInput, hostFound := findFlagValue(args, "host", "h")
	if !hostFound {
		return config, "", "", errors.New("required flag --host (-h) is missing")
	}

	outputFilePath, outputFound := findFlagValue(args, "output", "o")
	if !outputFound {
		baseName := strings.TrimSuffix(hostInput, filepath.Ext(hostInput))
		outputFilePath = baseName + "_results.txt"
		fmt.Printf("[*] No output file specified. Defaulting to: %s\n", outputFilePath)
	}
	
	config.DebugMode = hasBoolFlag(args, "debug", "d")

	if val, found := findFlagValue(args, "workers", "w"); found {
		if i, err := strconv.Atoi(val); err == nil { config.NumWorkers = i }
	}
	if val, found := findFlagValue(args, "rate", "r"); found {
		if i, err := strconv.Atoi(val); err == nil { config.ScansPerSecond = i }
	}
	if val, found := findFlagValue(args, "timeout", "t"); found {
		if i, err := strconv.Atoi(val); err == nil { config.ScanTimeout = time.Duration(i) * time.Second }
	}
	if val, found := findFlagValue(args, "start-port", "s"); found {
		if i, err := strconv.Atoi(val); err == nil { config.StartPort = i }
	}
	if val, found := findFlagValue(args, "end-port", "e"); found {
		if i, err := strconv.Atoi(val); err == nil { config.EndPort = i }
	}

	return config, hostInput, outputFilePath, nil
}

func PrintUsage() {
	fmt.Fprintf(os.Stderr, "A fast, concurrent Minecraft server port scanner.\n\n")
	fmt.Fprintf(os.Stderr, "Usage:\n  %s -h < host | hosts.txt > [flags]\n\n", os.Args[0])
	fmt.Fprintf(os.Stderr, "Example:\n  %s --host hosts.txt -s 25500 -e 25600 -o found.txt\n\n", os.Args[0])
	fmt.Fprintf(os.Stderr, "Flags:\n")
	fmt.Fprintf(os.Stderr, "  -h, --host <string>       (Required) A single domain or a path to a .txt file with hosts.\n")
	fmt.Fprintf(os.Stderr, "  -o, --output <string>     Path to save the results. (default: <host>_results.txt)\n")
	fmt.Fprintf(os.Stderr, "  -w, --workers <int>       Number of concurrent scan threads (default 150)\n")
	fmt.Fprintf(os.Stderr, "  -r, --rate <int>          Max scans per second (default 200)\n")
	fmt.Fprintf(os.Stderr, "  -t, --timeout <int>       Connection timeout in seconds (default 5)\n")
	fmt.Fprintf(os.Stderr, "  -s, --start-port <int>    Port to start scanning from (default 25000)\n")
	fmt.Fprintf(os.Stderr, "  -e, --end-port <int>      Port to end scanning at (default 30000)\n")
	fmt.Fprintf(os.Stderr, "  -d, --debug               Enable detailed debug logging.\n")
	fmt.Fprintf(os.Stderr, "      --help, -?            Show this help message.\n")
}