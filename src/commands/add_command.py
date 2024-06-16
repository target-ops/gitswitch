import getpass
import click
from click_help_colors import HelpColorsCommand,HelpColorsGroup
from configs.config import load_config
from configs.git import add_user, upload_ssh_key_to_vendor
from configs.ssh import generate_ssh_key


@click.group(cls=HelpColorsGroup,help_headers_color='white',help_options_color='green')
def add():
    """Commands to manage Git users and their SSH keys.                                 
    Use these commands to add new Git user profiles, configure user details,
    and upload SSH keys to GitHub or GitLab for secure authentication.
    """
    pass

@add.command(name='uploadkey',
    cls=HelpColorsCommand,
    help_options_color='green'
)
@click.option('-v','--vendor',prompt='Vendor name', required=True, type=click.Choice(["github", "gitlab"]),help='Vendor name')
@click.option('-u','--username',prompt='Username',required=True, help='Username of the user')
@click.option('-pk','--pub_key_path', prompt='Public Key Path',required=True ,help='Path to the public key file')
def uploadKey(vendor,username, pub_key_path):
    """Upload the SSH key to the vendor.

    This command allows you to upload an SSH key to GitHub or GitLab for the specified user.
    You will be prompted to enter your personal access token for authentication.

    Example usage:\n
    - gitswitch add uploadkey -v github -u username -pk /path/to/public/key
    """
    token = getpass.getpass(f"Enter your {vendor} personal access token: ")
    upload_ssh_key_to_vendor(vendor, username, pub_key_path, token)

@add.command(name='user',
    cls=HelpColorsCommand,
    help_options_color='green'
)
@click.option('-v','--vendor',prompt='Vendor name', required=True, type=click.Choice(["github", "gitlab"]),help='Vendor name')
@click.option('-u','--username',prompt='Username',required=True, help='Username of the user')
@click.option('-e','--email',prompt='Email',required=True, help='Email address of the user')
@click.option('-g','--generate',prompt='Would you like to generate an SSH key?',help='Generate SSH Key', default=False,is_flag=True)
def user(vendor, username, email, generate):
    """Add a new user.

    This command allows you to add a new Git user profile for a specified vendor.                                                                    
    You will be prompted to enter the vendor name (GitHub or GitLab), the username,
    and the email address. If you choose to generate an SSH key,                                            
    it will be done automatically.

    Example usage:\n
    - gitswitch add user -v github -u username -e email@example.com
    """
    config = load_config()
    print(generate)
    if generate:
        print("inside generate")
        generate_ssh_key(email)
    add_user(config, vendor, username, email)
    click.secho(f"User: {username} added for vendor {vendor}.", fg='green')
