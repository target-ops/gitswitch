from config.config import Config
from logging_config import get_logger

logger = get_logger(__name__)

class CurrentCommand:
    config = Config()

    @staticmethod
    def add_arguments(subparsers):
        parser = subparsers.add_parser('current', help='Show current active user')
        parser.set_defaults(command_class=CurrentCommand)

    def execute(self, args):
        config = self.config.load()
        vendor, username = self.config.get_current_user(config)
        if vendor and username:
            logger.info(f"Current active user: {username} for vendor {vendor}")
        else:
            logger.info("No active user set.")
