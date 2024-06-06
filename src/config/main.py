import argparse
import getpass
import argcomplete
import logging
from config import load_config, set_current_user, get_current_user
from ssh import generate_ssh_key, update_ssh_config
from git import set_global_git_user, add_user, delete_user, list_users, upload_ssh_key_to_vendor

# Configure logging
logging.basicConfig(level=logging.INFO)
logger = logging.getLogger(__name__)

class GitSwitch:
    def __init__(self):
        self.config = load_config()
        self.parser = self.create_parser()

    def create_parser(self):
        parser = argparse.ArgumentParser(
            description='Manage multiple Git users for different vendors.',
            epilog="""
Examples:
    # Add a new user
    gitswitch add github myusername myemail@example.com /path/to/key

    # Generate a new SSH key
    gitswitch generate-key myemail@example.com /path/to/key

    # List all users
    gitswitch list

    # Switch to a different user
    gitswitch switch github myusername

    # Delete a user
    gitswitch delete github myusername

    # Show current active user
    gitswitch current
            """,
            formatter_class=argparse.RawDescriptionHelpFormatter
        )
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

        argcomplete.autocomplete(parser)
        return parser

    def add_user_command(self, args):
        """
        Adds a new user and optionally uploads the SSH key to the vendor.
        
        :param args: Parsed command-line arguments
        """
        try:
            add_user(self.config, args.vendor, args.username, args.email, args.key_path)
            logger.info(f"User {args.username} added for vendor {args.vendor}.")
            if args.upload_key:
                token = getpass.getpass(f"Enter your {args.vendor} personal access token: ")
                upload_ssh_key_to_vendor(args.vendor, args.username, args.email, args.key_path, token)
        except Exception as e:
            logger.error(f"Error adding user: {e}")

    def generate_ssh_key_command(self, args):
        """
        Generates a new SSH key.
        
        :param args: Parsed command-line arguments
        """
        try:
            generate_ssh_key(args.email, args.key_path)
            logger.info(f"SSH key generated for {args.email}.")
        except Exception as e:
            logger.error(f"Error generating SSH key: {e}")

    def list_users_command(self):
        """
        Lists all users.
        """
        try:
            list_users(self.config)
        except Exception as e:
            logger.error(f"Error listing users: {e}")

    def switch_user_command(self, args):
        """
        Switches to a different user.
        
        :param args: Parsed command-line arguments
        """
        try:
            if args.vendor in self.config and args.username in self.config[args.vendor]:
                email, key_path = self.config[args.vendor][args.username].split(',')
                set_global_git_user(args.username, email)
                update_ssh_config(args.vendor, key_path)
                set_current_user(self.config, args.vendor, args.username)
                logger.info(f"Switched to user {args.username} for vendor {args.vendor}.")
            else:
                logger.error(f"User {args.username} not found for vendor {args.vendor}.")
        except Exception as e:
            logger.error(f"Error switching user: {e}")

    def delete_user_command(self, args):
        """
        Deletes a user.
        
        :param args: Parsed command-line arguments
        """
        try:
            delete_user(self.config, args.vendor, args.username)
            logger.info(f"User {args.username} deleted for vendor {args.vendor}.")
        except Exception as e:
            logger.error(f"Error deleting user: {e}")

    def current_user_command(self):
        """
        Shows the current active user.
        """
        try:
            vendor, username = get_current_user(self.config)
            if vendor and username:
                logger.info(f"Current active user: {username} for vendor {vendor}")
            else:
                logger.info("No active user set.")
        except Exception as e:
            logger.error(f"Error getting current user: {e}")

    def execute(self):
        """
        Parses arguments and executes the corresponding command.
        """
        args = self.parser.parse_args()
        if args.command is None:
            self.parser.print_help()
        elif args.command == 'add':
            self.add_user_command(args)
        elif args.command == 'generate-key':
            self.generate_ssh_key_command(args)
        elif args.command == 'list':
            self.list_users_command()
        elif args.command == 'switch':
            self.switch_user_command(args)
        elif args.command == 'delete':
            self.delete_user_command(args)
        elif args.command == 'current':
            self.current_user_command()

if __name__ == "__main__":
    GitSwitch().execute()
