# dpc2
DoorParty Connector v2

A small server that connects local RLogin clients to a remote RLogin server via an SSH tunnel.

This is a replacement for the original [doorparty-connector](https://github.com/echicken/doorparty-connector).  It's smaller and in some cases easier to install, but there is no compelling reason to switch if you already have the original one working.

This program was written in [Go](https://golang.org/), and you can build it from source if you want.  [Prebuilt executables](build/) are available for [Linux X86_64](build/linux_x64/), [Linux ARM6](build/linux_arm6/) (for use on a Raspberry Pi, for example), and [Win32](build/win32/).

## Installation

* Create a directory for doorparty-connector on your system, and place the following two files in it:
	* The appropriate executable for your system from one of the [build directories](build/)
	* A copy of [doorparty-connector.ini](doorparty-connector.ini)
* Edit your copy of [doorparty-connector.ini](doorparty-connector.ini)
	* Fill out `system_tag`, `ssh_username`, and `ssh_password` with the values you were given by the DoorParty administrator (you can omit the [square brackets] from your `system_tag`)
	* Edit `local_interface` and `local_port` if the default values don't suit your system
	* Leave the other values at in their default state.  They have been made customizable just in case the remote server's details change in the future.
	
### Linux

A [systemd unit file](doorparty-connector.service) is available.  It assumes that `doorparty-connector` and `doorparty-connector.ini` reside within `/opt/doorparty-connector`, so edit the unit file if you placed this program elsewhere.
	
## Usage

`doorparty-connector` runs as a server, by default listening on TCP port `9999` of `localhost` for connections from RLogin clients.  You must start it up and leave it running.  A typical installation has `doorparty-connector` running on the same machine that hosts your BBS - but this isn't a requirement.

When a user on your BBS chooses to connect to DoorParty, your BBS should open an RLogin connection to `doorparty-connector`, likely running on the same machine. `doorparty-connector` then connects the user to DoorParty's RLogin server via an SSH tunnel.
