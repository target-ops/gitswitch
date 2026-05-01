# Launch materials — GIF storyboard, voice guide, channels checklist

This file is a working doc for the 1.0 launch. It exists so the look, voice, and rollout stay coherent across whoever ends up touching the README, the demo, the blog post, and the PR descriptions.

---

## 1. Demo GIF — storyboard

**Goal:** the GIF is the marketing. It plays autoplaying-muted on the README, on Twitter, embedded in the Show HN post. A first-time visitor decides whether to install based on this 15-second loop.

**Hard constraints:**

- ≤ 15 seconds total (Twitter autoplays 30s; we want a tight loop)
- ≤ 2 MB (renders fast on slow connections; embedded inline)
- 800–960 px wide, 2× pixel density for retina
- Real shell, real terminal, real font (Berkeley Mono, JetBrains Mono, or Iosevka). **Not** a screen-recorded SVG, **not** asciinema. People can tell.
- Dark terminal background with high contrast. Red and green must read instantly.
- The cursor visible, the prompt visible. People want to feel like they're watching a real session.

**Recording stack (recommended):**

- Terminal: Ghostty or iTerm2, font size ~16pt
- Recorder: `vhs` (charm.sh — scriptable, deterministic, ships GIF and MP4)
- Window chrome: minimal — no traffic-light buttons, no menubar
- Save as both `demo.gif` (for README) and `demo.mp4` (for Twitter — much smaller)

### Three scenes, ~15 seconds total

#### Scene 1 (0:00 – 0:04) — the setup

```
~/work/some-repo  ❯ git status
On branch main
Changes not staged for commit:
  modified:   src/handler.go

~/work/some-repo  ❯ git add -A && git commit -m "fix the thing"
```

The viewer sees a normal-looking commit about to happen. *Quiet, neutral.* No hint that anything is wrong yet.

#### Scene 2 (0:04 – 0:11) — the block

The commit fires. Output appears in red, slightly larger than ambient text:

```
✗ gitswitch guard: blocked commit

  in directory:   ~/work/some-repo/
  expected:       you@company.com   (bound identity: work)
  got:            you@gmail.com

  fix:            gitswitch use work
                  (or: git commit --no-verify to override this once)
```

This is the **emotional event**. Hold the frame for a beat — let the viewer read it. The whole point of the GIF is that this moment makes them go *"oh — I want that."*

#### Scene 3 (0:11 – 0:15) — the recovery

```
~/work/some-repo  ❯ gitswitch use work
✓ active identity: work

~/work/some-repo  ❯ git commit -m "fix the thing"
[main 3a5ad97] fix the thing
 1 file changed, 2 insertions(+)
```

Green check on `gitswitch use work`. Commit succeeds cleanly. Loop.

### What NOT to put in the GIF

- Multiple commands the viewer has to read carefully — they're scanning, not reading
- Full identity tables, multi-line `gitswitch doctor` output — too dense for a GIF
- Marketing taglines overlaid on the video — let the terminal speak
- Slow typing animation — feels fake; use realistic paste-then-execute pacing
- Music or sound — README GIFs play muted; nobody will ever hear it

---

## 2. Voice guide

How gitswitch talks. Used for: README, error messages, PR descriptions, issue responses, blog posts, social copy.

### Principles

1. **First person where appropriate, never marketing-plural.** The maintainer is one person. Don't say "we're excited to announce" — say "I built this because I got bitten." Belonging is built from honesty, not corporate plurals.

2. **Specific over general.** *"forty-seven commits over three weeks"* beats *"a lot of commits."* *"the trailing-slash gotcha in `includeIf`"* beats *"git's quirks."* Specific is memorable; general is forgettable.

3. **Show what the user sees, not what we did.** Errors say *"fix: gitswitch use work"* not *"identity validation failed."* Output paths matter; abstract status codes don't.

4. **Avoid praise of yourself in the docs.** Don't say *"a powerful tool that makes managing identities easy."* Let the user decide if it's powerful or easy. Just describe what it does and let the work speak.

5. **Acknowledge the prior art.** `includeIf`, `gh auth switch`, `IdentitiesOnly yes` — these all existed before gitswitch. The contribution is the integration, not the invention. Saying so makes us trustworthy.

### Tone calibration

| Don't say | Say instead |
|---|---|
| "We're proud to announce…" | "I built this because…" |
| "Powerful identity management" | "Bind an identity to a directory." |
| "Seamlessly switch between accounts" | "Every `cd` is a switch." |
| "Robust security" | "SSH commit signing on by default." |
| "Easy to use" | (nothing — show the quickstart instead) |
| "Best-in-class" | (just delete the sentence) |
| "User-friendly error messages" | (give an example error message) |

### Error message format

Every error gets three lines:
1. **What went wrong** (red, terse)
2. **Why it went wrong** (the relevant facts: file paths, expected vs actual, identity names)
3. **What to do next** (an exact command they can run)

