import os
import click
import requests
from configs.utils import run_command
from configs.config import save_config

def set_global_git_user(username, email):
    """Function set_global_git_user."""
    run_command(f'git config --global user.name "{username}"')
    run_command(f'git config --global user.email "{email}"')

def add_user(config, vendor, username, email, key_path):
    """Function add user."""
    if vendor not in config:
        config[vendor] = {}
    config[vendor][username] = f"{email},{key_path}"
    save_config(config)

def delete_user(config, vendor, username):
    """Function delete_user."""
    if vendor in config and username in config[vendor]:
        del config[vendor][username]
        if not config[vendor]:
            del config[vendor]
        save_config(config)
    else:
        raise Exception(f"User {username} not found for vendor {vendor}")

def list_users(config):
    for vendor in config.sections():
        if vendor == "current":
            continue
        for username in config[vendor]:
            email, key_path = config[vendor][username].split(',')
            click.echo(f"vendor: " + click.style(vendor, fg="green") + " username: " + click.style(username, fg="green"))

def upload_ssh_key_to_vendor(vendor, username, key_path, token):
    public_key_path = f"{key_path}.pub"
    if not os.path.isfile(public_key_path):
        raise FileNotFoundError(f"The public key file {public_key_path} does not exist.")

    with open(public_key_path, 'r') as f:
        public_key = f.read()

    headers = {
        "Authorization": f"Bearer {token}",
        "Content-Type": "application/json"
    }
    data = {
        "title": f"{username}'s key",
        "key": public_key
    }

    if vendor == 'github':
        response = requests.post("https://api.github.com/user/keys", headers=headers, json=data)
    elif vendor == 'gitlab':
        response = requests.post("https://gitlab.com/api/v4/user/keys", headers=headers, json=data)
    else:
        raise Exception(f"Unsupported vendor: {vendor}")

    if response.status_code in [201, 200]:
        click.secho("Public key successfully uploaded.", fg='green')
    else:
        click.secho(f"Failed to upload public key: {response.status_code}", fg='red')
        click.secho(response.json(), fg='red')
