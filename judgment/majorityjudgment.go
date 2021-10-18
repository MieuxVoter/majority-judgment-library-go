package judgment

import (
	"fmt"
	"sort"
)

// MajorityJudgment is one of the deliberators ; it implements DeliberatorInterface.
type MajorityJudgment struct {
	favorContestation bool // strategy for evenness of judgments ; defaults to true
}

// Deliberate is part of the DeliberatorInterface
func (mj *MajorityJudgment) Deliberate(tally *PollTally) (_ *PollResult, err error) {

	amountOfProposals := len(tally.Proposals)
	if 0 == amountOfProposals {
		return &PollResult{Proposals: []*ProposalResult{}}, nil
	}

	amountOfGrades := len(tally.Proposals[0].Tally)
	for _, proposalTally := range tally.Proposals {
		if amountOfGrades != len(proposalTally.Tally) {
			return nil, fmt.Errorf("mishaped tally: " +
				"some proposals hold more grades than others ; " +
				"please provide tallies of the same shape")
		}
	}

	maximumAmountOfJudgments := uint64(0)
	for _, proposalTally := range tally.Proposals {
		amountOfJudgments := proposalTally.CountJudgments()
		if amountOfJudgments > maximumAmountOfJudgments {
			maximumAmountOfJudgments = amountOfJudgments
		}
	}

	amountOfJudges := tally.AmountOfJudges
	if 0 == amountOfJudges {
		amountOfJudges = tally.GuessAmountOfJudges()
	}
	if amountOfJudges < maximumAmountOfJudgments {
		return nil, fmt.Errorf("incoherent tally: " +
			"some proposals hold more judgments than the specified amount of judges ; " +
			"perhaps you forgot to set PollTally.AmountOfJudges " +
			"or to call PollTally.GuessAmountOfJudges()")
	}

	for proposalIndex, proposalTally := range tally.Proposals {
		if amountOfJudges != proposalTally.CountJudgments() {
			return nil, fmt.Errorf("unbalanced tally: "+
				"a proposal (#%d) holds less judgments than there are judges ; "+
				"use one of the PollTally.Balance() methods first", proposalIndex)
		}
	}

	proposalsResults := make(ProposalsResults, 0, 16)
	proposalsResultsSorted := make(ProposalsResults, 0, 16)
	for proposalIndex, proposalTally := range tally.Proposals {
		score, scoreErr := mj.ComputeScore(proposalTally, true)
		if nil != scoreErr {
			return nil, scoreErr
		}
		proposalResult := &ProposalResult{
			Index:    proposalIndex,
			Score:    score,
			Analysis: proposalTally.Analyze(),
			Rank:     0, // we set it below after the sort
		}
		proposalsResults = append(proposalsResults, proposalResult)
		proposalsResultsSorted = append(proposalsResultsSorted, proposalResult)
	}

	sort.Sort(sort.Reverse(proposalsResultsSorted))

	// Rule: Multiple Proposals may have the same Rank in case of perfect equality.
	previousScore := ""
	for proposalIndex, proposalResult := range proposalsResultsSorted {
		rank := proposalIndex + 1
		if (previousScore == proposalResult.Score) && (proposalIndex > 0) {
			rank = proposalsResultsSorted[proposalIndex-1].Rank
		}
		proposalResult.Rank = rank
		previousScore = proposalResult.Score
	}

	result := &PollResult{
		Proposals:       proposalsResults,
		ProposalsSorted: proposalsResultsSorted,
	}

	return result, nil
}

// ComputeScore is the heart of our MajorityJudgment Deliberator.
// Not sure it should be exported, though.
// See docs/score-calculus-flowchart.png
func (mj *MajorityJudgment) ComputeScore(tally *ProposalTally, favorContestation bool) (_ string, err error) {
	score := ""

	analysis := &ProposalAnalysis{}
	amountOfGrades := tally.CountAvailableGrades()
	amountOfJudgments := tally.CountJudgments()
	amountOfDigitsForGrade := countDigitsUint8(amountOfGrades)
	amountOfDigitsForAdhesionScore := countDigitsUint64(amountOfJudgments * 2)

	amountOfJudgmentsInt := int(amountOfJudgments)
	if amountOfJudgmentsInt < 0 {
		return "", fmt.Errorf("too many judgments ; see branch|fork using math/big")
	}

	mutatedTally := tally.Copy()
	for i := uint8(0); i < amountOfGrades; i++ {
		analysis.Run(mutatedTally, favorContestation)
		score += fmt.Sprintf("%0"+fmt.Sprintf("%d", amountOfDigitsForGrade)+"d", analysis.MedianGrade)
		//adhesionScore := amountOfJudgments) + analysis.SecondGroupSize * analysis.SecondGroupSign
		adhesionScore := amountOfJudgments
		if analysis.SecondGroupSign > 0 {
			adhesionScore = adhesionScore + analysis.SecondGroupSize
		} else if analysis.SecondGroupSign < 0 {
			adhesionScore = adhesionScore - analysis.SecondGroupSize
		}
		score += fmt.Sprintf("%0"+fmt.Sprintf("%d", amountOfDigitsForAdhesionScore)+"d", adhesionScore)
		regradingErr := mutatedTally.RegradeJudgments(analysis.MedianGrade, analysis.SecondMedianGrade)
		if nil != regradingErr {
			return "", regradingErr // 仕方がない – C'est la vie ! (see issue #4)
		}
	}

	return score, nil
}

func countDigitsUint8(i uint8) (count uint8) {
	for i > 0 {
		i = i / 10 // euclidean division
		count++
	}

	return
}

func countDigitsUint64(i uint64) (count uint8) {
	for i > 0 {
		i = i / 10 // Euclid wuz hear
		count++
	}

	return
}