Bad:
```
Error: identity mismatch.
```

Good:
```
✗ gitswitch guard: blocked commit
  expected:  you@company.com  (bound identity: work)
  got:       you@gmail.com
  fix:       gitswitch use work
```

### Common phrasings to keep consistent

- "**bind**" an identity to a directory (not "register", not "associate")
- "**guard**" the pre-commit hook (always lowercase, never "the Guard")
- "**identity**" not "profile", not "account", not "user" (it's the thing that travels across git/ssh/gh)
- "**directory**" not "folder" (we're a CLI tool; speak Unix)
- "**bound**" or "**unbound**" to describe whether an identity has a directory associated

---

## 3. Launch channels — checklist

In the order they go out, not all on the same day. Spacing matters; people share less when they see the same project everywhere at once.

### Day 0 — the foundation

- [ ] README rewritten and merged (this PR)
- [ ] `LICENSE` in place (done)
- [ ] Demo GIF recorded and embedded in the README
- [ ] `1.0.0` tag cut (only after `init`, `use`, `guard`, `doctor`, `why` ship)
- [ ] Homebrew formula switched to prebuilt binary download
- [ ] `curl | sh` installer at `gitswitch.dev/install` (or a stable raw GitHub URL)
- [ ] GitHub Discussions enabled with a pinned "introduce yourself" thread

### Day 1 — owned channels

- [ ] Personal blog / dev.to: publish the launch post (`docs/launch-post.md`)
- [ ] Twitter/X thread: 6–8 tweets, GIF in the first tweet, install command in the last
- [ ] LinkedIn (if relevant audience): the dev.to post, abridged
- [ ] Mastodon / fediverse: same as Twitter, smaller audience but engaged

### Day 2 — Show HN

- [ ] Tuesday or Wednesday, **9am Pacific** (peak HN traffic)
- [ ] Title: `Show HN: Gitswitch – Stop committing as the wrong person`
- [ ] Body: 2 short paragraphs. Personal pain (1) + what gitswitch does (1). Link the GIF.
- [ ] **Sit in the comments for the first 4 hours.** Reply to every comment within 15 minutes. This is what determines whether the post takes off — engaged authors who answer questions get upvotes.

### Day 3+ — long tail

- [ ] PR to `awesome-cli-apps`
- [ ] PR to `awesome-shell`
- [ ] PR to `awesome-git`
- [ ] Reddit `/r/git` (probably the warmest reception)
- [ ] Reddit `/r/programming` (cooler — only if HN goes well)
- [ ] Reddit `/r/commandline` (small but engaged)
- [ ] Lobsters (friendly to CLI tools; one post, no spamming)
- [ ] Hacker Newsletter, TLDR, etc. — submit via their normal channels
- [ ] Reach out to one or two adjacent maintainers (`gh` CLI? `lazygit`?) with a heads-up and a thanks for the prior art

### Forever — community

- [ ] Respond to every issue within 48 hours, even if just to acknowledge
- [ ] Triage labels in good order (`bug`, `enhancement`, `docs`, `discussion`)
- [ ] Once / month: write a small Discussions post with what changed
- [ ] Once / quarter: write a "where we are" post with metrics, lessons, what's next

### Don't do (yet)

- ❌ Product Hunt — wrong audience for dev tools, low ROI
- ❌ Paid promotion — wastes money at this stage
- ❌ "Featured on" badges in the README — adds noise, looks insecure

---

## 4. The first-issue safety nets

Before launch, make sure `gitswitch init` works on these realistic states without producing a horrible first impression:

- A user with no identity configured at all (fresh laptop)
- A user with one identity (the common case)
- A user with two GitHub accounts and one GitLab (the target audience)
- A user with an existing `~/.ssh/config` that already has bastion / jump-host entries (must preserve)
- A user with `gh` not installed (graceful skip, friendly hint)
- A user with `gh` installed but not authenticated
- A user with their key in 1Password SSH agent (no key file on disk, just agent)
- A user behind a corporate proxy (network calls must time out, never hang)

Every one of these is a real person who will install gitswitch within the first month. Each one of them having a clean experience is worth more than any new feature.

---

## 5. Once we're past launch

When 1.0 is out and the dust settles, the next conversation is what 1.x looks like — the features that turn gitswitch from "a tool people install" into "a tool people belong to":

- `gitswitch presence` — shell prompt segment
- `gitswitch share` — onboarding link generator
- `gitswitch agent` — AI-coding-agent identity registry
- `.gitswitch.yml` — repo-side identity policy
- `gitswitch lockdown` — panic button
- `gitswitch audit` — privacy / leak scan
- `gitswitch credential-helper` — fix the macOS keychain bug

Each one its own small launch. Each one a reason for the people who already have gitswitch installed to tell someone new about it.
