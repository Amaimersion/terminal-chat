# Terminal Chat

## Content

- [Content](#content)
- [Overview](#overview)
- [Download](#download)
- [Commands](#commands)
- [Protocol](#protocol)
- [Room URL](#room-url)
- [Platforms](#platforms)
- [Examples](#examples)

## Overview

It is a multiroom terminal chat. It allows communicating with UTF-8 text between different computers.

User can have multiple rooms at a same time. Room consists of another users. There can be only one active room at a time. When you type a message, all users from active room will receive it. You will see arrived messages from all rooms, not only from active room.

Room organization is independent. It is means that room organization (room name, receivers in that room, messages in that room that are visible to you) will not match on receiver side. Every receiver will have its own room organization.

Sender will be not able to contact receiver until receiver allows accepting of messages from this sender. Every chat side (sender or receiver) may stop receiving of messages from another side without any notifications.

## Download

You can download executable files in "Releases" tab.

## Commands

All interaction with the program is performed using specific commands with possible positional arguments. Every command starts with `/` sign. Everything else, including not known commands, will be treated as UTF-8 messages and will be sended to all receivers in active room. To see available commands, run the program and type `/help`.

## Protocol

Custom protocol named STTP is used for this application. See [protocol documentation](STTP.md) for more.

## Room URL

To start communication between two computers, every computer should know IP of recipient computer. At room start, the program will print possible URLs by which this room can be reached by another computer. However, it works mostly in LAN. For WAN you may need to obtain public IP by your own.

## Platforms

Supported and tested platforms are:
- Linux AMD64
- Windows AMD64

For other platforms you will need to build the program and test it by your own.

## Examples

Here is example of communication between two computers in the same LAN.

**Computer № 1:**

```
Welcome to multiroom chat.
Type "/help" for more information.

Users can reach this room using following URLs.
Pick "outbound" one if available.
URLs:
sttp://192.168.1.235:4444/0 (outbound)
sttp://172.17.0.1:4444/0 (local)

[main]: /user comp_2 sttp://192.168.1.249:4444/1
[main]: hello!
< 12:11 main: hello!
> 12:11 main comp_2: hi!
```

**Computer № 2:**

```
Welcome to multiroom chat.
Type "/help" for more information.

Users can reach this room using following URLs.
Pick "outbound" one if available.
URLs:
sttp://192.168.1.249:4444/0 (outbound)
sttp://172.17.0.1:3333/0 (local)

[main]: /room chat
Users can reach this room using following URLs.
Pick "outbound" one if available.
URLs:
sttp://192.168.1.249:4444/1 (outbound)
sttp://172.17.0.1:3333/1 (local)

[chat]: /user comp_1 sttp://192.168.1.235:4444/0
> 12:11 chat comp_1: hello!
[chat]: hi!
< 12:11 chat: hi!
```
