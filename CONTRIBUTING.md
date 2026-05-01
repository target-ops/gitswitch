# Contributing to gitswitch

PRs welcome. Two conventions matter — the rest is just good taste.

## Commit messages: Conventional Commits

We use [Conventional Commits](https://www.conventionalcommits.org/) so
[release-please](https://github.com/googleapis/release-please) can
auto-version releases from the commit history. Your commit subject
should start with one of:

| Prefix | When | Effect on the next release |
|---|---|---|
| `feat:` | new user-visible capability | minor bump (e.g., 1.0.x → 1.1.0) |
| `fix:` | bug fix | patch bump (e.g., 1.1.0 → 1.1.1) |
| `feat!:` or `fix!:` *or* `BREAKING CHANGE:` in body | incompatible change | major bump (e.g., 1.x → 2.0.0) |
| `docs:` | documentation only | no version bump (still appears in changelog) |
| `chore:`, `refactor:`, `test:`, `ci:`, `build:`, `perf:` | internal | no version bump |

Examples:

```
feat: add gitswitch why command
fix: handle missing ~/.ssh/config in init
docs: clarify the includeIf trailing-slash gotcha
chore: bump cobra to v1.10.3
```

If a PR squashes multiple commits, set the **PR title** to the
release-relevant Conventional Commit form — that's the line
release-please reads when squash-merging.

## Releases happen via PR, not by hand

You don't tag releases. After your PR lands on `main`,
release-please opens (or updates) a *"chore: release vX.Y.Z"* PR.
Merge that PR when you're ready to ship — that's what tags the
release, fires the cross-compile workflow, and updates the brew
formula. Fully automated downstream.

If you want to test a `main` build locally before the release PR
ships, just:

```bash
git pull origin main
go build -o /tmp/gitswitch ./cmd/gitswitch
/tmp/gitswitch --help
```

## Code conventions

- **One package per concern.** `internal/identity/` owns the JSON
  store; `internal/git/` wraps the git CLI; `internal/cmd/` is
  cobra commands only. Don't reach across.
- **Errors say what to do next.** Every error string ends with the
  exact command that fixes the situation. See `internal/cmd/use.go`
  for the pattern.
- **No raw ANSI escapes in `cmd/`.** Use `internal/style` so
  `--no-color` and `NO_COLOR` keep working. Same for `green`,
  `red`, etc. — they're routed through the style package.
- **Smoke-test against your own setup before opening a PR.**
  `gitswitch doctor` is the cheapest sanity check.

## Reviewing your own PR

A short checklist that catches 90% of mistakes:

- [ ] PR title is a Conventional Commit subject.
- [ ] `go vet ./...` is clean.
- [ ] Personal email / username didn't sneak into a doc or comment.
- [ ] Behavior change includes a test plan in the PR body, even if
      not automated yet.

That's it. The bar is "does it work, is it readable, and would
someone reading `git log` understand why it landed?"
