package models

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/sqids/sqids-go"
)

// PollType is the type of a poll.
type PollType int

// Supported Poll types.
const (
	MultipleChoicePoll PollType = iota
	WordCloudPoll
)

// Poll is an interface for a poll.
type Poll interface {
	ID() string
	Question() string
	Code() string
	CreatedAt() time.Time
	Type() PollType
}

// CommonPollData is a struct with common poll data
type CommonPollData struct {
	id        string
	idgen     *sqids.Sqids
	question  string
	createdAt time.Time
	polltype  PollType
}

func initPoll(question string, polltype PollType) CommonPollData {
	code := newPollCode()
	idgen := newIDGenerator()
	id, err := idgen.Encode([]uint64{code})
	if err != nil {
		panic(err)
	}
	cpd := CommonPollData{
		id:        id,
		idgen:     idgen,
		question:  question,
		createdAt: time.Now(),
		polltype:  polltype,
	}
	return cpd
}

// ID returns the id of the poll.
func (cpd CommonPollData) ID() string {
	if cpd.id == "" {
		panic("ID is empty")
	}
	return cpd.id
}

// Question returns the question of the poll.
func (cpd CommonPollData) Question() string {
	return cpd.question
}

// Code returns the code of the poll.
func (cpd CommonPollData) Code() string {
	return fmt.Sprintf("%d", cpd.idgen.Decode(cpd.id)[0])
}

// CreatedAt returns the creation time of the poll.
func (cpd CommonPollData) CreatedAt() time.Time {
	return cpd.createdAt
}

// Type returns the type of the poll.
func (cpd CommonPollData) Type() PollType {
	return cpd.polltype
}

func newPollCode() uint64 {
	return uint64(10000000+rand.Intn(90000000))
}

func newIDGenerator() *sqids.Sqids {
	s, _ := sqids.New()
	return s
}
