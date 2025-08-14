package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func ResolveHosts(input string) ([]string, error) {
	if strings.HasSuffix(strings.ToLower(input), ".txt") {
		fmt.Printf("[*] Reading hosts from file: %s\n", input)
		file, err := os.Open(input)
		if err != nil {
			return nil, fmt.Errorf("could not open host file: %w", err)
		}
		defer file.Close()

		var hosts []string
		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			line := strings.TrimSpace(scanner.Text())
			if line != "" && !strings.HasPrefix(line, "#") {
				hosts = append(hosts, line)
			}
		}
		if err := scanner.Err(); err != nil {
			return nil, fmt.Errorf("error reading host file: %w", err)
		}
		return hosts, nil
	}

	fmt.Printf("[*] Target is a single host: %s\n", input)
	return []string{input}, nil
}