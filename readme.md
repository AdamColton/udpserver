## UDP Server
This is a very simple UDP server. When it starts, it runs one Go routine to
receive messages and provides helper functions to send packets.

When starting the server, the caller must provide an object that implements the
PacketHandler interface. Received messages will be sent to the PacketHandler.
My [Packeter](https://github.com/AdamColton/packeter) project is one option, but
building something to do this is easy.

It also will find the the local network address of the server if it's behind a
router.