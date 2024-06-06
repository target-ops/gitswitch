from cli import CLI
from logging_config import setup_logging

def main():
    setup_logging()
    CLI().execute()

if __name__ == "__main__":
    main()
