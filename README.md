# Majority Judgment for Golang

[![MIT](https://img.shields.io/github/license/MieuxVoter/majority-judgment-library-go?style=for-the-badge)](LICENSE.md)
[![Release](https://img.shields.io/github/v/release/MieuxVoter/majority-judgment-library-go?include_prereleases&style=for-the-badge)](https://github.com/MieuxVoter/majority-judgment-library-go/releases)
[![Build Status](https://img.shields.io/github/workflow/status/MieuxVoter/majority-judgment-library-go/Go?style=for-the-badge)](https://github.com/MieuxVoter/majority-judgment-library-go/actions/workflows/go.yml)
[![Coverage](https://img.shields.io/codecov/c/github/MieuxVoter/majority-judgment-library-go?style=for-the-badge&token=FEUB64HRNM)](https://app.codecov.io/gh/MieuxVoter/majority-judgment-library-go/)
[![Code Quality](https://img.shields.io/codefactor/grade/github/MieuxVoter/majority-judgment-library-go?style=for-the-badge)](https://www.codefactor.io/repository/github/mieuxvoter/majority-judgment-library-go)
[![A+](https://img.shields.io/badge/go%20report-A+-brightgreen.svg?style=for-the-badge)](https://goreportcard.com/report/github.com/mieuxvoter/majority-judgment-library-go)
![LoC](https://img.shields.io/tokei/lines/github/MieuxVoter/majority-judgment-library-go?style=for-the-badge)
[![Discord Chat https://discord.gg/rAAQG9S](https://img.shields.io/discord/705322981102190593.svg?style=for-the-badge)](https://discord.gg/rAAQG9S)

A Golang module to deliberate using Majority Judgment.


## Features

- **score-based algorithm**, for performance and scalability
- supports billions of judgments with almost the same cost as dozens
- supports thousands of proposals per poll
- default judgment balancing tools: static grade, median grade


## Installation

    go get -u github.com/mieuxvoter/majority-judgment-library-go

It exposes the package `judgment`, for concision, since the repo name itself is quite long and we can't rename it.

> It's a pre-release.  We still have some int types that may change until v1.


## Usage

Say you have the following tally:

![Example of a merit profile](./docs/2-2-2-2-2_2-1-1-1-5_2-1-1-2-4_2-1-5-0-2_2-2-2-2-2.png)

You can compute out the majority judgment rank of each proposal like so:

```go
package main

import (
    "fmt"
    "github.com/mieuxvoter/majority-judgment-library-go"
    "log"
)

func main() {

    pollTally := &(judgment.PollTally{
        AmountOfJudges: 10,
        Proposals: []*judgment.ProposalTally{
            {Tally: []uint64{2, 2, 2, 2, 2}}, // Proposal A   Amount of judgments received for each grade,
            {Tally: []uint64{2, 1, 1, 1, 5}}, // Proposal B   from "worst" grade to "best" grade.
            {Tally: []uint64{2, 1, 1, 2, 4}}, // Proposal C   Make sure all tallies are balanced, that is they
            {Tally: []uint64{2, 1, 5, 0, 2}}, // Proposal D   hold the same total amount of judgments.
            {Tally: []uint64{2, 2, 2, 2, 2}}, // Proposal E   Equal proposals share the same rank.
            // …
        },
    })
    deliberator := &(judgment.MajorityJudgment{})
    result, err := deliberator.Deliberate(pollTally)

    if nil != err {
        log.Fatalf("Deliberation failed: %v", err)
    }

    // Proposals results are ordered like tallies, but Rank is available. 
    // result.Proposals[0].Rank == 4 // Proposal A
    // result.Proposals[1].Rank == 1 // Proposal B
    // result.Proposals[2].Rank == 2 // Proposal C
    // result.Proposals[3].Rank == 3 // Proposal D
    // result.Proposals[4].Rank == 4 // Proposal E

    // You may also use proposals sorted by Rank ; their initial Index is available
    // result.ProposalsSorted[0].Index == 1 // Proposal B
    // result.ProposalsSorted[1].Index == 2 // Proposal C
    // result.ProposalsSorted[2].Index == 3 // Proposal D
    // result.ProposalsSorted[3].Index == 0 // Proposal A
    // result.ProposalsSorted[4].Index == 4 // Proposal E
    
    fmt.Printf("Best Proposal Index: %d\n", result.ProposalsSorted[0].Index)

}
```

### Balancing uneven proposals

Sometimes, some proposals receive more judgments than others, and the tallies are unbalanced.
In those cases, a _default judgment_ strategy has to be picked:

#### Static Default Grade

Missing (or "I don't know") judgments are considered of a grade defined in advance,
usually `to reject` in order to incentivize proposals to be explicit and clear.

You may use `PollTally.BalanceWithStaticDefault(defaultGrade)` to that effect:

```go
pollTally := &PollTally{
    AmountOfJudges: 10,
    Proposals: []*ProposalTally{
        {Tally: []uint64{2, 1, 2, 2, 1}}, // Proposal A
        {Tally: []uint64{3, 1, 3, 1, 1}}, // Proposal B
        {Tally: []uint64{0, 1, 1, 0, 0}}, // Proposal C
        // …
    },
}
defaultGrade := 0
pollTally.BalanceWithStaticDefault(defaultGrade)

// pollTally was mutated and now contains balanced proposals' tallies
// pollTally.Proposals[0].Tally == {4, 1, 2, 2, 1}  // Proposal A is now balanced 
// pollTally.Proposals[1].Tally == {4, 1, 3, 1, 1}  // Proposal B is now balanced
// pollTally.Proposals[2].Tally == {8, 1, 1, 0, 0}  // Proposal C is now balanced
```

#### Median Default Grade

Same behavior as static, but the default grade for each proposal is its median grade.

Use `PollTally.BalanceWithMedianDefault()`.


#### Normalization

> Not implemented yet ; would require either `math/big` for LCM or floating-point arithmetic
> This is part of why deciding on int types is so tricky.


## License

`MIT` 🐜


## Contribute

This project needs a review by `Go` devs.
Feel free to suggest changes, report issues, make improvements, etc.

Some more information is available in [`docs/`](./docs).

