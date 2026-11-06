# Issue tracker: GitHub

Issues and specs live in GitHub Issues for `parklogic/go`.
Use the `gh` CLI with `--repo parklogic/go`.

## Operations

- Read: `gh issue view <number> --repo parklogic/go --comments`
- List: `gh issue list --repo parklogic/go --state open`
- Create: `gh issue create --repo parklogic/go --title "..." --body-file <file>`
- Comment: `gh issue comment <number> --repo parklogic/go --body-file <file>`
- Label: `gh issue edit <number> --repo parklogic/go --add-label "..."`
- Remove label: `gh issue edit <number> --repo parklogic/go --remove-label "..."`
- Close: post the explanation as a comment, then
  `gh issue close <number> --repo parklogic/go`.

When a skill says to publish to the issue tracker, create a GitHub issue.
When it says to fetch a ticket, read the issue and its comments.

## Relationships

Use native GitHub sub-issues and issue dependencies when supported.
Otherwise, put `Part of #<parent>` and `Blocked by: #<number>` at the
top of child issue bodies, and maintain a task list in the parent.
A ticket is unblocked when every blocker is closed.

## Pull requests as a triage surface

PRs as a request surface: no.

## Wayfinding

Use one issue labelled `wayfinder:map` for the map, with linked child
tickets labelled `wayfinder:research`, `wayfinder:prototype`,
`wayfinder:grilling`, or `wayfinder:task`.

The frontier is the first open child in map order with no assignee
and no open blockers. Claim it by assigning yourself.
Resolve it by commenting with the answer, closing it, and adding
a summary and link to the map's Decisions-so-far.
