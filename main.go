package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/kofuk/haniho/generator"
	"github.com/kofuk/haniho/tokenizer"
)

func runHttpServer() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if err := r.ParseForm(); err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		if !r.Form.Has("text") {
			log.Println("no text specified")
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		text := r.Form.Get("text")

		rawData, err := tokenizer.Tokenize([]rune(text))
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "audio/x-wav")

		if err := generator.Generate(rawData, w); err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	})

	fmt.Println("Listening :8000")
	log.Fatal(http.ListenAndServe(":8000", nil))
}

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "usage: haniho <str>\n       haniho --server")
		os.Exit(1)
	}

	if os.Args[1] == "--server" {
		runHttpServer()
		return
	}

	rawData, err := tokenizer.Tokenize([]rune(os.Args[1]))
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
