package pkg

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func ReadInputLines(linesCount int) ([]string, error) {
	var lines []string

	reader := bufio.NewReader(os.Stdin)
	for i := 0; i < linesCount; i++ {
		fmt.Print("> ")
		line, err := reader.ReadString('\n')
		line = strings.TrimSpace(line)
		if err != nil {
			return nil, err
		}
		lines = append(lines, line)
	}

	return lines, nil
}

func ReadLinesFromFile(filePath string) ([]string, error) {
	var lines []string

	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return lines, nil
}
