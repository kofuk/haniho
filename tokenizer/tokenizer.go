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
	rest := 0
	for i := 0; i < len(text); i++ {
		note := 0
		noteValue := 0

		switch text[i] {
		case 'ド':
			note = 24
			noteValue = 2

		case 'ﾄ':
			if i+1 < len(text) && text[i+1] == 'ﾞ' {
				i++
			}
			note = 24
			noteValue = 1

		case 'レ':
			note = 26
			noteValue = 2

		case 'ﾚ':
			note = 26
			noteValue = 1

		case 'ミ':
			note = 28
			noteValue = 2

		case 'ﾐ':
			note = 28
			noteValue = 1

		case 'フ':
			if i+1 < len(text) && text[i+1] == 'ァ' {
				i++
			}
			note = 29
			noteValue = 2

		case 'ﾌ':
			if i+1 < len(text) && text[i+1] == 'ｱ' {
				i++
			}
			note = 29
			noteValue = 1

		case 'ソ':
			note = 31
			noteValue = 2

		case 'ｿ':
			note = 31
			noteValue = 1

		case 'ラ':
			note = 33
			noteValue = 2

		case 'ﾗ':
			note = 33
			noteValue = 1

		case 'シ':
			note = 35
			noteValue = 2

		case 'ｼ':
			note = 35
			noteValue = 1

		case '^':
			octaveShift++
			continue
		case 'v':
			octaveShift--
			continue

		case 'ッ':
			rest = 2
			continue

		case 'ｯ':
			rest = 1
			continue
		}

		for j := i + 1; j < len(text); j++ {
			switch text[j] {
			case 'ー':
				noteValue += 2

			case '-':
				noteValue += 1

			case '#', '♯':
				note++

			case 'b', '♭':
				note--

			default:
				goto out
			}

			i++
		}
	out:

		note += octaveShift * 12

		track.Events = append(track.Events,
			Event{Type: EventNoteOn, Note: note, DeltaTime: rest},
			Event{Type: EventNoteOff, Note: note, DeltaTime: noteValue},
		)
		rest = 0
	}
	result.Resolution = 0.125
	result.BPM = 120
	result.Tracks = append(result.Tracks, track)

	return result, nil
}
