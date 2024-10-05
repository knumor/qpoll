package models

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"sort"
	"sync"
	"time"
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
func NewWordCloud(question, owner string) *WordCloud {
	return &WordCloud{
		CommonPollData: initPoll(question, owner, WordCloudPoll),
		words:          make(map[string]int),
	}
}

// WordCloudFromJSON makes a word cloud with a given id/code and data from JSON.
func WordCloudFromJSON(id string, data []byte) (*WordCloud, error) {
	var loadStruct struct {
		Question     string         `json:"question"`
		Words        map[string]int `json:"words"`
		Owner        string         `json:"owner"`
		CreatedAt    time.Time      `json:"createdAt"`
		PollType     PollType       `json:"polltype"`
		NumResponses int            `json:"numResponses"`
		NumVotes     int            `json:"numVotes"`
	}

	err := json.Unmarshal(data, &loadStruct)
	if err != nil {
		slog.Error("Failed to unmarshal word cloud", "error", err)
	}

	var wc WordCloud
	wc.id = id
	wc.idgen = newIDGenerator()
	wc.words = loadStruct.Words
	wc.owner = loadStruct.Owner
	wc.createdAt = loadStruct.CreatedAt
	wc.question = loadStruct.Question
	wc.polltype = loadStruct.PollType
	wc.numResponses = loadStruct.NumResponses
	wc.numVotes = loadStruct.NumVotes

	return &wc, nil
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

// MarshalJSON marshals the word cloud to JSON.
func (wc *WordCloud) MarshalJSON() ([]byte, error) {
	ts, err := wc.CreatedAt().MarshalJSON()
	words, err := json.Marshal(wc.words)
	if err != nil {
		return nil, err
	}
	return []byte(fmt.Sprintf(
		`{"id":"%s","question":"%s","createdAt":%s,"owner":"%s","polltype":%d,"numResponses":%d,"numVotes":%d,"words":%s}`,
		wc.ID(),
		wc.Question(),
		string(ts),
		wc.owner,
		wc.Type(),
		wc.numResponses,
		wc.numVotes,
		string(words),
	)), nil
}
