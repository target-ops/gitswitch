import os
import configparser

CONFIG_FILE = os.path.join(os.path.dirname(__file__), '..', 'config.ini')

def load_config():
    config = configparser.ConfigParser()
    if os.path.exists(CONFIG_FILE):
        config.read(CONFIG_FILE)
    if not os.path.exists(CONFIG_FILE):
        # Create the file if it doesn't exist
        with open(CONFIG_FILE, 'w') as file:
            config.write(file)
    # Read the config file
    config.read(CONFIG_FILE)
    return config

def save_config(config):
    with open(CONFIG_FILE, 'w') as configfile:
        config.write(configfile)

def set_current_user(config, vendor, username):
    if 'current' not in config:
        config['current'] = {}
    config['current']['vendor'] = vendor
    config['current']['username'] = username
    save_config(config)

def get_current_user(config):
    if 'current' in config:
        vendor = config['current']['vendor']
        username = config['current']['username']
        return vendor, username
    else:
        return None, None
