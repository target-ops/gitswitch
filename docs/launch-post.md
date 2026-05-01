# I committed to a client's repo as my personal email — for three weeks. So I built this.

*Draft for dev.to / personal blog. Voice: first person, vulnerable, specific. ~950 words. Title is the hook; don't dilute it.*

---

It was a Saturday, around 11pm, and I was idly looking at my GitHub contribution graph for no particular reason. I noticed something was off. The green squares for the past three weeks were on my personal account, but I'd been doing client work that whole time. I clicked one of the squares.

It was a commit on the client's repo. Authored by `ofir474@gmail.com`. My personal email. On the company project.

I scrolled. There were forty-seven of them.

I sat there for a minute trying to figure out how I'd let this happen. The truth is depressingly mundane: two weeks earlier I'd `git config --global user.email` to my personal address to push a fix to a side project, and I'd forgotten to switch it back. Forty-seven commits. Three weeks. Nobody noticed except me, that Saturday, by accident.

I didn't tell the client. I rewrote the history on a branch I hadn't pushed yet, force-pushed the rest, sent a humiliating message to a senior engineer I respected explaining what had happened, and went to bed. It bothered me for days. Not the mistake — the fact that I had no idea it was happening. There was no warning. No friction. Git happily wrote the wrong email on every commit and then GitHub happily accepted them.

So I went looking for the tool that prevents this. I assumed it existed.

---

## What I found instead

I found that git has a feature called `includeIf`, which lets you load different config based on which directory you're in. It works. It's also obscure — added in git 2.13 (2017), barely mentioned in the official docs, and every blog post about it has a comment from someone who got bitten by the trailing-slash gotcha. (The pattern after `gitdir:` is a glob, the trailing slash is significant, and there's no error if you get it wrong — just silent failure.)

I found that managing two SSH keys for two GitHub accounts is a separate problem with its own thirty-step tutorial, none of which are quite right, and that without `IdentitiesOnly yes` your SSH agent broadcasts every key you've ever loaded to every server you connect to — which is both a privacy concern and the source of the famous "Permission denied (publickey)" error nobody can debug.

I found that the GitHub CLI has its own idea of who you are, completely independent of git. You can be `git config user.email = work@company.com` while `gh pr create` runs as your personal account.

I found that the macOS keychain stores one credential per host, period, so if you use HTTPS with two GitHub accounts you're fundamentally fighting the OS.

I found dozens of tutorials, all subtly wrong, all assuming you'd remember to manually switch on every project change. None of them prevented the failure I'd just experienced. None of them refused to let me commit as the wrong person.

So I built the thing.

---

## What gitswitch does

`gitswitch` is the tool I needed at 11pm that Saturday.

It binds identities to directories. After a one-time setup, every `cd` is a switch — git, SSH, the GitHub CLI, and commit signing all line up automatically based on where you are.

It installs a pre-commit hook that **refuses commits where the identity is wrong**. If you forget to switch, you don't ship a wrong-author commit; you get a clear error message telling you what to do, and the commit doesn't happen. This is the part I needed three weeks ago.

It comes with a `doctor` command that proves, in one screen, that all the layers agree about who you are right now. And a `why` command that explains the magic — because automatic tools you can't inspect are just a different kind of bug waiting to happen.

It's a single binary. No Python runtime, no package manager spaghetti. Install in two seconds, works on every Mac/Linux/Windows machine, never rots.

```bash
brew install gitswitch
gitswitch init        # autodetect what you already have
gitswitch use work     ~/work
gitswitch use personal ~/code
gitswitch guard install
```

Done. From there, your identity is correct or your commits don't happen. There is no third option, which is exactly what I wanted.

---

## What I learned building this

Two things.

**First: most "user error" is tooling failure with extra steps.** When I committed forty-seven times as the wrong person, the system that allowed that to happen is the failure. I could call myself careless, but the next person to make this mistake won't be careless either — they'll be busy, or distracted, or in a hurry, or doing exactly what every git tutorial told them to do. A tool that prevents a class of error is a better answer than a habit you have to maintain.

**Second: the existing solutions are mostly true but practically false.** `includeIf` works. Per-account SSH host aliases work. `gh auth switch` works. The combination of all three, configured correctly, signed commits and all — that's what works. But assembling it manually is a 90-minute setup with four places to break, and you'll forget how when you onboard your next laptop. That assembly *is* the product.

---

## Try it

If you have multiple GitHub accounts, or you've ever committed under the wrong email and felt that specific kind of stomach-drop when you noticed — `gitswitch` is for you.

```bash
brew install gitswitch
gitswitch init
```

Source: [github.com/target-ops/gitswitch](https://github.com/target-ops/gitswitch)
Discussions: [github.com/target-ops/gitswitch/discussions](https://github.com/target-ops/gitswitch/discussions)

If it saves you from a moment like mine, tell me — I read every issue, and that's the kind of report that keeps me building.

— Ofir

---

*P.S. — If you've already committed as the wrong person and want to fix the history, `git filter-repo --email-callback` is the modern answer. Be careful, force-pushing rewrites is its own kind of disaster. The point of `gitswitch guard` is that you never have to.*
