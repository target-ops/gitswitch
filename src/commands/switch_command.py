import click
import inquirer
from click_help_colors import HelpColorsCommand
from configs.config import load_config,set_current_user
from configs.git import set_global_git_user
from configs.ssh import update_ssh_config

# @click.command(cls=HelpColorsCommand,help_options_color='green')
# @click.option('-v','--vendor',required=True,type=click.Choice(["github", "gitlab"]), help='Vendor name')
# @click.option('-u','--username',required=True, help='Username')
# def switch(vendor, username):
#     """Switch to a different user.

#     This command allows you to switch between different Git user profiles for a specified vendor.
#     You will be prompted to enter the vendor name (GitHub or GitLab) and the username you want to switch to.                                
#     It updates the global Git user configuration and SSH settings to match the selected user.

#     Example usage:\n
#     - gitswitch switch -v github -u username
#     """
#     config = load_config()
#     if vendor in config and username in config[vendor]:
#         email, key_path = config[vendor][username].split(',')
#         set_global_git_user(username, email)
#         update_ssh_config(vendor, key_path)
#         set_current_user(config, vendor, username)
#         click.secho("Switched to: " + click.style(username, fg="green"))
#     else:
#         click.secho(f"User {username} not found for vendor {vendor}.", fg='red')

@click.command(cls=HelpColorsCommand, help_options_color='green')
@click.option('-v','--vendor', help='Vendor name')
@click.option('-u','--username', help='Username')
def switch(vendor, username):
    """Switch to a different user.

    This command allows you to switch between different Git user profiles for a specified vendor.
    You will be prompted to enter the vendor name (GitHub or GitLab) and the username you want to switch to.                                
    It updates the global Git user configuration and SSH settings to match the selected user.

    Example usage:\n
    - gitswitch switch
    """
    
    config = load_config()

    if vendor is None:
        vendors = list(config.keys())
        if 'current' in vendors:
            vendors.remove('DEFAULT')
            vendors.remove('current')

        vendor_questions = [inquirer.List('vendor',message="Select a vendor",choices=vendors)]
        vendor_answers = inquirer.prompt(vendor_questions)
        vendor = vendor_answers['vendor']

    if vendor in config:
        if username is None:
            usernames = list(config[vendor].keys())
            username_questions = [inquirer.List('username',message="Select a username",choices=usernames)]
            username_answers = inquirer.prompt(username_questions)
            username = username_answers['username']

        if username in config[vendor]:
            email, key_path = config[vendor][username].split(',')
            set_global_git_user(username, email)
            update_ssh_config(vendor, key_path)
            set_current_user(config, vendor, username)
            click.secho("Switched to: " + click.style(username, fg="green"))
        else:
            click.secho(f"Username {username} not found for vendor {vendor}.", fg='red')
    else:
        click.secho(f"Vendor {vendor} not found.", fg='red')