package rabinkarpsimilarity

import "sync"

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
