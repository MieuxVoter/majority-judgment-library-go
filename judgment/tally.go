package judgment

import (
	"fmt"
)

// PollTally describes the amount of judgments received by each proposal on each grade.
type PollTally struct {
	AmountOfJudges uint64           `json:"amountOfJudges"` // Helps balancing tallies using default judgments.
	Proposals      []*ProposalTally `json:"proposals"`      // Tallies of each proposal.  Its order is preserved in the result.
}

// GuessAmountOfJudges also mutates the PollTally by filling the AmountOfJudges property
func (pollTally *PollTally) GuessAmountOfJudges() (_ uint64) {
	pollTally.AmountOfJudges = 0
	for _, proposalTally := range pollTally.Proposals {
		amountOfJudges := proposalTally.CountJudgments()
		if pollTally.AmountOfJudges < amountOfJudges {
			pollTally.AmountOfJudges = amountOfJudges
		}
	}
	return pollTally.AmountOfJudges
}

// BalanceWithStaticDefault mutates the PollTally
func (pollTally *PollTally) BalanceWithStaticDefault(defaultGrade uint8) (err error) {
	for _, proposalTally := range pollTally.Proposals {
		proposalErr := proposalTally.FillWithStaticDefault(pollTally.AmountOfJudges, defaultGrade)
		if proposalErr != nil {
			return proposalErr
		}
	}
	return nil
}

// BalanceWithMedianDefault mutates the PollTally
func (pollTally *PollTally) BalanceWithMedianDefault() (err error) {
	for _, proposalTally := range pollTally.Proposals {
		proposalErr := proposalTally.FillWithMedianDefault(pollTally.AmountOfJudges)
		if proposalErr != nil {
			return proposalErr
		}
	}
	return nil
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

// ProposalTally holds the amount of judgments received per Grade for a single Proposal
type ProposalTally struct {
	Tally []uint64 `json:"tally"` // Amount of judgments received for each grade, from "worst" grade to "best" grade.
}

// Analyze a ProposalTally and return its ProposalAnalysis
func (proposalTally *ProposalTally) Analyze() (_ *ProposalAnalysis) {
	analysis := &ProposalAnalysis{}
	analysis.Run(proposalTally, true)
	return analysis
}

// Copy a ProposalTally (deeply)
func (proposalTally *ProposalTally) Copy() (_ *ProposalTally) {
	// There might exist an elegant one-liner to copy a slice of uint64
	intTally := make([]uint64, 0, 8)
	for _, gradeTally := range proposalTally.Tally {
		intTally = append(intTally, gradeTally) // uint64 is copied, hopefully
	}
	return &ProposalTally{
		Tally: intTally,
	}
}

// CountJudgments tallies the received judgments by a Proposal
func (proposalTally *ProposalTally) CountJudgments() (_ uint64) {
	amountOfJudgments := uint64(0)
	for _, gradeTally := range proposalTally.Tally {
		amountOfJudgments += gradeTally
	}
	return amountOfJudgments
}

// CountAvailableGrades returns the amount of available grades in the poll (usually 7 or so).
func (proposalTally *ProposalTally) CountAvailableGrades() (_ uint8) {
	return uint8(len(proposalTally.Tally))
}

// RegradeJudgments mutates the proposalTally by moving judgments from one grade to another.
// Useful when computing the score ; perhaps this method should not be exported, though.
func (proposalTally *ProposalTally) RegradeJudgments(fromGrade uint8, intoGrade uint8) (err error) {

	if fromGrade == intoGrade {
		return nil
	}

	amountOfGrades := proposalTally.CountAvailableGrades()
	if fromGrade >= amountOfGrades {
		return fmt.Errorf("RegradeJudgments() fromGrade is too high")
	}
	if intoGrade >= amountOfGrades {
		return fmt.Errorf("RegradeJudgments() intoGrade is too high")
	}

	proposalTally.Tally[intoGrade] += proposalTally.Tally[fromGrade]
	proposalTally.Tally[fromGrade] = 0

	return nil
}

// FillWithStaticDefault mutates the proposalTally
func (proposalTally *ProposalTally) FillWithStaticDefault(upToAmount uint64, defaultGrade uint8) (err error) {
	// More silent integer casting awkwardness… ; we need to fix this
	missingAmount := int(upToAmount) - int(proposalTally.CountJudgments())
	if missingAmount < 0 {
		return fmt.Errorf("FillWithStaticDefault() upToAmount is lower than the actual amount of judgments")
	} else if missingAmount == 0 {
		return nil
	}

	if defaultGrade >= proposalTally.CountAvailableGrades() {
		return fmt.Errorf("FillWithStaticDefault() defaultGrade is higher than the amount of available grades")
	}

	proposalTally.Tally[defaultGrade] += uint64(missingAmount)

	return nil
}

// FillWithMedianDefault mutates the proposalTally
func (proposalTally *ProposalTally) FillWithMedianDefault(upToAmount uint64) (err error) {
	analysis := proposalTally.Analyze()
	fillErr := proposalTally.FillWithStaticDefault(upToAmount, analysis.MedianGrade)
	if fillErr != nil {
		return fillErr
	}
	return nil
}
