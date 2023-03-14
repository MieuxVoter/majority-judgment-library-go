# Contributing


## Dependencies

No dependencies would be best.
We do have _one_ dev dependency on an assertions library for testing.
Should be shaken out automatically if tests are not exported, right?


## Visualize

[Gource](https://gource.io/) is very useful to quickly browse the history and structure of a project:

    git log --pretty=format:"%at|%s" --reverse --no-merges > commitmsg.txt
    gource --font-scale 2.0 --highlight-dirs --filename-time 7.0 --caption-file commitmsg.txt --caption-size 26 --realtime


## Tests

    cd judgment
    go test


## Publishing

https://golang.org/doc/modules/publishing

    GOPROXY=proxy.golang.org go list -m github.com/mieuxvoter/judgment@v0.1.1


## Score Calculus

See the page about [Majority Judgment Score Calculus](./SCORE.md).
