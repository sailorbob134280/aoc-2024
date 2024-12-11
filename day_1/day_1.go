package main

import (
	"bufio"
	_ "embed"
	"errors"
	"flag"
	"fmt"
	"io"
	"log/slog"
	"os"
	"slices"
	"sort"
	"strconv"
	"strings"
)

// The file is in a file, with each list separated by a space on separate line. E.g.
// Left_1 Right_1
// Left_2 Right_2
// etc

//go:embed data.txt
var data []byte

var ErrBadFile = errors.New("input file is badly formatted")

type Locations struct {
	Left  []int
	Right []int
}

// NewLocations creates a new Locations struct from a reader.
func NewLocations(data io.Reader) (*Locations, error) {
	l := Locations{}
	lines := bufio.NewScanner(data)
	lines.Split(bufio.ScanLines)

	for lines.Scan() {
		line := lines.Text()
		words := bufio.NewScanner(strings.NewReader(line))
		words.Split(bufio.ScanWords)

		// Left one first
		if !words.Scan() {
			return nil, ErrBadFile
		}

		var err error
		l.Left, err = appendFromString(l.Left, words.Text())
		if err != nil {
			return nil, err
		}

		// Now the right
		if !words.Scan() {
			return nil, ErrBadFile
		}

		l.Right, err = appendFromString(l.Right, words.Text())
		if err != nil {
			return nil, err
		}
	}

	return &l, nil
}

func main() {
	slog.SetDefault(slog.New(slog.NewTextHandler(os.Stdout, nil)))

	filename := flag.String("datafile", "data.txt", "The input data file to parse")
	flag.Parse()

	if filename == nil || *filename == "" {
		slog.Error("Filename is invalid, exiting")
		os.Exit(1)
	}

	// Ingest each list
	f, err := os.Open(*filename)
	if err != nil {
		slog.Error("Failed to open file", "err", err)
		os.Exit(1)
	}

	l, err := NewLocations(f)
	if err != nil {
		slog.Error("Failed to parse input file", "err", err)
		os.Exit(1)
	}

	// Sort the lists
	sort.Ints(l.Left)
	sort.Ints(l.Right)

	// For each, compare the diference and calculate the similatity
	// We assume the lists are the same length because otherwise the file is malformed
	var res, sim int
	for i := range l.Left {
		res += distance(l.Left[i], l.Right[i])
		sim += l.Left[i] * numOccurances(l.Right, l.Left[i])
	}

	fmt.Printf("\nTotal Distance: %v", res)
	fmt.Printf("\nTotal Similarity: %v\n", sim)
}

func numOccurances(arr []int, n int) int {
	i, found := slices.BinarySearch(arr, n)
	if !found {
		return 0
	}

	count := 0
	for arr[i] == n {
		count++
		i++
	}

	return count
}

func appendFromString(arr []int, s string) ([]int, error) {
	n, err := strconv.Atoi(s)
	if err != nil {
		return arr, ErrBadFile
	}
	arr = append(arr, n)
	return arr, nil
}

// Go inexplicably does not have an abs function for ints, and I have
// been bitten by float casting far too many times to do it that way,
// so we write our own.
func distance(a, b int) int {
	if a < b {
		return b - a
	}
	return a - b
}
