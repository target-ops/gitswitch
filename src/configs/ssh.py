import os
import re
import click

from configs.utils import run_command

SSH_CONFIG_FILE = os.path.expanduser('~/.ssh/config')


def default_key_path(vendor, username):
    """Resolve the key path for a (vendor, username).

    Prefers the namespaced path id_rsa_{vendor}_{username}. Falls back to the
    legacy id_rsa_{username} when only that exists, so users who upgraded from
    earlier versions are not broken.
    """
    ssh_dir = os.path.expanduser('~/.ssh')
    namespaced = os.path.join(ssh_dir, f'id_rsa_{vendor}_{username}') if vendor else None
    legacy = os.path.join(ssh_dir, f'id_rsa_{username}')
    if namespaced and os.path.exists(namespaced):
        return namespaced
    if os.path.exists(legacy):
        return legacy
    return namespaced or legacy


def generate_ssh_key(email, vendor=None, username=None):
    """Function generate_ssh_key."""
    if not re.match(r"[^@]+@[^@]+\.[^@]+", email):
        click.secho(f'Invalid email: {email}', fg='red')
        exit(0)

    if not username:
        username = email.split('@')[0]
    ssh_dir = os.path.expanduser('~/.ssh')
    if vendor:
        key_path = os.path.join(ssh_dir, f'id_rsa_{vendor}_{username}')
    else:
        key_path = os.path.join(ssh_dir, f'id_rsa_{username}')

    if not os.path.exists(ssh_dir):
        os.makedirs(ssh_dir)

    if os.path.exists(key_path):
        click.secho(f'SSH key already exists with the name: {key_path}', fg='red')
        exit(0)

    run_command(f'ssh-keygen -b 4096 -t rsa -C "{email}" -f {key_path} -N ""')
    return key_path


def _build_host_block(vendor, key_path):
    host = f"{vendor}.com"
    return (
        f"Host {host}\n"
        f"    HostName {host}\n"
        f"    PreferredAuthentications publickey\n"
        f"    IdentityFile {key_path}\n"
        f"    IdentitiesOnly yes\n"
    )


def update_ssh_config(vendor, key_path):
    """Replace just the Host block for this vendor; preserve everything else."""
    host = f"{vendor}.com"
    new_block = _build_host_block(vendor, key_path)

    ssh_dir = os.path.dirname(SSH_CONFIG_FILE)
    if not os.path.exists(ssh_dir):
        os.makedirs(ssh_dir, mode=0o700)

    if not os.path.exists(SSH_CONFIG_FILE):
        with open(SSH_CONFIG_FILE, 'w') as f:
            f.write(new_block)
        os.chmod(SSH_CONFIG_FILE, 0o600)
        return

    with open(SSH_CONFIG_FILE, 'r') as f:
        content = f.read()

    # Strip any existing block whose Host line matches this vendor.
    # A block runs from its `Host` line up to (but not including) the next
    # top-level `Host` line or end of file.
    pattern = re.compile(
        r'^[ \t]*Host[ \t]+[^\n]*\b' + re.escape(host) + r'\b[^\n]*\n'
        r'(?:(?![ \t]*Host[ \t])[^\n]*\n?)*',
        re.MULTILINE,
    )
    content = pattern.sub('', content)

    if content and not content.endswith('\n'):
        content += '\n'
    if content and not content.endswith('\n\n'):
        content += '\n'
    content += new_block

    with open(SSH_CONFIG_FILE, 'w') as f:
        f.write(content)
    os.chmod(SSH_CONFIG_FILE, 0o600)
