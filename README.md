# Auto-Approve-Bot

Telegram bot to auto approve chat join requests.

## Prerequisites

- **Go** version 1.23.0 or later
- **MongoDB** Database Url
- **Telegram Bot Token** from [BotFather](https://t.me/botfather)

## Installation

### 1. Install Go

Follow the instructions to install Go on your system: [Go Installation Guide](https://golang.org/doc/install)

<details>
<summary>Easy Way:</summary>

```shell
git clone https://github.com/udhos/update-golang dlgo && cd dlgo && sudo ./update-golang.sh && source /etc/profile.d/golang_path.sh
```

Exit the terminal and open the terminal to check the installation.
</details>

Verify the installation by running:

```shell
go version
```

### 2. Clone the repository

```shell
git clone https://github.com/AshokShau/Auto-Approve-Bot.git ApproveBot && cd ApproveBo
```

### 3. Set up the environment

Copy the sample environment file and edit it as needed:

```shell
cp sample.env .env
vi .env
```

### 4. Build the project

```shell
go build 
```

## Start the bot

```shell
sudo ./Auto-Approve-Bot
```

## Deploy on Vercel

* Deploy the bot on Vercel with the following steps:
* Fork this repository 🍴
* Login your [Vercel](https://vercel.com/) account
* Go to your [Add New Project](https://vercel.com/new)
* Choose the repository you forked
* Configure the environment variables: `DB_URI` [MongoDB](https://www.mongodb.com/)
* Tap on Deploy
* After deployment, visit the deployed URL & Connect your bot with the deployed URL.

* Start the bot by sending `/start` command to the bot.
* Congratulations 🎉 Enjoy the bot 🌟 if you have any questions, join the [support Channel](https://t.me/FallenProjects)
  🤗

## Usage

1. **Start the bot**: Start the bot by sending `/start` command to the bot.
2. **Add the bot to your group**: Add the bot to your group and make it an admin with permission to approve new members.

> **Note**: The bot must be an admin with <u>Invite Users</u> permission to approve new members to auto-approve new
> members.

### Commands

- `/start`: Start the bot.
- `/ping`: Check if the bot is running.
- `/autoApprove`: Enable/Disable auto-approve mode (Admin Only).
- `/stats`: Get the bot's statistics (Bot Owner Only).
- `/broadcast`: Broadcast a message to all users (Bot Owner Only).

## Contributing

<details>
<summary>Contributing Guidelines</summary>

Contributions are welcome! Follow these steps to contribute:

1. **Fork the repository**: Click the "Fork" button at the top right of this page to create a copy of this repository in
   your GitHub account.

2. **Clone the repository**: Clone your forked repository to your local machine.
    ```shell
    git clone https://github.com/your-username/Auto-Approve-Bot.git
    cd Auto-Approve-Bot
    ```

3. **Create a branch**: Create a new branch for your changes.
    ```shell
    git checkout -b feature-branch
    ```

4. **Make your changes**: Make your changes to the codebase.

5. **Commit your changes**: Commit your changes with a descriptive commit message.
    ```shell
    git add .
    git commit -m "Description of your changes"
    ```

6. **Push to your branch**: Push your changes to your forked repository.
    ```shell
    git push origin feature-branch
    ```

7. **Submit a pull request**: Go to the original repository on GitHub and create a pull request from your forked
   repository.

Please ensure your code follows the project's coding standards and includes appropriate tests.

Thank you for contributing!
</details>

## License

This project is licensed under the MIT License—see the [LICENSE](LICENSE) file for details.
