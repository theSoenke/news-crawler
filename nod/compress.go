package nod

import (
	"os"

	"github.com/dsnet/compress/bzip2"
)

func compressBz2(output string, filename string) error {
	f, err := os.Create("out/nod/" + filename + ".bz2")
	if err != nil {
		return err
	}

	bzWriter, err := bzip2.NewWriter(f, &bzip2.WriterConfig{Level: 2})
	if err != nil {
		return err
	}

	_, err = bzWriter.Write([]byte(output))
	if err != nil {
		return err
	}

	defer bzWriter.Close()
	if err != nil {
		return err
	}

	return nil
}
