import os
import re
import click

from configs.utils import run_command

SSH_CONFIG_FILE = os.path.expanduser('~/.ssh/config')

def generate_ssh_key(email):
    """Function generate_ssh_key."""
    if not re.match(r"[^@]+@[^@]+\.[^@]+", email):
        click.secho(f'Invalid email: {email}', fg='red')
        exit(0)
    
    username = email.split('@')[0]
    ssh_dir = os.path.expanduser('~/.ssh')
    key_path = os.path.join(ssh_dir, f'id_rsa_{username}')

    if not os.path.exists(ssh_dir):
        os.makedirs(ssh_dir)

    if os.path.exists(key_path):
        click.secho(f'SSH key already exists with the name: {key_path}', fg='red')
        exit(0)

    run_command(f'ssh-keygen -b 4096 -t rsa -C "{email}" -f {key_path} -N ""')


def update_ssh_config(vendor, key_path):
    """Function update_ssh_config."""
    host_entry = f"""
Host {vendor}.com
    HostName {vendor}.com
    PreferredAuthentications publickey
    IdentityFile {key_path}
    """

    with open(SSH_CONFIG_FILE, 'w') as f:
        f.write(host_entry)
