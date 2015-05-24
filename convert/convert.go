package convert

import (
	"fmt"
	"io"
	"os"
	"strconv"

	"golang.org/x/text/transform"
)

const (
	// TmpMaxTrial is max count of try to generate temporary file path.
	TmpMaxTrial = 10
)

func isNotExist(path string) bool {
	_, err := os.Stat(path)
	if err != nil {
		return os.IsNotExist(err)
	}
	return false
}

func tmpPath(path, suffix string) (string, error) {
	p := path + "." + suffix
	if isNotExist(p) {
		return p, nil
	}
	for i := 1; i < TmpMaxTrial; i++ {
		p := path + "." + strconv.Itoa(i) + "." + suffix
		if isNotExist(p) {
			return p, nil
		}
	}
	return "", fmt.Errorf("can't genretate temporary path for %s", path)
}

func convertAll(dstpath, srcpath string, t transform.Transformer) error {
	// Open write/read files.
	w, err := os.OpenFile(dstpath, os.O_RDWR|os.O_CREATE|os.O_EXCL, 0666)
	if err != nil {
		return err
	}
	defer w.Close()
	r, err := os.Open(srcpath)
	if err != nil {
		return err
	}
	defer r.Close()
	// Convert
	w2 := transform.NewWriter(w, t)
	_, err = io.Copy(w2, r)
	return err
}

func swapFiles(path, pathNew string) error {
	pathOld, err := tmpPath(path, "old")
	if err != nil {
		return err
	}
	if err := os.Rename(path, pathOld); err != nil {
		return nil
	}
	if err := os.Rename(pathNew, path); err != nil {
		return nil
	}
	return os.Remove(pathOld)
}

// Convert converts a file with Transformer in-place.
func Convert(path string, t transform.Transformer) error {
	dstpath, err := tmpPath(path, "new")
	if err != nil {
		return err
	}
	if err := convertAll(path, dstpath, t); err != nil {
		os.Remove(dstpath)
		return err
	}
	return swapFiles(path, dstpath)
}
