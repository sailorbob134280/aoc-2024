package main

import (
	"bufio"
	"bytes"
	_ "embed"
	"fmt"
	"log/slog"
	"os"
	"slices"
	"sort"
	"strconv"
	"strings"

	"github.com/sailorbob134280/aoc-2024/utils"
)

// Cheating with a little knowledge about the data, we never see a line longer than 8
const maxLineLength = 10

//go:embed data.txt
var data []byte

func main() {
	slog.SetDefault(slog.New(slog.NewTextHandler(os.Stdout, nil)))

	safeReports := 0
	safeWithDamper := 0

	s := bufio.NewScanner(bytes.NewBuffer(data))
	s.Split(bufio.ScanLines)
	for s.Scan() {
		line := s.Text()
		report, err := parseLine(line)
		if err != nil {
			slog.Error("Error parsing line", "line", line, "error", err.Error())
			os.Exit(1)
		}

		if validateReport(report) {
			safeReports++
		} else if problemDamper(report) {
			safeWithDamper++
		}
	}

	if err := s.Err(); err != nil {
		slog.Error("Error scanning data", "error", err.Error())
		os.Exit(1)
	}

	fmt.Printf("\nSafe reports: %v\n", safeReports)
	fmt.Printf("Safe reports with Problem Damper: %v\n", safeReports+safeWithDamper)
	os.Exit(0)
}

// problemDamper applies the Problem Damper to the report via a brute force method and returns true if the report is safe.
func problemDamper(r []int) bool {
	for i := range r {
		// Must clone to avoid clobbering the original report
		report := slices.Clone(r)
		report = slices.Delete(report, i, i+1)
		if validateReport(report) {
			return true
		}
	}
	return false
}

// validateReport returns true if the report is safe, and false otherwise.
func validateReport(r []int) bool {
	if !sort.IsSorted(sort.IntSlice(r)) && !sort.IsSorted(sort.Reverse(sort.IntSlice(r))) {
		return false
	}

	// Tiny edge case if there's only one element, can't be unsafe
	if len(r) == 1 {
		return true
	}

	// Intentionally start one after the start
	for i := 1; i < len(r); i++ {
		dist := utils.Distance(r[i-1], r[i])
		if dist < 1 || dist > 3 {
			return false
		}
	}

	return true
}

// parseLine expects a string that contains a space-separated list of ints of any length.
func parseLine(s string) ([]int, error) {
	w := bufio.NewScanner(strings.NewReader(s))
	w.Split(bufio.ScanWords)
	res := make([]int, 0, maxLineLength)
	for w.Scan() {
		val, err := strconv.ParseInt(w.Text(), 10, 0)
		if err != nil {
			return nil, err
		}
		res = append(res, int(val))
	}

	if err := w.Err(); err != nil {
		return nil, err
	}

	return res, nil
}
