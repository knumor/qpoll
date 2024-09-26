package models

import (
	"sort"
	"sync"
)

// Word is a word in a word cloud.
type Word struct {
	Index  int
	Text   string
	Weight float64
	Freq   int
}

// WordCloud is a word cloud model.
type WordCloud struct {
	CommonPollData
	sync.RWMutex
	words map[string]int
}

// NewWordCloud creates a new word cloud.
func NewWordCloud(question string) *WordCloud {
	return &WordCloud{
		CommonPollData: initPoll(question, WordCloudPoll),
		words: make(map[string]int),
	}
}

// AddWord adds a word to the word cloud.
func (wc *WordCloud) AddWord(word string) {
	wc.Lock()
	wc.words[word]++
	wc.Unlock()
}

// GetWords returns the words in the word cloud.
func (wc *WordCloud) GetWords() []Word {
	var maxCount int

	wc.RLock()
	words := make([]Word, 0, len(wc.words))
	for word, count := range wc.words {
		if count > maxCount {
			maxCount = count
		}
		words = append(words, Word{
			Text:   word,
			Weight: float64(count),
			Freq:   count,
		})
	}
	wc.RUnlock()

	sort.SliceStable(words, func(i, j int) bool {
		return words[i].Text < words[j].Text
	})

	for i := range words {
		words[i].Index = i
		words[i].Weight /= float64(maxCount)
	}
	return words
}
