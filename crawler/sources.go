package crawler

import (
	"encoding/json"
	"io/ioutil"
)

func readSourcesFile(path string) ([]string, error) {
	sourceFile, err := ioutil.ReadFile(path)

	if err != nil {
		return nil, err
	}

	var sources = make([]string, 0)
	err = json.Unmarshal(sourceFile, &sources)

	if err != nil {
		return nil, err
	}

	return sources, nil
}
