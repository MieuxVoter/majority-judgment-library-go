package judgment

// DeliberatorInterface ought to be implemented by all deliberators ; not overly useful for now, but hey.
type DeliberatorInterface interface {
	Deliberate(tally *PollTally) (result *PollResult, err error)
}
