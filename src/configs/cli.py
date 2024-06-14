import click
from click_help_colors import HelpColorsGroup

# Import the command groups from the other files
from commands.add_command import add, uploadKey
from commands.generate_command import generate
from commands.list_command import list_cm
from commands.switch_command import switch
from commands.delete_command import delete
from commands.current_command import current


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

cli.add_command(add)
cli.add_command(generate)
cli.add_command(list_cm)
cli.add_command(switch)
cli.add_command(delete)
cli.add_command(current)
generate.add_command(uploadKey)
