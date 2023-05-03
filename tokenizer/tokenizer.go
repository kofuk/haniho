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
	for i := 0; i < len(text); i++ {
		note := 0

		switch text[i] {
		case 'ド':
			note = 60
		case 'レ':
			note = 62
		case 'ミ':
			note = 64
		case 'フ':
			if i+1 < len(text) && text[i+1] == 'ァ' {
				i++
			}
			note = 65
		case 'ソ':
			note = 67
		case 'ラ':
			note = 69
		case 'シ':
			note = 71
		}

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
