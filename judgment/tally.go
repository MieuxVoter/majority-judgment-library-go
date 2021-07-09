package judgment

import "fmt"

type PollTally struct {
	AmountOfJudges uint64           // Helps balancing tallies using default judgments.
	Proposals      []*ProposalTally // Tallies of each proposal.  Its order is preserved in the result.
}

// Mutates the PollTally
func (pollTally *PollTally) BalanceWithStaticDefault(defaultGrade uint8) (err error) {
	for _, proposalTally := range pollTally.Proposals {
		proposalErr := proposalTally.FillWithStaticDefault(pollTally.AmountOfJudges, defaultGrade)
		if proposalErr != nil {
			return proposalErr
		}
	}
	return nil
}

// Mutates the PollTally
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

type ProposalTally struct {
	Tally []uint64 // Amount of judgments received for each grade, from "worst" grade to "best" grade.
}

func (proposalTally *ProposalTally) Analyze() (_ *ProposalAnalysis) {
	analysis := &ProposalAnalysis{}
	analysis.Run(proposalTally, true)
	return analysis
}

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

func (proposalTally *ProposalTally) CountJudgments() (_ uint64) {
	amountOfJudgments := uint64(0)
	for _, gradeTally := range proposalTally.Tally {
		amountOfJudgments += gradeTally
	}
	return amountOfJudgments
}

func (proposalTally *ProposalTally) CountAvailableGrades() (_ uint8) {
	return uint8(len(proposalTally.Tally))
}

// Mutates the proposalTally.
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

// Mutates the proposalTally
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

// Mutates the proposalTally
func (proposalTally *ProposalTally) FillWithMedianDefault(upToAmount uint64) (err error) {
	analysis := proposalTally.Analyze()
	fillErr := proposalTally.FillWithStaticDefault(upToAmount, analysis.MedianGrade)
	if fillErr != nil {
		return fillErr
	}
	return nil
}
