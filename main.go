package main

import (
	"crypto/sha256"
	"fmt"
	"github.com/jessevdk/go-flags"
	"log"
	"os"
)

type FileSearcher interface {
	GetFilePaths() []string
	GetBytes(path string) ([]byte, error)
}

type Deleter interface {
	Delete(path string) error
}

func main() {
	_, err := flags.ParseArgs(&opts, os.Args)
	if err != nil {
		return
	}

	var deleter Deleter

	if opts.Trial {
		deleter = &NoopDeleter{}
	} else {
		deleter = &StandardDeleter{}
	}

	if opts.Backup {
		deleter = NewBackupDeleter(opts.Input, deleter)
	}

	var fileSearcher FileSearcher
	fileSearcher, err = NewStandardFileSearcher(opts.Input)

	if err != nil {
		log.Fatal(err)
	}

	itemsToDelete, err := processChecks(&fileSearcher)

	if err != nil {
		log.Fatal(err)
	}

	deleteAll(itemsToDelete, &deleter)
	verboseLn("Done.")
}

func processChecks(searcher *FileSearcher) ([]string, error) {
	return sameBytes(searcher), nil
}

func sameBytes(searcher *FileSearcher) []string {
	out := make([]string, 0)
	hasher := sha256.New()
	s := *searcher
	set := NewSet()

	for _, path := range s.GetFilePaths() {
		verbose("Checking " + path + "... ")
		bytes, err := s.GetBytes(path)

		if err != nil {
			log.Println(err)
			continue
		}

		hash := string(hasher.Sum(bytes))

		if set.Contains(hash) {
			verbose("Matching hash has been found!")
			out = append(out, path)
		}
		verbose("\n")
		
		set.Add(hash)
	}

	return out
}

func deleteAll(items []string, deleter *Deleter) {
	for _, item := range items {
		verboseLn("Deleting " + item)
		err := (*deleter).Delete(item)

		if err != nil {
			log.Println(err)
		}
	}
}

func verbose(msg string) {
	if opts.Verbose {
		fmt.Print(msg)
	}
}

func verboseLn(msg string) {
	if opts.Verbose {
		fmt.Println(msg)
	}
}
