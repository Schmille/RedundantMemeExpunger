package main

import "testing"

func TestSameBytes(t *testing.T) {
	var searcher FileSearcher = GetTestSearcher()

	checks := sameBytes(&searcher)

	if checks[0] != "/mnt/g/test" {
		t.Fail()
	}
}

func TestProcessChecks(t *testing.T) {
	var searcher FileSearcher = GetTestSearcher()

	checks, err := processChecks(&searcher)

	if checks[0] != "/mnt/g/test" {
		t.Fail()
	}

	if err != nil {
		t.Fail()
	}
}

func TestConvertBytes(t *testing.T) {
	result, err := ParseSizeString("1GB")

	if err != nil {
		t.Fail()
	}

	if result != 1073741824 {
		t.Fail()
	}
}

func TestConvertBytesFloat(t *testing.T) {
	result, err := ParseSizeString("2.5MB")

	if err != nil {
		t.Fail()
	}

	if result != 2621440 {
		t.Fail()
	}
}

func TestConvertBytesNegative(t *testing.T) {
	_, err := ParseSizeString("-1B")

	if err == nil {
		t.Fail()
	}
}

func TestConvertBytesFloatComma(t *testing.T) {
	result, err := ParseSizeString("2,5MB")

	if err != nil {
		t.Fail()
	}

	if result != 2621440 {
		t.Fail()
	}
}

func TestConvertBytesShouldFail(t *testing.T) {
	_, err := ParseSizeString("fail")

	if err == nil {
		t.Fail()
	}
}

func GetTestSearcher() *MockFileSearcher {

	s1 := "~/golang/test"
	s2 := "C:\\Users\\Test"
	s3 := "/mnt/g/test"

	var p = []string{s1, s2, s3}

	b1 := []byte{0xFF, 0x32, 0x45, 0x78, 0x15}
	b2 := []byte{0xDF, 0x3F, 0x25, 0x70, 0x10}
	b3 := []byte{0xFF, 0x32, 0x45, 0x78, 0x15}

	m := make(map[string][]byte)
	m[s1] = b1
	m[s2] = b2
	m[s3] = b3

	mock := MockFileSearcher{
		paths: p,
		bytes: m,
	}

	return &mock
}
