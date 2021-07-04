package judgment

import (
	"github.com/stretchr/testify/assert"
	"math"
	"testing"
)

func TestReadmeDemo(t *testing.T) {
	poll := &PollTally{
		AmountOfJudges: 10,
		Proposals: []*ProposalTally{
			{Tally: []uint64{2, 2, 2, 2, 2}},
			{Tally: []uint64{2, 1, 1, 1, 5}},
			{Tally: []uint64{2, 1, 1, 2, 4}},
			{Tally: []uint64{2, 1, 5, 0, 2}},
			{Tally: []uint64{2, 2, 2, 2, 2}},
		},
	}
	deliberator := &MajorityJudgment{}
	result, err := deliberator.Deliberate(poll)
	assert.NoError(t, err, "Deliberation should succeed")
	assert.Len(t, result.Proposals, len(poll.Proposals), "There should be as many results as there are tallies.")
	assert.Equal(t, 4, result.Proposals[0].Rank, "Rank of proposal A")
	assert.Equal(t, 1, result.Proposals[1].Rank, "Rank of proposal B")
	assert.Equal(t, 2, result.Proposals[2].Rank, "Rank of proposal C")
	assert.Equal(t, 3, result.Proposals[3].Rank, "Rank of proposal D")
	assert.Equal(t, 4, result.Proposals[4].Rank, "Rank of proposal E")
}

func TestNoProposals(t *testing.T) {
	poll := &PollTally{
		AmountOfJudges: 10,
		Proposals:      []*ProposalTally{},
	}
	deliberator := &MajorityJudgment{}
	result, err := deliberator.Deliberate(poll)
	assert.NoError(t, err, "Deliberation should succeed")
	assert.Len(t, result.Proposals, len(poll.Proposals), "There should be as many results as there are tallies.")
}

func TestIncoherentTally(t *testing.T) {
	poll := &PollTally{
		AmountOfJudges: 2, // not 8 as it should
		Proposals: []*ProposalTally{
			{Tally: []uint64{4, 4}},
			{Tally: []uint64{2, 6}},
		},
	}
	deliberator := &MajorityJudgment{}
	result, err := deliberator.Deliberate(poll)
	assert.Error(t, err, "Deliberation should fail")
	assert.Nil(t, result, "Deliberation result should be nil")
}

func TestMishapedTally(t *testing.T) {
	poll := &PollTally{
		AmountOfJudges: 10,
		Proposals: []*ProposalTally{
			{Tally: []uint64{2, 2, 2, 2, 2}},
			{Tally: []uint64{2, 2, 2, 2}},
		},
	}
	deliberator := &MajorityJudgment{}
	result, err := deliberator.Deliberate(poll)
	assert.Error(t, err, "Deliberation should fail")
	assert.Nil(t, result, "Deliberation result should be nil")
}

func TestUnbalancedTally(t *testing.T) {
	poll := &PollTally{
		AmountOfJudges: 10,
		Proposals: []*ProposalTally{
			{Tally: []uint64{2, 2, 2, 2, 2}},
			{Tally: []uint64{2, 0, 0, 0, 2}},
		},
	}
	deliberator := &MajorityJudgment{}
	result, err := deliberator.Deliberate(poll)
	assert.Error(t, err, "Deliberation should fail")
	assert.Nil(t, result, "Deliberation result should be nil")
}

func TestExcessivelyBigTally(t *testing.T) {
	// math.MaxInt64 + 1 is a valid uint64 value, but we need to cast to int internally so it overflows
	poll := &PollTally{
		AmountOfJudges: math.MaxInt64 + 1,
		Proposals: []*ProposalTally{
			{Tally: []uint64{math.MaxInt64, 1}},
			{Tally: []uint64{math.MaxInt64, 1}},
		},
	}
	deliberator := &MajorityJudgment{}
	result, err := deliberator.Deliberate(poll)
	if assert.Error(t, err, "Deliberation should fail") {
		//println(err.Error())
	}
	assert.Nil(t, result, "Deliberation result should be nil")
}
