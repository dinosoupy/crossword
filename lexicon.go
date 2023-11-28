package crossword

import (
	"fmt"
	"os"
	"strings"

	bitarray "github.com/Workiva/go-datastructures/bitarray"
)

// Number of words in the dictionary
const NUM_WORDS = 16

// AllWords just stores all words in a slice
// Lexicon points to indices within allWords to retrieve words
var AllWords = make(map[int][]string)

// Lexicon maps a word-length int to a slice of Layers;
// for every word length `wlen` there are `wlen` layers
var Lexicon = make(map[int][]Layer)

// A Layer contains 26 bitarrays for every letter in the alphabet
type Layer struct {
	bitarrays [26]bitarray.BitArray // uint8 is wasteful as we are just storing 0 or 1
}

func NewLayer(bitarrSize uint64) Layer {
	layer := Layer{}
	for i := 0; i < 26; i++ {
		bitarr := bitarray.NewBitArray(bitarrSize)
		layer.bitarrays[i] = bitarr
	}
	return layer
}

// ReadAllWords reads the contents of file at `filename` and initialises the AllWords map.
// This function should be called before calling NewLexicon()
func ReadAllWords(filename string) {
	data, err := os.ReadFile(filename)
	if err != nil {
		fmt.Printf("lexicon: %v", err)
		os.Exit(1)
	}

	lines := strings.Split(string(data), "\n")
	for _, line := range lines {
		word := strings.Split(line, ";")[0]
		wlen := len(word)
		AllWords[wlen] = append(AllWords[wlen], word)
	}
}

// NewLexicon initialises a lexicon from AllWords
func NewLexicon() error {
	for wlen, words := range AllWords {
		// Initialise Lexicon[wlen] layers by append new layers with empty bitarrays
		for i := 0; i < wlen; i++ {
			Lexicon[wlen] = append(Lexicon[wlen], NewLayer(uint64(len(words))))
		}
		// For every letter in every word, set the bit in the layer's appropriate bitarray
		for wordIdx, word := range words {
			for letterIdx, letter := range word {
				ch := letter - 'a'
				layer := Lexicon[wlen][letterIdx]
				if err := layer.bitarrays[ch].SetBit(uint64(wordIdx)); err != nil {
					return err
				}
			}
		}
	}
	return nil
}
