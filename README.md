## Rabin–Karp Similarity

A small Go library that uses the Rabin–Karp rolling hash over k-grams to compute Jaccard similarity between strings (or text corpora).

The core idea:

- Split texts into overlapping k-grams.
- Hash each k-gram using a Rabin–Karp rolling hash.
- Compare the sets of hashes using Jaccard similarity.

## Installation

Make sure your module path matches the one declared in `go.mod`. For this repo it is:

```bash
go get github.com/hanzkid/rabinkarp-similarity-check@v0.1.0
```

You can also use the latest version:

```bash
go get github.com/hanzkid/rabinkarp-similarity-check@latest
```

Then import it:

```go
import rabinkarpsimilarity "github.com/hanzkid/rabinkarp-similarity-check"
```

## Quick start

```go
package main

import (
	"fmt"

	rabinkarpsimilarity "github.com/hanzkid/rabinkarp-similarity-check"
)

func main() {
	// Create a similarity calculator with k-gram size 3.
	rk := rabinkarpsimilarity.NewRabinKarpSimilarity(2, 3)

	// Register corpus texts and precompute their hashes.
	rk.CalculateCorpusHashes("the quick brown fox", "text-1")
	rk.CalculateCorpusHashes("the quick blue fox", "text-2")

	// Compare a query against a specific corpus key.
	query := "the quick brown fox jumps"
	sim1 := rk.CalculateSimilarity(query, "text-1")
	sim2 := rk.CalculateSimilarity(query, "text-2")

	fmt.Printf("Similarity(query, text-1) = %.4f\n", sim1)
	fmt.Printf("Similarity(query, text-2) = %.4f\n", sim2)
}
```

## API overview

- **`type RabinKarpSimilarity`**
  - Holds:
    - `CorpusHashes map[string][]int64`: precomputed k-gram hashes per corpus key.
    - `QueryHashes []int64`: k-gram hashes for the last query.
    - `KGram int`: chosen k-gram size.

- **`NewRabinKarpSimilarity(numberOfCorpus int, kgram int) *RabinKarpSimilarity`**
  - Creates a new similarity calculator.
  - `numberOfCorpus` is an initial capacity hint for the internal `map`.
  - `kgram` is the k-gram size used for hashing (e.g. 2, 3, 5).

- **`(*RabinKarpSimilarity) CalculateCorpusHashes(corpus, key string)`**
  - Computes all k-gram hashes for `corpus` and stores them under `key`.
  - Safe to call from multiple goroutines; internal access is guarded by a mutex.

- **`(*RabinKarpSimilarity) CalculateSimilarity(query, key string) float64`**
  - Computes Jaccard similarity between the k-gram hash set of `query` and the corpus identified by `key`.
  - Returns a value in \[0, 1\]:
    - `1.0` means identical k-gram sets.
    - `0.0` means no shared k-grams (under the hash function).

## How similarity is computed

For a chosen `K`:

1. Both `corpus` and `query` are converted to overlapping k-grams of size `K` (using runes).
2. Each k-gram is hashed via a Rabin–Karp rolling hash.
3. Duplicates are removed by converting the slices of hashes into sets.
4. The similarity is:

\[
\text{similarity} = \frac{|A \cap B|}{|A \cup B|}
\]

where:

- \(A\) is the set of hashes for the query.
- \(B\) is the set of hashes for the corpus.

## Running tests

From the project root:

```bash
go test ./...
```

