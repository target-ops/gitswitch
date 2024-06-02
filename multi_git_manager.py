import os
import subprocess
import configparser
import argparse
import requests
import getpass

CONFIG_FILE = 'config.ini'

def run_command(command):
    process = subprocess.Popen(command, shell=True, stdout=subprocess.PIPE, stderr=subprocess.PIPE)
    stdout, stderr = process.communicate()
    if process.returncode != 0:
        raise Exception(f"Command failed: {stderr.decode('utf-8')}")
    return stdout.decode('utf-8').strip()

def load_config():
    config = configparser.ConfigParser()
    if os.path.exists(CONFIG_FILE):
        config.read(CONFIG_FILE)
    return config

def save_config(config):
    with open(CONFIG_FILE, 'w') as configfile:
        config.write(configfile)

def add_user(config, vendor, username, email, key_path):
    if vendor not in config:
        config[vendor] = {}
    config[vendor][username] = f"{email},{key_path}"
    save_config(config)

def delete_user(config, vendor, username):
    if vendor in config and username in config[vendor]:
        del config[vendor][username]
        if not config[vendor]:
            del config[vendor]
        save_config(config)
    else:
        raise Exception(f"User {username} not found for vendor {vendor}")

def generate_ssh_key(email, key_path):
    run_command(f'ssh-keygen -t ed25519 -C "{email}" -f {key_path} -N ""')

def add_ssh_key_to_agent(key_path):
    run_command("eval $(ssh-agent -s)")
    run_command(f'ssh-add {key_path}')

def configure_repository(repo_path, host, username, email):
    os.chdir(repo_path)
    run_command(f'git remote set-url origin git@{host}:{username}/repo.git')
    run_command(f'git config user.name "{username}"')
    run_command(f'git config user.email "{email}"')

def list_users(config):
    for vendor in config.sections():
        print(f"{vendor}:")
        for username in config[vendor]:
            email, key_path = config[vendor][username].split(',')
            print(f"  Username: {username}, Email: {email}, SSH Key: {key_path}")

def upload_ssh_key_to_vendor(vendor, username, email, key_path, token):
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

def main():
    parser = argparse.ArgumentParser(description='Manage multiple Git users for different vendors.')
    subparsers = parser.add_subparsers(dest='command')

    # Add user command
    parser_add = subparsers.add_parser('add', help='Add a new user')
    parser_add.add_argument('vendor', type=str, help='Git vendor (e.g., github, gitlab)')
    parser_add.add_argument('username', type=str, help='Git username')
    parser_add.add_argument('email', type=str, help='User email')
    parser_add.add_argument('key_path', type=str, help='Path to SSH key')
    parser_add.add_argument('--upload-key', action='store_true', help='Upload the SSH key to the vendor')

    # Generate SSH key command
    parser_generate = subparsers.add_parser('generate-key', help='Generate a new SSH key')
    parser_generate.add_argument('email', type=str, help='Email for the SSH key')
    parser_generate.add_argument('key_path', type=str, help='Path to store the SSH key')

    # List users command
    parser_list = subparsers.add_parser('list', help='List all users')

    # Switch user command
    parser_switch = subparsers.add_parser('switch', help='Switch to a different user')
    parser_switch.add_argument('vendor', type=str, help='Git vendor (e.g., github, gitlab)')
    parser_switch.add_argument('username', type=str, help='Git username')
    parser_switch.add_argument('repo_path', type=str, help='Path to the repository')

    # Delete user command
    parser_delete = subparsers.add_parser('delete', help='Delete a user')
    parser_delete.add_argument('vendor', type=str, help='Git vendor (e.g., github, gitlab)')
    parser_delete.add_argument('username', type=str, help='Git username')

    args = parser.parse_args()

    config = load_config()

    if args.command == 'add':
        add_user(config, args.vendor, args.username, args.email, args.key_path)
        print(f"User {args.username} added for vendor {args.vendor}.")
        if args.upload_key:
            token = getpass.getpass(f"Enter your {args.vendor} personal access token: ")
            upload_ssh_key_to_vendor(args.vendor, args.username, args.email, args.key_path, token)

    elif args.command == 'generate-key':
        generate_ssh_key(args.email, args.key_path)
        add_ssh_key_to_agent(args.key_path)
        print(f"SSH key generated and added to agent for {args.email}.")

    elif args.command == 'list':
        list_users(config)

    elif args.command == 'switch':
        if args.vendor in config and args.username in config[args.vendor]:
            email, key_path = config[vendor][args.username].split(',')
            add_ssh_key_to_agent(key_path)
            configure_repository(args.repo_path, f"{args.vendor}.com", args.username, email)
            print(f"Switched to user {args.username} for vendor {args.vendor}.")
        else:
            print(f"User {args.username} not found for vendor {args.vendor}.")

    elif args.command == 'delete':
        try:
            delete_user(config, args.vendor, args.username)
            print(f"User {args.username} deleted for vendor {args.vendor}.")
        except Exception as e:
            print(e)

if __name__ == "__main__":
    main()
