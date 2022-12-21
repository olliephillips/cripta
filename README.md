# Cripta Messenger

## Peer to peer encrypted & signed messaging via MQTT pub/sub.

[![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/olliephillips/cripta?style=flat-square)](https://go.dev/)
[![Go Report Card](https://goreportcard.com/badge/github.com/olliephillips/cripta?style=flat-square)](https://goreportcard.com/report/github.com/olliephillips/cripta)
[![GitHub Workflow Status](https://img.shields.io/github/actions/workflow/status/olliephillips/cripta/build.yml?branch=main&style=flat-square)](https://github.com/olliephillips/cripta/actions/workflows/build.yml)
[![GitHub Workflow Status (branch)](https://img.shields.io/github/actions/workflow/status/olliephillips/cripta/unit_test.yml?branch=main&label=tests&style=flat-square)](https://github.com/olliephillips/cripta/actions/workflows/unit_test.yml)

[![GitHub Release Date](https://img.shields.io/github/release-date/olliephillips/cripta?style=flat-square)](https://github.com/olliephillips/cripta/releases)
[![GitHub commits since latest release (by date)](https://img.shields.io/github/commits-since/olliephillips/cripta/latest?style=flat-square)](https://github.com/olliephillips/cripta/commits)
[![GitHub](https://img.shields.io/github/license/olliephillips/cripta?label=license&style=flat-square)](LICENSE)

With Cripta you can relay secure messages instantly using free (and mostly available) MQTT brokers. The broker is
configurable, to use a different broker edit, rebuild the app and distribute that.

## Tests

There are no tests in this project. It started as an afternoon hack. I hope to add at some point.

## Features

- Asymmetric encryption using recipient public key
- Inbox management (list, read, delete & empty)
- Mail groups
- Message signing with private key

## Releases

Download an executable for your OS [here](https://github.com/olliephillips/cripta/releases).

## Usage

Create a `config.env` file in same directory as the cripta executable. Add your Twitter username.

```bash
## Enter your twitter username without the @ symbol
TWITTER_USERNAME=olliephillips
```

Note: There is currently no validation of the username against Twitter. You could pick any username.

Run it on Mac/Linux.

```bash
./cripta

Cripta Messenger
---------------------
Version: v0.0.1
Go version: 1.19
---------------------

Enter 'help' for command list.

```

Show command list with `help`.

```bash
> help
Showing available commands..

Command List
------------
help                    Show available commands
list                    List all messages in mailbox
friends                 List available friends (with public keys)
groups                  List available groups
read <msgId>            Print the message with id <msgId> to the console
delete <msgId>          Delete the messsage with id <msgId> from mailbox
empty                   Delete all messages in the mailbox
quit                    Quit the program

To send a message to a single recipient:-
> @<username> <subject>::<message>

To send to a group of recipients (group must exist):-
> -><groupname> <subject>::<message>

The message(s) will be queue and published.
``` 

### Receiving messages

To receive a message from someone else you must have shared your public key in `my_public_key.txt` with the sender.

### Sending messages

You can send messages to any user who's public key you have saved in your friends folder. By default the folder
is `friend_keys`. A user's key is saved in a `.txt` file named with their username e.g. `username.txt`.

You can send test messages to yourself. On first start your public key is created and copied to the `friend_keys`
folder.

Use this format to send to a user.

```bash
> @<username> <subject>::<message>
```

### Creating groups

You can also send to groups. To create a group, create a named `.txt` file in the `groups` folder (default),
e.g. `mygroup.txt`. Add usernames (for whom you hold a public key) one per line, including the `@` symbol.

```bash
@user1
@user2
```

Use this format to send to a group:

```bash
> -><groupname> <subject>::<message>
```

## License

MIT