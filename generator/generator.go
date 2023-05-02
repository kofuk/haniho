package generator

import (
	"bufio"
	"math"
	"os"

	"github.com/kofuk/haniho/tokenizer"
	"github.com/youpy/go-wav"
)

func Generate(data tokenizer.RawData) error {
	out, err := os.Create("/tmp/hoge.wav")
	if err != nil {
		return err
	}
	defer out.Close()

	bufWriter := bufio.NewWriter(out)
	defer bufWriter.Flush()

	encoder := wav.NewWriter(out, 44100, 2, 44100, 16)

	for i := 0; i < 44100; i++ {
		sample := wav.Sample{}

		sampleVal := int(math.Sin((44100%440)/440.0*2*math.Pi) * 32768)

		sample.Values[0] = sampleVal
		sample.Values[1] = sampleVal

		if err := encoder.WriteSamples([]wav.Sample{sample}); err != nil {
			return err
		}
	}

	return nil
}
