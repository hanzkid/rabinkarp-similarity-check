package rabinkarpsimilarity

import (
	"sync"
)

type RabinKarpSimilarity struct {
	CorpusHashes map[string][]int64
	QueryHashes  []int64
	KGram        int
	Mu           sync.Mutex
}

// NewRabinKarpSimilarity creates a new RabinKarpSimilarity instance
// numberOfCorpus is the number of corpus to calculate the hashes for
// kgram is the k-gram size to use for the hashes
func NewRabinKarpSimilarity(numberOfCorpus int, kgram int) *RabinKarpSimilarity {
	rabinKarp := &RabinKarpSimilarity{
		CorpusHashes: make(map[string][]int64, numberOfCorpus),
		KGram:        kgram,
	}
	return rabinKarp
}

func (r *RabinKarpSimilarity) calculateQueryHashes(query string) {
	runes := []rune(query)
	if r.KGram <= 0 || len(runes) < r.KGram {
		r.QueryHashes = nil
		return
	}

	rabinKarp := newRabinKarp(query, r.KGram)
	hashes := []int64{rabinKarp.Hash}
	for rabinKarp.nextWindow() {
		hashes = append(hashes, rabinKarp.Hash)
	}
	r.QueryHashes = hashes
}

func uniqueHashes(hashes []int64) map[int64]struct{} {
	set := make(map[int64]struct{}, len(hashes))
	for _, h := range hashes {
		set[h] = struct{}{}
	}
	return set
}

// intersectHashes returns the number of unique hashes that appear in both sets.
func (r *RabinKarpSimilarity) intersectHashes(set1, set2 map[int64]struct{}) int {
	if len(set1) == 0 || len(set2) == 0 {
		return 0
	}

	// Iterate over the smaller set for efficiency.
	if len(set1) > len(set2) {
		set1, set2 = set2, set1
	}

	intersection := 0
	for h := range set1 {
		if _, ok := set2[h]; ok {
			intersection++
		}
	}

	return intersection
}

// CalculateCorpusHashes calculates the hashes of the corpus and stores them in the CorpusHashes map
// using a mutex to prevent race conditions
// key is the key of the corpus in the CorpusHashes map, example: "text-1", "text-2", etc.
func (r *RabinKarpSimilarity) CalculateCorpusHashes(corpus string, key string) {
	r.Mu.Lock()
	defer r.Mu.Unlock()

	runes := []rune(corpus)
	if r.KGram <= 0 || len(runes) < r.KGram {
		r.CorpusHashes[key] = nil
		return
	}

	rabinKarp := newRabinKarp(corpus, r.KGram)
	hashes := []int64{rabinKarp.Hash}
	for rabinKarp.nextWindow() {
		hashes = append(hashes, rabinKarp.Hash)
	}
	r.CorpusHashes[key] = hashes
}

// CalculateSimilarity calculates the similarity between the query and the corpus
// using the intersectHashes function to find the intersection of the hashes
// key is the key of the corpus in the CorpusHashes map, example: "text-1", "text-2", etc.
func (r *RabinKarpSimilarity) CalculateSimilarity(query string, key string) float64 {
	r.calculateQueryHashes(query)
	corpusHashes := r.CorpusHashes[key]

	if len(r.QueryHashes) == 0 || len(corpusHashes) == 0 {
		return 0
	}

	querySet := uniqueHashes(r.QueryHashes)
	corpusSet := uniqueHashes(corpusHashes)

	intersect := r.intersectHashes(querySet, corpusSet)
	totalHashes := len(querySet) + len(corpusSet) - intersect
	if totalHashes == 0 {
		return 0
	}

	return float64(intersect) / float64(totalHashes)
}
