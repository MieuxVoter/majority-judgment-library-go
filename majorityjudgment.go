package judgment

import (
	"fmt"
	"sort"
)

type DeliberatorInterface interface {
	Deliberate(tally *PollTally) (result *PollResult, err error)
}

type MajorityJudgment struct {
	favorContestation bool // strategy for evenness of judgments ; defaults to true
}

func (mj *MajorityJudgment) Deliberate(tally *PollTally) (_ *PollResult, err error) {

	// TODO: preset balancers : static, median, normalization
	//// Consider missing judgments as TO_REJECT judgments
	//// One reason why we need many algorithms, and options on each
	//for _, candidateTally := range pollTally.Candidates {
	//	missing := pollTally.MaxJudgmentsAmount - candidateTally.JudgmentsAmount
	//	if 0 < missing {
	//		candidateTally.JudgmentsAmount = pollTally.MaxJudgmentsAmount
	//		candidateTally.Grades[0].Amount += missing
	//		//println("Added the missing TO REJECT judgments", missing)
	//	}
	//}

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

	amountOfJudgments := tally.Proposals[0].CountJudgments()
	for _, proposalTally := range tally.Proposals {
		if amountOfJudgments != proposalTally.CountJudgments() {
			return nil, fmt.Errorf("unbalanced tally: " +
				"some proposals hold more judgments than others ; " +
				"use one of the tally balancers or make your own")
		}
	}

	proposalsResults := make(ProposalsResults, 0, 16)
	proposalsResultsSorted := make(ProposalsResults, 0, 16)
	for _, proposalTally := range tally.Proposals {
		score, scoreErr := mj.ComputeScore(proposalTally, true)
		if nil != scoreErr {
			return nil, scoreErr
		}
		proposalResult := &ProposalResult{
			Rank:  0, // we set it below after the sort
			Score: score,
			//Analysis: proposalTally.Analyze(),
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
		Proposals: proposalsResults,
	}

	return result, nil
}

// See docs/score-calculus-flowchart.png
func (mj *MajorityJudgment) ComputeScore(tally *ProposalTally, favorContestation bool) (_ string, err error) {
	score := ""

	analysis := &ProposalAnalysis{}
	amountOfJudgments := tally.CountJudgments()
	amountOfDigitsForGrade := 3          // fixme: compute
	amountOfDigitsForAdhesionScore := 12 // fixme: compute

	mutatedTally := tally.Copy()
	for range tally.Tally { // loop once per grade
		analysis.Run(mutatedTally, favorContestation)
		score += fmt.Sprintf("%0"+fmt.Sprintf("%d", amountOfDigitsForGrade)+"d", analysis.MedianGrade)
		// fixme: BAD → uint64 to int conversion — either move to int everywhere, or use whatever bigint Go has
		score += fmt.Sprintf("%0"+fmt.Sprintf("%d", amountOfDigitsForAdhesionScore)+"d", int(amountOfJudgments)+int(analysis.SecondGroupSize)*analysis.SecondGroupSign)
		err := mutatedTally.RegradeJudgments(analysis.MedianGrade, analysis.SecondMedianGrade)
		if nil != err {
			return "", err
		}
	}

	return score, nil
}
