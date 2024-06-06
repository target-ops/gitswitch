import os
import subprocess
from src.utils import run_command

class SSHConfig:
    SSH_CONFIG_FILE = os.path.expanduser('~/.ssh/config')

    @staticmethod
    def generate_ssh_key(email, key_path):
        """Generate an SSH key."""
        run_command(f'ssh-keygen -b 4096 -t rsa -C "{email}" -f {key_path} -N ""')

    @staticmethod
    def update_ssh_config(vendor, key_path):
        """Update the SSH config file."""
        host_entry = f"""
    Host {vendor}.com
        HostName {vendor}.com
        PreferredAuthentications publickey
        IdentityFile {key_path}
        """

        with open(SSHConfig.SSH_CONFIG_FILE, 'w') as f:
            f.write(host_entry)