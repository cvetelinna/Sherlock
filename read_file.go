package github_api

import (
	"bufio"
	"os"
)

func Read(filename string) ([]string, error) {
	usernames := make([]string, 0)

	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if line != "" {
			usernames = append(usernames, line)
		}
	}

	if scanner.Err() != nil {
		return nil, err
	}

	return usernames, nil
}
