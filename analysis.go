package judgment

// Searching for good names
// ------------------------
// Rebuttal == Contestation == ???
// Rebuttal has a nice length isonomy with Adhesion
// Contestation ain't in some dictionaries
// Contestation feels more natural to a french thinker
// Rebuttal may be a little (too) intense
//
// uint8 for grades
// ----------------
// Not sure about that.  Is it worth the hassle?  May change.  Advice welcome.
//

type ProposalAnalysis struct {
	TotalSize              uint64 // total amount of judges|judgments across all grades
	MedianGrade            uint8  // 0 == "worst" grade, goes up to the amount of grades minus one
	MedianGroupSize        uint64 // in judges|judgments
	SecondMedianGrade      uint8  // used in Majority Judgment deliberation
	SecondGroupSize        uint64 // either adhesion or contestation, whichever is bigger (contestation prevails)
	SecondGroupSign        int    // -1 for contestation group, +1 for adhesion group
	AdhesionGroupGrade     uint8
	AdhesionGroupSize      uint64
	ContestationGroupGrade uint8
	ContestationGroupSize  uint64
	// Can't decide between Rebuttal and Contestation…  Help!
	//RebuttalGroupGrade uint8
	//RebuttalGroupSize  uint64
}

// This MUTATES THE ANALYSIS, but leaves the proposalTally intact, unchanged.
// MJ uses the low median by default (favors contestation), but there's a parameter if need be.
func (analysis *ProposalAnalysis) Run(proposalTally *ProposalTally, favorContestation bool) {
	analysis.TotalSize = proposalTally.CountJudgments()
	analysis.MedianGrade = 0
	analysis.MedianGroupSize = 0
	analysis.SecondMedianGrade = 0
	analysis.SecondGroupSize = 0
	analysis.SecondGroupSign = 0
	analysis.AdhesionGroupGrade = 0
	analysis.AdhesionGroupSize = 0
	analysis.ContestationGroupGrade = 0
	analysis.ContestationGroupSize = 0

	if 0 == analysis.TotalSize {
		return
	}

	adjustedTotal := analysis.TotalSize - 1
	//if ! favorContestation {
	//	adjustedTotal = analysis.TotalSize + 1
	//}
	medianIndex := adjustedTotal / 2 // Euclidean division
	startIndex := uint64(0)
	cursorIndex := uint64(0)
	for gradeIndex, gradeTally := range proposalTally.Tally {
		if 0 == gradeTally {
			continue
		}

		startIndex = cursorIndex
		cursorIndex += gradeTally
		if (startIndex < medianIndex) && (cursorIndex <= medianIndex) {
			analysis.ContestationGroupSize += gradeTally
			analysis.ContestationGroupGrade = uint8(gradeIndex)
		} else if (startIndex <= medianIndex) && (medianIndex < cursorIndex) {
			analysis.MedianGroupSize = gradeTally
			analysis.MedianGrade = uint8(gradeIndex)
		} else if (startIndex > medianIndex) && (medianIndex < cursorIndex) {
			analysis.AdhesionGroupSize += gradeTally
			if 0 == analysis.AdhesionGroupGrade {
				analysis.AdhesionGroupGrade = uint8(gradeIndex)
			}
		}
	}

	// How to multiline-format this big condition?  I failed.
	if (favorContestation && analysis.AdhesionGroupSize <= analysis.ContestationGroupSize) || (!favorContestation && analysis.AdhesionGroupSize < analysis.ContestationGroupSize) {
		analysis.SecondMedianGrade = analysis.ContestationGroupGrade
		analysis.SecondGroupSize = analysis.ContestationGroupSize
		if 0 < analysis.SecondGroupSize {
			analysis.SecondGroupSign = -1
		}
	} else {
		analysis.SecondMedianGrade = analysis.AdhesionGroupGrade
		analysis.SecondGroupSize = analysis.AdhesionGroupSize
		if 0 < analysis.SecondGroupSize {
			analysis.SecondGroupSign = 1
		}
	}

}