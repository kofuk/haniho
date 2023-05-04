package main

//go:generate go run .

import (
	"bytes"
	"fmt"
	"go/format"
	"io"
	"log"
	"math"
	"os"
)

var (
	keyNames = []string{"C", "C#", "D", "D#", "E", "F", "F#", "G", "G#", "A", "A#", "B"}
)

func generateSinWaveTable(out io.Writer) {
	fmt.Fprintln(out, "\"sin\": []float64{")
	for i := 0; i < 2048; i++ {
		fmt.Fprintf(out, "%f,", math.Sin(float64(i)/2048*2*math.Pi))
	}
	fmt.Fprintln(out, "},")
}

func generateSawWaveTable(out io.Writer) {
	fmt.Fprintln(out, "\"saw\": []float64{")
	for i := 0; i < 2048; i++ {
		fmt.Fprintf(out, "%f,", float64(i-1024)/1024)
	}
	fmt.Fprintln(out, "},")
}

func generateRectWaveTable(out io.Writer) {
	fmt.Fprintln(out, "\"rect\": []float64{")
	for i := 0; i < 2048; i++ {
		val := -1.0
		if i < 1024 {
			val = 1.0
		}
		fmt.Fprintf(out, "%f,", val)
	}
	fmt.Fprintln(out, "},")
}

func generateWaveTables(out io.Writer) {
	fmt.Fprintln(out, "var wavetables = map[string][]float64{")
	generateSinWaveTable(out)
	generateSawWaveTable(out)
	generateRectWaveTable(out)
	fmt.Fprintln(out, "}")
}

func generateNoteNoData(out io.Writer) {
	fmt.Fprintln(out, "var noteNo = []float64{")
	for d := 0; d < 128; d++ {
		freq := math.Pow(2, (float64(d)-69.0)/12.0) * 440
		keyName := keyNames[d%12]
		octave := d/12 - 2
		fmt.Fprintf(out, "%.4f, // %d: %s%d\n", freq, d, keyName, octave)
	}
	fmt.Fprintln(out, "}")
}

func main() {
	out := bytes.NewBuffer(nil)

	fmt.Fprintln(out, "// Code generated by codegen.go; DO NOT EDIT.")
	fmt.Fprintln(out, "package generator;")

	generateNoteNoData(out)
	generateWaveTables(out)

	outFile, err := os.Create("../data.go")
	if err != nil {
		log.Fatal(err)
	}
	defer outFile.Close()

	formatted, err := format.Source(out.Bytes())
	if err != nil {
		log.Fatal(err)
	}

	if _, err := outFile.Write(formatted); err != nil {
		log.Fatal(err)
	}
}
