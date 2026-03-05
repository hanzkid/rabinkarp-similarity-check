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

func NewRabinKarpSimilarity(numberOfCorpus int, kgram int) *RabinKarpSimilarity {
	rabinKarp := &RabinKarpSimilarity{
		CorpusHashes: make(map[string][]int64, numberOfCorpus),
		KGram:        kgram,
	}

	return rabinKarp
}

func (r *RabinKarpSimilarity) CalculateCorpusHashes(corpus string, key string) {
	r.Mu.Lock()
	defer r.Mu.Unlock()

	rabinKarp := NewRabinKarp(corpus, r.KGram)
	hash := []int64{}

	for i := 0; i <= len(corpus)-r.KGram+1; i++ {
		hash = append(hash, rabinKarp.Hash)

		if !rabinKarp.NextWindow() {
			break
		}
	}
	r.CorpusHashes[key] = hash
}

func (r *RabinKarpSimilarity) CalculateQueryHashes(query string) {
	rabinKarp := NewRabinKarp(query, r.KGram)
	hash := []int64{}
	for i := 0; i <= len(query)-r.KGram+1; i++ {
		hash = append(hash, rabinKarp.Hash)
	}
	r.QueryHashes = hash
}

func (r *RabinKarpSimilarity) CalculateSimilarity(query string, key string) float64 {
	r.CalculateQueryHashes(query)
	intersect := r.intersectHashes(r.QueryHashes, r.CorpusHashes[key])
	totalHashes := len(r.QueryHashes) + len(r.CorpusHashes[key])
	return float64(intersect) / float64(totalHashes)
}

func (r *RabinKarpSimilarity) intersectHashes(hashes1, hashes2 []int64) int {
	intersection := 0

	set := make(map[int64]bool)
	set2 := make(map[int64]bool)

	for _, num := range hashes1 {
		set[num] = true
	}

	for _, num := range hashes2 {
		if set[num] {
			intersection++
		}
		set2[num] = true
	}

	for _, num := range hashes1 {
		if set2[num] {
			intersection++
		}
	}

	return intersection
}
