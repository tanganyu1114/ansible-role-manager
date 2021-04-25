package roles

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"github.com/alibaba/sentinel-golang/util"
	"io"
	"os"
	"path/filepath"
)

type Decompressor interface {
	Decompress(exDir string, compressedData []byte) error
}

type tgzArchiver struct {
}

func NewDecompressor() Decompressor {
	decompressor := new(tgzArchiver)
	return Decompressor(decompressor)
}

func (g tgzArchiver) Decompress(exDir string, compressedData []byte) error {
	err := util.CreateDirIfNotExists(exDir)
	if err != nil {
		return err
	}
	buf := bytes.NewReader(compressedData)
	gzipReader, err := gzip.NewReader(buf)
	if err != nil {
		return err
	}
	defer gzipReader.Close()

	tReader := tar.NewReader(gzipReader)
	for {
		tHeader, err := tReader.Next()
		if err == io.EOF {
			break
		}

		if err != nil {
			return err
		}

		exF, err := os.Create(filepath.Join(exDir, tHeader.Name))
		if err != nil {
			return err
		}

		_, err = io.Copy(exF, tReader)
		if err != nil {
			exF.Close()
			return err
		}
		exF.Close()
	}
	return nil
}
