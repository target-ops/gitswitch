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
            print(f"{vendor}:")
            for username in config[vendor]:
                email, key_path = config[vendor][username].split(',')
                print(f"  Username: {username}, Email: {email}, SSH Key: {key_path}")

    @staticmethod
    def upload_ssh_key_to_vendor(vendor, username, email, key_path, token):
        """Upload SSH key to the vendor's platform."""
        public_key_path = f"{key_path}.pub"
        if not os.path.isfile(public_key_path):
            raise FileNotFoundError(f"The public key file {public_key_path} should exist.")
