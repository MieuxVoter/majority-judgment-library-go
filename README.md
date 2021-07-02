# Majority Judgment for Golang

> WORK IN PROGRESS

A Go library to deliberate using Majority Judgment.

We use a **score-based algorithm**, for performance and scalability.

Supports billions of judgments and thousands of proposals per poll, if need be.


## Installation

    go get -u github.com/mieuxvoter/judgment


## Usage

```go

package main

import (
	"fmt"
	"log"

	"github.com/mieuxvoter/judgment"
)

func main() {

    poll := &judgment.PollTally{
        AmountOfJudges: 10,
        Proposals: []*judgment.ProposalTally{
            {Tally: []uint64{2, 2, 2, 2, 2}}, // Proposal A   Amount of judgments received for each grade,
            {Tally: []uint64{2, 1, 1, 1, 5}}, // Proposal B   from "worst" grade to "best" grade.
            {Tally: []uint64{2, 1, 1, 2, 4}}, // Proposal C   Make sure all tallies are balanced, that is they
            {Tally: []uint64{2, 1, 5, 0, 2}}, // Proposal D   hold the same total amount of judgments.
            // …
        },
    }
    deliberator := &judgment.MajorityJudgment{}
    result, err := deliberator.Deliberate(poll)

    if nil != err {
        log.Fatalf("Deliberation failed: %v", err)
    }

    fmt.Printf("Result: %v\n", result) // proposals results are ordered like tallies, but Rank is available 

}

```


## License

`MIT` 🐜


## Contribute

A review by a seasoned `Go` veteran would be appreciated.


