package generator

import (
	"bufio"
	"io"
	"log"
	"math"

	"github.com/kofuk/haniho/tokenizer"
	"github.com/youpy/go-wav"
)

func getSampleValue(ratio float64) float64 {
	return math.Sin(ratio * math.Pi * 2.0)
}

func calculateRatio(sampleNo, samplePerSec, freq int) float64 {
	samplePerCycle := samplePerSec / freq
	return float64(sampleNo%samplePerCycle) / float64(samplePerCycle)
}

func calculateLength(data *tokenizer.RawData) float64 {
	if len(data.Tracks) != 1 {
		log.Fatal("Multi-track support has not implemented yet!")
	}

	result := 0.0

	for _, ev := range data.Tracks[0].Events {
		result += float64(ev.DeltaTime) * data.Resolution * 4
	}

	result *= 60.0 / float64(data.BPM)

	return result
}

func Generate(data *tokenizer.RawData, w io.Writer) (err error) {
	bufWriter := bufio.NewWriter(w)
	defer func() {
		if err != nil {
			err = bufWriter.Flush()
		}
	}()

	length := calculateLength(data)

	encoder := wav.NewWriter(bufWriter, uint32(length)*44100, 2, 44100, 16)

	for i := 0; i < int(length)*44100; i++ {
		sample := wav.Sample{}

		freq := int(noteNo[60])

		sampleVal := (int(getSampleValue(calculateRatio(i, 44100, freq))*32767) + int(getSampleValue(calculateRatio(i, 44100, freq/2))*32767) + int(getSampleValue(calculateRatio(i, 44100, freq*2))*32767)) / 3
		sampleVal = int(float64(sampleVal) * (getSampleValue(calculateRatio(i, 44100, 6))*0.3 + 1.7) / 2.0)

		sample.Values[0] = sampleVal
		sample.Values[1] = sampleVal

		if err = encoder.WriteSamples([]wav.Sample{sample}); err != nil {
			return
		}
	}

	return
}
