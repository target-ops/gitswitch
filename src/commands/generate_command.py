import click
from configs.ssh import generate_ssh_key
from configs.config import load_config

@click.group()
def generate():
    """Generate SSH keys for different Git vendors.

    Use these commands to create SSH keys for secure access to repositories on GitHub and GitLab.
    This simplifies authentication and enhances security for your Git operations.
    """
    pass

@generate.command()
@click.option('-e','--email',required=True, help='Email address of the user')
@click.option('-pk','--pub_key_path', required=True ,help='Path to the public key file')
def key(email, pub_key_path):
    """Generate a new SSH key.

    This command creates a new SSH key pair using the provided email address and 
    saves the public key to the specified file path. The SSH key is essential for 
    secure communication with Git repositories on platforms like GitHub and GitLab.

    Example usage:\n
    - gitswitch generate key -e email@example.com -pk /path/pubkey
    """
    generate_ssh_key(email, pub_key_path)
    click.secho(f"SSH key generated for {email}.", fg='green')
