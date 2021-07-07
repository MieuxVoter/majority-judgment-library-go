# Contributing


## Dependencies

No dependencies would be best.
We do have a dev dependency on an assertions lib for testing.
Should be shaken out if tests are not exported, right?


## Tests

    go test


## Publishing

https://golang.org/doc/modules/publishing

    GOPROXY=proxy.golang.org go list -m github.com/mieuxvoter/judgment@v0.1.0


## Score Calculus

See the page about [Majority Judgment Score Calculus](./SCORE.md).
