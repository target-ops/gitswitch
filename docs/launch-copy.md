# Launch copy — paste-ready posts

Four pieces, each tuned to its target audience. Replace any details
(numbers, tone) that aren't *your* real story before publishing —
specificity is what makes this work.

**Stagger the posting**, don't fire them all the same hour:
- Day 1: LinkedIn
- Day 2 (Tue/Wed, ~9am the audience's time): Reddit /r/git
- Day 3 (or same day, different time): /r/commandline + /r/golang
- Sit in the comments for the first 4 hours after each Reddit post

---

## 1. LinkedIn

**Why this format:** LinkedIn truncates with "see more" around the 200-character mark on mobile. The hook has to land in the first 2 lines. Plain text + line breaks beat formatting; this isn't a blog post. Single video/GIF at the end (LinkedIn embeds GIFs natively from URL, but uploading the MP4 is more reliable — use `docs/demo.mp4`).

**Post body (paste as-is, replace the GIF placeholder with your upload):**

```
I committed forty-seven times to a client's repo as my personal email. Over three weeks. Nobody noticed except me, on a Saturday night, scrolling my contribution graph by accident.

The mistake was depressingly mundane: I'd `git config --global user.email` to my personal address two weeks earlier to push a fix to a side project, and forgot to switch back.

Git doesn't know you're two different people depending on which folder you're in. SSH happily sends every key in your agent to every host. The GitHub CLI has its own opinion about who you are, independent of `git config`. None of these layers talk to each other, and any of them can silently disagree.

So I built gitswitch.

It binds an identity to a directory. After a one-time setup, every `cd` becomes the switch — git config, SSH key, gh CLI account, and commit signing all line up automatically.

The killer feature: a pre-commit hook that refuses commits where the active identity is wrong for the directory. Forget to switch? You don't ship the bad commit. You get a one-line error telling you what to do.

[ATTACH demo.mp4 OR docs/demo.gif HERE]

Single Go binary. MIT licensed. Install in 8 seconds:

  brew install target-ops/tap/gitswitch
  gitswitch init

If you've ever felt that specific stomach-drop when you noticed your last commit went out as the wrong person — this is for you.

→ github.com/target-ops/gitswitch

What's your worst "I committed as the wrong person" story?

#git #DeveloperTools #OpenSource #DX #GoLang
```

**Length:** ~280 words. The "see more" cutoff hides everything after the second paragraph; that's why the hook leads.

**One change you should make before posting:** the `forty-seven commits` number is mine, not yours — replace with your real number, even if it's just "a few times" or "twice last year." Specificity beats made-up specificity.

---

## 2. Reddit — /r/git

**Why this audience:** Git enthusiasts. They'll appreciate the technical depth (`includeIf`, `IdentitiesOnly yes`, the trailing-slash gotcha). They have *opinions* — engage in the comments, don't just drop and run.

**Title:**

```
I built gitswitch — per-directory git identity binding with a pre-commit hook that refuses wrong-author commits
```

**Body:**

```
After committing to a client repo with my personal email for three weeks (forty-seven commits, only noticed by accident looking at my contribution graph), I went looking for the tool that prevents this. The pieces existed:

- `includeIf` for per-directory gitconfig (added in git 2.13, barely documented; every blog post about it has a comment about the trailing-slash gotcha)
- per-account SSH host aliases with `IdentitiesOnly yes`
- `gh auth switch` for the GitHub CLI
- SSH commit signing (git 2.34+, `gpg.format=ssh`)

None of them refused to let me commit when the email was wrong, and configuring all four to work together is a 90-minute yak shave I'd forget by the next laptop.

`gitswitch` ties them together:

```
gitswitch init                       # autodetect existing identities
gitswitch use work     ~/work        # bind work identity to a directory
gitswitch use personal ~/code        # ditto for personal
gitswitch guard install              # the killer feature
```

`guard` installs a pre-commit hook that refuses commits where `git config user.email` doesn't match the directory's bound identity. The error includes the exact one-line fix.

[Demo GIF: https://github.com/target-ops/gitswitch/raw/main/docs/demo.gif]

Source: https://github.com/target-ops/gitswitch (Go, MIT, brew install, single 2MB binary)

Curious: how do you handle multiple identities today? `includeIf` alone? Manual switching? 1Password SSH? Something else?
```

---

## 3. Reddit — /r/commandline

**Why this audience:** CLI-tool collectors. Terse, no-nonsense. Doesn't tolerate marketing copy. Wants to know what it does in two sentences.

**Title:**

```
gitswitch — per-directory git identity binding, with a pre-commit hook (Go, single binary)
```

**Body:**

```
Built after I committed to a client's repo as my personal email for three weeks straight without noticing.

gitswitch binds an identity (email, signing key, SSH key, gh account) to a directory and installs a pre-commit hook that refuses commits where the active identity is wrong for that directory. After setup, every `cd` is the switch — no manual step.

[Demo GIF: https://github.com/target-ops/gitswitch/raw/main/docs/demo.gif]

Setup:

```
brew install target-ops/tap/gitswitch
gitswitch init
gitswitch use work ~/work
gitswitch guard install
```

5 headline commands: `init`, `use`, `guard`, `doctor`, `why`. Plus `list`, `add`, `delete`, `rename` for lifecycle.

Single Go binary, ~2MB, no runtime deps. macOS arm64/x64, Linux x64/arm64, Windows. MIT.

https://github.com/target-ops/gitswitch
```

---

## 4. Reddit — /r/golang

**Why this audience:** Go developers will care about the Go-specific things (single binary, GoReleaser, the lipgloss panels). They'll also notice if you're showing off — keep the tone matter-of-fact.

**Title:**

```
gitswitch — Go CLI for per-directory git identity, with a pre-commit hook
```

**Body:**

```
Shipped my first end-to-end Go project this week — a CLI that binds a git identity (email, signing key, SSH key, gh account) to a directory via `includeIf`, and installs a pre-commit hook that refuses commits with the wrong active identity.

Stack:

- cobra for the CLI surface
- charmbracelet/lipgloss for the styled panels (the `guard` blocked-commit message has a red rounded border — surprisingly satisfying)
- charmbracelet/huh for interactive prompts (`init` walks you through detected identities)
- GoReleaser cross-compiles 5 platform binaries on every release tag
- release-please drives semver bumps from Conventional Commits, opens release PRs that merge to tag
- `brews:` block in goreleaser auto-pushes the formula to a homebrew tap

Single binary, ~2MB, `CGO_ENABLED=0` so no glibc surprises.

The release pipeline is fully hands-off: `feat:` PR merged → release-please opens `chore: release` PR → merge → 5 binaries published + brew formula auto-updated. End-to-end in ~90s. Cut three releases this week without manually tagging anything.

[Demo GIF: https://github.com/target-ops/gitswitch/raw/main/docs/demo.gif]

```
brew install target-ops/tap/gitswitch
```

https://github.com/target-ops/gitswitch (MIT)

Honest critique welcome — especially around (a) lack of tests yet (working on it next), (b) where I drew the line between a `style` package and reaching for lipgloss directly.
```

---

## Posting checklist

For each Reddit post:

- [ ] Read the subreddit's rules first (sticky post or sidebar). Some require a `[Show]` or `[Project]` prefix in the title.
- [ ] Post Tuesday or Wednesday, 9–11am the audience's local time (US-east morning catches both US and Europe).
- [ ] Don't post the same content to multiple subreddits within an hour — looks like spam.
- [ ] **Sit in the comments for the first 4 hours.** Reply to every comment in <15 minutes. This determines whether the post takes off.
- [ ] If someone asks "what about [other tool]?", respond with respect — the prior art mention from your /r/git body covers most of these (`includeIf`, `gh auth switch`).
- [ ] Don't crosspost. Each subreddit gets its own purpose-built post.

For the LinkedIn post:

- [ ] Upload the MP4, not just a link to the GIF — LinkedIn renders MP4 inline; GIF links are sometimes flattened to a thumbnail.
- [ ] Post on a Tuesday-Thursday morning your time zone (LinkedIn engagement skews business-hours).
- [ ] Reply to comments within the first 4 hours — same engagement-loop dynamic as Reddit.
- [ ] Don't add too many hashtags (LinkedIn caps useful reach around 3–5).

## What NOT to do

- Don't post to /r/programming yet. It's hostile to self-promotion unless you have HN-level traction first. Wait until the post on /r/git or HN does well, then crosspost there with a different angle.
- Don't post to Hacker News from a brand-new account or as "Show HN: My App". Read the [Show HN guidelines](https://news.ycombinator.com/showhn.html) first; the title format matters.
- Don't write a Medium post. The dev.to draft (`docs/launch-post.md`) is the blog version.
