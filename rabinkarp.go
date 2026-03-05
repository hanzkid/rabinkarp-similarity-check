package rabinkarpsimilarity

import "strings"

type RabinKarp struct {
	Base        int
	Text        string
	PatternSize int
	Start       int
	End         int
	Mod         int64
	Hash        int64
	basePow     int64 // precomputed Base^(PatternSize-1) % Mod
}

func NewRabinKarp(text string, patternSize int) *RabinKarp {
	rb := &RabinKarp{
		Base:        26,
		PatternSize: patternSize,
		Start:       0,
		End:         0,
		Mod:         16777619,
		Text:        strings.ToLower(text),
	}

	rb.init()

	return rb
}

// init precomputes the first window hash and Base^(PatternSize-1) % Mod.
func (rb *RabinKarp) init() {
	textRunes := []rune(rb.Text)
	if rb.PatternSize <= 0 || rb.PatternSize > len(textRunes) {
		rb.Start = 0
		rb.End = 0
		rb.Hash = 0
		rb.basePow = 1
		return
	}

	// Precompute basePow = Base^(PatternSize-1) % Mod
	rb.basePow = 1
	for i := 0; i < rb.PatternSize-1; i++ {
		rb.basePow = (rb.basePow * int64(rb.Base)) % rb.Mod
	}

	// Compute hash for the first window.
	hashValue := int64(0)
	for i := 0; i < rb.PatternSize; i++ {
		// Map rune to a positive value; assume lowercase letters.
		value := int64(textRunes[i]) - 96
		if value < 0 {
			value = 0
		}
		hashValue = (hashValue*int64(rb.Base) + value) % rb.Mod
	}

	rb.Start = 0
	rb.End = rb.PatternSize
	rb.Hash = hashValue
}

// NextWindow returns boolean and create new hash for the next window
// using rolling hash technique
func (rb *RabinKarp) NextWindow() bool {
	textRunes := []rune(rb.Text)

	// No more full windows possible.
	if rb.End >= len(textRunes) {
		return false
	}

	// Rolling hash: remove leading rune, multiply, add trailing rune.
	leading := int64(textRunes[rb.Start]) - 96
	if leading < 0 {
		leading = 0
	}
	trailing := int64(textRunes[rb.End]) - 96
	if trailing < 0 {
		trailing = 0
	}

	rb.Hash = (rb.Hash - (leading*rb.basePow)%rb.Mod + rb.Mod) % rb.Mod
	rb.Hash = (rb.Hash*int64(rb.Base) + trailing) % rb.Mod

	rb.Start++
	rb.End++
	return true
}

// CurrentWindowText return the current window text
func (rb *RabinKarp) CurrentWindowText() string {
	return rb.Text[rb.Start:rb.End]
}
