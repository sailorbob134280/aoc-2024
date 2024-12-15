package main

import (
	_ "embed"
	"fmt"
	"log/slog"
	"os"
	"regexp"
	"strconv"
	"strings"
)

//go:embed data.txt
var data string

func main() {
	slog.SetDefault(slog.New(slog.NewTextHandler(os.Stdout, nil)))

	instructions := regexp.MustCompile(`(mul\(\d+,\d+\)|do\(\)|don't\(\))`).FindAllString(data, -1)

	if len(instructions) == 0 {
		slog.Error("No multiplication operations found")
		os.Exit(1)
	}

	var res int
	enabled := true
	for _, r := range instructions {
		if enabled && strings.Contains(r, "mul") {
			valStrings := regexp.MustCompile(`\d+`).FindAllString(r, 2)
			vals, err := parseIntList(valStrings)
			if err != nil {
				slog.Error("Error parsing int strings", "error", err.Error(), "vals", valStrings)
				os.Exit(1)
			}
			res += vals[0] * vals[1]
		} else if strings.Contains(r, "don't") {
			enabled = false
		} else if strings.Contains(r, "do") {
			enabled = true
		}
	}

	fmt.Printf("\nResult: %v\n", res)

	os.Exit(0)
}

func parseIntList(vals []string) ([]int, error) {
	res := make([]int, len(vals))
	for i, s := range vals {
		r, err := strconv.ParseInt(s, 10, 0)
		if err != nil {
			return nil, err
		}
		res[i] = int(r)
	}
	return res, nil
}
