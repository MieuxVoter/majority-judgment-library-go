package judgment

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestProposalTally_RegradeJudgments(t *testing.T) {
	proposalTally := ProposalTally{Tally: []uint64{1, 2, 3, 4, 5, 6, 7}}
	err := proposalTally.RegradeJudgments(0, 6)
	assert.NoError(t, err, "Regrading should succeed")
	assert.Equal(t, uint64(0), proposalTally.Tally[0])
	assert.Equal(t, uint64(2), proposalTally.Tally[1])
	assert.Equal(t, uint64(3), proposalTally.Tally[2])
	assert.Equal(t, uint64(4), proposalTally.Tally[3])
	assert.Equal(t, uint64(5), proposalTally.Tally[4])
	assert.Equal(t, uint64(6), proposalTally.Tally[5])
	assert.Equal(t, uint64(8), proposalTally.Tally[6])
}

func TestProposalTally_RegradeJudgments_Failure1(t *testing.T) {
	proposalTally := ProposalTally{Tally: []uint64{1, 2, 3, 4, 5, 6, 7}}
	err := proposalTally.RegradeJudgments(0, 60)
	assert.Error(t, err, "Regrading should fail")
	assert.Equal(t, uint64(1), proposalTally.Tally[0])
	assert.Equal(t, uint64(2), proposalTally.Tally[1])
	assert.Equal(t, uint64(3), proposalTally.Tally[2])
	assert.Equal(t, uint64(4), proposalTally.Tally[3])
	assert.Equal(t, uint64(5), proposalTally.Tally[4])
	assert.Equal(t, uint64(6), proposalTally.Tally[5])
	assert.Equal(t, uint64(7), proposalTally.Tally[6])
}

func TestProposalTally_RegradeJudgments_Failure2(t *testing.T) {
	proposalTally := ProposalTally{Tally: []uint64{1, 2, 3, 4, 5, 6, 7}}
	err := proposalTally.RegradeJudgments(60, 0)
	assert.Error(t, err, "Regrading should fail")
	assert.Equal(t, uint64(1), proposalTally.Tally[0])
	assert.Equal(t, uint64(2), proposalTally.Tally[1])
	assert.Equal(t, uint64(3), proposalTally.Tally[2])
	assert.Equal(t, uint64(4), proposalTally.Tally[3])
	assert.Equal(t, uint64(5), proposalTally.Tally[4])
	assert.Equal(t, uint64(6), proposalTally.Tally[5])
	assert.Equal(t, uint64(7), proposalTally.Tally[6])
}
