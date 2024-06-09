import click
import getpass
from click_help_colors import HelpColorsGroup, HelpColorsCommand
from src.config import load_config, set_current_user, get_current_user
from src.ssh import generate_ssh_key, update_ssh_config
from src.git import set_global_git_user, add_user, delete_user, list_users, upload_ssh_key_to_vendor

version = "0.0.1"

@click.version_option(version, "--version", "-V", message='%(version)s')

@click.group(
    cls=HelpColorsGroup,
    help_headers_color='white',
    help_options_color='green',
)

def cli():
    """gitswitch cli \n
    Easily manage multiple Git user profiles for different vendors. 
    Seamlessly switch between configurations, avoid commit errors, and streamline your workflow. 
    Perfect for developers juggling various projects and clients.."""
    pass

@cli.group()
def add():
    """Commands to manage Git users and their SSH keys.
    Use these commands to add new Git user profiles, configure user details,
    and upload SSH keys to GitHub or GitLab for secure authentication.
    """
    pass

@cli.group()
def generate():
    """Generate SSH keys for different Git vendors.

    Use these commands to create SSH keys for secure access to repositories on GitHub and GitLab.
    This simplifies authentication and enhances security for your Git operations.
    """
    pass


@cli.command()
def list():
    """List all configured Git users.

    This command displays a list of all Git user profiles configured in the system.
    It includes details for each user, such as the vendor (GitHub or GitLab) and the username.

    Example usage:\n
    - gitswitch list
    """
    config = load_config()
    list_users(config)
    # pass

@add.command()
@click.option('-v','--vendor',prompt='Vendor name', required=True, type=click.Choice(["github", "gitlab"]),help='Vendor name')
@click.option('-u','--username',prompt='Username',required=True, help='Username of the user')
@click.option('-e','--email',prompt='Email',required=True, help='Email address of the user')
@click.option('-pk','--pub_key_path', prompt='Public Key Path',required=True ,help='Path to the public key file')
def user(vendor, username, email, pub_key_path):
    """Add a new user.

    This command allows you to add a new Git user profile for a specified vendor.
    You will be prompted to enter the vendor name (GitHub or GitLab), the username,
    email address, and the path to the public key file for SSH authentication.

    Example usage:\n
    - gitswitch add user -v github -u username -e email@example.com -pk /path/to/public/key
    """
    config = load_config()
    add_user(config, vendor, username, email, pub_key_path)
    click.secho(f"User: {username} added for vendor {vendor}.", fg='green')

@add.command()
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

@cli.command()
@click.option('-v','--vendor', prompt='Vendor name', type=click.Choice(["github", "gitlab"]), help='Vendor name')
@click.option('-u','--username',prompt='Username',required=True, help='Username of the user')
def switch(vendor, username):
    """Switch to a different user.

    This command allows you to switch between different Git user profiles for a specified vendor.
    You will be prompted to enter the vendor name (GitHub or GitLab) and the username you want to switch to.
    It updates the global Git user configuration and SSH settings to match the selected user.

    Example usage:\n
    - gitswitch switch -v github -u username
    """
    config = load_config()
    if vendor in config and username in config[vendor]:
        email, key_path = config[vendor][username].split(',')
        set_global_git_user(username, email)
        update_ssh_config(vendor, key_path)
        set_current_user(config, vendor, username)
        click.secho(f"Switched to: " + click.style(username, fg="green"))
    else:
        click.secho(f"User {username} not found for vendor {vendor}.", fg='red')

@cli.command()
@click.option('-v','--vendor', prompt='Vendor name', required=True, type=click.Choice(["github", "gitlab"]),help='Vendor name')
@click.option('-u','--username',prompt='Username', required=True, help='Username of the user')
def delete(vendor, username):
    """Delete a user from the configuration file.

    This command removes a specified Git user profile from the configuration file.
    You will be prompted to enter the vendor name (GitHub or GitLab) and the username to be deleted.
    It ensures the selected user is no longer available for switching or other operations.

    Example usage:\n
    - gitswitch delete -v github -u username
    """
    config = load_config()
    try:
        delete_user(config, vendor, username)
        click.secho(f"User {username} deleted for vendor {vendor}.", fg='green')
    except Exception as e:
        click.secho(str(e), fg='red')

@cli.command()
def current():
    """Show the current active Git user.

    This command displays the currently active Git user profile, including the vendor (GitHub or GitLab)
    and the username. If no user is set as active, it informs you accordingly.

    Example usage:\n
    - gitswitch current
    """
    config = load_config()
    vendor, username = get_current_user(config)
    if vendor and username:
        click.secho(f"Active user: "+ click.style(username, fg="green")+ " for vendor: "+ click.style(vendor, fg="green"))
    else:
        click.secho("No active user set.", fg='yellow')


# add.add_command(uploadKey)
# add.add_command(current)
generate.add_command(uploadKey)

if __name__ == "__main__":
    cli()
