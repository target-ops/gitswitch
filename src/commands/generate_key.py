from config.ssh import SSHConfig

from logging_config import get_logger

logger = get_logger(__name__)

class GenerateKeyCommand:
    ssh_config = SSHConfig()
    
    @staticmethod
    def add_arguments(subparsers):
        parser = subparsers.add_parser('generate-key', help='Generate a new SSH key')
        parser.add_argument('email', type=str, help='Email for the SSH key')
        parser.add_argument('key_path', type=str, help='Path to store the SSH key')
        parser.set_defaults(command_class=GenerateKeyCommand)

    def execute(self, args):
        try:
            self.ssh_config.generate_ssh_key(args.email, args.key_path)
            logger.info(f"SSH key generated for {args.email}.")
        except Exception as e:
            logger.error(f"Error generating SSH key: {e}")
