package parser

import "fmt"

func Parse(input string) (int, error) {
	if input == "" {
		return 0, fmt.Errorf("empty input")
	}

	return 1, nil
}
