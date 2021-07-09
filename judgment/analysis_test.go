package judgment

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestProposalAnalysis_Run(t *testing.T) {
	type expectations struct {
		TotalSize              uint64
		MedianGrade            uint8
		MedianGroupSize        uint64
		SecondMedianGrade      uint8
		SecondGroupSize        uint64
		SecondGroupSign        int
		AdhesionGroupGrade     uint8
		AdhesionGroupSize      uint64
		ContestationGroupGrade uint8
		ContestationGroupSize  uint64
	}
	type args struct {
		proposalTally     *ProposalTally
		favorContestation bool
	}
	tests := []struct {
		name         string
		expectations expectations
		args         args
	}{
		{
			name: "All zeroes",
			args: args{
				proposalTally:     &ProposalTally{Tally: []uint64{0, 0, 0, 0, 0, 0, 0}},
				favorContestation: true,
			},
			expectations: expectations{
				MedianGrade:            0,
				MedianGroupSize:        0,
				SecondMedianGrade:      0,
				SecondGroupSize:        0,
				SecondGroupSign:        0,
				AdhesionGroupGrade:     0,
				AdhesionGroupSize:      0,
				ContestationGroupGrade: 0,
				ContestationGroupSize:  0,
			},
		},
		{
			name: "Single Grade (absurd! perhaps we should return err instead)",
			args: args{
				proposalTally:     &ProposalTally{Tally: []uint64{777}},
				favorContestation: true,
			},
			expectations: expectations{
				MedianGrade:            0,
				MedianGroupSize:        777,
				SecondMedianGrade:      0,
				SecondGroupSize:        0,
				SecondGroupSign:        0,
				AdhesionGroupGrade:     0,
				AdhesionGroupSize:      0,
				ContestationGroupGrade: 0,
				ContestationGroupSize:  0,
			},
		},
		{
			name: "Two grades (approbation poll)",
			args: args{
				proposalTally:     &ProposalTally{Tally: []uint64{421, 124}},
				favorContestation: true,
			},
			expectations: expectations{
				MedianGrade:            0,
				MedianGroupSize:        421,
				SecondMedianGrade:      1,
				SecondGroupSize:        124,
				SecondGroupSign:        1,
				AdhesionGroupGrade:     1,
				AdhesionGroupSize:      124,
				ContestationGroupGrade: 0,
				ContestationGroupSize:  0,
			},
		},
		{
			name: "Single judgment",
			args: args{
				proposalTally:     &ProposalTally{Tally: []uint64{0, 0, 0, 0, 1, 0, 0}},
				favorContestation: true,
			},
			expectations: expectations{
				MedianGrade:            4,
				MedianGroupSize:        1,
				SecondMedianGrade:      0,
				SecondGroupSize:        0,
				SecondGroupSign:        0,
				AdhesionGroupGrade:     0,
				AdhesionGroupSize:      0,
				ContestationGroupGrade: 0,
				ContestationGroupSize:  0,
			},
		},
		{
			name: "All ones",
			args: args{
				proposalTally:     &ProposalTally{Tally: []uint64{1, 1, 1, 1, 1, 1, 1}},
				favorContestation: true,
			},
			expectations: expectations{
				MedianGrade:            3,
				MedianGroupSize:        1,
				SecondMedianGrade:      2,
				SecondGroupSize:        3,
				SecondGroupSign:        -1,
				AdhesionGroupGrade:     4,
				AdhesionGroupSize:      3,
				ContestationGroupGrade: 2,
				ContestationGroupSize:  3,
			},
		},
		{
			name: "All ones, favoring adhesion",
			args: args{
				proposalTally:     &ProposalTally{Tally: []uint64{1, 1, 1, 1, 1, 1, 1}},
				favorContestation: false,
			},
			expectations: expectations{
				MedianGrade:            3,
				MedianGroupSize:        1,
				SecondMedianGrade:      4,
				SecondGroupSize:        3,
				SecondGroupSign:        1,
				AdhesionGroupGrade:     4,
				AdhesionGroupSize:      3,
				ContestationGroupGrade: 2,
				ContestationGroupSize:  3,
			},
		},
		{
			name: "All ones (even total)",
			args: args{
				proposalTally:     &ProposalTally{Tally: []uint64{1, 1, 1, 1, 1, 1}},
				favorContestation: true,
			},
			expectations: expectations{
				MedianGrade:            2,
				MedianGroupSize:        1,
				SecondMedianGrade:      3,
				SecondGroupSize:        3,
				SecondGroupSign:        1,
				AdhesionGroupGrade:     3,
				AdhesionGroupSize:      3,
				ContestationGroupGrade: 1,
				ContestationGroupSize:  2,
			},
		},
		{
			name: "All ones (even total), favoring adhesion",
			args: args{
				proposalTally:     &ProposalTally{Tally: []uint64{1, 1, 1, 1, 1, 1}},
				favorContestation: false,
			},
			expectations: expectations{
				MedianGrade:            3,
				MedianGroupSize:        1,
				SecondMedianGrade:      2,
				SecondGroupSize:        3,
				SecondGroupSign:        -1,
				AdhesionGroupGrade:     4,
				AdhesionGroupSize:      2,
				ContestationGroupGrade: 2,
				ContestationGroupSize:  3,
			},
		},
		{
			name: "Basic usage 1",
			args: args{
				proposalTally:     &ProposalTally{Tally: []uint64{3, 2, 3, 1, 3}},
				favorContestation: true,
			},
			expectations: expectations{
				MedianGrade:            2,
				MedianGroupSize:        3,
				SecondMedianGrade:      1,
				SecondGroupSize:        5,
				SecondGroupSign:        -1,
				AdhesionGroupGrade:     3,
				AdhesionGroupSize:      4,
				ContestationGroupGrade: 1,
				ContestationGroupSize:  5,
			},
		},
		{
			name: "Basic usage with zeroes 1",
			args: args{
				proposalTally:     &ProposalTally{Tally: []uint64{3, 2, 0, 0, 5}},
				favorContestation: true,
			},
			expectations: expectations{
				MedianGrade:            1,
				MedianGroupSize:        2,
				SecondMedianGrade:      4,
				SecondGroupSize:        5,
				SecondGroupSign:        1,
				AdhesionGroupGrade:     4,
				AdhesionGroupSize:      5,
				ContestationGroupGrade: 0,
				ContestationGroupSize:  3,
			},
		},
		{
			name: "Basic usage with zeroes 2",
			args: args{
				proposalTally:     &ProposalTally{Tally: []uint64{2, 0, 0, 0, 0, 0, 5}},
				favorContestation: true,
			},
			expectations: expectations{
				MedianGrade:            6,
				MedianGroupSize:        5,
				SecondMedianGrade:      0,
				SecondGroupSize:        2,
				SecondGroupSign:        -1,
				AdhesionGroupGrade:     0,
				AdhesionGroupSize:      0,
				ContestationGroupGrade: 0,
				ContestationGroupSize:  2,
			},
		},
		{
			name: "Basic usage with zeroes 3",
			args: args{
				proposalTally:     &ProposalTally{Tally: []uint64{20, 0, 0, 0, 0, 0, 5}},
				favorContestation: true,
			},
			expectations: expectations{
				MedianGrade:            0,
				MedianGroupSize:        20,
				SecondMedianGrade:      6,
				SecondGroupSize:        5,
				SecondGroupSign:        1,
				AdhesionGroupGrade:     6,
				AdhesionGroupSize:      5,
				ContestationGroupGrade: 0,
				ContestationGroupSize:  0,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			expected := &ProposalAnalysis{
				TotalSize:              tt.expectations.TotalSize,
				MedianGrade:            tt.expectations.MedianGrade,
				MedianGroupSize:        tt.expectations.MedianGroupSize,
				SecondMedianGrade:      tt.expectations.SecondMedianGrade,
				SecondGroupSize:        tt.expectations.SecondGroupSize,
				SecondGroupSign:        tt.expectations.SecondGroupSign,
				AdhesionGroupGrade:     tt.expectations.AdhesionGroupGrade,
				AdhesionGroupSize:      tt.expectations.AdhesionGroupSize,
				ContestationGroupGrade: tt.expectations.ContestationGroupGrade,
				ContestationGroupSize:  tt.expectations.ContestationGroupSize,
			}
			analysis := &ProposalAnalysis{}
			analysis.Run(tt.args.proposalTally, tt.args.favorContestation)
			assert.Equal(t, expected.MedianGrade, analysis.MedianGrade, "MedianGrade")
			assert.Equal(t, expected.MedianGroupSize, analysis.MedianGroupSize, "MedianGroupSize")
			assert.Equal(t, expected.SecondMedianGrade, analysis.SecondMedianGrade, "SecondMedianGrade")
			assert.Equal(t, expected.SecondGroupSize, analysis.SecondGroupSize, "SecondGroupSize")
			assert.Equal(t, expected.SecondGroupSign, analysis.SecondGroupSign, "SecondGroupSign")
			assert.Equal(t, expected.AdhesionGroupGrade, analysis.AdhesionGroupGrade, "AdhesionGroupGrade")
			assert.Equal(t, expected.AdhesionGroupSize, analysis.AdhesionGroupSize, "AdhesionGroupSize")
			assert.Equal(t, expected.ContestationGroupGrade, analysis.ContestationGroupGrade, "ContestationGroupGrade")
			assert.Equal(t, expected.ContestationGroupSize, analysis.ContestationGroupSize, "ContestationGroupSize")

		})
	}
}
