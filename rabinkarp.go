package rabinkarpsimilarity

import (
	"math"
	"strings"
)

type RabinKarp struct {
	Base        int
	Text        string
	PatternSize int
	Start       int
	End         int
	Mod         int64
	Hash        int64
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

	rb.GetHash()

	return rb
}

// GetHash creates hash for the first window
func (rb *RabinKarp) GetHash() {
	hashValue := int64(0)
	for n := 0; n < rb.PatternSize; n++ {
		runeText := []rune(rb.Text)
		value := int64(runeText[n])
		mathpower := int64(math.Pow(float64(rb.Base), float64(rb.PatternSize-n-1)))
		hashValue = (hashValue + (value-96)*(mathpower)) % rb.Mod
	}
	rb.Start = 0
	rb.End = rb.PatternSize
	rb.Hash = hashValue
}

// NextWindow returns boolean and create new hash for the next window
// using rolling hash technique
func (rb *RabinKarp) NextWindow() bool {
	textRunes := []rune(rb.Text)

	if rb.End < len(textRunes)-1 {
		mathpower := int64(math.Pow(float64(rb.Base), float64(rb.PatternSize-1)))
		rb.Hash -= (int64(textRunes[rb.Start]) - 96) * mathpower
		rb.Hash *= int64(rb.Base)
		rb.Hash += int64(textRunes[rb.End] - 96)
		rb.Hash = rb.Hash % rb.Mod
		rb.Start++
		rb.End++
		return true
	}
	return false
}

// CurrentWindowText return the current window text
func (rb *RabinKarp) CurrentWindowText() string {
	return rb.Text[rb.Start:rb.End]
}
