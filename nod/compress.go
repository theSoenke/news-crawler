package nod

import (
	"os"
	"path"

	"github.com/dsnet/compress/bzip2"
)

func (corpus *dayCorpus) compress(output string, dir string, filename string) error {
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		err := os.MkdirAll(dir, os.ModePerm)
		if err != nil {
			return err
		}
	}

	file := path.Join(dir, filename+".bz2")
	f, err := os.Create(file)
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

	err = bz2.Close()
	return err
}
