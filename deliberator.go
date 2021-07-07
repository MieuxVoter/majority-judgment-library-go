package judgment

type DeliberatorInterface interface {
	Deliberate(tally *PollTally) (result *PollResult, err error)
}
