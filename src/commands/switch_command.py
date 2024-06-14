import click
from configs.config import load_config,set_current_user
from configs.git import set_global_git_user
from configs.ssh import update_ssh_config

@click.group()
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