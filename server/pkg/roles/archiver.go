package roles

import (
	"bytes"
	"fmt"
	"github.com/c4milo/packit"
	"github.com/c4milo/unpackit"
	"strings"
)

type Archiver interface {
	Decompress(exDir string, compressedData []byte) error
	Compress(exDir string) ([]byte, error)
}

type archiver struct {
}

func newArchiver() Archiver {
	decompressor := new(archiver)
	return Archiver(decompressor)
}

func (a archiver) Compress(exDir string) ([]byte, error) {
	//err := os.Chdir(exDir)
	//if err != nil {
	//	return nil, err
	//}
	buf := bytes.NewBuffer(nil)
	packit.Zip(exDir, buf)
	data := buf.Bytes()
	if data == nil || len(data) == 0 {
		return nil, fmt.Errorf("failed to compress the '%s' directroy file", exDir)
	}
	return data, nil
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
