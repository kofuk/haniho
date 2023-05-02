package tokenizer

import (
	"testing"
)

func Test_Tokenize(t *testing.T) {
	result, _ := Tokenize([]rune("ドレミファソ"))
	if len(result.Tracks) != 1 {
		t.Fatalf("len(DDRawData.Tracks) must be 1, but was %d\n", len(result.Tracks))
	}
	if len(result.Tracks[0].Events) != 10 {
		t.Fatalf("len(DDRawData.Tracks[9].Events) must be 10, but was %d\n", len(result.Tracks[0].Events))
	}
}
