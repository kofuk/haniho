package main

import (
	"log"

	"github.com/kofuk/haniho/generator"
	"github.com/kofuk/haniho/tokenizer"
)

func main() {
	rawData, err := tokenizer.Tokenize([]rune("ドレミファソ"))
	if err != nil {
		log.Fatal(err)
	}

	if err := generator.Generate(*rawData); err != nil {
		log.Fatal(err)
	}
}
