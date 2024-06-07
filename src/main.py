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
    """Manage multiple Git users for different vendors."""
    pass

@cli.group()
def add():
    """Manage Git users."""
    pass

@cli.group()
def generate():
    """Generate Key for different vendors."""
    pass


@cli.command()
def list():
    """List all users."""
    config = load_config()
    list_users(config)
    # pass

@add.command()
@click.option('-v','--vendor',prompt='Vendor name', required=True, type=click.Choice(["github", "gitlab"]),help='Vendor name')
@click.option('-u','--username',prompt='Username',required=True, help='Username of the user')
@click.option('-e','--email',prompt='Email',required=True, help='Email address of the user')
@click.option('-pk','--pub_key_path', prompt='Public Key Path',required=True ,help='Path to the public key file')
def user(vendor, username, email, pub_key_path):
    """Add a new user."""
    config = load_config()
    add_user(config, vendor, username, email, pub_key_path)
    click.secho(f"User: {username} added for vendor {vendor}.", fg='green')

@add.command()
@click.option('-v','--vendor',prompt='Vendor name', required=True, type=click.Choice(["github", "gitlab"]),help='Vendor name')
@click.option('-u','--username',prompt='Username',required=True, help='Username of the user')
@click.option('-pk','--pub_key_path', prompt='Public Key Path',required=True ,help='Path to the public key file')
def uploadKey(vendor,username, pub_key_path):
    """Upload the SSH key to the vendor.""" 
    token = getpass.getpass(f"Enter your {vendor} personal access token: ")
    upload_ssh_key_to_vendor(vendor, username, pub_key_path, token)

@generate.command()
@click.option('-e','--email',required=True, help='Email address of the user')
@click.option('-pk','--pub_key_path', required=True ,help='Path to the public key file')
def key(email, pub_key_path):
    """Generate a new SSH key."""
    generate_ssh_key(email, pub_key_path)
    click.secho(f"SSH key generated for {email}.", fg='green')

@cli.command()
@click.option('-v','--vendor', prompt='Vendor name', type=click.Choice(["github", "gitlab"]), help='Vendor name')
@click.option('-u','--username',prompt='Username',required=True, help='Username of the user')
def switch(vendor, username):
    """Switch to a different user."""
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
    """Delete user from config file."""
    config = load_config()
    try:
        delete_user(config, vendor, username)
        click.secho(f"User {username} deleted for vendor {vendor}.", fg='green')
    except Exception as e:
        click.secho(str(e), fg='red')

@cli.command()
def current():
    """Show current active user."""
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
