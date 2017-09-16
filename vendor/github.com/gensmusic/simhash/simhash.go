package simhash

import (
	"hash/fnv"

	"github.com/gensmusic/jiebago/analyse"
)

var (
	_HASH_BIT_LENGTH uint8 = 64
	extracter        analyse.TagExtracter
)

type hashWeigth struct {
	Hash   uint64
	Weight float64
}

// load jieba, idf, stop words dictionaries
func LoadDictionary(jiebapath, idfpath, stopwords string) error {
	if err := extracter.LoadDictionary(jiebapath); err != nil {
		return err
	}
	if err := extracter.LoadIdf(idfpath); err != nil {
		return err
	}
	if err := extracter.LoadStopWords(stopwords); err != nil {
		return err
	}

	return nil
}

/*
calculate simhash with top n keywords
calculate with all words if topN < 0
*/
func Simhash(s *string, topN int) uint64 {
	if s == nil {
		panic("simhash cannot hash nil string")
	}

	hashes := extractAndHash(s, topN)
	if len(hashes) == 0 {
		return 0
	}

	weights := calWeights(hashes)
	return fingerprint(weights)
}

func hasher(s string) uint64 {
	h := fnv.New64a()
	h.Write([]byte(s))
	return h.Sum64()
}

func extractAndHash(s *string, topN int) []hashWeigth {
	if s == nil {
		panic("cannot extract nil string")
	}

	words := extracter.ExtractTags(*s, topN)
	wordsLen := len(words)
	if wordsLen == 0 {
		return []hashWeigth{}
	}

	result := make([]hashWeigth, wordsLen)
	for i, w := range words {
		hash := hasher(w.Text())
		result[i] = hashWeigth{hash, w.Weight()}
	}
	return result
}

func calWeights(hashes []hashWeigth) [64]float64 {
	var weights [64]float64
	for _, v := range hashes {
		for i := uint8(0); i < 64; i++ {
			weight := v.Weight
			if (1 << i & v.Hash) == 0 {
				weight *= -1
			}
			weights[i] += weight
		}
	}
	return weights
}

func fingerprint(weights [64]float64) uint64 {
	var f uint64
	for i := uint8(0); i < 64; i++ {
		if weights[i] >= 0.0 {
			f |= (1 << i)
		}
	}
	return f
}
