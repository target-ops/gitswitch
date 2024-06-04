# 🚀 Target-Ops: GitSwitch

Welcome to **GitSwitch**, the ultimate solution for managing multiple Git users across different vendors like GitHub and GitLab. Whether you're a developer juggling multiple identities or a team lead needing to streamline SSH key management, Target-Ops has got you covered.

## 🌟 Features
- **Add Users**: Easily add new Git users for various vendors.
- **Switch Users**: Quickly switch between different Git user profiles.
- **List Users**: View all configured Git users.
- **Delete Users**: Remove Git users when they're no longer needed.
- **Current User**: Display the currently active user.
- **Generate SSH Keys**: Create and manage SSH keys seamlessly.
- **Upload SSH Key**: Upload the SSH key.pub to your github/gitlab account

## 📦 Installation
Install the required dependencies:
```
pip install -r requirements.txt
```
## 🚀 Usage
GitSwitch provides a command-line interface for managing Git users. Here are the available commands:
#### Add User
Add a new Git user for a specific vendor:
```python main.py add <vendor> <username> <email> <key_path> --upload-key```
#### Generate SSH Key
Generate a new SSH key:
```python main.py generate-key <email> <key_path>```
#### List Users
List all configured Git users:
```python main.py list```
#### Switch User
Switch to a different Git user:
```python main.py switch <vendor> <username>```
#### Delete User
Delete a configured Git user:
``` python main.py delete <vendor> <username>```
#### Current User
Display the currently active Git user:
```python main.py current```

## ⚙️ Configuration
Configuration is handled through a configuration file, typically located in your home directory. This file keeps track of all users and their associated details.

## 💡 Contributing
We welcome contributions from the community! Feel free to fork the repository and submit pull requests.
```
Fork the repository.
Create a new branch: git checkout -b my-feature-branch.
Make your changes and commit them: git commit -m 'Add some feature'.
Push to the branch: git push origin my-feature-branch.
Submit a pull request.
```