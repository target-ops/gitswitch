import shutil
import subprocess
import click


def is_gh_installed():
    """Return True if the GitHub CLI (`gh`) is on PATH."""
    return shutil.which('gh') is not None


def gh_logged_in_users():
    """Return the set of usernames currently authenticated with `gh`."""
    if not is_gh_installed():
        return set()
    result = subprocess.run(
        ['gh', 'auth', 'status'],
        capture_output=True, text=True,
    )
    # `gh auth status` prints lines like "  - Logged in to github.com account <user>"
    # on stderr (older versions) or stdout (newer). Parse both.
    output = (result.stdout or '') + '\n' + (result.stderr or '')
    users = set()
    for line in output.splitlines():
        line = line.strip()
        if 'account' in line:
            parts = line.split('account')
            tail = parts[-1].strip().split()
            if tail:
                users.add(tail[0].strip('()'))
    return users


def switch_gh_user(username):
    """Switch the active `gh` CLI account to `username`.

    Quietly skips if `gh` is not installed. If the user is not yet logged in,
    surfaces a hint to run `gh auth login`.
    """
    if not is_gh_installed():
        click.secho(
            "gh CLI not installed; skipping GitHub CLI account switch.",
            fg='yellow',
        )
        return False

    result = subprocess.run(
        ['gh', 'auth', 'switch', '--user', username],
        capture_output=True, text=True,
    )
    if result.returncode == 0:
        click.secho(f"GitHub CLI: switched to {username}.", fg='green')
        return True

    click.secho(
        f"GitHub CLI: could not switch to '{username}'. "
        f"Run `gh auth login` to add this account, then retry.",
        fg='yellow',
    )
    if result.stderr.strip():
        click.secho(result.stderr.strip(), fg='yellow')
    return False


def login_gh_user():
    """Run an interactive `gh auth login` so the new account is added to gh."""
    if not is_gh_installed():
        click.secho(
            "gh CLI not installed; skipping `gh auth login`.",
            fg='yellow',
        )
        return False
    # Use subprocess.call so stdin/stdout stay attached to the user's terminal.
    return subprocess.call(['gh', 'auth', 'login']) == 0
