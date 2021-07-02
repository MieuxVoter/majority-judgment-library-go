package judgment

type PollResult struct {
	Proposals ProposalsResults // Matches the order of the input proposals' tallies
}

type ProposalResult struct {
	Index int    // index of the proposal in the input proposals' tallies
	Rank  int    // Rank starts at 1 (best) and goes upwards.  Equal Proposals share the same rank.
	Score string // higher lexicographically → better rank
}

// ProposalsResults implements sort.Interface based on the Score field.
type ProposalsResults []*ProposalResult

func (a ProposalsResults) Len() int           { return len(a) }
func (a ProposalsResults) Less(i, j int) bool { return a[i].Score < a[j].Score }
func (a ProposalsResults) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
