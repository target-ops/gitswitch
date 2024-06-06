import os
import configparser

CONFIG_FILE = os.path.join(os.path.dirname(__file__), 'config.ini')

class Config:
    @staticmethod
    def load():
        """Load config file. If the file does not exist, create it."""
        config = configparser.ConfigParser()
        if os.path.exists(CONFIG_FILE):
            config.read(CONFIG_FILE)
        else:
            # Create the file if it doesn't exist
            with open(CONFIG_FILE, 'w') as file:
                config.write(file)
        return config

    @staticmethod
    def save(config):
        """Save the configuration to the config file."""
        with open(CONFIG_FILE, 'w') as configfile:
            config.write(configfile)

    @staticmethod
    def set_current_user(config, vendor, username):
        """Set the current active user."""
        if 'current' not in config:
            config['current'] = {}
        config['current']['vendor'] = vendor
        config['current']['username'] = username
        Config.save(config)

    @staticmethod
    def get_current_user(config):
        """Get the current active user."""
        if 'current' in config:
            vendor = config['current']['vendor']
            username = config['current']['username']
            return vendor, username
        else:
            return None, None
