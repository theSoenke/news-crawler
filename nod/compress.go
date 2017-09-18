package nod

import (
	"os"

	"github.com/dsnet/compress/bzip2"
)

func compressBz2(output string, filename string) error {
	outDir := "out/nod/"
	if _, err := os.Stat(outDir); os.IsNotExist(err) {
		err := os.MkdirAll(outDir, os.ModePerm)
		if err != nil {
			return err
		}
	}
	f, err := os.Create(outDir + filename + ".bz2")
	if err != nil {
		return err
	}

	bz2, err := bzip2.NewWriter(f, &bzip2.WriterConfig{Level: 2})
	if err != nil {
		return err
	}

	_, err = bz2.Write([]byte(output))
	if err != nil {
		return err
	}

	defer bz2.Close()
	if err != nil {
		return err
	}

	return nil
}
