# TCP Server

## Spec
Client / server library that listens on a given port, and allows a user to send
arbitary messages to it.

The server will log the message with any of the following modifications:
- Only including alphanumeric characters
- Remove all alphanumeric characters

The action should be defined when the client is called.

## Requirements

- User defined port
- One call client API
- Messages are sent using custom binary encoding
