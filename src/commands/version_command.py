import click

version = "0.0.1"

@click.command()
@click.version_option(version, "--version", "-v", message='%(version)s')
def version_command():
    pass