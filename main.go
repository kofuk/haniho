package main

import (
	"log"
	"os"

	"github.com/kofuk/haniho/generator"
	"github.com/kofuk/haniho/tokenizer"
)

func main() {
	rawData, err := tokenizer.Tokenize([]rune("ドレミファソラシド"))
	if err != nil {
		log.Fatal(err)
	}

	out, err := os.Create("/tmp/hoge.wav")
	if err != nil {
		log.Fatal(err)
	}
	defer out.Close()

	if err := generator.Generate(rawData, out); err != nil {
		log.Fatal(err)
	}
}
