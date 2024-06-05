import argparse
import getpass
from config import load_config, set_current_user, get_current_user
from ssh import generate_ssh_key, update_ssh_config
from git import set_global_git_user, add_user, delete_user, list_users, upload_ssh_key_to_vendor

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

    # Delete user command
    parser_delete = subparsers.add_parser('delete', help='Delete a user')
    parser_delete.add_argument('vendor', type=str, help='Git vendor (e.g., github, gitlab)')
    parser_delete.add_argument('username', type=str, help='Git username')

    # Current user command
    parser_current = subparsers.add_parser('current', help='Show current active user')

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
        print(f"SSH key generated for {args.email}.")

    elif args.command == 'list':
        list_users(config)

    elif args.command == 'switch':
        if args.vendor in config and args.username in config[args.vendor]:
            email, key_path = config[args.vendor][args.username].split(',')
            set_global_git_user(args.username, email)
            update_ssh_config(args.vendor, key_path)
            set_current_user(config, args.vendor, args.username)
            print(f"Switched to user {args.username} for vendor {args.vendor}.")
        else:
            print(f"User {args.username} not found for vendor {args.vendor}.")

    elif args.command == 'delete':
        try:
            delete_user(config, args.vendor, args.username)
            print(f"User {args.username} deleted for vendor {args.vendor}.")
        except Exception as e:
            print(e)

    elif args.command == 'current':
        vendor, username = get_current_user(config)
        if vendor and username:
            print(f"Current active user: {username} for vendor {vendor}")
        else:
            print("No active user set.")

if __name__ == "__main__":
    main()
