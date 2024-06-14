import click
from configs.config import load_config
from configs.git import list_users

@click.command(name='list')
def list_cm():
    """List all configured Git users.

    This command displays a list of all Git user profiles configured in the system.
    It includes details for each user, such as the vendor (GitHub or GitLab) and the username.

    Example usage:\n
    - gitswitch list
    """
    config = load_config()
    list_users(config)
    # pass