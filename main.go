package main

import (
	"log"
)

func main() {
	_, err := Parse("http://spiegel.de")

	if err != nil {
		log.Fatal(err)
		return
	}
}
