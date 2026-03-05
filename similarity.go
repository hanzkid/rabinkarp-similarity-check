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
	rabinKarp := newRabinKarp(query, r.KGram)
	hash := []int64{}
	for i := 0; i <= len(query)-r.KGram+1; i++ {
		hash = append(hash, rabinKarp.Hash)
	}
	r.QueryHashes = hash
}

func (r *RabinKarpSimilarity) intersectHashes(hashes1, hashes2 []int64) int {
	set := make(map[int64]struct{})

	for _, h := range hashes1 {
		set[h] = struct{}{}
	}

	intersection := 0
	for _, h := range hashes2 {
		if _, ok := set[h]; ok {
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

	rabinKarp := newRabinKarp(corpus, r.KGram)
	hash := []int64{}

	for i := 0; i <= len(corpus)-r.KGram+1; i++ {
		hash = append(hash, rabinKarp.Hash)

		if !rabinKarp.nextWindow() {
			break
		}
	}
	r.CorpusHashes[key] = hash
}

// CalculateSimilarity calculates the similarity between the query and the corpus
// using the intersectHashes function to find the intersection of the hashes
// key is the key of the corpus in the CorpusHashes map, example: "text-1", "text-2", etc.
func (r *RabinKarpSimilarity) CalculateSimilarity(query string, key string) float64 {
	r.calculateQueryHashes(query)
	intersect := r.intersectHashes(r.QueryHashes, r.CorpusHashes[key])
	totalHashes := len(r.QueryHashes) + len(r.CorpusHashes[key]) - intersect
	return float64(intersect) / float64(totalHashes)
}
