package judgment

import (
	"fmt"
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

func TestProposalTally_RegradeJudgmentsIntoSelf(t *testing.T) {
	proposalTally := ProposalTally{Tally: []uint64{1, 2, 3, 4, 5, 6, 7}}
	err := proposalTally.RegradeJudgments(2, 2)
	assert.NoError(t, err, "Regrading should succeed")
	assert.Equal(t, uint64(1), proposalTally.Tally[0])
	assert.Equal(t, uint64(2), proposalTally.Tally[1])
	assert.Equal(t, uint64(3), proposalTally.Tally[2])
	assert.Equal(t, uint64(4), proposalTally.Tally[3])
	assert.Equal(t, uint64(5), proposalTally.Tally[4])
	assert.Equal(t, uint64(6), proposalTally.Tally[5])
	assert.Equal(t, uint64(7), proposalTally.Tally[6])
}

func TestProposalTally_RegradeJudgments_Failure1(t *testing.T) {
	proposalTally := ProposalTally{Tally: []uint64{1, 2, 3, 4, 5, 6, 7}}
	err := proposalTally.RegradeJudgments(0, uint8(len(proposalTally.Tally)))
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

func TestProposalTally_FillWithStaticDefault(t *testing.T) {
	proposalTally := ProposalTally{Tally: []uint64{1, 2, 3, 4, 5, 6, 7}}
	expectedTally := ProposalTally{Tally: []uint64{3, 2, 3, 4, 5, 6, 7}}
	err := proposalTally.FillWithStaticDefault(30, 0)
	assert.NoError(t, err, "Filling should succeed")
	for i := 0; i < 7; i++ {
		assert.Equal(t, expectedTally.Tally[i], proposalTally.Tally[i], fmt.Sprintf("Grade #%d", i))
	}
}

func TestProposalTally_FillWithStaticDefaultNotZero(t *testing.T) {
	proposalTally := ProposalTally{Tally: []uint64{0, 1, 0, 1, 0, 0, 1}}
	expectedTally := ProposalTally{Tally: []uint64{0, 1, 1, 1, 0, 0, 1}}
	err := proposalTally.FillWithStaticDefault(4, 2)
	assert.NoError(t, err, "Filling should succeed")
	for i := 0; i < 7; i++ {
		assert.Equal(t, expectedTally.Tally[i], proposalTally.Tally[i], fmt.Sprintf("Grade #%d", i))
	}
}

func TestProposalTally_FillWithStaticDefaultNoop(t *testing.T) {
	proposalTally := ProposalTally{Tally: []uint64{0, 1, 2, 1, 0, 0, 1}}
	expectedTally := ProposalTally{Tally: []uint64{0, 1, 2, 1, 0, 0, 1}}
	err := proposalTally.FillWithStaticDefault(5, 5)
	assert.NoError(t, err, "Filling should succeed")
	for i := 0; i < 7; i++ {
		assert.Equal(t, expectedTally.Tally[i], proposalTally.Tally[i], fmt.Sprintf("Grade #%d", i))
	}
}

func TestProposalTally_FillWithStaticDefaultFailureGradeTooHigh(t *testing.T) {
	proposalTally := ProposalTally{Tally: []uint64{0, 1, 0, 1, 2, 3, 4}}
	expectedTally := ProposalTally{Tally: []uint64{0, 1, 0, 1, 2, 3, 4}}
	err := proposalTally.FillWithStaticDefault(12, 200)
	assert.Error(t, err, "Filling should fail")
	for i := 0; i < 7; i++ {
		assert.Equal(t, expectedTally.Tally[i], proposalTally.Tally[i], fmt.Sprintf("Grade #%d", i))
	}
}

func TestProposalTally_FillWithStaticDefaultFailureAmountToLow(t *testing.T) {
	proposalTally := ProposalTally{Tally: []uint64{0, 1, 0, 1, 2, 3, 4}}
	expectedTally := ProposalTally{Tally: []uint64{0, 1, 0, 1, 2, 3, 4}}
	err := proposalTally.FillWithStaticDefault(5, 0)
	assert.Error(t, err, "Filling should fail")
	for i := 0; i < 7; i++ {
		assert.Equal(t, expectedTally.Tally[i], proposalTally.Tally[i], fmt.Sprintf("Grade #%d", i))
	}
}

func TestProposalTally_FillWithStaticDefaultSuccesses(t *testing.T) {
	type test struct {
		name           string
		amountOfJudges uint64
		defaultGrade   uint8
		input          ProposalTally
		expected       ProposalTally
	}
	tests := []test{
		{
			name:           "Basic usage",
			amountOfJudges: 30,
			defaultGrade:   0,
			input:          ProposalTally{Tally: []uint64{1, 2, 3, 4, 5, 6, 7}},
			expected:       ProposalTally{Tally: []uint64{3, 2, 3, 4, 5, 6, 7}},
		},
		// …
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			proposalTally := tt.input
			expectedTally := tt.expected
			err := proposalTally.FillWithStaticDefault(tt.amountOfJudges, tt.defaultGrade)
			assert.NoError(t, err, "Filling should succeed")
			for i := 0; i < len(proposalTally.Tally); i++ {
				assert.Equal(t, expectedTally.Tally[i], proposalTally.Tally[i], fmt.Sprintf("Grade #%d", i))
			}
		})
	}
}

func TestProposalTally_FillWithMedianDefaultSuccesses(t *testing.T) {
	type test struct {
		name           string
		amountOfJudges uint64
		input          ProposalTally
		expected       ProposalTally
	}
	tests := []test{
		{
			name:           "Basic usage",
			amountOfJudges: 30,
			input:          ProposalTally{Tally: []uint64{1, 2, 3, 4, 5, 6, 7}},
			expected:       ProposalTally{Tally: []uint64{1, 2, 3, 4, 7, 6, 7}},
		},
		{
			name:           "Lots of zeroes",
			amountOfJudges: 5,
			input:          ProposalTally{Tally: []uint64{0, 0, 0, 1, 0, 0, 1}},
			expected:       ProposalTally{Tally: []uint64{0, 0, 0, 4, 0, 0, 1}},
		},
		// …
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			proposalTally := tt.input
			expectedTally := tt.expected
			err := proposalTally.FillWithMedianDefault(tt.amountOfJudges)
			assert.NoError(t, err, "Filling should succeed")
			for i := 0; i < len(proposalTally.Tally); i++ {
				assert.Equal(t, expectedTally.Tally[i], proposalTally.Tally[i], fmt.Sprintf("Grade #%d", i))
			}
		})
	}
}

func TestProposalTally_FillWithMedianDefaultFailureAmountToLow(t *testing.T) {
	proposalTally := ProposalTally{Tally: []uint64{0, 1, 0, 1, 2, 3, 4}}
	expectedTally := ProposalTally{Tally: []uint64{0, 1, 0, 1, 2, 3, 4}}
	err := proposalTally.FillWithMedianDefault(5)
	assert.Error(t, err, "Filling should fail")
	for i := 0; i < 7; i++ {
		assert.Equal(t, expectedTally.Tally[i], proposalTally.Tally[i], fmt.Sprintf("Grade #%d", i))
	}
}
