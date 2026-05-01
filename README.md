<div align="center">

# gitswitch

**Stop committing as the wrong person.**

[![Latest release](https://img.shields.io/github/v/release/target-ops/gitswitch?style=flat-square)](https://github.com/target-ops/gitswitch/releases/latest)
[![License: MIT](https://img.shields.io/badge/license-MIT-blue.svg?style=flat-square)](LICENSE)
[![Build](https://img.shields.io/github/actions/workflow/status/target-ops/gitswitch/release.yml?style=flat-square&label=release)](https://github.com/target-ops/gitswitch/actions)
[![Stars](https://img.shields.io/github/stars/target-ops/gitswitch?style=flat-square)](https://github.com/target-ops/gitswitch/stargazers)

</div>

---

```text
$ git commit -m "fix the auth flow"

✗ gitswitch guard: blocked commit

  in directory:   ~/work/some-repo/
  expected:       you@company.com   (bound identity: work)
  got:            you@gmail.com

  fix:            gitswitch use work
                  (or: git commit --no-verify to override this once)
```

## The problem

You set your `git user.email` to your work address last Tuesday. Forgot to switch back. For the next three weeks every commit to your personal side project went out under your employer's email. You found out by accident, looking at your contribution graph.

Git has no idea you're two different people depending on which folder you're in. SSH cheerfully sends every key in your agent to every host. The GitHub CLI is logged into a different account than your last `git config`. None of these layers talk to each other, and any of them can silently disagree.

`gitswitch` makes them agree, and **refuses to let you commit when they don't.**

---

## Install

```bash
brew install target-ops/tap/gitswitch
```

Single binary. ~2 MB. macOS arm64/x64, Linux x64/arm64. No Python. No runtime dependencies. Installs in 8 seconds.

Windows: grab the `.zip` from [releases](https://github.com/target-ops/gitswitch/releases/latest).

---

## Quickstart — 30 seconds

```bash
gitswitch init                       # auto-detect what's already on this machine
gitswitch use work     ~/work        # bind work identity to a directory
gitswitch use personal ~/code        # bind personal identity to another
gitswitch guard install              # refuse wrong-author commits at the source
```

That's it. Every `cd` is now a switch. Forget to switch — `gitswitch` won't let you ship the bad commit.

---

## What it does

| Command | What it does |
|---|---|
| `gitswitch init` | Auto-detect identities from your `~/.gitconfig`, `~/.ssh/config`, GitHub CLI, ssh keys |
| `gitswitch use <id> [<dir>]` | Bind an identity to a directory (writes a sentinel-marked `includeIf` block) |
| `gitswitch guard install` | Install the global pre-commit hook that blocks wrong-author commits |
| `gitswitch doctor` | Verify all layers — git, ssh, gh, signing — agree on who you are right now |
| `gitswitch why` | Explain, in plain English, why your active identity is what it is |

Run any command with `--help` for the full reference.

---

## How it compares

|  | Manual `includeIf` | `gh auth switch` | **gitswitch** |
|---|:---:|:---:|:---:|
| Auto-switch by directory | ✓ (if you nail the trailing-slash gotcha) | ✗ | **✓** |
| Per-account SSH key isolation | manual | ✗ | **✓** |
| Keeps `gh` in lockstep with `git` | ✗ | partial | **✓** |
| SSH commit signing on by default | ✗ | ✗ | **✓** |
| Refuses wrong-author commits | ✗ | ✗ | **✓** |
| One-command verify all layers | ✗ | ✗ | **`gitswitch doctor`** |

---

## How it works

Three things happen the first time you `gitswitch use <id> <dir>`:

1. **Per-identity gitconfig** at `~/.config/gitswitch/identities/<id>.gitconfig` — sets `user.name`, `user.email`, signing key, and `core.sshCommand` with `IdentitiesOnly=yes`.
2. **Conditional include** — a sentinel-marked block in `~/.gitconfig` that loads the per-identity file only when you're inside the bound directory:

   ```
   # >>> gitswitch:work
   [includeIf "gitdir:~/work/"]
       path = ~/.config/gitswitch/identities/work.gitconfig
   # <<< gitswitch:work
   ```

3. **Binding record** in `~/.config/gitswitch/config.json` so `gitswitch why` can explain itself later.

Once `guard` is installed, every `git commit` runs a ~5 ms check: does `user.email` match the identity bound to this directory? Yes → commit. No → refuse with a one-line fix. The dev.to story *"I used the wrong git email for two weeks and no one noticed"* — gitswitch makes that story impossible.

---

## Philosophy

Git identity should be **infrastructure**, not something you remember.

The tool is small, the binary is single, the only state on your machine lives in `~/.config/gitswitch/` and the directories *you* tell it to manage. No service. No cloud sync of your keys. No telemetry. The maintainer is one developer who got bitten and built this; the issue tracker responds in days, not weeks; the roadmap is public; the tests pass.

---

## Community

- **[Issues](https://github.com/target-ops/gitswitch/issues)** — bugs and feature requests
- **[Discussions](https://github.com/target-ops/gitswitch/discussions)** — questions, design conversations
- PRs welcome. Run `go build ./cmd/gitswitch` and try the binary against your own setup before opening one.

---

## License

MIT. See [LICENSE](LICENSE).
