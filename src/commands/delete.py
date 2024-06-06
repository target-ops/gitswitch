from config.config import Config
from config.git import GitConfig
from logging_config import get_logger

logger = get_logger(__name__)

class DeleteCommand:
    config = Config()
    git_config = GitConfig()
    
    @staticmethod
    def add_arguments(subparsers):
        parser = subparsers.add_parser('delete', help='Delete a user')
        parser.add_argument('vendor', type=str, help='Git vendor (e.g., github, gitlab)')
        parser.add_argument('username', type=str, help='Git username')
        parser.set_defaults(command_class=DeleteCommand)

    def execute(self, args):
        try:
            self.git_config.delete_user(self.config.load(), args.vendor, args.username)
            logger.info(f"User {args.username} deleted for vendor {args.vendor}.")
        except Exception as e:
            logger.error(f"Error deleting user: {e}")
