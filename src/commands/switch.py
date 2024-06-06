from config.config import Config
from config.git import GitConfig
from config.ssh import SSHConfig
from logging_config import get_logger

logger = get_logger(__name__)

class SwitchCommand:
    @staticmethod
    def add_arguments(subparsers):
        parser = subparsers.add_parser('switch', help='Switch to a different user')
        parser.add_argument('vendor', type=str, help='Git vendor (e.g., github, gitlab)')
        parser.add_argument('username', type=str, help='Git username')
        parser.set_defaults(command_class=SwitchCommand)

    def execute(self, args):
        config = Config.load()
        try:
            if args.vendor in config and args.username in config[args.vendor]:
                email, key_path = config[args.vendor][args.username].split(',')
                GitConfig.set_global_git_user(args.username, email)
                SSHConfig.update_ssh_config(args.vendor, key_path)
                Config.set_current_user(config, args.vendor, args.username)
                logger.info(f"Switched to user {args.username} for vendor {args.vendor}.")
            else:
                logger.warning(f"User {args.username} not found for vendor {args.vendor}.")
        except Exception as e:
            logger.error(f"Error switching user: {e}")
