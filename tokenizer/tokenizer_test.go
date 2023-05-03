package tokenizer

import (
	"math"
	"testing"
)

const (
	epsilon = 1e-5
)

func Test_Tokenize(t *testing.T) {
	testcases := []struct {
		text   string
		expect *RawData
	}{
		{
			text: "ドﾚミ",
			expect: &RawData{
				Resolution: 0.125,
				BPM:        120.0,
				Tracks: []Track{
					{
						Events: []Event{
							{
								Type:      EventNoteOn,
								Note:      60,
								DeltaTime: 0,
							},
							{
								Type:      EventNoteOff,
								Note:      60,
								DeltaTime: 2,
							},
							{
								Type:      EventNoteOn,
								Note:      62,
								DeltaTime: 0,
							},
							{
								Type:      EventNoteOff,
								Note:      62,
								DeltaTime: 1,
							},
							{
								Type:      EventNoteOn,
								Note:      64,
								DeltaTime: 0,
							},
							{
								Type:      EventNoteOff,
								Note:      64,
								DeltaTime: 2,
							},
						},
					},
				},
			},
		},
		{
			text: "ド-^ドvvドー",
			expect: &RawData{
				Resolution: 0.125,
				BPM:        120.0,
				Tracks: []Track{
					{
						Events: []Event{
							{
								Type:      EventNoteOn,
								Note:      60,
								DeltaTime: 0,
							},
							{
								Type:      EventNoteOff,
								Note:      60,
								DeltaTime: 3,
							},
							{
								Type:      EventNoteOn,
								Note:      72,
								DeltaTime: 0,
							},
							{
								Type:      EventNoteOff,
								Note:      72,
								DeltaTime: 2,
							},
							{
								Type:      EventNoteOn,
								Note:      48,
								DeltaTime: 0,
							},
							{
								Type:      EventNoteOff,
								Note:      48,
								DeltaTime: 4,
							},
						},
					},
				},
			},
		},
		{
			text: "レbレレ#",
			expect: &RawData{
				Resolution: 0.125,
				BPM:        120.0,
				Tracks: []Track{
					{
						Events: []Event{
							{
								Type:      EventNoteOn,
								Note:      61,
								DeltaTime: 0,
							},
							{
								Type:      EventNoteOff,
								Note:      61,
								DeltaTime: 2,
							},
							{
								Type:      EventNoteOn,
								Note:      62,
								DeltaTime: 0,
							},
							{
								Type:      EventNoteOff,
								Note:      62,
								DeltaTime: 2,
							},
							{
								Type:      EventNoteOn,
								Note:      63,
								DeltaTime: 0,
							},
							{
								Type:      EventNoteOff,
								Note:      63,
								DeltaTime: 2,
							},
						},
					},
				},
			},
		},
	}

	for _, tc := range testcases {
		result, _ := Tokenize([]rune(tc.text))

		if math.Abs(result.Resolution-tc.expect.Resolution) > epsilon {
			t.Fatalf("Resolution must be %v, but was %v\n", result.Resolution, tc.expect.Resolution)
		}
		if math.Abs(result.BPM-tc.expect.BPM) > epsilon {
			t.Fatalf("BPM must be %v, but was %v\n", result.BPM, tc.expect.BPM)
		}

		if len(result.Tracks) != len(tc.expect.Tracks) {
			t.Fatalf("len(Tracks) must be %v, but was %v\n", len(result.Tracks), len(tc.expect.Tracks))
		}

		for i := 0; i < len(result.Tracks); i++ {
			if len(result.Tracks[i].Events) != len(tc.expect.Tracks[i].Events) {
				t.Fatalf("len(Tracks[%d].Events) must be %v, but was %v\n", i, len(result.Tracks[i].Events), len(tc.expect.Tracks[i].Events))
			}

			for j := 0; j < len(result.Tracks[i].Events); j++ {
				ev1 := result.Tracks[i].Events[j]
				ev2 := tc.expect.Tracks[i].Events[j]

				if ev1.Type != ev2.Type {
					t.Fatalf("Tracks[%d].Events[%d].Type must be %v, but was %v\n", i, j, ev1.Type, ev2.Type)
				}
				if ev1.DeltaTime != ev2.DeltaTime {
					t.Fatalf("Tracks[%d].Events[%d].DeltaTime must be %v, but was %v\n", i, j, ev1.DeltaTime, ev2.DeltaTime)
				}
				if ev1.Note != ev2.Note {
					t.Fatalf("Tracks[%d].Events[%d].Note must be %v, but was %v\n", i, j, ev1.Note, ev2.Note)
				}
			}
		}
	}
}
