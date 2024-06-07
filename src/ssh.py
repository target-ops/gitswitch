import os
from utils import run_command

SSH_CONFIG_FILE = os.path.expanduser('~/.ssh/config')

def generate_ssh_key(email, key_path):
    """Function generate_ssh_key."""
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
