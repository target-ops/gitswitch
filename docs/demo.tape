# gitswitch — README demo recording.
#
# Run with:
#   vhs docs/demo.tape
#
# Produces docs/demo.gif (embedded in the README) plus
# docs/demo.mp4 (smaller, for Twitter / Show HN / LinkedIn).
#
# The tape is fully self-contained: it runs setup invisibly,
# then records four scenes (state, commit attempt, blocked,
# fixed + verified) in a sandbox $HOME so it never touches
# the maintainer's real config or credentials.

Output docs/demo.gif
Output docs/demo.mp4

Set FontSize 18
Set Width 1100
Set Height 700
Set Theme "Catppuccin Mocha"
Set TypingSpeed 60ms
Set PlaybackSpeed 1.0
Set Padding 24

# ---------------------------------------------------------------
# Hidden setup. Sandbox $HOME, seed two identities + bindings, set
# up a work repo whose local user.email is deliberately wrong so
# the next commit will trip the guard.
# ---------------------------------------------------------------
Hide
Type `export HOME=/tmp/gitswitch-demo && rm -rf $HOME && mkdir -p $HOME/.config/gitswitch $HOME/.ssh /tmp/demo/work /tmp/demo/personal`
Enter
Type `printf '%s\n' '{"version":1,"identities":[{"name":"work","email":"you@company.com","git_name":"Work","ssh_key":"/tmp/demo/key","signing_key":"/tmp/demo/key.pub","gh_account":"you-work","vendor":"github"},{"name":"personal","email":"you@gmail.com","git_name":"Personal","ssh_key":"/tmp/demo/key","signing_key":"/tmp/demo/key.pub","gh_account":"you-personal","vendor":"github"}],"bindings":[{"directory":"/tmp/demo/work","identity":"work"},{"directory":"/tmp/demo/personal","identity":"personal"}]}' > $HOME/.config/gitswitch/config.json`
Enter
Type `gitswitch use work /tmp/demo/work > /dev/null 2>&1 && gitswitch use personal /tmp/demo/personal > /dev/null 2>&1 && gitswitch guard install > /dev/null 2>&1`
Enter
# The per-identity gitconfig that 'gitswitch use' just wrote turns on
# SSH commit signing by default — but the sandbox SSH key is fake, so
# real `git commit` would fail. Strip the signing/sshCommand bits for
# the demo so the second-commit-succeeds frame actually succeeds.
Type `for f in $HOME/.config/gitswitch/identities/*.gitconfig; do printf '[user]\nname = %s\nemail = %s\n' "$(grep -m1 '^    name' $f | sed 's/^    name = //')" "$(grep -m1 '^    email' $f | sed 's/^    email = //')" > $f; done`
Enter
Type `cd /tmp/demo/work && git init -q && git config user.email you@gmail.com && git config user.name "Personal" && echo "auth flow change" > handler.go && git add handler.go`
Enter
Type `PS1='$ '`
Enter
Sleep 200ms
Ctrl+L
Show

# ---------------------------------------------------------------
# Scene 1 — establish: the user has two identities configured.
# ---------------------------------------------------------------
Type "gitswitch list"
Enter
Sleep 3500ms

# ---------------------------------------------------------------
# Scene 2 — about to commit, looks normal.
# ---------------------------------------------------------------
Type "git status"
Enter
Sleep 2200ms

Type "git commit -m 'fix the auth flow'"
Enter

# ---------------------------------------------------------------
# Scene 3 — the marquee moment. Hold long enough for a viewer to
# actually read the red panel.
# ---------------------------------------------------------------
Sleep 5500ms

# ---------------------------------------------------------------
# Scene 4 — the fix. Unset the local override so the includeIf
# block from `gitswitch use` provides the right email; then retry.
# ---------------------------------------------------------------
Type "git config --unset user.email"
Enter
Sleep 1200ms

Type "git commit -m 'fix the auth flow'"
Enter
Sleep 2500ms

# ---------------------------------------------------------------
# Scene 5 — victory lap. `why` confirms the active identity for
# this directory. (Skipping `doctor` here because it depends on
# real ssh-auth + gh — neither work in the sandbox; `why` works
# entirely from local config.)
# ---------------------------------------------------------------
Type "gitswitch why"
Enter
Sleep 4500ms
