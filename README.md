# 🚀 Target-Ops: GitSwitch

Welcome to **GitSwitch**, the ultimate solution for managing multiple Git users across different vendors like GitHub and GitLab. Whether you're a developer juggling multiple identities or a team lead needing to streamline SSH key management, Target-Ops has got you covered.

## 🌟 Features
- **User Management**: Add, switch, list, and delete Git users across various vendors with ease.
- **SSH Key Management**: Generate and manage SSH keys seamlessly.
- **SSH Key Upload**: Upload the SSH key.pub to your GitHub/GitLab account directly.
- **Active User Display**: Easily view the currently active user.

## 📦 Installation
GitSwitch can be installed using Homebrew with the following commands:

```sh
brew tap target-ops/homebrew-tap 
brew install target-ops/tap/gitswitch
```

## 🚀 Usage
GitSwitch provides a command-line interface for managing Git users. Here are the available commands:
#### Add User
```
gitswitch add user --vendor <vendor> --username <username> --email <email> --pub_key_path <path_to_public_key>
```

#### Generate SSH Key
```
gitswitch generate key --email <email> --pub_key_path <path_to_public_key>
```
#### List Users
```
gitswitch list
```
#### Switch User
```
gitswitch switch --vendor <vendor> --username <username>
```
#### Delete User
```
gitswitch delete --vendor <vendor> --username <username>
```
#### Current User
```
gitswitch current
```

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