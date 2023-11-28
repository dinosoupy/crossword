package crossword

import (
	"fmt"
	"testing"
)

func TestLexicon(t *testing.T) {
	savedAllWords, savedLexicon := AllWords, Lexicon
	defer func() {
		AllWords, Lexicon = savedAllWords, savedLexicon
	}()

	AllWords = map[int][]string{
		3:  {"aaa", "aba"},
		4:  {"aaah", "aaas"},
		12: {"aaabatteries"},
	}

	if err := NewLexicon(); err != nil {
		t.Errorf("lexicon: %v", err)
		return
	}

	var tests = []struct {
		wlen     int
		letter   rune
		position int
		want     []string
	}{
		{4, 's', 3, []string{"aaas"}},
		{4, 'a', 0, []string{"aaah", "aaas"}},
		{3, 'a', 0, []string{"aaa", "aba"}},
	}

	for _, test := range tests {
		descr := fmt.Sprintf("Words of length %d, that have letter %c at position %d:\n", test.wlen, test.letter, test.position)

		numWords := len(AllWords[test.wlen])
		buf := make([]uint64, numWords)
		got := make([]string, 0)
		ch := test.letter - 'a'
		l := Lexicon[test.wlen][test.position].bitarrays[ch]
		out := l.GetSetBits(0, buf)

		for _, i := range out {
			word := AllWords[test.wlen][i]
			got = append(got, word)
		}

		if !isEqual(got, test.want) {
			t.Errorf("%s\nGot:\t%v\nWant:\t%v", descr, got, test.want)
		}
	}
}

func isEqual(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}
