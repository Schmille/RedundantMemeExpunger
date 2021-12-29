package main

import (
	"errors"
	"io"
	"io/fs"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"time"
)

var opts struct {
	Verbose bool   `short:"v" long:"verbose" description:"Show verbose outputs"`
	Trial   bool   `short:"t" long:"trail-run" description:"Perform all actions, but do not delete anything"`
	Input   string `short:"i" long:"input-path" description:"Select the input folder"`
	Backup  bool   `short:"b" long:"backup" description:"Instead of deleting, copy all files into a sub-folder"`
	SizeLimit string `short:"s" long:"size-limit" description:"Limit on the maximum size of files that should be processed" default:"3MB"`
}

type NoopDeleter struct {
}

type MockFileSearcher struct {
	paths []string
	bytes map[string][]byte
}

func (n *NoopDeleter) Delete(_ string) error {
	// Do nothing.
	return nil
}

func (m *MockFileSearcher) GetFilePaths() []string {
	return m.paths
}

func (m *MockFileSearcher) GetBytes(path string) ([]byte, error) {
	return m.bytes[path], nil
}

type StandardFileSearcher struct {
	basePath   string
	maxFileSize int64
	childPaths map[string]fs.FileInfo
}

type StandardDeleter struct {
}

type BackupDeleter struct {
	inner      Deleter
	backupPath string
}

func NewStandardFileSearcher(basePath string, maxFileSize int64) (*StandardFileSearcher, error) {
	if len(opts.Input) <= 0 {
		return nil, errors.New("please provide an input path")
	}

	paths := make(map[string]fs.FileInfo)

	err := filepath.Walk(basePath, func(path string, info fs.FileInfo, err error) error {

		if info.IsDir() {
			if path != basePath {
				return fs.SkipDir
			}
			return nil
		}

		paths[path] = info
		return nil
	})

	if err != nil {
		return nil, err
	}

	return &StandardFileSearcher{
		basePath:   basePath,
		maxFileSize: maxFileSize,
		childPaths: paths,
	}, nil
}

func (s *StandardFileSearcher) GetFilePaths() []string {
	keys := make([]string, 0, len(s.childPaths))
	for k, v := range s.childPaths {

		if v.Size() > s.maxFileSize {
			verboseLn(k + " too large to process. Skipping.")
			continue
		}

		keys = append(keys, k)
	}

	return keys
}

func (s *StandardFileSearcher) GetBytes(path string) ([]byte, error) {
	return ioutil.ReadFile(path)
}

func (d *StandardDeleter) Delete(path string) error {
	return os.Remove(path)
}

func NewBackupDeleter(base string, deleter Deleter) *BackupDeleter {
	timestamp := time.Now().Format("2006-01-02_15-04-05")
	path := filepath.Join(base, ".backup", timestamp)

	err := os.MkdirAll(path, 777)

	if err != nil {
		log.Fatal(err)
	}

	return &BackupDeleter{
		inner:      deleter,
		backupPath: path,
	}
}

func (d *BackupDeleter) Delete(path string) error {
	reader, err := os.Open(path)
	if err != nil {
		return err
	}

	_, name := filepath.Split(path)
	writer, err := os.Create(filepath.Join(d.backupPath, name))
	if err != nil {
		return err
	}

	_, err = io.Copy(writer, reader)
	if err != nil {
		return err
	}

	err = reader.Close()
	if err != nil {
		return err
	}

	err = writer.Close()
	if err != nil {
		return err
	}

	return d.inner.Delete(path)
}