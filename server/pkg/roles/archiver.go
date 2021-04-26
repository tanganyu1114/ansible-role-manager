package roles

import (
	"bytes"
	"fmt"
	"github.com/c4milo/unpackit"
	"strings"
)

type Decompressor interface {
	Decompress(exDir string, compressedData []byte) error
}

type archiver struct {
}

func NewDecompressor() Decompressor {
	decompressor := new(archiver)
	return Decompressor(decompressor)
}

func (a archiver) Decompress(exDir string, compressedData []byte) error {
	buf := bytes.NewReader(compressedData)
	afterExDir, err := unpackit.Unpack(buf, exDir)
	if err != nil {
		return err
	}
	if !strings.EqualFold(strings.TrimSpace(afterExDir), strings.TrimSpace(exDir)) {
		return fmt.Errorf("the exDir is inconsistent before and after decompression. before exDir = '%s', after exDir = '%s'", exDir, afterExDir)
	}
	return nil
}
