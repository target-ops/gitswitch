from config.config import Config
from config.git import GitConfig
from logging_config import get_logger

logger = get_logger(__name__)

class ListCommand:
    config = Config()
    git_config = GitConfig()
    
    @staticmethod
    def add_arguments(subparsers):
        parser = subparsers.add_parser('list', help='List all users')
        parser.set_defaults(command_class=ListCommand)

    def execute(self, args):
        try:
            self.git_config.list_users(self.config.load())
        except Exception as e:
            logger.error(f"Error listing users: {e}")
