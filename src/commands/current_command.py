import click
from configs.config import load_config, get_current_user

@click.command(name='current')
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
        click.secho("Active user: " + click.style(username, fg="yellow") + " for vendor: " + click.style(vendor, fg="yellow"))
    else:
        click.secho("No active user set.", fg='yellow')
