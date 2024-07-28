# TCP with custom protocol

Create a TCP server / client which encodes / decodes a custom binary protocol.

The binary protocol can be very simple:
- e.g. 1B version, 2B length, etc

All encoding and decoding must be done by hand.

The server should treat every message individually, regardless of how they come in ie.
multiple messages in a single connection, or the whole message via a single connection.
