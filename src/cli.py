import argparse
import argcomplete
from commands import add,current,delete,generate_key,list,switch
# from commands.add import AddCommand
# from commands.current import CurrentCommand
# from commands.delete import DeleteCommand
# from commands.generate_key import GenerateKeyCommand
# from commands.list import ListCommand
# from commands.switch import SwitchCommand

class CLI:
    addCmd=add.AddCommand()
    currentCmd=current.CurrentCommand()
    deleteCmd=delete.DeleteCommand()
    generateKeyCmd=generate_key.GenerateKeyCommand()
    listCmd=list.ListCommand()
    switchCmd=switch.SwitchCommand()


    def __init__(self):
        self.parser = self.create_parser()
    
    def create_parser(self) -> argparse.ArgumentParser:
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

        self.addCmd.add_arguments(subparsers)
        self.generateKeyCmd.add_arguments(subparsers)
        self.listCmd.add_arguments(subparsers)
        self.switchCmd.add_arguments(subparsers)
        self.deleteCmd.add_arguments(subparsers)
        self.currentCmd.add_arguments(subparsers)

        argcomplete.autocomplete(parser)
        return parser

    def execute(self):
        args = self.parser.parse_args()
        if args.command is None:
            self.parser.print_help()
        else:
            command = args.command_class()
            command.execute(args)
