package judgment

import (
	"fmt"
	"github.com/stretchr/testify/assert"
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
	fmt.Printf("Result: %v\n", result)
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
