[![Go Report Card](https://goreportcard.com/badge/github.com/echicken/dpc2?style=flat-square)](https://goreportcard.com/report/github.com/echicken/dpc2)


# dpc2
DoorParty Connector v2

A small server that connects local RLOGIN clients to a remote RLOGIN server via an SSH tunnel.

This is a replacement for the original [doorparty-connector](https://github.com/echicken/doorparty-connector).  It's smaller and in some cases easier to install, but there is no compelling reason to switch if you already have the original one working.

This program was written in [Go](https://golang.org/), and you can build it from source if you want.  [Prebuilt executables](build/) are available for [Linux X86_64](build/linux_x64/), [Linux ARM6](build/linux_arm6/) (for use on a Raspberry Pi, for example), and [Win32](build/win32/).

## Installation

* Create a directory for doorparty-connector on your system, and place the following two files in it:
	* The appropriate executable for your system from one of the [build directories](build/)
	* A copy of [doorparty-connector.ini](configs/doorparty-connector.ini)
* Edit your copy of [doorparty-connector.ini](configs/doorparty-connector.ini)
	* Fill out `system_tag`, `ssh_username`, and `ssh_password` with the values you were given by the DoorParty administrator (you can omit the [square brackets] from your `system_tag`)
	* Edit `local_interface` and `local_port` if the default values don't suit your system
	* Leave the other values at in their default state.  They have been made customizable just in case the remote server's details change in the future.
	
### Linux

For automatic startup, a [systemd unit file](init/doorparty-connector.service) is available.  It assumes that `doorparty-connector` and `doorparty-connector.ini` reside within `/opt/doorparty-connector`, so edit the path as needed.

If you use some other init system, you're on your own - but feel free to share your init script / config file / whatever and I can add it to this repo.

### Windows

I dunno, put a shortcut in your Startup menu or some shit like that.
	
## Usage

`doorparty-connector` runs as a server, by default listening on TCP port `9999` of `localhost` for connections from RLOGIN clients.  You must start it up and leave it running.  A typical installation has `doorparty-connector` running on the same machine that hosts your BBS - but this isn't a requirement.

When a user on your BBS chooses to connect to DoorParty:

* Your BBS should open an RLOGIN connection to `doorparty-connector` (eg. port `9999` on `localhost`)
	* The RLOGIN "server-user-name" is the user's alias prefixed with your DoorParty "system tag", eg. `[ec]echicken`
	* The RLOGIN "client-user-name" is the user's DoorParty password
		* You can use whatever you like for this value, and even use the same value for all users if you must, but you must always send the same password for the same user every time they connect
		* The user does not need to know this value
	* `doorparty-connector` then connects the user to DoorParty's RLOGIN server via an SSH tunnel, passing the necessary user details along
		* An account is automatically created for the user on the DoorParty server if it doesn't already exist
			* This account is linked to this user on your BBS
			* The user can only access it from your BBS
			* The user does not need to know the password
			* The user may have several DoorParty accounts, one for each BBS they connect from; they cannot use the same DoorParty account from multiple BBSs

On Mystic, for example, this is menu command `IR`, with a `DATA` field like:
* `/ADDR=localhost:9999 /USER=[system_tag]@USER@ /PASS=some_password`
	* Mind that `system_tag` and `some_password` must be replaced with your own values

A [script](third_party/synchronet/doorparty.js) and [instructions for Synchronet](third_party/synchronet/) are available [here](third_party/synchronet/).