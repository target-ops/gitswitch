from config.config import Config
from config.git import GitConfig
import getpass
from logging_config import get_logger

logger = get_logger(__name__)

class AddCommand:    
    config = Config()
    git_config = GitConfig()

    @staticmethod
    def add_arguments(subparsers):
        parser = subparsers.add_parser('add', help='Add a new user')
        parser.add_argument('vendor', type=str, help='Git vendor (e.g., github, gitlab)')
        parser.add_argument('username', type=str, help='Git username')
        parser.add_argument('email', type=str, help='User email')
        parser.add_argument('key_path', type=str, help='Path to SSH key')
        parser.add_argument('--upload-key', action='store_true', help='Upload the SSH key to the vendor')
        parser.set_defaults(command_class=AddCommand)

    def execute(self, args):
        try:
            self.git_config.add_user(self.config.load, args.vendor, args.username, args.email, args.key_path)
            logger.info(f"User {args.username} added for vendor {args.vendor}.")
            if args.upload_key:
                token = getpass.getpass(f"Enter your {args.vendor} personal access token: ")
                self.git_config.upload_ssh_key_to_vendor(args.vendor, args.username, args.email, args.key_path, token)
        except Exception as e:
            logger.error(f"Error adding user: {e}")
