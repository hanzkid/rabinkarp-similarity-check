package rabinkarpsimilarity

import (
	"math"
	"testing"
)

func almostEqual(a, b, eps float64) bool {
	return math.Abs(a-b) <= eps
}

func TestSimilarityIdenticalTexts(t *testing.T) {
	rk := NewRabinKarpSimilarity(1, 3)

	corpus := "the quick brown fox"
	key := "text-1"
	rk.CalculateCorpusHashes(corpus, key)

	sim := rk.CalculateSimilarity(corpus, key)
	if !almostEqual(sim, 1.0, 1e-9) {
		t.Fatalf("expected similarity 1.0 for identical texts, got %f", sim)
	}
}

func TestSimilarityCompletelyDifferent(t *testing.T) {
	rk := NewRabinKarpSimilarity(2, 3)

	// Use clearly disjoint character sets to minimize accidental hash collisions.
	corpus := "aaaaaa"
	query := "bbbbbb"
	key := "text-1"

	rk.CalculateCorpusHashes(corpus, key)
	sim := rk.CalculateSimilarity(query, key)

	if sim != 0 {
		t.Fatalf("expected similarity 0 for disjoint texts, got %f", sim)
	}
}

func TestSimilarityPartialOverlap(t *testing.T) {
	rk := NewRabinKarpSimilarity(1, 2)

	corpus := "abcd"
	query := "bc"
	key := "text-1"

	// For k=2:
	// corpus k-grams: "ab", "bc", "cd"  -> 3 unique
	// query  k-grams: "bc"              -> 1 unique
	// intersection: {"bc"} = 1
	// union: 3 + 1 - 1 = 3
	// expected similarity = 1/3

	rk.CalculateCorpusHashes(corpus, key)
	sim := rk.CalculateSimilarity(query, key)

	expected := 1.0 / 3.0
	if !almostEqual(sim, expected, 1e-9) {
		t.Fatalf("expected similarity %f, got %f", expected, sim)
	}
}

