package nod

import (
	"io/ioutil"
	"os"

	"gopkg.in/neurosnap/sentences.v1"
)

// NewSentenceTokenizer creates a new sentence tokenizer
func NewSentenceTokenizer(language string) (sentences.SentenceTokenizer, error) {
	// not a perfect solution to rely on a installed package from the binary
	gopath := os.Getenv("GOPATH")
	trainingFile := gopath + "/src/gopkg.in/neurosnap/sentences.v1/data/" + language + ".json"
	f, err := os.Open(trainingFile)
	if err != nil {
		return nil, err
	}
	b, err := ioutil.ReadAll(f)
	if err != nil {
		return nil, err
	}

	training, err := sentences.LoadTraining(b)
	if err != nil {
		return nil, err
	}

	tokenizer := sentences.NewSentenceTokenizer(training)
	return tokenizer, nil
}
