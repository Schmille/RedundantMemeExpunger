package main

import (
	"errors"
	"math"
	"strconv"
	"strings"
)

const (
	GIGABYTE = "GB"
	MEGABYTE = "MB"
	KILOBYTE = "KB"
	BYTE = "B"
	EMPTY = ""
)

func ParseSizeString(input string) (int64, error) {
	var multiplier float64
	input = strings.ReplaceAll(input, ",", ".")

	if strings.HasSuffix(input, GIGABYTE) {
		multiplier = 1024 * 1024 * 1024
		input = strings.ReplaceAll(input, GIGABYTE, EMPTY)
	} else if strings.HasSuffix(input, MEGABYTE) {
		multiplier = 1024 * 1024
		input = strings.ReplaceAll(input, MEGABYTE, EMPTY)
	} else if strings.HasSuffix(input, KILOBYTE) {
		multiplier = 1024
		input = strings.ReplaceAll(input, KILOBYTE, EMPTY)
	} else if strings.HasSuffix(input, BYTE) {
		multiplier = 1
		input = strings.ReplaceAll(input, BYTE, EMPTY)
	} else {
		return -1, errors.New("unrecognised size pattern")
	}

	input = strings.Trim(input, "%q%s")

	size, err := strconv.ParseFloat(input, 64)

	if err != nil {
		return -1, err
	}

	if size < 0 {
		return -1, errors.New("size must not be negative")
	}

	size = size * multiplier
	size = math.RoundToEven(size)

	return int64(size), nil
}

type Set struct {
	content map[string]bool
}

func NewSet() *Set {
	return &Set{
		content: make(map[string]bool),
	}
}

func (s *Set) Add(entry string) {
	s.content[entry] = true
}

func (s *Set) Remove(entry string) {
	delete(s.content, entry)
}

func (s *Set) Contains(entry string) bool {
	return s.content[entry]
}
