# Contributing to `go-apt-transport`

First, thanks for the interest in contributing! Before contributing, be sure to
read the [Code of Conduct](CODE_OF_CONDUCT.md)

# License

By contributing to `go-apt-transport`, you agree to license your commits under
the [MIT License](LICENSE.md)

# Git Commit Messages

Git commit messages must utlize [gitm😍ji](https://gitmoji.dev) or the actual
emoji themselves for each line item. If a line requires more information place
it in a paragraph directly below the line item (unless the description is for
the subject line of the commit). Additionally, messages should use markdown, as
most tooling can render it. For example:

```gitcommit
⚡ Improve speed by fooing bars instead of bazzes.

Fooing a baz has been proven to be slow in certain cases due to string
expansion. By using a bar, we're able to reduce this performance issue.

🐛 Resolve incorrect usage of cmake_minimum_required
After consulting with 13 oracles, 12 wizards, and a wise Kiwi goat named
"Harold", we've finally figured out when to call `cmake_minimum_required`.

♻ Refactor several internal functions just to keep people guessing.
When writing lists, make sure to indent as needed for readability. Think of
explanations for a commit message as "a tweet but with less anger and more
technical reasoning"

  - Please make nice clean lists
  - spacing between elements matters
  - thank you

```

Please note that the [gitm😍ji](https://gitmoji.dev) list has changed over
time, and will continue to change or evolve over time, but always assume that
the currently published list is to be used.

NOTE: We do not currently have a conventional commit standard used to generate
changelogs, but once implemented this will be an automated operation that will
be enforced, and hopefully still rely on gitmoji.
