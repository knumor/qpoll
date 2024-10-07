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
	MultipleChoicePoll PollType = iota + 1
	WordCloudPoll
)

func (pt PollType) String() string {
	switch pt {
	case MultipleChoicePoll:
		return "Multiple Choice"
	case WordCloudPoll:
		return "Word Cloud"
	default:
		return fmt.Sprintf("PollType(%d)", pt)
	}
}

// Poll is an interface for a poll.
type Poll interface {
	ID() string
	Question() string
	Code() string
	CreatedAt() time.Time
	Owner() string
	Type() PollType
	ResponseCount() int
	VoteCount() int
	AddVote(int)
	MarshalJSON() ([]byte, error)
}

// CommonPollData is a struct with common poll data
type CommonPollData struct {
	id           string
	idgen        *sqids.Sqids
	question     string
	owner        string
	createdAt    time.Time
	polltype     PollType
	numResponses int
	numVotes     int
}

func initPoll(question, owner string, polltype PollType) CommonPollData {
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
		owner:     owner,
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

// ResponseCount returns the number of responses to the poll.
func (cpd CommonPollData) ResponseCount() int {
	return cpd.numResponses
}

// VoteCount returns the number of clients that have accessed the poll.
func (cpd CommonPollData) VoteCount() int {
	return cpd.numVotes
}

// AddVote increments the number of votes and responses.
func (cpd *CommonPollData) AddVote(responseCount int) {
	cpd.numVotes++
	cpd.numResponses += responseCount
}

func (cpt *CommonPollData) Owner() string {
	return cpt.owner
}

func newPollCode() uint64 {
	return uint64(10000000 + rand.Intn(90000000))
}

func newIDGenerator() *sqids.Sqids {
	s, _ := sqids.New()
	return s
}
