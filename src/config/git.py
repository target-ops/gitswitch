import os
import requests
from src.utils import run_command
from .config import Config

class GitConfig:
    @staticmethod
    def set_global_git_user(username, email):
        """Set the global Git user."""
        run_command(f'git config --global user.name "{username}"')
        run_command(f'git config --global user.email "{email}"')

    @staticmethod
    def add_user(config, vendor, username, email, key_path):
        """Add a new user to the config."""
        if vendor not in config:
            config[vendor] = {}
        config[vendor][username] = f"{email},{key_path}"
        Config.save(config)

    @staticmethod
    def delete_user(config, vendor, username):
        """Delete a user from the config."""
        if vendor in config and username in config[vendor]:
            del config[vendor][username]
            if not config[vendor]:
                del config[vendor]
            Config.save(config)
        else:
            raise Exception(f"User {username} not found for vendor {vendor}")

    @staticmethod
    def list_users(config):
        """List all users from the config."""
        for vendor in config.sections():
            if vendor == "current":
                continue
            for username in config[vendor]:
                print(f"vendor: {vendor}, username: {username}")

    @staticmethod
    def upload_ssh_key_to_vendor(vendor, username, email, key_path, token):
        """Function upload_ssh_key_to_vendor."""
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
            print("Public key successfully uploaded.")
        else:
            print(f"Failed to upload public key: {response.status_code}")
            print(response.json())