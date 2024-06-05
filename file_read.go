package main

import "os"

func readFile(filePath string) ([]byte, error) {
	fileContent, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}
	return fileContent, nil
}
