<p align="center">
<img src="https://i.postimg.cc/YSXdd1np/Untitled-2.png">
</p>
<p align="center">
<a href="https://join.slack.com/t/target-ops/shared_invite/zt-2kxdr9djp-YoQSCoRzARa9psxO8aYoaQ"><img src="https://img.shields.io/badge/Slack-4A154B?style=for-the-badge&logo=slack&logoColor=white"></a>
<a href="https://www.linkedin.com/company/target-ops"><img src="https://img.shields.io/badge/LinkedIn-0077B5?style=for-the-badge&logo=linkedin&logoColor=white"></a>
<a href="https://dev.to/target-ops"><img src="https://img.shields.io/badge/dev.to-0A0A0A?style=for-the-badge&logo=devdotto&logoColor=white"></a>
<a href="https://dly.to/13bF3DMZs9K"><img src="https://img.shields.io/badge/daily.dev-CE3DF3?style=for-the-badge&logo=dailydotdev&logoColor=white"></a>
<a href="https://t.me/targetops"><img src="https://img.shields.io/badge/Telegram-2CA5E0?style=for-the-badge&logo=telegram&logoColor=white"></a>
<a href="https://www.patreon.com/target_ops"><img src="https://img.shields.io/badge/Patreon-F96854?style=for-the-badge&logo=patreon&logoColor=white"></a>
<img alt="GitHub Org's stars" src="https://img.shields.io/github/stars/target-ops?style=for-the-badge&logoColor=green&cacheSeconds=3600?style=for-the-badge">
</p>
</br>

# 🚀 GitSwitch

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
