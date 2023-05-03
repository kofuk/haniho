package tokenizer

const (
	EventNoteOn = iota
	EventNoteOff
)

type EventType int

type Event struct {
	Type EventType

	// Note number
	Note int

	// Time after previous event
	DeltaTime int
}

type Track struct {
	Events []Event
}

type RawData struct {
	// Note type assigned to a token.
	// For example, 0.25 means quarter notes.
	Resolution float64
	// Beats per minute
	BPM    float64
	Tracks []Track
}

func Tokenize(text []rune) (*RawData, error) {
	result := &RawData{}

	track := Track{}
	octaveShift := 3
	for i := 0; i < len(text); i++ {
		note := 0

		switch text[i] {
		case 'ド':
			note = 24
		case 'レ':
			note = 26
		case 'ミ':
			note = 28
		case 'フ':
			if i+1 < len(text) && text[i+1] == 'ァ' {
				i++
			}
			note = 29
		case 'ソ':
			note = 31
		case 'ラ':
			note = 33
		case 'シ':
			note = 35

		case '^':
			octaveShift++
			continue
		case 'v':
			octaveShift--
			continue
		}

		if i+1 < len(text) {
			next := text[i+1]
			if next == '#' {
				note++
				i++
			} else if next == 'b' {
				note--
				i++
			}
		}

		note += octaveShift * 12

		track.Events = append(track.Events,
			Event{Type: EventNoteOn, Note: note, DeltaTime: 0},
			Event{Type: EventNoteOff, Note: note, DeltaTime: 1},
		)
	}
	result.Resolution = 0.25
	result.BPM = 120
	result.Tracks = append(result.Tracks, track)

	return result, nil
}
