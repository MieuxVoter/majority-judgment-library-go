package judgment

type PollResult struct {
	Proposals       ProposalsResults // matches the order of the input proposals' tallies
	ProposalsSorted ProposalsResults // same Results, but sorted by Rank this time
}

type ProposalResult struct {
	Index    int    // index of the proposal in the input proposals' tallies.  Useful with ProposalSorted
	Rank     int    // Rank starts at 1 (best) and goes upwards.  Equal Proposals share the same rank.
	Score    string // higher lexicographically → better rank
	Analysis *ProposalAnalysis
}

// ProposalsResults implements sort.Interface based on the Score field.
type ProposalsResults []*ProposalResult

func (a ProposalsResults) Len() int           { return len(a) }
func (a ProposalsResults) Less(i, j int) bool { return a[i].Score < a[j].Score }
func (a ProposalsResults) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
