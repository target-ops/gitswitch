# gitswitch — README demo recording.
#
# Run with:
#   vhs docs/demo.tape
#
# Produces docs/demo.gif (embedded in the README) plus
# docs/demo.mp4 (smaller, for Twitter / Show HN).
#
# The tape is fully self-contained: it runs setup invisibly,
# then records three scenes (about-to-commit, blocked, fixed)
# in a sandbox HOME so it never touches the maintainer's real
# config or credentials.

Output docs/demo.gif
Output docs/demo.mp4

Set FontSize 18
Set Width 1000
Set Height 540
Set Theme "Catppuccin Mocha"
Set TypingSpeed 35ms
Set PlaybackSpeed 1.0
Set Padding 24

# ----- hidden setup -----
# Sandbox $HOME so we don't touch real ~/.config or ~/.gitconfig.
# Seed a gitswitch identity bound to /tmp/demo and a wrong local
# email so the next commit attempt fails the guard check.
Hide
Type `export HOME=/tmp/gitswitch-demo && rm -rf $HOME && mkdir -p $HOME/.config/gitswitch $HOME/.ssh /tmp/demo`
Enter
Type `printf '%s\n' '{"version":1,"identities":[{"name":"work","email":"you@company.com","git_name":"Work","ssh_key":"/tmp/demo/key","signing_key":"/tmp/demo/key.pub","gh_account":"you-work","vendor":"github"}],"bindings":[{"directory":"/tmp/demo","identity":"work"}]}' > $HOME/.config/gitswitch/config.json`
Enter
Type `gitswitch use work /tmp/demo > /dev/null 2>&1 && gitswitch guard install > /dev/null 2>&1`
Enter
Type `cd /tmp/demo && git init -q && git config user.email you@gmail.com && git config user.name "Personal" && echo "auth flow change" > handler.go && git add handler.go && clear`
Enter
Type `PS1='$ '`
Enter
Sleep 200ms
Ctrl+L
Show

# ----- scene 1: about to commit -----
Type "git status"
Enter
Sleep 1500ms

Type "git commit -m 'fix the auth flow'"
Enter
Sleep 4500ms

# ----- scene 2: hold the blocked panel -----
# (no input — the red panel from `gitswitch guard` is on screen)

# ----- scene 3: fix and retry -----
Type "git config --unset user.email"
Enter
Sleep 800ms

Type "git commit -m 'fix the auth flow'"
Enter
Sleep 2500ms
