package broker

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/ml444/gomega/config"
)

func getFile(path string) (*os.File, error) {
	file, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		return nil, err
	}
	return file, nil
}

func getIdxFile(fileDir string, sequence uint64) (*os.File, error) {
	path := filepath.Join(fileDir, fmt.Sprintf("%d%s", sequence, config.IdxFileSuffix))
	return getFile(path)
}

func getDataFile(fileDir string, sequence uint64) (*os.File, error) {
	path := filepath.Join(fileDir, fmt.Sprintf("%d%s", sequence, config.DataFileSuffix))
	return getFile(path)
}

func searchIdxFileWithSequence(fileDir string, sequence uint64) (string, uint64, error) {
	var foundFile string
	var foundSequence uint64
	err := filepath.Walk(fileDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			fmt.Println(err)
			return err
		}
		if info.IsDir() {
			return nil
		}
		fname := info.Name()
		if filepath.Ext(fname) == ".idx" {
			seqStr := strings.TrimSuffix(fname, ".idx")
			seq, err := strconv.ParseUint(seqStr, 10, 64)
			if err != nil {
				return err
			}
			if seq <= sequence {
				foundSequence = seq
				foundFile = path
				return fmt.Errorf("found")
			}
		}
		return nil
	})
	if foundFile != "" {
		return foundFile, foundSequence, nil
	}
	if err == nil {
		return "", 0, fmt.Errorf("not found")
	}
	return "", 0, err
}
