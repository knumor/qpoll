package models

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"sync"
	"time"
)

// Option represents an option in a multiple choice poll.
type Option struct {
	Index int
	Text  string
	Count int
	Weight float64
}

// MultipleChoice is a multiple choice poll model.
type MultipleChoice struct {
	CommonPollData
	sync.RWMutex
	options []string
	counts  []int
}

// NewMultipleChoice creates a new multiple choice poll.
func NewMultipleChoice(question string, options []string) *MultipleChoice {
	return &MultipleChoice{
		CommonPollData: initPoll(question, MultipleChoicePoll),
		options:        options,
		counts:         make([]int, len(options)),
	}
}

// MultipleChoiceFromJSON makes a multiple choice poll with a given id/code and data from JSON.
func MultipleChoiceFromJSON(id string, data []byte) (*MultipleChoice, error) {
	var loadStruct struct {
		Question     string    `json:"question"`
		Options      []string  `json:"options"`
		Counts       []int     `json:"counts"`
		CreatedAt    time.Time `json:"createdAt"`
		PollType     PollType  `json:"polltype"`
		NumResponses int       `json:"numResponses"`
		NumVotes     int       `json:"numVotes"`
	}

	err := json.Unmarshal(data, &loadStruct)
	if err != nil {
		slog.Error("Failed to unmarshal multiple choice poll", "error", err)
	}

	var mc MultipleChoice
	mc.id = id
	mc.idgen = newIDGenerator()
	mc.options = loadStruct.Options
	mc.counts = loadStruct.Counts
	mc.createdAt = loadStruct.CreatedAt
	mc.question = loadStruct.Question
	mc.polltype = loadStruct.PollType
	mc.numResponses = loadStruct.NumResponses
	mc.numVotes = loadStruct.NumVotes

	return &mc, nil
}

// AddOption adds an option to the poll.
func (mc *MultipleChoice) AddOption(option string) {
	mc.Lock()
	mc.options = append(mc.options, option)
	mc.counts = append(mc.counts, 0)
	mc.Unlock()
}

// AddVoteForOption adds a vote for an option.
func (mc *MultipleChoice) AddVoteForOption(index int) {
	mc.Lock()
	mc.counts[index]++
	mc.Unlock()
	mc.AddVote(1)
}

// GetOptions returns the options and their counts from the poll.
func (mc *MultipleChoice) GetOptions() []Option {
	mc.RLock()
	options := make([]Option, len(mc.options))
	sum := sumInts(mc.counts)
	for i := range mc.options {
		weight := 0.00
		if sum > 0 {
			weight = float64(mc.counts[i]) / float64(sum)
		}
		options[i] = Option{
			Index: i,
			Text:  mc.options[i],
			Count: mc.counts[i],
			Weight: weight,
		}
	}
	mc.RUnlock()

	return options
}

func sumInts(ints []int) int {
	sum := 0
	for _, i := range ints {
		sum += i
	}
	return sum
}

// MarshalJSON marshals the multiple choice poll to JSON.
func (mc *MultipleChoice) MarshalJSON() ([]byte, error) {
	ts, err := mc.CreatedAt().MarshalJSON()
	options, err := json.Marshal(mc.options)
	counts, err := json.Marshal(mc.counts)
	if err != nil {
		return nil, err
	}
	return []byte(fmt.Sprintf(
		`{"id":"%s","question":"%s","createdAt":%s,"polltype":%d,"numResponses":%d,"numVotes":%d,"options":%s,"counts":%s}`,
		mc.ID(),
		mc.Question(),
		string(ts),
		mc.Type(),
		mc.numResponses,
		mc.numVotes,
		string(options),
		string(counts),
	)), nil
}
