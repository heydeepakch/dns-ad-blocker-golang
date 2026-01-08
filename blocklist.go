package main

import (
	"bufio"
	"os"
	"strings"
)

func loadBlocklist(path string) (map[string]bool, error) {
	blocked := make(map[string]bool)

	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		// Ignore comments and empty lines
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		fields := strings.Fields(line)
		if len(fields) < 2 {
			continue
		}

		domain := strings.TrimSuffix(fields[1], ".")
		blocked[domain] = true
	}

	return blocked, scanner.Err()
}
