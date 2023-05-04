package generator

import (
	"bufio"
	"io"
	"log"
	"math"

	"github.com/kofuk/haniho/tokenizer"
	"github.com/youpy/go-wav"
)

const (
	sampleRate = 44100
)

type sinOscillator struct {
	frame int
	freq  float64
}

func (self *sinOscillator) oscillate() float64 {
	curFrame := self.frame
	self.frame++
	samplePerCycle := sampleRate / self.freq
	return math.Sin(math.Remainder(float64(curFrame), samplePerCycle) / samplePerCycle * math.Pi * 2)
}

func (self *sinOscillator) handleEvent(ev tokenizer.Event) {
	if ev.Type == tokenizer.EventNoteOn {
		self.frame = 0
		self.freq = noteNo[ev.Note]
	}
}

type envState int

const (
	envNone envState = iota
	envAttack
	envDecay
	envSustaion
	envRelease
)

type envelopeGenerator struct {
	state        envState
	frame        int
	attackFrame  int
	decayFrame   int
	sustainLevel float64
	releaseFrame int
}

func newEnvelope(attack, decay, sustain, release float64) envelopeGenerator {
	return envelopeGenerator{
		state:        envNone,
		frame:        0,
		attackFrame:  int(attack * sampleRate),
		decayFrame:   int(decay * sampleRate),
		sustainLevel: sustain,
		releaseFrame: int(release * sampleRate),
	}
}

func (self *envelopeGenerator) handleEvent(ev tokenizer.Event) {
	if ev.Type == tokenizer.EventNoteOn {
		self.state = envAttack
		self.frame = 0
	} else if ev.Type == tokenizer.EventNoteOff {
		self.state = envRelease
		self.frame = 0
	}
}

func (self *envelopeGenerator) filter(src float64) float64 {
	if self.state == envAttack {
		p := float64(self.frame) / float64(self.attackFrame)

		if self.frame < self.attackFrame {
			self.frame++
		} else {
			self.frame = 0
			self.state = envDecay
		}
		return src * p
	} else if self.state == envDecay {
		p := float64(self.frame) / float64(self.decayFrame)

		if self.frame < self.decayFrame {
			self.frame++
		} else {
			self.frame = 0
			self.state = envSustaion
		}
		return src - src*(1.0-self.sustainLevel)*p
	} else if self.state == envSustaion {
		return src * self.sustainLevel
	} else if self.state == envRelease {
		p := float64(self.frame) / float64(self.decayFrame)

		if self.frame < self.releaseFrame {
			self.frame++
		} else {
			self.frame = 0
			self.state = envNone
		}
		return src * self.sustainLevel * p
	}

	return 0
}

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

	encoder := wav.NewWriter(bufWriter, uint32(length*44100), 2, 44100, 16)

	osc := sinOscillator{}
	env := newEnvelope(0.0, 0.2, 0.0, 0.1)

	framesElapsed := 0
	curEventNo := 0

	for i := 0; i < int(length*44100.0); i++ {
		if curEventNo < len(data.Tracks[0].Events) {
			ev := data.Tracks[0].Events[curEventNo]
			deltaFrame := float64(ev.DeltaTime) * data.Resolution * 4
			deltaFrame *= 60.0 / float64(data.BPM) * 44100
			if float64(framesElapsed) >= deltaFrame {
				osc.handleEvent(ev)
				env.handleEvent(ev)

				curEventNo++
				framesElapsed = 0
			} else {
				framesElapsed++
			}
		}

		rawSampleVal := env.filter(osc.oscillate())
		sampleVal := int(rawSampleVal * 32767)

		sample := wav.Sample{}

		sample.Values[0] = sampleVal
		sample.Values[1] = sampleVal

		if err = encoder.WriteSamples([]wav.Sample{sample}); err != nil {
			return
		}
	}

	return
}
