import click
from click_help_colors import HelpColorsCommand,HelpColorsGroup
from configs.ssh import generate_ssh_key
from commands.add_command import uploadKey


@click.group(cls=HelpColorsGroup,help_headers_color='white',help_options_color='green')
def generate():
    """Generate SSH keys for different Git vendors.

    Use these commands to create SSH keys for secure access to repositories on GitHub and GitLab.                                                      
    This simplifies authentication and enhances security for your Git operations.
    """
    pass

@click.command(cls=HelpColorsCommand,help_options_color='green')
@click.option('-e','--email',required=True, help='Email address of the user')
@click.option('-v','--vendor', type=click.Choice(["github", "gitlab"]), help='Vendor name (namespaces the key file per vendor)')
@click.option('-u','--username', help='Username (defaults to the email local-part)')
def key(email, vendor, username):
    """Generate a new SSH key.

    This command creates a new SSH key pair using the provided email address.
    The SSH key is essential for secure communication with Git repositories
    on platforms like GitHub and GitLab.

    The key pair is saved in the ~/.ssh directory. When --vendor and
    --username are supplied, the file is named id_rsa_{vendor}_{username};
    otherwise it falls back to id_rsa_{username} (where {username} is the
    part of the email before the '@' symbol).

    Example usage:\n
    - gitswitch generate key -e email@example.com\n
    - gitswitch generate key -v github -u work -e work@company.com
    """
    generate_ssh_key(email, vendor=vendor, username=username)
    click.secho(f"SSH key generated for {email}.", fg='green')

generate.add_command(key)
generate.add_command(uploadKey)
