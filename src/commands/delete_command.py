import click
from click_help_colors import HelpColorsCommand
from configs.config import load_config
from configs.git import delete_user

@click.command(name='del',
    cls=HelpColorsCommand,
    help_options_color='green'
)
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