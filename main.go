package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/kofuk/haniho/generator"
	"github.com/kofuk/haniho/tokenizer"
)

func runHttpServer() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		log.Println("/")

		text := ""

		if r.Method == http.MethodGet {
			if err := r.ParseForm(); err != nil {
				log.Println(err)
				w.WriteHeader(http.StatusBadRequest)
				return
			}

			if !r.Form.Has("text") {
				log.Println("text not specified")
				w.WriteHeader(http.StatusBadRequest)
				return
			}

			text = r.Form.Get("text")
		} else if r.Method == http.MethodPost {
			log.Println("reading text from body")

			limited := io.LimitReader(r.Body, 4096)
			data, err := io.ReadAll(limited)
			if err != nil {
				log.Println(err)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			text = string(data)
		}

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

	http.HandleFunc("/play", func(w http.ResponseWriter, r *http.Request) {
		if err := r.ParseForm(); err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "text/html")

		text := r.Form.Get("text")
		textHtml := strings.Replace(text, "<", "&lt;", -1)
		textUrl := url.QueryEscape(text)

		w.Write([]byte(fmt.Sprintf("<!doctype html><html><head><meta charset=\"utf-8\"><title>Haniho Playground</title><body><form action=\"/play\"><textarea name=\"text\">%s</textarea><input type=\"submit\"><br><audio controls><source src=\"/?text=%s\">", textHtml, textUrl)))
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
