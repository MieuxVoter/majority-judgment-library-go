package judgment

import (
	"github.com/stretchr/testify/assert"
	"math"
	"testing"
)

func TestBunchOfWorkingScenarios(t *testing.T) {
	type rawRanks []int
	type rawTally []uint64
	type test struct {
		Name           string
		AmountOfJudges uint64
		Proposals      []rawTally
		Ranks          rawRanks
	}
	tests := []test{
		{
			Name:           "Basic 01",
			AmountOfJudges: 3,
			Proposals: []rawTally{
				[]uint64{1, 1, 1},
				[]uint64{1, 2, 0},
				[]uint64{0, 2, 1},
			},
			Ranks: []int{2, 3, 1},
		},
		{
			Name:           "Billions of participants",
			AmountOfJudges: 20e9, // 20 billion
			Proposals: []rawTally{
				[]uint64{10e9, 10e9},
				[]uint64{9999999999, 10000000001},
			},
			Ranks: []int{2, 1},
		},
		// …
	}

	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			proposalsTallies := make([]*ProposalTally, 0, 10)
			for _, p := range tt.Proposals {
				proposalTally := &ProposalTally{Tally: p}
				proposalsTallies = append(proposalsTallies, proposalTally)
			}
			poll := &PollTally{
				AmountOfJudges: tt.AmountOfJudges,
				Proposals:      proposalsTallies,
			}
			deliberator := &MajorityJudgment{}
			result, err := deliberator.Deliberate(poll)
			assert.NoError(t, err, "Deliberation should succeed")
			for proposalResultIndex, proposalResult := range result.Proposals {
				assert.Equal(t, tt.Ranks[proposalResultIndex], proposalResult.Rank, "Rank of proposal")
			}

		})
	}
}

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

	assert.Equal(t, 1, result.ProposalsSorted[0].Rank, "Rank of sorted proposal A")
	assert.Equal(t, 2, result.ProposalsSorted[1].Rank, "Rank of sorted proposal 1")
	assert.Equal(t, 3, result.ProposalsSorted[2].Rank, "Rank of sorted proposal 2")
	assert.Equal(t, 4, result.ProposalsSorted[3].Rank, "Rank of sorted proposal 3")
	assert.Equal(t, 4, result.ProposalsSorted[4].Rank, "Rank of sorted proposal 4")

	assert.Equal(t, 1, result.ProposalsSorted[0].Index, "Index of sorted proposal A")
	assert.Equal(t, 2, result.ProposalsSorted[1].Index, "Index of sorted proposal 1")
	assert.Equal(t, 3, result.ProposalsSorted[2].Index, "Index of sorted proposal 2")
	assert.Equal(t, 0, result.ProposalsSorted[3].Index, "Index of sorted proposal 3")
	assert.Equal(t, 4, result.ProposalsSorted[4].Index, "Index of sorted proposal 4")

}

func TestGuessingAmountOfJudges(t *testing.T) {
	poll := &PollTally{
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

	assert.Equal(t, 1, result.ProposalsSorted[0].Rank, "Rank of sorted proposal A")
	assert.Equal(t, 2, result.ProposalsSorted[1].Rank, "Rank of sorted proposal 1")
	assert.Equal(t, 3, result.ProposalsSorted[2].Rank, "Rank of sorted proposal 2")
	assert.Equal(t, 4, result.ProposalsSorted[3].Rank, "Rank of sorted proposal 3")
	assert.Equal(t, 4, result.ProposalsSorted[4].Rank, "Rank of sorted proposal 4")

	assert.Equal(t, 1, result.ProposalsSorted[0].Index, "Index of sorted proposal A")
	assert.Equal(t, 2, result.ProposalsSorted[1].Index, "Index of sorted proposal 1")
	assert.Equal(t, 3, result.ProposalsSorted[2].Index, "Index of sorted proposal 2")
	assert.Equal(t, 0, result.ProposalsSorted[3].Index, "Index of sorted proposal 3")
	assert.Equal(t, 4, result.ProposalsSorted[4].Index, "Index of sorted proposal 4")

}

func TestNoProposals(t *testing.T) {
	poll := &PollTally{
		AmountOfJudges: 0,
		Proposals:      []*ProposalTally{},
	}
	deliberator := &MajorityJudgment{}
	result, err := deliberator.Deliberate(poll)
	assert.NoError(t, err, "Deliberation should succeed")
	assert.Len(t, result.Proposals, len(poll.Proposals), "There should be as many results as there are tallies.")
}

func TestIncoherentTally(t *testing.T) {
	poll := &PollTally{
		AmountOfJudges: 2, // lower than expected 8
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
