package generator

import (
	"bufio"
	"math"
	"os"

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

func Generate(data tokenizer.RawData) error {
	out, err := os.Create("/tmp/hoge.wav")
	if err != nil {
		return err
	}
	defer out.Close()

	bufWriter := bufio.NewWriter(out)
	defer bufWriter.Flush()

	encoder := wav.NewWriter(out, 88200, 2, 44100, 16)

	for i := 0; i < 88200; i++ {
		sample := wav.Sample{}

		freq := int(noteNo[60])

		sampleVal := (int(getSampleValue(calculateRatio(i, 44100, freq))*32767) + int(getSampleValue(calculateRatio(i, 44100, freq / 2))*32767) + int(getSampleValue(calculateRatio(i, 44100, freq * 2))*32767)) / 3
		sampleVal = int(float64(sampleVal) * (getSampleValue(calculateRatio(i, 44100, 6))*0.3 + 1.7) / 2.0)

		sample.Values[0] = sampleVal
		sample.Values[1] = sampleVal

		if err := encoder.WriteSamples([]wav.Sample{sample}); err != nil {
			return err
		}
	}

	return nil
}
